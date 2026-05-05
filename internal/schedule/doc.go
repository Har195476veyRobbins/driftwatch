// Package schedule provides a periodic scheduler used by driftwatch to
// trigger drift detection at a configurable interval.
//
// Basic usage:
//
//	s := schedule.New(schedule.Config{Interval: 30 * time.Second})
//	err := s.Run(ctx, func(ctx context.Context) {
//		// perform drift check
//	})
//
// The scheduler fires the callback immediately on start and then on every
// tick. Cancelling the context stops the loop and Run returns ctx.Err().
package schedule
