package forge

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	docker "github.com/docker/docker/client"

	"github.com/sclevine/forge/engine"
	"github.com/sclevine/forge/internal"
)

const stagerScript = `
	set -e
	{{- range .BuildpackMD5s}}
	su root -c "unzip -qq /tmp/{{.}}.zip -d /tmp/buildpacks/{{.}}" && rm /tmp/{{.}}.zip
	{{- end}}

	{{if not .RSync}}exec {{end}}su root -p -c "PATH=$PATH exec /lifecycle/builder -outputDroplet /droplet -buildpackOrder '$0' -skipDetect=$1"
	{{- if .RSync}}
	rsync -a /tmp/app/ /tmp/local/
	{{- end}}
`

type Stager struct {
	ImageTag         string
	SystemBuildpacks Buildpacks
	Logs             io.Writer
	Loader           Loader
	versioner        versioner
	engine           forgeEngine
	image            forgeImage
}

type StageConfig struct {
	AppTar        io.Reader
	Cache         ReadResetWriter
	CacheEmpty    bool
	BuildpackZips map[string]engine.Stream
	Stack         string
	AppDir        string
	ForceDetect   bool
	RSync         bool
	Color         Colorizer
	AppConfig     *AppConfig
}

type ReadResetWriter interface {
	io.ReadWriter
	Reset() error
}

func NewStager(client *docker.Client, httpClient *http.Client, exit <-chan struct{}) *Stager {
	return &Stager{
		ImageTag: "forge",
		Logs:     os.Stdout,
		Loader:   noopLoader{},
		versioner: &internal.Version{
			Client: httpClient,
		},
		engine: &dockerEngine{
			Docker: client,
			Exit:   exit,
		},
		image: &engine.Image{
			Docker: client,
			Exit:   exit,
		},
	}
}

func (s *Stager) Stage(config *StageConfig) (droplet engine.Stream, err error) {
	if err := s.buildDockerfile(config.Stack); err != nil {
		return engine.Stream{}, err
	}

	var buildpackMD5s []string
	for checksum := range config.BuildpackZips {
		buildpackMD5s = append(buildpackMD5s, checksum)
	}
	sort.Strings(buildpackMD5s)
	containerConfig, err := s.buildContainerConfig(config.AppConfig, buildpackMD5s, config.ForceDetect, config.RSync)
	if err != nil {
		return engine.Stream{}, err
	}
	remoteDir := "/tmp/app"
	if config.RSync {
		remoteDir = "/tmp/local"
	}
	hostConfig := s.buildHostConfig(config.AppDir, remoteDir)
	contr, err := s.engine.NewContainer(config.AppConfig.Name+"-staging", containerConfig, hostConfig)
	if err != nil {
		return engine.Stream{}, err
	}
	defer contr.CloseAfterStream(&droplet)
	for checksum, zip := range config.BuildpackZips {
		if err := contr.StreamFileTo(zip, fmt.Sprintf("/tmp/%s.zip", checksum)); err != nil {
			return engine.Stream{}, err
		}
	}

	if err := contr.ExtractTo(config.AppTar, "/tmp/app"); err != nil {
		return engine.Stream{}, err
	}
	if !config.CacheEmpty {
		if err := contr.ExtractTo(config.Cache, "/tmp/cache"); err != nil {
			return engine.Stream{}, err
		}
	}

	status, err := contr.Start(config.Color("[%s] ", config.AppConfig.Name), s.Logs, nil)
	if err != nil {
		return engine.Stream{}, err
	}
	if status != 0 {
		return engine.Stream{}, fmt.Errorf("container exited with status %d", status)
	}

	if err := config.Cache.Reset(); err != nil {
		return engine.Stream{}, err
	}
	if err := streamOut(contr, config.Cache, "/tmp/output-cache"); err != nil {
		return engine.Stream{}, err
	}

	return contr.StreamFileFrom("/droplet")
}

func (s *Stager) buildContainerConfig(config *AppConfig, buildpackMD5s []string, forceDetect, rsync bool) (*container.Config, error) {
	var (
		buildpacks []string
		detect     bool
	)
	if config.Buildpack == "" && len(config.Buildpacks) == 0 {
		buildpacks = s.SystemBuildpacks.names()
		detect = true
	} else if len(config.Buildpacks) > 0 {
		buildpacks = config.Buildpacks
	} else {
		buildpacks = []string{config.Buildpack}
	}
	detect = detect || forceDetect
	if detect {
		fmt.Fprintln(s.Logs, "Buildpack: will detect")
	} else {
		var plurality string
		if len(buildpacks) > 1 {
			plurality = "s"
		}
		fmt.Fprintf(s.Logs, "Buildpack%s: %s\n", plurality, strings.Join(buildpacks, ", "))
	}

	// TODO: fill with real information -- get/set container limits
	vcapApp, err := json.Marshal(&vcapApplication{
		ApplicationID:      "01d31c12-d066-495e-aca2-8d3403165360",
		ApplicationName:    config.Name,
		ApplicationURIs:    []string{"localhost"},
		ApplicationVersion: "2b860df9-a0a1-474c-b02f-5985f53ea0bb",
		Limits:             map[string]int64{"fds": 16384, "mem": 1024, "disk": 4096},
		Name:               config.Name,
		SpaceID:            "18300c1c-1aa4-4ae7-81e6-ae59c6cdbaf1",
		SpaceName:          config.Name + "-space",
		URIs:               []string{"localhost"},
		Version:            "18300c1c-1aa4-4ae7-81e6-ae59c6cdbaf1",
	})
	if err != nil {
		return nil, err
	}

	services := config.Services
	if services == nil {
		services = Services{}
	}
	vcapServices, err := json.Marshal(services)
	if err != nil {
		return nil, err
	}

	env := map[string]string{
		"CF_INSTANCE_ADDR":  "",
		"CF_INSTANCE_IP":    "0.0.0.0",
		"CF_INSTANCE_PORT":  "",
		"CF_INSTANCE_PORTS": "[]",
		"CF_STACK":          "cflinuxfs2",
		"HOME":              "/root",
		"LANG":              "en_US.UTF-8",
		"MEMORY_LIMIT":      "1024m",
		"PATH":              "/usr/local/bin:/usr/bin:/bin",
		"USER":              "root",
		"STACK":             "heroku-16",
		"VCAP_APPLICATION":  string(vcapApp),
		"VCAP_SERVICES":     string(vcapServices),
	}

	scriptBuf := &bytes.Buffer{}
	tmpl := template.Must(template.New("").Parse(stagerScript))
	if err := tmpl.Execute(scriptBuf, struct {
		RSync         bool
		BuildpackMD5s []string
	}{
		rsync,
		buildpackMD5s,
	}); err != nil {
		return nil, err
	}

	return &container.Config{
		Hostname:   config.Name,
		User:       "root",
		Env:        mapToEnv(mergeMaps(env, config.StagingEnv, config.Env)),
		Image:      s.ImageTag,
		WorkingDir: "/root",
		Entrypoint: strslice.StrSlice{
			"/bin/bash", "-c", scriptBuf.String(),
			strings.Join(buildpacks, ","),
			strconv.FormatBool(!detect),
		},
	}, nil
}

func (*Stager) buildHostConfig(appDir, remoteDir string) *container.HostConfig {
	if appDir == "" || remoteDir == "" {
		return nil
	}
	return &container.HostConfig{Binds: []string{appDir + ":" + remoteDir}}
}

func streamOut(contr Container, out io.Writer, path string) error {
	stream, err := contr.StreamFileFrom(path)
	if err != nil {
		return err
	}
	return stream.Out(out)
}

func (s *Stager) Download(path, stack string) (stream engine.Stream, err error) {
	contr, err := s.container(stack)
	if err != nil {
		return engine.Stream{}, err
	}
	defer contr.CloseAfterStream(&stream)
	return contr.StreamFileFrom(path)
}

func (s *Stager) DownloadTar(path, stack string) (stream engine.Stream, err error) {
	contr, err := s.container(stack)
	if err != nil {
		return engine.Stream{}, err
	}
	defer contr.CloseAfterStream(&stream)
	return contr.StreamTarFrom(path)
}

func (s *Stager) container(stack string) (Container, error) {
	if err := s.buildDockerfile(stack); err != nil {
		return nil, err
	}
	containerConfig := &container.Config{
		Hostname:   "download",
		User:       "root",
		Image:      s.ImageTag,
		Entrypoint: strslice.StrSlice{"read"},
	}
	contr, err := s.engine.NewContainer("download", containerConfig, nil)
	if err != nil {
		return nil, err
	}
	return contr, nil
}

func (s *Stager) buildDockerfile(stack string) error {
	buildpacks, err := s.buildpacks()
	if err == internal.ErrNetwork || err == internal.ErrUnavailable {
		fmt.Fprintln(s.Logs, "Warning: cannot build image: ", err)
		return nil
	}
	if err != nil {
		return err
	}
	dockerfileBuf := &bytes.Buffer{}
	dockerfileData := struct {
		Stack      string
		Buildpacks []buildpackInfo
	}{
		stack,
		buildpacks,
	}
	dockerfileTmpl := template.Must(template.New("Dockerfile").Parse(dockerfile))
	if err := dockerfileTmpl.Execute(dockerfileBuf, dockerfileData); err != nil {
		return err
	}
	dockerfileStream := engine.NewStream(ioutil.NopCloser(dockerfileBuf), int64(dockerfileBuf.Len()))
	return s.Loader.Loading("Image", s.image.Build(s.ImageTag, dockerfileStream))
}

func (s *Stager) buildpacks() ([]buildpackInfo, error) {
	var buildpacks []buildpackInfo
	for _, buildpack := range s.SystemBuildpacks {
		url, err := s.versioner.Build(buildpack.URL, buildpack.VersionURL)
		if err != nil {
			return nil, err
		}
		checksum := fmt.Sprintf("%x", md5.Sum([]byte(buildpack.Name)))
		info := buildpackInfo{buildpack.Name, url, checksum}
		buildpacks = append(buildpacks, info)
	}
	return buildpacks, nil
}

type buildpackInfo struct {
	Name, URL, MD5 string
}
