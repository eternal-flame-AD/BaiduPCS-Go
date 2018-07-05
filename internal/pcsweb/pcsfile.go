package pcsweb

import (
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"github.com/iikira/BaiduPCS-Go/internal/pcscommand"
	"github.com/iikira/BaiduPCS-Go/internal/pcsconfig"
)

const (
	statOngoing = iota
	statSuccess = iota
	statFailure = iota
)

var (
	downloadManager = &downloadController{
		downloadOptions: &pcscommand.DownloadOptions{
			IsTest:               false,
			IsPrintStatus:        false,
			IsExecutedPermission: false && runtime.GOOS != "windows",
			IsOverwrite:          false,
			IsShareDownload:      false,
			IsLocateDownload:     true,
			IsStreaming:          false,
			SaveTo:               ".",
			Parallel:             100,
			Load:                 1,
		},
		tasks: make([]*downloadTask, 0),
	}
	currid = int64(0)
)

func fileList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fpath := r.Form.Get("path")
	dataReadCloser, err := pcsconfig.Config.ActiveUserBaiduPCS().PrepareFilesDirectoriesList(fpath, baidupcs.DefaultOrderOptions)
	if err != nil {
		w.Write((&ErrInfo{
			ErrroCode: 1,
			ErrorMsg:  err.Error(),
		}).JSON())
		return
	}

	defer dataReadCloser.Close()
	io.Copy(w, dataReadCloser)
}

func downloadQuery(w http.ResponseWriter, r *http.Request) {
	w.Write(toJSON(downloadManager.getDownloadTask()))
}

func downloadStart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	path := r.Form.Get("path")
	if path == "" {
		w.WriteHeader(400)
		w.Write((&ErrInfo{ErrroCode: 31023, ErrorMsg: "incorrect parameter"}).JSON())
		return
	}

	downloadManager.addDownloadTask(path)
	w.Write((&ErrInfo{ErrroCode: 0, ErrorMsg: "no error"}).JSON())
}

func downloadCancel(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write((&ErrInfo{ErrroCode: 31210, ErrorMsg: "not implemented"}).JSON())
}

func downloadDir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write(toJSON(downloadManager.downloadOptions.SaveTo))
	case "POST":
		r.ParseForm()

		path := r.Form.Get("path")

		dir, err := os.Stat(path)
		if path == "" || err != nil || !dir.IsDir() {
			w.WriteHeader(400)
			w.Write((&ErrInfo{ErrroCode: 31023, ErrorMsg: "incorrect parameter"}).JSON())
			return
		}

		downloadManager.downloadOptions.SaveTo = path
		w.Write((&ErrInfo{ErrroCode: 0, ErrorMsg: "no error"}).JSON())
	}
}

func downloadProcess(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		downloadQuery(w, r)
	case "POST":
		downloadStart(w, r)
	case "DELETE":
		downloadCancel(w, r)
	}
}

type downloadController struct {
	downloadOptions *pcscommand.DownloadOptions
	tasks           []*downloadTask
}

type downloadTask struct {
	ID     int64
	Path   []string
	Status int
}

func (c *downloadController) addDownloadTask(path string) error {
	//pcscommand.RunBgDownload([]string{path}, c.downloadOptions)
	go func() {
		currid++
		task := downloadTask{
			ID:     currid,
			Path:   []string{path},
			Status: statOngoing,
		}
		c.tasks = append(c.tasks, &task)
		pcscommand.RunDownload([]string{path}, c.downloadOptions)
		task.Status = statSuccess
	}()
	return nil
}

func (c *downloadController) getDownloadTask() []*downloadTask {
	return c.tasks
}
