// Package webhook provides HTTP webhook delivery for cronwrap job lifecycle
// events.
//
// # Overview
//
// A Client handles low-level HTTP POST delivery to a configured URL, while a
// Dispatcher wraps the Client with named helpers for common events:
//
//	- JobStarted  — fired immediately before a job begins execution.
//	- JobFinished — fired on successful completion, includes exit code and wall
//	                clock duration.
//	- JobFailed   — fired when a job exits with a non-zero status or encounters
//	                a runtime error, includes the error message.
//
// # Authentication
//
// When a non-empty secret is provided to New, every outgoing request carries
// an X-Webhook-Secret header so that receivers can verify the source.
//
// # Error Handling
//
// Dispatcher.send is best-effort: delivery failures are logged via the
// standard library logger but never propagate to the caller, ensuring that
// webhook outages cannot interrupt job execution.
package webhook
