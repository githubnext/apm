package tokenmanager

import "time"

// timerAfter returns a channel that closes after n seconds.
func timerAfter(n int) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(n) * time.Second)
		close(ch)
	}()
	return ch
}
