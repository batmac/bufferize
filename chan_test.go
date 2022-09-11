// Copyright 2022 Baptiste Canton.
// SPDX-License-Identifier: MIT

package bufferize_test

import (
	"context"
	"testing"
	"time"

	"github.com/batmac/bufferize"
)

var (
	maxDurationSeconds = 3 // must be at least 3
	bufSize            = 10
)

func TestDontPanicIfZeroSize(t *testing.T) {
	orig := make(chan int, 1)
	channel := bufferize.New(orig, 0)
	go func() {
		orig <- 1
		close(orig)
	}()
	for {
		v, ok := <-channel
		if !ok {
			break
		}
		t.Log(v)
	}
	t.Log("done")
}

// almost the same as NewTicker, but simpler and with known values
func tick(n int, ch chan int) {
	for i := 0; i < n; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	close(ch)
}

func TestNew(t *testing.T) {
	last := -1
	orig := make(chan int, 1)
	channel := bufferize.New(orig, bufSize)

	go tick(maxDurationSeconds, orig)

	time.Sleep(2 * time.Second)
	for {
		v, ok := <-channel
		if !ok {
			break
		}
		last = v
		t.Log(v)
	}
	t.Log("done")
	if last != maxDurationSeconds-1 {
		t.Errorf("expected %d, got %d", maxDurationSeconds-1, last)
	}
}

func TestNewCtx(t *testing.T) {
	last := -1
	orig := make(chan int, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0.7*float32(maxDurationSeconds))*time.Second)
	defer cancel()
	channel := bufferize.NewCtx(ctx, orig, 10)

	go tick(10, orig)
	time.Sleep(2 * time.Second)

loop:
	for {
		select {
		case v, ok := <-channel:
			if !ok {
				break loop
			}
			last = v
			t.Log(v)
		case <-ctx.Done():
			t.Log("done (ctx)")
			break loop
		}
	}
	if last > maxDurationSeconds/2 {
		t.Errorf("last value (%d) is too high", last)
	}
	t.Logf("last value: %d", last)
}
