package app

import (
	"time"
)

// Worker is a simple event loop that processes tasks in FIFO policy.
// Provides a few queues, all of which will be processed in order.
type Worker struct {
	// debounce pair queue
	dpq chan *DebounceLoad

	// timed debounce queue
	dq chan func()

	// simple queue
	q chan func()

	// delay to be used in both debounce queues
	d time.Duration
}

// DebounceLoad holds 3 functions that are used to debounce incoming events.
type DebounceLoad struct {
	// refresh will be called until there are no more events left
	refresh func()

	// final will be called after there are no more events
	final func()

	// loading function, after a delay at each event if debounce has not finished yet
	onLoad func()
}

func newRefreshQueue(d time.Duration, qSize int) *Worker {
	r := &Worker{}
	r.q = make(chan func(), qSize)
	r.dpq = make(chan *DebounceLoad, qSize)
	r.d = d
	return r
}

// Starts the event loop. Should be called as/in a goroutine.
func (r *Worker) Start() {
	for {
		select {
		case f := <-r.q:
			r.handleSimple(f)
		case f := <-r.dq:
			r.handleDebounce(f)
		case dp := <-r.dpq:
			r.handleDebounceLoad(dp)
		}
	}
}

// Queues function f in the simple queue and will be ran by the worker when it gets the chance.
func (r *Worker) Queue(f func()) {
	r.q <- f
}

// Queues f to run by the worker and will get ran after (last call to Debounce()) + d.
func (r *Worker) Debounce(f func()) {
	r.dq <- f
}

// Queues the triplet to run by the worker as a debounce load.
// Continuously consumes events produced by DebounceLoad() and runs refresh. Runs final after the last event.
// Runs onLoad if refresh takes longer that d (delay).
func (r *Worker) DebounceLoad(refresh func(), final func(), onLoad func()) {
	r.dpq <- &DebounceLoad{
		refresh: refresh,
		final:   final,
		onLoad:  onLoad,
	}
}

func (r *Worker) handleSimple(f func()) {
	f()
}

func (r *Worker) handleDebounce(f func()) {
	for {
		// start timer
		timer := time.NewTimer(r.d)
		defer timer.Stop()

		// wait for the next event or timer to elapse
		select {
		case f = <-r.dq:
			// if we get next event we start a new timer and loop back
			timer.Stop()
			timer.Reset(r.d)
		case <-timer.C:
			// if the timer fires we execute the function and stop the debounce
			f()
			return
		}
	}
}

func (r *Worker) handleDebounceLoad(dp *DebounceLoad) {
	// keep track if a event was seen (debounce was called)
	seen := false
	for {
		// drain the channel
		func() {
			for {
				select {
				case dp = <-r.dpq:
					seen = false
				default:
					return
				}
			}
		}()

		// if there was no new event since debounce call, we call final and stop
		if seen {
			dp.final()
			return
		}

		// call onload after d if refresh is still running
		onLoadTimer := time.AfterFunc(r.d, func() {
			dp.onLoad()
		})

		// call refresh
		dp.refresh()

		// cancel timer so that if refresh finishes before, onLoad wont fire
		onLoadTimer.Stop()

		// we called debounce, so we could mark it as seen
		seen = true
	}
}
