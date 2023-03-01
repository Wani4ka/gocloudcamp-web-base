package playlist

import (
	"time"
)

type Timer interface {
	Schedule(length time.Duration, callback func())
	Stop()
	Pause()
	IsPaused() bool
	Resume()
}

type timerImpl struct {
	length    time.Duration
	startedAt time.Time
	pausedAt  time.Time
	callback  func()
	object    *time.Timer
}

func NewTimer() Timer {
	return &timerImpl{}
}

func (timer *timerImpl) Stop() {
	timer.startedAt = time.Time{}
	timer.pausedAt = time.Time{}
	timer.length = 0
	timer.callback = nil
	if timer.object != nil {
		timer.object.Stop()
		timer.object = nil
	}
}

func work(timer *time.Timer, callback, stop func()) {
	<-timer.C
	stop()
	callback()
}

func (timer *timerImpl) Schedule(length time.Duration, callback func()) {
	timer.Stop()
	timer.startedAt = time.Now()
	timer.length = length
	timer.callback = callback
	timer.object = time.NewTimer(length)
	go work(timer.object, callback, timer.Stop)
}

func (timer *timerImpl) Pause() {
	if timer.IsPaused() || timer.object == nil {
		return
	}
	timer.pausedAt = time.Now()
	timer.object.Stop()
}

func (timer *timerImpl) Resume() {
	if !timer.IsPaused() {
		return
	}
	elapsed := timer.pausedAt.Sub(timer.startedAt)
	timer.startedAt = time.Now().Add(-elapsed)
	timer.pausedAt = time.Time{}
	timer.object = time.NewTimer(timer.length - elapsed)
	go work(timer.object, timer.callback, timer.Stop)
}

func (timer *timerImpl) IsPaused() bool {
	return !timer.startedAt.IsZero() && !timer.pausedAt.IsZero()
}
