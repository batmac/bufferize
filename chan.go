// Copyright 2022 Baptiste Canton.
// SPDX-License-Identifier: MIT

package bufferize

import "context"

const DEF_BUF_SIZE = 10

func New[T any](orig <-chan T, size int) chan T {
	return NewCtx(context.Background(), orig, size)
}

func NewCtx[T any](ctx context.Context, orig <-chan T, size int) chan T {
	if size <= 0 {
		size = DEF_BUF_SIZE
	}
	ch := make(chan T, size)
	go func() {
		for {
			select {
			case v, ok := <-orig:
				if !ok {
					close(ch)
					return
				}
				ch <- v
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
