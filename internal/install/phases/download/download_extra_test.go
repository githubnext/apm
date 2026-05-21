package download_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/download"
)

type countingDownloader struct {
	mu    sync.Mutex
	calls int
}

func (d *countingDownloader) DownloadPackage(ref, path string) (interface{}, error) {
	d.mu.Lock()
	d.calls++
	d.mu.Unlock()
	return ref, nil
}

type partialDownloader struct {
	failKeys map[string]bool
}

func (d *partialDownloader) DownloadPackage(ref, path string) (interface{}, error) {
	if d.failKeys[ref] {
		return nil, errors.New("partial failure")
	}
	return ref, nil
}

func TestRunParallelDownload_SingleWorker(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "a", DownloadRef: "r-a", ShortName: "a"},
		{DepKey: "b", DownloadRef: "r-b", ShortName: "b"},
		{DepKey: "c", DownloadRef: "r-c", ShortName: "c"},
	}
	d := &countingDownloader{}
	result := download.RunParallelDownload(tasks, 1, d, nil)
	if len(result) != 3 {
		t.Errorf("expected 3 results, got %d", len(result))
	}
	if d.calls != 3 {
		t.Errorf("expected 3 download calls, got %d", d.calls)
	}
}

func TestRunParallelDownload_AllFail(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "x", DownloadRef: "rx", ShortName: "x"},
		{DepKey: "y", DownloadRef: "ry", ShortName: "y"},
	}
	d := &partialDownloader{failKeys: map[string]bool{"rx": true, "ry": true}}
	result := download.RunParallelDownload(tasks, 2, d, nil)
	if len(result) != 0 {
		t.Errorf("expected 0 results when all fail, got %d", len(result))
	}
}

func TestRunParallelDownload_PartialSuccess(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "ok", DownloadRef: "r-ok", ShortName: "ok"},
		{DepKey: "bad", DownloadRef: "r-bad", ShortName: "bad"},
	}
	d := &partialDownloader{failKeys: map[string]bool{"r-bad": true}}
	result := download.RunParallelDownload(tasks, 2, d, nil)
	if len(result) != 1 {
		t.Errorf("expected 1 success, got %d", len(result))
	}
	if _, ok := result["ok"]; !ok {
		t.Error("expected ok in results")
	}
	if _, ok := result["bad"]; ok {
		t.Error("bad should not be in results")
	}
}

func TestRunParallelDownload_ManyTasksManyWorkers(t *testing.T) {
	const n = 20
	tasks := make([]download.DownloadTask, n)
	for i := 0; i < n; i++ {
		key := "task-" + string(rune('a'+i%26))
		if i >= 26 {
			key = key + string(rune('0'+i/26))
		}
		tasks[i] = download.DownloadTask{DepKey: key, DownloadRef: "ref-" + key, ShortName: key}
	}
	d := &countingDownloader{}
	result := download.RunParallelDownload(tasks, 10, d, nil)
	if len(result) != n {
		t.Errorf("expected %d results, got %d", n, len(result))
	}
}

func TestRunParallelDownload_DepKeyInResult(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "mykey", DownloadRef: "ref1", ShortName: "mykey"},
	}
	d := &countingDownloader{}
	result := download.RunParallelDownload(tasks, 1, d, nil)
	if _, ok := result["mykey"]; !ok {
		t.Error("expected result keyed by DepKey")
	}
}

func TestRunParallelDownload_NilProgress(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "p", DownloadRef: "r", ShortName: "p"},
	}
	d := &countingDownloader{}
	// Should not panic with nil progress
	result := download.RunParallelDownload(tasks, 1, d, nil)
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}

func TestDownloadTask_Fields(t *testing.T) {
	task := download.DownloadTask{
		DepKey:      "key",
		DownloadRef: "ref",
		InstallPath: "/path/to/install",
		DisplayName: "My Package",
		ShortName:   "pkg",
	}
	if task.DepKey != "key" || task.ShortName != "pkg" {
		t.Error("DownloadTask fields not set correctly")
	}
}

func TestDownloadResult_Fields(t *testing.T) {
	res := download.DownloadResult{
		DepKey: "key",
		Info:   "info-data",
		Err:    nil,
	}
	if res.DepKey != "key" || res.Err != nil {
		t.Error("DownloadResult fields not set correctly")
	}
}

func TestRunParallelDownload_ReturnType(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "a", DownloadRef: "ra", ShortName: "a"},
	}
	d := &countingDownloader{}
	result := download.RunParallelDownload(tasks, 1, d, nil)
	// Result must be a map (not nil even on partial failure)
	if result == nil {
		t.Error("result should never be nil")
	}
}
