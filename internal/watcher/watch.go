package watcher

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Watch(dir string, onChange func()) *Batcher {
	watcher, err := NewBatchWatcher(1 * time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	watcher.Add(dir)

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				// fmt.Println("received system events:", evs)
				for _, ev := range evs {
					// sometimes during rm -rf operations a '"": REMOVE' is triggered, just ignore these
					if ev.Name == "" {
						continue
					}
					// events to watch
					importantEvent := (ev.Op == fsnotify.Create || ev.Op == fsnotify.Write || ev.Op == fsnotify.Rename || ev.Op == fsnotify.Remove)
					if importantEvent {
						onChange()
						break
					}
				}
			}
		}
	}()

	return watcher
}
