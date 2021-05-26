package commands

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/opencredo/venafi-vault-wizard/app/config"
	"github.com/opencredo/venafi-vault-wizard/app/plugins"
	"github.com/opencredo/venafi-vault-wizard/app/plugins/venafi"
	pki_backend "github.com/opencredo/venafi-vault-wizard/app/plugins/venafi/pki-backend"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
	mocks "github.com/opencredo/venafi-vault-wizard/mocks/app/questions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateConfig(t *testing.T) {
	testCases := map[string]struct {
		questions []struct {
			question     string
			answer       string
			questionType QuestionType
		}
		expectedConfig *config.Config
	}{
		"normal VM pki-backend": {
			questions: []struct {
				question     string
				answer       string
				questionType QuestionType
			}{
				{
					question:     "What is Vault's API address?",
					answer:       "http://localhost:8200",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What token should be used to authenticate with Vault?",
					answer:       "root",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "Is Vault running in a VM or a container",
					answer:       "VM",
					questionType: ClosedQuestion,
				},
				{
					question:     "Do you have SSH access to the Vault server(s)",
					answer:       "Yes",
					questionType: ClosedQuestion,
				},
				{
					question:     "Is Vault running in High-Availability (HA) mode",
					answer:       "No, just one node",
					questionType: ClosedQuestion,
				},
				{
					question:     "What is the hostname of the Vault server?",
					answer:       "localhost",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What is the SSH username to log into the Vault server?",
					answer:       "vagrant",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What is the SSH password to log into the Vault server?",
					answer:       "vagrant",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What is the SSH port for logging into the Vault server?",
					answer:       "22",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "Which plugin would you like to configure",
					answer:       "venafi-pki-backend",
					questionType: ClosedQuestion,
				},
				{
					question:     "Which version of the plugin would you like to use?",
					answer:       "v0.9.0",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "Which Vault path should the plugin be mounted at?",
					answer:       "pki",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What should the role be called?",
					answer:       "web",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What type of Venafi instance will be used?",
					answer:       "Venafi-as-a-Service",
					questionType: ClosedQuestion,
				},
				{
					question:     "What is the Venafi-as-a-Service API Key?",
					answer:       "venafiAPIKey",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "What project zone should be used for issuing certificates?",
					answer:       "projectzoneID",
					questionType: OpenEndedQuestion,
				},
				{
					question:     "You have configured 1 roles, are there more",
					answer:       "No that's it",
					questionType: ClosedQuestion,
				},
				{
					question:     "You have configured 1 plugins, are there more",
					answer:       "No that's it",
					questionType: ClosedQuestion,
				},
			},
			expectedConfig: &config.Config{
				Vault: config.VaultConfig{
					VaultAddress: "http://localhost:8200",
					VaultToken:   "root",
					SSHConfig: []config.SSH{
						{
							Hostname: "localhost",
							Username: "vagrant",
							Password: "vagrant",
							Port:     22,
						},
					},
				},
				Plugins: []plugins.Plugin{
					{
						Type:      "venafi-pki-backend",
						Version:   "v0.9.0",
						MountPath: "pki",
						Impl: &pki_backend.VenafiPKIBackendConfig{
							MountPath: "pki",
							Version:   "v0.9.0",
							Roles: []pki_backend.Role{
								{
									Name: "web",
									Secret: venafi.VenafiSecret{
										Name: "vaas",
										Cloud: &venafi.VenafiCloudConnection{
											APIKey: "venafiAPIKey",
											Zone:   "projectzoneID",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testFileName := fmt.Sprintf("TestGenerateConfig_vvw_%d.hcl", time.Now().Unix())
			defer os.Remove(testFileName)

			questioner := new(mocks.Questioner)
			defer questioner.AssertExpectations(t)
			for _, question := range tc.questions {
				expectQuestion(questioner, question.question, question.answer, question.questionType)
			}
			expectUnansweredQuestions(questioner)

			GenerateConfig(testFileName, questioner)

			actualConfig, err := config.NewConfigFromFile(testFileName)
			assert.NoError(t, err)

			diff := cmp.Diff(
				tc.expectedConfig,
				actualConfig,
				cmp.FilterPath(func(path cmp.Path) bool {
					return path.String() == "Plugins.Config"
				}, cmp.Ignore()),
			)
			if diff != "" {
				t.Errorf("GenerateConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type QuestionType int

const (
	OpenEndedQuestion = iota
	ClosedQuestion
)

func expectQuestion(questioner *mocks.Questioner, question, answer string, questionType QuestionType) {
	mockQuestion := new(mocks.Question)
	mockQuestion.On("Ask").Return(nil)
	mockQuestion.On("Answer").Return(questions.Answer(answer))
	switch questionType {
	case OpenEndedQuestion:
		questioner.On(
			"NewOpenEndedQuestion",
			mock.MatchedBy(func(q *questions.OpenEndedQuestion) bool { return q.Question == question }),
		).Return(mockQuestion)
	case ClosedQuestion:
		questioner.On(
			"NewClosedQuestion",
			mock.MatchedBy(func(q *questions.ClosedQuestion) bool { return q.Question == question }),
		).Return(mockQuestion)
	}
}

func expectUnansweredQuestions(questioner *mocks.Questioner) {
	mockQuestion := new(mocks.Question)
	questioner.On("NewOpenEndedQuestion", mock.Anything).Maybe().Return(mockQuestion)
	questioner.On("NewClosedQuestion", mock.Anything).Maybe().Return(mockQuestion)
}
