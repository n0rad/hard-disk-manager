package hdm

import (
	"github.com/ghodss/yaml"
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"io/ioutil"
	"os"
	"sync"
)

type BackupDB struct {
	path    string
	lock    sync.RWMutex
	servers Servers
}

func (db *BackupDB) Init(path string, servers Servers) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errs.WithEF(err, data.WithField("path", path), "Failed to create db path")
	}

	db.servers = servers
	db.path = path
	return nil
}

func (db *BackupDB) GetOrCreateBackup(diskLabel string, path string) (Backup, error) {
	file := db.path + "/" + diskLabel + ".yaml"
	_, err := os.Stat(file)
	if err == nil {
		backups, err := db.readFile(diskLabel)
		if err != nil {
			return Backup{}, errs.WithE(err, "Failed to load read backup db file")
		}

		for _, backup := range backups {
			if backup.SourcePath == path {
				return backup, nil
			}
		}
	} else if os.IsNotExist(err) {
		// not in db, preparing
		backup := Backup{
			SourceDiskLabel: diskLabel,
			SourcePath:      path,
		}
		//return backup, backup.Init("", block, db.servers)
		return backup, nil
	}
	return Backup{}, errs.WithEF(err, data.WithField("file", file), "Failed to read disk db file")
}

func (db *BackupDB) SaveBackup(backup Backup) error {
	backups, err := db.readFile(backup.SourceDiskLabel)
	if err != nil {
		return errs.WithE(err, "Failed to load read backup db file")
	}

	found := false
	for i := range backups {
		if backups[i].SourcePath != backup.SourcePath {
			continue
		}

		backups[i] = backup
		found = true
	}

	if !found {
		backups = append(backups, backup)
	}

	return db.WriteFile(backups, backup.SourceDiskLabel)
}

/////////////////////////////////////////////

func (db *BackupDB) readFile(label string) ([]Backup, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var backups []Backup
	filePath := db.path + "/" + label + ".yaml"
	bytes, err := ioutil.ReadFile(filePath)
	fileField := data.WithField("file", filePath)
	if err != nil {
		return backups, errs.WithEF(err, fileField, "Failed to read disk db file")
	}

	if err := yaml.Unmarshal(bytes, &backups); err != nil {
		return backups, errs.WithE(err, "Failed to marshal blockDevice")
	}

	//for i := range backups {
	//	if err := backups[i].Init(); err != nil {
	//		return backups, errs.WithEF(err, data.WithField("file", filePath), "Failed to init backup")
	//	}
	//}

	return backups, nil
}

func (db *BackupDB) WriteFile(backups []Backup, label string) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	backupsBytes, err := yaml.Marshal(backups)
	if err != nil {
		return errs.WithE(err, "Failed to marshal backups")
	}

	filePath := db.path + "/" + label + ".yaml"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errs.WithEF(err, data.WithField("file", filePath), "Failed to open backup file")
	}
	defer file.Close()

	if _, err := file.Write(backupsBytes); err != nil {
		return errs.WithEF(err, data.WithField("path", filePath), "Failed to write backup file")
	}
	return nil
}
