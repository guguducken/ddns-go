package ddns

import "context"

type Stopper interface {
	// Stop must execute asynchronously and return a channel which report shutdown finished
	Stop() chan struct{}
	SetCancelFunc(cancel context.CancelFunc)
}
