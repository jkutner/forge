package cmd_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/sclevine/cflocal/cf/cmd"
	"github.com/sclevine/cflocal/cf/cmd/mocks"
	"github.com/sclevine/cflocal/local"
	sharedmocks "github.com/sclevine/cflocal/mocks"
	"github.com/sclevine/cflocal/remote"
)

var _ = Describe("Pull", func() {
	var (
		mockCtrl   *gomock.Controller
		mockUI     *sharedmocks.MockUI
		mockApp    *mocks.MockApp
		mockFS     *mocks.MockFS
		mockHelp   *mocks.MockHelp
		mockConfig *mocks.MockConfig
		cmd        *Pull
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUI = sharedmocks.NewMockUI()
		mockApp = mocks.NewMockApp(mockCtrl)
		mockFS = mocks.NewMockFS(mockCtrl)
		mockHelp = mocks.NewMockHelp(mockCtrl)
		mockConfig = mocks.NewMockConfig(mockCtrl)
		cmd = &Pull{
			UI:     mockUI,
			App:    mockApp,
			FS:     mockFS,
			Help:   mockHelp,
			Config: mockConfig,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("#Match", func() {
		It("should return true when the first argument is pull", func() {
			Expect(cmd.Match([]string{"pull"})).To(BeTrue())
			Expect(cmd.Match([]string{"not-pull"})).To(BeFalse())
			Expect(cmd.Match([]string{})).To(BeFalse())
			Expect(cmd.Match(nil)).To(BeFalse())
		})
	})

	Describe("#Run", func() {
		It("should download a droplet and save its env vars", func() {
			droplet := mocks.NewMockBuffer("some-droplet")
			file := mocks.NewMockBuffer("")
			env := &remote.AppEnv{
				Staging: map[string]string{"a": "b"},
				Running: map[string]string{"c": "d"},
				App:     map[string]string{"e": "f"},
			}
			oldLocalYML := &local.LocalYML{
				Applications: []*local.AppConfig{
					{Name: "some-other-app"},
					{
						Name:       "some-app",
						Command:    "some-old-command",
						StagingEnv: map[string]string{"g": "h"},
						RunningEnv: map[string]string{"i": "j"},
						Env:        map[string]string{"k": "l"},
					},
				},
			}
			newLocalYML := &local.LocalYML{
				Applications: []*local.AppConfig{
					{Name: "some-other-app"},
					{
						Name:       "some-app",
						Command:    "some-command",
						StagingEnv: map[string]string{"a": "b"},
						RunningEnv: map[string]string{"c": "d"},
						Env:        map[string]string{"e": "f"},
					},
				},
			}
			mockApp.EXPECT().Droplet("some-app").Return(droplet, int64(100), nil)
			mockFS.EXPECT().WriteFile("./some-app.droplet").Return(file, nil)
			mockConfig.EXPECT().Load().Return(oldLocalYML, nil)
			mockApp.EXPECT().Env("some-app").Return(env, nil)
			mockApp.EXPECT().Command("some-app").Return("some-command", nil)
			mockConfig.EXPECT().Save(newLocalYML)

			Expect(cmd.Run([]string{"pull", "some-app"})).To(Succeed())
			Expect(file.Result()).To(Equal("some-droplet"))
			Expect(droplet.Result()).To(BeEmpty())
			Expect(mockUI.Out).To(gbytes.Say("Successfully downloaded: some-app"))
		})

		// TODO: test when app isn't in local.yml
	})
})
