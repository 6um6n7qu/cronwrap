// Package retry implements configurable retry logic for cron job execution.
//
// It supports:
//   - A maximum attempt count
//   - A base delay between attempts
//   - An optional exponential back-off multiplier
//   - Context-aware cancellation between retries
//
// Basic usage:
//
//	p := retry.DefaultPolicy()
//	attempts, err := retry.Do(ctx, p, func(ctx context.Context) error {
//		return runMyJob(ctx)
//	})
//	if err != nil {
//		log.Printf("job failed after %d attempts: %v", len(attempts), err)
//	}
package retry
