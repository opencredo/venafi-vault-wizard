package questions

type Answer string

// AnswerQueue is a FIFO collection of answers where Push adds them to the queue, and Pop retrieves them in the order
// they were added
type AnswerQueue []*Answer

// NewAnswerQueue returns a new empty AnswerQueue
func NewAnswerQueue() *AnswerQueue {
	return &AnswerQueue{}
}

// Push adds an Answer to the back of the queue
func (q *AnswerQueue) Push(answer Answer) {
	*q = append(*q, &answer)
}

// Pop retrieves an Answer from the front of the queue and deletes it
func (q *AnswerQueue) Pop() *Answer {
	if len(*q) == 0 {
		return nil
	}
	// take first value of slice as front of queue
	value := (*q)[0]
	// shrink slice by chopping off first value
	*q = (*q)[1:]

	return value
}

// PeekLast returns the Answer at the back of the queue (most recently added) without deleting it
func (q *AnswerQueue) PeekLast() *Answer {
	if len(*q) == 0 {
		return nil
	}
	return (*q)[len(*q)-1]
}
