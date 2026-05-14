// Package download implements the parallel package pre-download phase of the
// install pipeline. Mirrors src/apm_cli/install/phases/download.py.
package download

import (
	"sync"
)

// DownloadTask describes a single package that needs to be fetched.
type DownloadTask struct {
	DepKey      string
	DownloadRef string
	InstallPath string
	DisplayName string
	ShortName   string
}

// DownloadResult holds the outcome of one download task.
type DownloadResult struct {
	DepKey  string
	Info    interface{} // opaque PackageInfo returned by the downloader
	Err     error
}

// Downloader is implemented by the component that fetches packages.
type Downloader interface {
	DownloadPackage(downloadRef, installPath string) (interface{}, error)
}

// ProgressReporter is an optional TUI adapter.
type ProgressReporter interface {
	TaskStarted(depKey, label string)
	TaskCompleted(depKey string)
	TaskFailed(depKey string)
}

// RunParallelDownload executes all tasks concurrently up to maxWorkers.
// Returns a map[depKey]PackageInfo for successful downloads; failures are
// silently dropped so the sequential integration loop retries and reports.
func RunParallelDownload(
	tasks []DownloadTask,
	maxWorkers int,
	downloader Downloader,
	progress ProgressReporter,
) map[string]interface{} {
	if len(tasks) == 0 || maxWorkers <= 0 {
		return map[string]interface{}{}
	}

	workers := maxWorkers
	if workers > len(tasks) {
		workers = len(tasks)
	}

	resultsCh := make(chan DownloadResult, len(tasks))
	tasksCh := make(chan DownloadTask, len(tasks))

	for _, t := range tasks {
		tasksCh <- t
	}
	close(tasksCh)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasksCh {
				if progress != nil {
					progress.TaskStarted(t.DepKey, "fetch "+t.ShortName)
				}
				info, err := downloader.DownloadPackage(t.DownloadRef, t.InstallPath)
				resultsCh <- DownloadResult{DepKey: t.DepKey, Info: info, Err: err}
				if err != nil {
					if progress != nil {
						progress.TaskFailed(t.DepKey)
					}
				} else {
					if progress != nil {
						progress.TaskCompleted(t.DepKey)
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	results := make(map[string]interface{}, len(tasks))
	for r := range resultsCh {
		if r.Err == nil {
			results[r.DepKey] = r.Info
		}
	}
	return results
}
