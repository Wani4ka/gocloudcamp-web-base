package playlist

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	got := NewTimer()
	_, ok := got.(Timer)
	if !ok {
		t.Fatalf("got %v, want Timer", reflect.TypeOf(got))
	}
}

func TestTimerSchedule(t *testing.T) {
	timer := NewTimer()
	dur := 500 * time.Millisecond
	grace := dur / 20
	ch := make(chan time.Time)
	timer.Schedule(dur, func() {
		ch <- time.Now()
	})
	time.Sleep(dur)
	want := time.Now()
	got := <-ch
	if got.After(want) && got.Sub(want) > grace {
		t.Fatalf("Timer finished %v late, only %v allowed", got.Sub(want), grace)
	} else if want.After(got) && want.Sub(got) > grace {
		t.Fatalf("Timer finished %v early, only %v allowed", want.Sub(got), grace)
	}
}

func TestTimerStop(t *testing.T) {
	timer := NewTimer()
	dur := 500 * time.Millisecond
	wait := 100 * time.Millisecond
	timer.Schedule(dur, func() {
		t.Fatal("Timer callback called besides being stopped")
	})
	time.Sleep(wait)
	timer.Stop()
	time.Sleep(dur)
}

func TestTimerPauseResume(t *testing.T) {
	timer := NewTimer()
	dur := 500 * time.Millisecond
	wait := 100 * time.Millisecond
	grace := dur / 20
	ch := make(chan time.Time)
	timer.Schedule(dur, func() {
		ch <- time.Now()
	})
	if timer.IsPaused() {
		t.Fatal("Timer is marked as paused right after being scheduled")
		return
	}
	time.Sleep(wait)
	timer.Pause()
	if !timer.IsPaused() {
		t.Fatal("Timer is not marked as paused but wasn't actually resumed")
		return
	}
	timer.Pause()
	if !timer.IsPaused() {
		t.Fatal("Timer is not marked as paused after calling Pause method again")
		return
	}
	timer.Resume()
	if timer.IsPaused() {
		t.Fatal("Timer is marked as paused right after being resumed")
		return
	}
	timer.Resume()
	if timer.IsPaused() {
		t.Fatal("Timer is marked as paused after calling Resume method again")
		return
	}
	time.Sleep(dur - wait)
	want := time.Now()
	got := <-ch
	if got.After(want) && got.Sub(want) > grace {
		t.Fatalf("Timer finished %v late, only %v allowed", got.Sub(want), grace)
	} else if want.After(got) && want.Sub(got) > grace {
		t.Fatalf("Timer finished %v early, only %v allowed", want.Sub(got), grace)
	}
}
