package integrity

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type Directory struct {
	Regex     *regexp.Regexp
	Exclusive bool
	Strategy  Strategy

	timers      map[string]*time.Timer
	timersMutex *sync.Mutex
}

type DirectoryAction func(p []byte) (n int, err error)

func (d Directory) List(path string) error {
	return d.directoryWalk(path, func(path string, info os.FileInfo) (error, error) {
		println(path)
		return nil, nil
	})
}

func (d Directory) Check(path string) error {
	return d.directoryWalk(path, func(path string, info os.FileInfo) (error, error) {
		set, err := d.Strategy.IsSet(path)
		if err != nil {
			return nil, errs.WithE(err, "Failed to check if sum is set")
		}
		if !set {
			logs.WithField("path", path).Warn("Missing sum")
			return nil, nil
		}

		logs.WithField("path", path).Info("Processing file")
		ok, err := d.Strategy.Check(path)
		if err != nil {
			return nil, errs.WithEF(err, data.WithField("path", path), "Failed to check file integrity")
		}
		if ok != nil {
			return errs.WithEF(ok, data.WithField("path", path), "File integrity failed"), nil
		}
		return nil, nil
	})
}

func (d Directory) Set(path string) error {
	return d.directoryWalk(path, func(path string, info os.FileInfo) (error, error) {
		if d.Strategy.IsSumFile(path) {
			return nil, nil
		}

		set, err := d.Strategy.IsSet(path)
		if err != nil {
			return nil, errs.WithE(err, "Failed to check if sum is set")
		}

		if !set {
			logs.WithField("path", path).Info("Processing file")
			if _, err := d.Strategy.SumAndSet(path); err != nil {
				return errs.WithE(err, "Failed to set sum"), nil
			}
		} else {
			logs.WithField("path", path).Debug("Sum already exists")
		}
		return nil, nil
	})
}

func (d Directory) Remove(path string) error {
	return d.directoryWalk(path, func(path string, info os.FileInfo) (error, error) {
		if d.Strategy.IsSumFile(path) {
			return nil, nil
		}

		set, err := d.Strategy.IsSet(path)
		if err != nil {
			return nil, errs.WithE(err, "Failed to check if sum is set")
		}

		if set {
			logs.WithField("path", path).Info("Processing file")
			if err := d.Strategy.Remove(path); err != nil {
				return nil, errs.WithEF(err, data.WithField("path", path), "Failed to remove integrity")
			}
		}
		return nil, nil
	})
}

func (d Directory) Watch(path string) error {
	d.timers = map[string]*time.Timer{}
	d.timersMutex = &sync.Mutex{}

	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if info.Mode().IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		return errs.WithE(err, "Failed to watch directory")
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if err := d.processEvent(event, watcher); err != nil {
					logs.WithE(err).Error("Failed to process event")
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
	<-done
	return nil
}

func (d Directory) processEvent(event fsnotify.Event, watcher *fsnotify.Watcher) error {
	logs.WithField("event", event).Trace("received fs event")
	if !d.matchesPattern(event.Name) {
		return nil
	}
	if d.Strategy.IsSumFile(event.Name) {
		return nil
	}

	switch event.Op {
	case fsnotify.Create:
		stat, err := os.Stat(event.Name)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			if err := watcher.Add(event.Name); err != nil {
				return errs.WithE(err, "Failed to watch directory")
			}
		} else {
			logs.WithField("file", event.Name).Info("Calculate sum of new file")
			if _, err := d.Strategy.SumAndSet(event.Name); err != nil {
				return errs.WithE(err, "Failed to sum new file")
			}
		}
	case fsnotify.Write:
		d.timersMutex.Lock()
		if timer, ok := d.timers[event.Name]; ok {
			timer.Stop()
		}
		d.timers[event.Name] = time.AfterFunc(1*time.Second, func() {
			logs.WithField("file", event.Name).Info("Recalculate sum after write")
			if _, err := d.Strategy.SumAndSet(event.Name); err != nil {
				logs.WithE(err).Error("Failed to sum new file")
			}

			d.timersMutex.Lock()
			defer d.timersMutex.Unlock()
			delete(d.timers, event.Name)
		})
		d.timersMutex.Unlock()
	case fsnotify.Remove:
		logs.WithField("file", event.Name).Info("Removing sum of deleted file")
		if err := d.Strategy.Remove(event.Name); err != nil {
			return errs.WithE(err, "Failed to remove sum file")
		}
	case fsnotify.Rename:
		logs.WithField("file", event.Name).Info("Removing sum of renamed file")
		if err := d.Strategy.Remove(event.Name); err != nil {
			return errs.WithE(err, "Failed to remove sum file")
		}
	case fsnotify.Chmod:
	}

	return nil
}

////////////////////

func (d Directory) directoryWalk(path string, f func(path string, info os.FileInfo) (error, error)) error {
	var fail bool
	if err := filepath.Walk(path, func(path string, info os.FileInfo, errIn error) error {
		if d.Strategy.IsSumFile(path) {
			return nil
		}

		if errIn != nil {
			logs.WithE(errIn).WithField("path", path).Error("Failed to process path")
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if d.matchesPattern(path) {
			logs.WithField("path", path).Debug("Processing file")
			ok, err := f(path, info)
			if err != nil {
				return err
			}
			if ok != nil {
				logs.WithError(ok).Error("File failed")
				fail = true
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if fail {
		return errs.With("Some file failed")
	}
	return nil
}

func (d Directory) matchesPattern(path string) bool {
	return !d.Exclusive && d.Regex.MatchString(path) ||
		d.Exclusive && !d.Regex.MatchString(path)
}
