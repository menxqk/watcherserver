package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const TIMEOUT = 100 * time.Millisecond

type Op uint

const (
	Created Op = iota
	Deleted
	Changed
)

var ops = map[Op]string{
	Created: "Created",
	Deleted: "Deleted",
	Changed: "Changed",
}

func (o Op) String() string {
	if op, found := ops[o]; found {
		return op
	}
	return "---"
}

type Event struct {
	Op
	Path string
}

type watcher struct {
	baseDir string
	evtCh   chan Event
	closeCh chan struct{}
}

type fileList map[string]os.FileInfo

func New(dir string) *watcher {
	w := &watcher{
		baseDir: dir,
	}
	return w
}

func (w *watcher) Start() (chan Event, error) {
	list, err := w.retrieveFileList()
	if err != nil {
		return nil, err
	}

	w.evtCh = make(chan Event)
	w.closeCh = make(chan struct{})

	go func(l fileList) {
	outer:
		for {
			timer := time.After(TIMEOUT)
			select {
			case <-timer:
				newList, err := w.retrieveFileList()
				if err != nil {
					fmt.Println(err)
					continue
				}
				w.checkForEvent(l, newList)
				l = newList
			case <-w.closeCh:
				break outer
			}
		}
		w.closeCh <- struct{}{}
	}(list)

	return w.evtCh, nil
}

func (w *watcher) Stop() {
	w.closeCh <- struct{}{}
	<-w.closeCh
	close(w.evtCh)
	close(w.closeCh)
}

func (w *watcher) retrieveFileList() (fileList, error) {
	entries, err := os.ReadDir(w.baseDir)
	if err != nil {
		return nil, err
	}

	list := make(fileList)

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			continue
		}
		name := filepath.Join(w.baseDir, info.Name())
		list[name] = info
	}

	return list, nil
}

func (w *watcher) checkForEvent(list fileList, newList fileList) {
	// check for created files
	for name := range newList {
		if _, found := list[name]; !found {
			evt := Event{
				Op:   Created,
				Path: name,
			}
			w.evtCh <- evt
			return
		}
	}

	// check for deleted files
	for name := range list {
		if _, found := newList[name]; !found {
			evt := Event{
				Op:   Deleted,
				Path: name,
			}
			w.evtCh <- evt
			return
		}
	}

	// check for changed files
	for name, info := range list {
		if _, found := newList[name]; found {
			newInfo := newList[name]
			if newInfo.ModTime().After(info.ModTime()) {
				evt := Event{
					Op:   Changed,
					Path: name,
				}
				w.evtCh <- evt
				return
			}
		}
	}
}
