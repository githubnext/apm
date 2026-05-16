package download_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/download"
)

// successDownloader always returns a static result.
type successDownloader struct {
	mu      sync.Mutex
	started []string
}

func (d *successDownloader) DownloadPackage(ref, path string) (interface{}, error) {
	d.mu.Lock()
	d.started = append(d.started, ref)
	d.mu.Unlock()
	return map[string]string{"ref": ref, "path": path}, nil
}

// failDownloader always returns an error.
type failDownloader struct{}

func (d *failDownloader) DownloadPackage(_, _ string) (interface{}, error) {
	return nil, errors.New("download failed")
}

// trackingProgress records events.
type trackingProgress struct {
	mu        sync.Mutex
	started   []string
	completed []string
	failed    []string
}

func (p *trackingProgress) TaskStarted(depKey, _ string) {
	p.mu.Lock()
	p.started = append(p.started, depKey)
	p.mu.Unlock()
}
func (p *trackingProgress) TaskCompleted(depKey string) {
	p.mu.Lock()
	p.completed = append(p.completed, depKey)
	p.mu.Unlock()
}
func (p *trackingProgress) TaskFailed(depKey string) {
	p.mu.Lock()
	p.failed = append(p.failed, depKey)
	p.mu.Unlock()
}

func TestRunParallelDownload_Empty(t *testing.T) {
	d := &successDownloader{}
	result := download.RunParallelDownload(nil, 4, d, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestRunParallelDownload_ZeroWorkers(t *testing.T) {
	tasks := []download.DownloadTask{{DepKey: "a", DownloadRef: "r1", ShortName: "a"}}
	d := &successDownloader{}
	result := download.RunParallelDownload(tasks, 0, d, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for 0 workers, got %d", len(result))
	}
}

func TestRunParallelDownload_Success(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "a", DownloadRef: "ref-a", InstallPath: "/tmp/a", ShortName: "a"},
		{DepKey: "b", DownloadRef: "ref-b", InstallPath: "/tmp/b", ShortName: "b"},
	}
	d := &successDownloader{}
	result := download.RunParallelDownload(tasks, 2, d, nil)
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
	if _, ok := result["a"]; !ok {
		t.Error("missing result for key a")
	}
	if _, ok := result["b"]; !ok {
		t.Error("missing result for key b")
	}
}

func TestRunParallelDownload_Failure_Excluded(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "bad", DownloadRef: "ref", ShortName: "bad"},
	}
	d := &failDownloader{}
	result := download.RunParallelDownload(tasks, 1, d, nil)
	if len(result) != 0 {
		t.Errorf("expected 0 results for failed download, got %d", len(result))
	}
}

func TestRunParallelDownload_Progress(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "x", DownloadRef: "ref-x", ShortName: "x"},
	}
	d := &successDownloader{}
	p := &trackingProgress{}
	download.RunParallelDownload(tasks, 1, d, p)
	if len(p.started) != 1 || p.started[0] != "x" {
		t.Errorf("expected started=[x], got %v", p.started)
	}
	if len(p.completed) != 1 || p.completed[0] != "x" {
		t.Errorf("expected completed=[x], got %v", p.completed)
	}
	if len(p.failed) != 0 {
		t.Errorf("expected no failures, got %v", p.failed)
	}
}

func TestRunParallelDownload_FailureProgress(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "y", DownloadRef: "ref-y", ShortName: "y"},
	}
	d := &failDownloader{}
	p := &trackingProgress{}
	download.RunParallelDownload(tasks, 1, d, p)
	if len(p.failed) != 1 || p.failed[0] != "y" {
		t.Errorf("expected failed=[y], got %v", p.failed)
	}
}

func TestRunParallelDownload_MoreWorkersThanTasks(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "a", DownloadRef: "r1", ShortName: "a"},
	}
	d := &successDownloader{}
	result := download.RunParallelDownload(tasks, 100, d, nil)
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}
