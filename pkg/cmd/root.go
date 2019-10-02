package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

//# compare all disks with unencrypted & mounted
//# run unencrypt
//# mount to /mnt
//# rebuild Merged
//# restart containers ?

// search for failing disks
// sync disks across servers

// diff of information from previous
// logs of events (diff)

//hdm scan // scan servers for disks
//hdm list // list disks
//hdm

//hdm disks scan			// get list of disks
//hdm logs

//hdm destroy
//hdm index		// index files from disks
//hdm search    	// search for a file on all disks
//hdm restore     // restore a backup file
//hdm backup		// scan files, for backup order and run backup
//hdm backupable  // check that backup orders can work (target disk is plugged, size of directory match disk size)

//hdm scan				// get information from the disks
//hdm repair			// scan for bad blocks, pending sectors, find affected files, etc...
//hdm location			// get location of disks

//hdm check					// check disks are prepared, mounted, repaired and backupable, backup up to date (period)
//hdm agent					// start an agent that ????????????????????????
//hdm disks sync  			// sync to other disks
//hdm load
//hdm unload

func RootCommand(Version string, BuildTime string) *cobra.Command {
	var logLevel string
	var hdmHome string

	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					logs.WithField("value", logLevel).Fatal("Unknown log level")
				}
				logs.SetLevel(level)
			} else {
				//logs.SetLevel(logs.WARN)
			}

			if err := hdm.HDM.Init(hdmHome); err != nil {
				logs.WithE(err).Fatal("Cannot start, failed to load configuration")
			}
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Display HDM version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hdm")
			fmt.Println("version : ", Version)
			fmt.Println("Build Time : ", BuildTime)
		},
	})

	// main
	cmd.AddCommand(
		command("agent", []string{}, Agent,
			"Run an agent that self handle disks"),
		passwordCmd(),
	)

	// state
	cmd.AddCommand(
		//command("list", []string{"ls"}, List,
		//	"List known disks (even unplugged)"),
		commandWithDiskSelector("index", []string{}, Index,
			"Index files from disks"),
		commandWithDiskSelector("location", []string{}, Location,
			"Get disk location"),
	)

	// cycle
	cmd.AddCommand(
		commandWithDiskSelector("add", []string{}, Add,
			"AddBlockDevice disks as usable (mdadm,crypt,mount,restart)"),
		commandWithRequiredDiskSelector("remove", []string{}, Remove,
			"Remove or cleanup removed disk (kill,umount,restart,mdadm,crypt)"),
		commandWithRequiredServerDiskAndLabel("prepare", []string{}, Prepare,
			"Make disks usable(luksOpen,mount,restart)"),
		commandWithRequiredServerDiskAndLabel("erase", []string{}, Erase,
			"securely erase disk"),
	)

	// heal

	// backup
	cmd.AddCommand(
		commandWithRequiredDiskSelector("backupable", []string{}, Backupable,
			"Find backup configs and run backups"),
		commandWithRequiredDiskSelector("backup", []string{}, BackupCmd,
			"Find backup configs and run backups"),
		command("backups", []string{}, Backups,
			"Find backup configs and run backups"),
	)

	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")
	cmd.PersistentFlags().StringVarP(&hdmHome, "home", "H", homeDotConfigPath()+"/hdm", "configFile")
	return cmd
}

func homeDotConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		logs.WithError(err).Fatal("Failed to find user home folder")
	}
	return home + "/.config"
}
