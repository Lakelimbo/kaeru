package jobs

import (
	"sync"
)

// Extremely basic pubsub for messaging
type JobMessaging struct {
	subs map[string][]chan string
	mu   sync.RWMutex
}

func NewPubSub() *JobMessaging {
	return &JobMessaging{
		subs: make(map[string][]chan string),
	}
}

// Pub publishes a message to the stack.
//
// Basically, when a Job is created, it is sent through a
// RWMutex to be executed inside of a goroutine.
func (j *JobMessaging) Pub(jobID string, msg string) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	for _, ch := range j.subs[jobID] {
		select {
		case ch <- msg:
		default:
		}
	}
}

// Sub subscribes to an existing message.
//
// After a message is published, you can subscribe to any eventual
// changes (such as in jobs.UpdateJob() ) that will then be re-streamed
// into here.
func (j *JobMessaging) Sub(jobID string) chan string {
	ch := make(chan string, 10)

	j.mu.Lock()
	j.subs[jobID] = append(j.subs[jobID], ch)
	j.mu.Unlock()

	return ch
}

// Unsub unsubscribes the current subscribed message.
//
// After the messaging is done (e.g. the job is completed/errored), you
// should unsubscribe, as there is nothing more to watch, and it will also
// remove it from the stack.
func (j *JobMessaging) Unsub(jobID string, ch chan string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	subs := j.subs[jobID]
	for i, c := range subs {
		if c == ch {
			j.subs[jobID] = append(subs[:i], subs[i+1:]...)
			close(c)
			break
		}
	}
}
