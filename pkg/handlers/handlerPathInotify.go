package handlers

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

func init() {
	handlers = append(handlers, handler{
		HandlerFilter{Type: "path"},
		func() Handler {
			return &HandlerBackup{
				CommonHandler: CommonHandler{
					handlerName: "handlerInotify",
				},
			}
		},
	})
}

type HandlerInotify struct {
	CommonHandler
}

func (h *HandlerInotify) Init() {

}

func (h *HandlerInotify) Start() {

}


var watcher *fsnotify.Watcher

// main
func sss() {

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk("/Users/skdomino/Desktop/test", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
