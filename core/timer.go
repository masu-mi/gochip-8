package core

import (
	"sync"
	"time"
)

type DelayedTimer struct {
	mux    sync.Mutex
	ticker *time.Ticker
	h      TimerHandler

	v uint8
}

type TimerHandler interface {
	Start()
	Stop()
}

func NewDelayedTimer(hz uint, h TimerHandler) *DelayedTimer {
	t := &DelayedTimer{
		ticker: time.NewTicker(time.Second / time.Duration(hz)),
		h:      h,
	}
	go func() {
		for {
			_ = <-t.ticker.C
			t.mux.Lock()
			if t.v > 0 {
				t.v--
				if t.v == 0 && t.h != nil {
					t.h.Stop()
				}
			}
			t.mux.Unlock()
		}
	}()
	return t
}

func (dt *DelayedTimer) GetV() uint8 {
	dt.mux.Lock()
	defer dt.mux.Unlock()
	return dt.v
}

func (dt *DelayedTimer) SetV(v uint8) {
	dt.mux.Lock()
	defer dt.mux.Unlock()
	if dt.v == 0 && dt.h != nil {
		defer dt.h.Start()
	}
	dt.v = v
}
