package questions

import "errors"

var ErrQuestionNotAnswered = errors.New("question.Answer() called without the question having been successfully asked")

type Questioner interface {
	NewOpenEndedQuestion(question *OpenEndedQuestion) Question
	NewClosedQuestion(question *ClosedQuestion) Question
}

// AskQuestions takes a list of questions and asks each one in sequence, returning the list of answers
func AskQuestions(questions []Question) error {
	for _, question := range questions {
		err := question.Ask()
		if err != nil {
			return err
		}
	}
	return nil
}

type Answer string

// Question represents a generic user prompt
type Question interface {
	Ask() error
	Answer() Answer
}

// OpenEndedQuestion is a prompt where the user can provide any answer
type OpenEndedQuestion struct {
	Question  string
	Default   string
	AllowEdit bool
	Validate  func(string) error
}

// ClosedQuestion is a prompt where the user can choose from a fixed set of answers
type ClosedQuestion struct {
	Question string
	Items    []string
}

// QuestionBranch is a sort of "meta-question" whose Ask method will delegate to ConditionQuestion and execute BranchA
// if ConditionAnswer is met, otherwise execute BranchB
type QuestionBranch struct {
	ConditionQuestion Question
	ConditionAnswer   Answer
	BranchA           []Question
	BranchB           []Question
}

func (q *QuestionBranch) Ask() error {
	err := q.ConditionQuestion.Ask()
	if err != nil {
		return err
	}

	answer := q.ConditionQuestion.Answer()

	if answer == q.ConditionAnswer {
		err := AskQuestions(q.BranchA)
		if err != nil {
			return err
		}
	} else {
		err := AskQuestions(q.BranchB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *QuestionBranch) Answer() Answer {
	return q.ConditionQuestion.Answer()
}
