package download_test

import (
	"errors"
	"testing"

	"github.com/githubnext/apm/internal/install/phases/download"
)

type alwaysSucceedDownloader struct{}

func (d *alwaysSucceedDownloader) DownloadPackage(ref, path string) (interface{}, error) {
	return ref + "-ok", nil
}

type alwaysFailDownloader struct{}

func (d *alwaysFailDownloader) DownloadPackage(ref, path string) (interface{}, error) {
	return nil, errors.New("always fail")
}

func TestRunParallelDownload_EmptyTasks(t *testing.T) {
	result := download.RunParallelDownload(nil, 4, &alwaysSucceedDownloader{}, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil tasks, got %v", result)
	}
}

func TestRunParallelDownload_ZeroWorkersExtra2(t *testing.T) {
	tasks := []download.DownloadTask{{DepKey: "a", DownloadRef: "ref-a", InstallPath: "/tmp/a"}}
	result := download.RunParallelDownload(tasks, 0, &alwaysSucceedDownloader{}, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for 0 workers, got %v", result)
	}
}

func TestRunParallelDownload_AllSucceed(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "pkg1", DownloadRef: "ref1", InstallPath: "/tmp/1"},
		{DepKey: "pkg2", DownloadRef: "ref2", InstallPath: "/tmp/2"},
		{DepKey: "pkg3", DownloadRef: "ref3", InstallPath: "/tmp/3"},
	}
	result := download.RunParallelDownload(tasks, 2, &alwaysSucceedDownloader{}, nil)
	if len(result) != 3 {
		t.Errorf("expected 3 results, got %d", len(result))
	}
	for _, task := range tasks {
		if _, ok := result[task.DepKey]; !ok {
			t.Errorf("expected key %s in result", task.DepKey)
		}
	}
}

func TestRunParallelDownload_AllFail_Extra2(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "pkg1", DownloadRef: "ref1", InstallPath: "/tmp/1"},
		{DepKey: "pkg2", DownloadRef: "ref2", InstallPath: "/tmp/2"},
	}
	result := download.RunParallelDownload(tasks, 4, &alwaysFailDownloader{}, nil)
	if len(result) != 0 {
		t.Errorf("expected 0 results for all-fail, got %d", len(result))
	}
}

func TestDownloadTask_FieldsSet(t *testing.T) {
	task := download.DownloadTask{
		DepKey:      "dep-key",
		DownloadRef: "v1.0",
		InstallPath: "/opt/pkg",
		DisplayName: "My Package",
		ShortName:   "mypkg",
	}
	if task.DepKey != "dep-key" {
		t.Error("DepKey field mismatch")
	}
	if task.ShortName != "mypkg" {
		t.Error("ShortName field mismatch")
	}
}

func TestDownloadResult_ZeroValue(t *testing.T) {
	var r download.DownloadResult
	if r.DepKey != "" || r.Info != nil || r.Err != nil {
		t.Error("DownloadResult zero value should have empty fields")
	}
}

func TestRunParallelDownload_WorkersCappedToTasks(t *testing.T) {
	tasks := []download.DownloadTask{
		{DepKey: "only", DownloadRef: "ref", InstallPath: "/tmp/only"},
	}
	// More workers than tasks — should still work
	result := download.RunParallelDownload(tasks, 100, &alwaysSucceedDownloader{}, nil)
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}
