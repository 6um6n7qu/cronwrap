package webhook

import (
	"context"
	"log"
	"time"
)

// Dispatcher wraps a Client and provides convenience methods for common
// cronwrap job lifecycle events.
type Dispatcher struct {
	client *Client
}

// NewDispatcher returns a Dispatcher backed by the given Client.
func NewDispatcher(c *Client) *Dispatcher {
	return &Dispatcher{client: c}
}

// JobStarted fires a "started" webhook for the named job.
func (d *Dispatcher) JobStarted(ctx context.Context, jobName string) {
	d.send(ctx, Payload{
		JobName: jobName,
		Event:   "started",
	})
}

// JobFinished fires a "finished" webhook including the exit code and duration.
func (d *Dispatcher) JobFinished(ctx context.Context, jobName string, exitCode int, dur time.Duration) {
	d.send(ctx, Payload{
		JobName:  jobName,
		Event:    "finished",
		ExitCode: exitCode,
		Duration: dur / time.Millisecond,
	})
}

// JobFailed fires a "failed" webhook including the error message.
func (d *Dispatcher) JobFailed(ctx context.Context, jobName string, exitCode int, errMsg string) {
	d.send(ctx, Payload{
		JobName:  jobName,
		Event:    "failed",
		ExitCode: exitCode,
		Error:    errMsg,
	})
}

// send is a best-effort delivery helper; errors are logged but not propagated
// so that webhook failures never block job execution.
func (d *Dispatcher) send(ctx context.Context, p Payload) {
	if err := d.client.Send(ctx, p); err != nil {
		log.Printf("webhook dispatch error (job=%s event=%s): %v", p.JobName, p.Event, err)
	}
}
