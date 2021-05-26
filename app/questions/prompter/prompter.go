package prompter

import (
	"github.com/manifoldco/promptui"
	"github.com/opencredo/venafi-vault-wizard/app/questions"
)

type prompter struct{}

// NewPrompter returns an instance of questions.Questioner that uses github.com/manifoldco/promptui to implement its
// questions
func NewPrompter() questions.Questioner {
	return &prompter{}
}

func (p *prompter) NewOpenEndedQuestion(question *questions.OpenEndedQuestion) questions.Question {
	return &openEndedQuestion{
		Question:  question.Question,
		Default:   question.Default,
		AllowEdit: question.AllowEdit,
		Validate:  question.Validate,
	}
}

func (p *prompter) NewClosedQuestion(question *questions.ClosedQuestion) questions.Question {
	return &closedQuestion{
		Question: question.Question,
		Items:    question.Items,
	}
}

// openEndedQuestion is a prompt where the user can provide any answer
type openEndedQuestion struct {
	Question  string
	Default   string
	AllowEdit bool
	Validate  promptui.ValidateFunc

	Result *questions.Answer
}

func (q *openEndedQuestion) Ask() error {
	prompt := promptui.Prompt{
		Label:       q.Question,
		Default:     q.Default,
		AllowEdit:   q.AllowEdit,
		HideEntered: true,
		Validate:    q.Validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return err
	}

	answer := questions.Answer(result)
	q.Result = &answer

	return nil
}

func (q *openEndedQuestion) Answer() questions.Answer {
	if q.Result == nil {
		panic(questions.ErrQuestionNotAnswered)
	}

	return *q.Result
}

// closedQuestion is a prompt where the user can choose from a fixed set of answers
type closedQuestion struct {
	Question string
	Items    []string

	Result *questions.Answer
}

func (q *closedQuestion) Ask() error {
	prompt := promptui.Select{
		Label:        q.Question,
		Items:        q.Items,
		HideSelected: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	answer := questions.Answer(result)
	q.Result = &answer

	return nil
}

func (q *closedQuestion) Answer() questions.Answer {
	if q.Result == nil {
		panic(questions.ErrQuestionNotAnswered)
	}

	return *q.Result
}
