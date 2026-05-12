// Package jobevents implements a lightweight publish/subscribe event bus
// for cronwrap job lifecycle events.
//
// # Overview
//
// A Bus allows multiple components (loggers, notifiers, webhooks, etc.) to
// react to job lifecycle transitions without being directly coupled to one
// another.
//
// # Usage
//
//	bus, err := jobevents.New(16)
//	if err != nil { ... }
//
//	ch, _ := bus.Subscribe(jobevents.EventFailed)
//	go func() {
//		for evt := range ch {
//			fmt.Println("job failed:", evt.Job)
//		}
//	}()
//
//	bus.Publish(jobevents.Event{
//		Type: jobevents.EventFailed,
//		Job:  "daily-backup",
//		Payload: map[string]any{"exit_code": 1},
//	})
//
// # Thread Safety
//
// All methods on Bus are safe for concurrent use.
package jobevents
