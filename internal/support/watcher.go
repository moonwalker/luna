// Copyright (c) The Hugo Authors. All rights reserved.
// https://github.com/gohugoio/hugo/blob/master/watcher/batcher.go

package support

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

// Batcher batches file watch events in a given interval.
type Batcher struct {
	*fsnotify.Watcher
	interval time.Duration
	done     chan struct{}

	Events chan []fsnotify.Event // Events are returned on this channel
}

// NewWatcher creates and starts a Batcher with the given time interval.
func NewWatcher(interval time.Duration) (*Batcher, error) {
	watcher, err := fsnotify.NewWatcher()

	batcher := &Batcher{}
	batcher.Watcher = watcher
	batcher.interval = interval
	batcher.done = make(chan struct{}, 1)
	batcher.Events = make(chan []fsnotify.Event, 1)

	if err == nil {
		go batcher.run()
	}

	return batcher, err
}

func (b *Batcher) run() {
	tick := time.Tick(b.interval)
	evs := make([]fsnotify.Event, 0)
OuterLoop:
	for {
		select {
		case ev := <-b.Watcher.Events:
			evs = append(evs, ev)
		case <-tick:
			if len(evs) == 0 {
				continue
			}
			b.Events <- evs
			evs = make([]fsnotify.Event, 0)
		case <-b.done:
			break OuterLoop
		}
	}
	close(b.done)
}

// Close stops the watching of the files.
func (b *Batcher) Close() {
	b.done <- struct{}{}
	b.Watcher.Close()
}
