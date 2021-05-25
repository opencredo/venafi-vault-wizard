package questions

import "github.com/manifoldco/promptui"

// AskQuestions takes a list of questions and asks each one in sequence, returning the list of answers
func AskQuestions(questions []Question, answers *AnswerQueue) error {
	for _, question := range questions {
		err := question.Ask(answers)
		if err != nil {
			return err
		}
	}
	return nil
}

func AskSingleQuestion(question Question) (Answer, error) {
	answerQueue := NewAnswerQueue()
	err := question.Ask(answerQueue)
	if err != nil {
		return "", err
	}
	return *answerQueue.Pop(), nil
}

// Question represents a generic user prompt
type Question interface {
	Ask(queue *AnswerQueue) error
}

// OpenEndedQuestion is a prompt where the user can provide any answer
type OpenEndedQuestion struct {
	Question  string
	Default   string
	AllowEdit bool
	Validate  promptui.ValidateFunc
}

func (q *OpenEndedQuestion) Ask(queue *AnswerQueue) error {
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

	queue.Push(Answer(result))
	return nil
}

// ClosedQuestion is a prompt where the user can choose from a fixed set of answers
type ClosedQuestion struct {
	Question string
	Items    []string
}

func (q *ClosedQuestion) Ask(queue *AnswerQueue) error {
	prompt := promptui.Select{
		Label:        q.Question,
		Items:        q.Items,
		HideSelected: true,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	queue.Push(Answer(result))
	return nil
}

// QuestionBranch is a sort of "meta-question" whose Ask method will delegate to ConditionQuestion and execute BranchA
// if ConditionAnswer is met, otherwise execute BranchB
type QuestionBranch struct {
	ConditionQuestion Question
	ConditionAnswer   Answer
	BranchA           []Question
	BranchB           []Question
}

func (q *QuestionBranch) Ask(queue *AnswerQueue) error {
	err := q.ConditionQuestion.Ask(queue)
	if err != nil {
		return err
	}

	if *queue.PeekLast() == q.ConditionAnswer {
		err := AskQuestions(q.BranchA, queue)
		if err != nil {
			return err
		}
	} else {
		err := AskQuestions(q.BranchB, queue)
		if err != nil {
			return err
		}
	}
	return nil
}
