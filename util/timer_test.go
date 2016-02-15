// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
// Author: Nathan VanBenschoten (nvanbenschoten@gmail.com)

package util

import (
	"testing"
	"time"
)

const timeStep = 10 * time.Millisecond

func TestTimerTimeout(t *testing.T) {
	var timer Timer
	defer func() {
		if stopped := timer.Stop(); stopped {
			t.Errorf("expected Stop to return false, got true")
		}
	}()
	timer.Reset(timeStep)

	<-timer.C
	timer.Read = true

	select {
	case <-timer.C:
		t.Errorf("expected timer to only timeout once after Reset; got two timeouts")
	case <-time.After(5 * timeStep):
	}
}

func TestTimerStop(t *testing.T) {
	var timer Timer
	timer.Reset(timeStep)
	if stopped := timer.Stop(); !stopped {
		t.Errorf("expected Stop to return true, got false")
	}

	select {
	case <-timer.C:
		t.Errorf("expected timer to stop after call to Stop; got timer that was not stopped")
	case <-time.After(5 * timeStep):
	}
}

func TestTimerUninitializedStopNoop(t *testing.T) {
	var timer Timer
	if stopped := timer.Stop(); stopped {
		t.Errorf("expected Stop to return false when the timer was never reset, got true")
	}
}

func TestTimerResetBeforeTimeout(t *testing.T) {
	var timer Timer
	defer timer.Stop()
	timer.Reset(timeStep)

	timer.Reset(timeStep)
	<-timer.C
	timer.Read = true

	select {
	case <-timer.C:
		t.Errorf("expected timer to only timeout once after Reset; got two timeouts")
	case <-time.After(5 * timeStep):
	}
}

func TestTimerResetAfterTimeoutAndNoRead(t *testing.T) {
	var timer Timer
	defer timer.Stop()
	timer.Reset(timeStep)

	time.Sleep(2 * timeStep)

	timer.Reset(timeStep)
	<-timer.C
	timer.Read = true

	select {
	case <-timer.C:
		t.Errorf("expected timer to only timeout once after Reset; got two timeouts")
	case <-time.After(5 * timeStep):
	}
}

func TestTimerResetAfterTimeoutAndRead(t *testing.T) {
	var timer Timer
	defer timer.Stop()
	timer.Reset(timeStep)

	<-timer.C
	timer.Read = true

	timer.Reset(timeStep)
	<-timer.C
	timer.Read = true

	select {
	case <-timer.C:
		t.Errorf("expected timer to only timeout once after Reset; got two timeouts")
	case <-time.After(5 * timeStep):
	}
}

func TestTimerMakesProgressInLoop(t *testing.T) {
	var timer Timer
	defer timer.Stop()
	for i := 0; i < 5; i++ {
		timer.Reset(timeStep)
		<-timer.C
		timer.Read = true
	}
}
