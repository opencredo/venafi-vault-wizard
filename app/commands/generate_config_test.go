package commands_test

import (
	"embed"
	"encoding/csv"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/opencredo/venafi-vault-wizard/app/commands"
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
		questionsCSVFilename string
		expectedConfig       *config.Config
	}{
		"one VM pki-backend": {
			questionsCSVFilename: "test_fixtures/one_vm_pki-backend.csv",
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
		"multi VM pki-backend": {
			questionsCSVFilename: "test_fixtures/multi_vm_pki-backend.csv",
			expectedConfig: &config.Config{
				Vault:   config.VaultConfig{
					VaultAddress: "http://localhost:8200",
					VaultToken:   "root",
					SSHConfig: []config.SSH{
						{
							Hostname: "localhost",
							Username: "vagrant",
							Password: "vagrant",
							Port:     22,
						},
						{
							Hostname: "localhost2",
							Username: "vagrant",
							Password: "vagrant",
							Port:     23,
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
		"container pki-backend": {
			questionsCSVFilename: "test_fixtures/container_pki-backend.csv",
			expectedConfig: &config.Config{
				Vault: config.VaultConfig{
					VaultAddress: "http://localhost:8200",
					VaultToken:   "root",
					SSHConfig:    nil,
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

			err := expectQuestionsInCSV(tc.questionsCSVFilename, questioner)
			assert.NoError(t, err)

			expectUnansweredQuestions(questioner)

			commands.GenerateConfig(testFileName, questioner)

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
		).Once().Return(mockQuestion)
	case ClosedQuestion:
		questioner.On(
			"NewClosedQuestion",
			mock.MatchedBy(func(q *questions.ClosedQuestion) bool { return q.Question == question }),
		).Once().Return(mockQuestion)
	}
}

func expectUnansweredQuestions(questioner *mocks.Questioner) {
	mockQuestion := new(mocks.Question)
	questioner.On("NewOpenEndedQuestion", mock.Anything).Maybe().Return(mockQuestion)
	questioner.On("NewClosedQuestion", mock.Anything).Maybe().Return(mockQuestion)
}

//go:embed test_fixtures
var questionsFiles embed.FS

func expectQuestionsInCSV(questionsFilename string, questioner *mocks.Questioner) error {
	file, err := questionsFiles.Open(questionsFilename)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	questionRows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, question := range questionRows {
		question, answer, questionTypeString := question[0], question[1], question[2]
		var questionType QuestionType
		switch questionTypeString {
		case "OpenEndedQuestion":
			questionType = OpenEndedQuestion
		case "ClosedQuestion":
			questionType = ClosedQuestion
		default:
			return fmt.Errorf("unexpected question type: %s", questionTypeString)
		}
		expectQuestion(questioner, question, answer, questionType)
	}

	return nil
}
