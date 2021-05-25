package questions

import "github.com/manifoldco/promptui"

// AskQuestions takes a list of questions and asks each one in sequence, returning the list of answers
func AskQuestions(questions []Question) ([]Answer, error) {
	var answers []Answer
	for _, question := range questions {
		answer, err := question.Ask()
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer...)
	}
	return answers, nil
}

// Question represents a generic user prompt
type Question interface {
	Ask() ([]Answer, error)
}

// OpenEndedQuestion is a prompt where the user can provide any answer
type OpenEndedQuestion struct {
	Question  string
	Default   string
	AllowEdit bool
	Validate  promptui.ValidateFunc
}

func (q *OpenEndedQuestion) Ask() ([]Answer, error) {
	prompt := promptui.Prompt{
		Label:       q.Question,
		Default:     q.Default,
		AllowEdit:   q.AllowEdit,
		HideEntered: true,
		Validate:    q.Validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return []Answer{Answer(result)}, nil
}

// ClosedQuestion is a prompt where the user can choose from a fixed set of answers
type ClosedQuestion struct {
	Question string
	Items    []string
}

func (q *ClosedQuestion) Ask() ([]Answer, error) {
	prompt := promptui.Select{
		Label:        q.Question,
		Items:        q.Items,
		HideSelected: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return []Answer{Answer(result)}, nil
}

// QuestionBranch is a sort of "meta-question" whose Ask method will delegate to ConditionQuestion and execute BranchA
// if ConditionAnswer is met, otherwise execute BranchB
type QuestionBranch struct {
	ConditionQuestion Question
	ConditionAnswer   Answer
	BranchA           []Question
	BranchB           []Question
}

func (q *QuestionBranch) Ask() ([]Answer, error) {
	result, err := q.ConditionQuestion.Ask()
	if err != nil {
		return nil, err
	}

	var branchResults []Answer
	if result[len(result)-1] == q.ConditionAnswer {
		branchResults, err = AskQuestions(q.BranchA)
		if err != nil {
			return nil, err
		}
	} else {
		branchResults, err = AskQuestions(q.BranchB)
		if err != nil {
			return nil, err
		}
	}
	branchResults = append(result, branchResults...)
	return branchResults, nil
}
