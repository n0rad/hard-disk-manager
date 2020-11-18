package backup

type Backup struct {
	// source
	SourceStrategy string // files, folder, leafFolder ?
	StartAtChild   bool
	// ToIgnore pattern ?

	// how
	ScanInterval string // cron / watch
	WatchChanges bool   //

	// target
	//append only
	//Versioning string // trash, periodic, ...
	Count      string // how many backup
}

type Disc struct {
	Usage string // backup, data
	//Type string // hdd, ssd
}
