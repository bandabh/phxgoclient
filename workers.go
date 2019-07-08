package phxgoclient

import (
	"time"
)

// https://bbengfort.github.io/snippets/2016/06/26/background-work-goroutines-timer.html

type Worker struct {
	Stopped         bool          // A flag determining the state of the worker
	ShutdownChannel chan string   // A channel to communicate to the routine
	Interval        time.Duration // The interval with which to run the Action
	period          time.Duration // The actual period of the wait
	action          ExecutionAction
}

func NewWorker(interval time.Duration, action ExecutionAction) *Worker {
	return &Worker{
		Stopped:         false,
		ShutdownChannel: make(chan string),
		Interval:        interval,
		period:          interval,
		action:          action,
	}
}

// Run starts the worker and listens for a shutdown call.
func (w *Worker) Run() {
	// Loop that runs forever
	for {
		select {
		case <-w.ShutdownChannel:
			w.ShutdownChannel <- "Down"
			return
		case <-time.After(w.period):
			break
		}

		started := time.Now()
		w.action()
		finished := time.Now()

		duration := finished.Sub(started)
		w.period = w.Interval - duration

	}

}

func (w *Worker) Shutdown() {
	w.Stopped = true

	w.ShutdownChannel <- "Down"
	<-w.ShutdownChannel

	close(w.ShutdownChannel)
}

type ExecutionAction func()

func (w *Worker) Action(action ExecutionAction) {
	action()
}
