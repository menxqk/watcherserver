package watcher

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_Watcher(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	w := New(curDir)
	evtCh, err := w.Start()
	if err != nil {
		t.Error(err)
	}
	defer w.Stop()

	testFilePath := filepath.Join(curDir, "test.file")

	// test created event
	f, err := os.Create(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	evt := <-evtCh
	if evt.Op != Created {
		t.Errorf("expected %s, got %s", Created, evt.Op)
	}
	if evt.Path != testFilePath {
		t.Errorf("wrong file path %s", evt.Path)
	}
	f.Close()

	// test changed event
	f, err = os.OpenFile(testFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	_, err = f.WriteString("test file")
	if err != nil {
		t.Error(err)
	}
	f.Close()
	evt = <-evtCh
	if evt.Op != Changed {
		t.Errorf("expected %s, got %s", Changed, evt.Op)
	}
	if evt.Path != testFilePath {
		t.Errorf("wrong file path %s", evt.Path)
	}

	// test deleted event
	os.Remove(testFilePath)
	evt = <-evtCh
	if evt.Op != Deleted {
		t.Errorf("expected %s, got %s", Deleted, evt.Op)
	}
	if evt.Path != testFilePath {
		t.Errorf("wrong file path %s", evt.Path)
	}
}
