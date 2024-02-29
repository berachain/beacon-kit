package async

import (
	"context"
	"sync"
)

// WaitForCondition waits for a given condition to be met.
// The condition is checked by the provided checkFunc.
// If the condition is not met immediately, it periodically calls the provided progressFunc to attempt to meet the condition.
func WaitForCondition(
	ctx context.Context, checkFunc func() error,
	progressFunc func(context.Context),
	stopCh chan struct{}, cond *sync.Cond,
) {
	cond.L.Lock()
	defer cond.L.Unlock()

	// If the condition is not met, we immediately request progress.
	if checkFunc() != nil {
		go progressFunc(ctx)
		for checkFunc() != nil {
			select {
			case <-stopCh:
				return
			case <-ctx.Done():
				return
			default:
				// Then we wait until the condition is met.
				cond.Wait()
			}
		}
	}
}
