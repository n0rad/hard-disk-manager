package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/logs"
	hdm2 "github.com/n0rad/hard-drive-manager/pkg/hdm"
	"github.com/spf13/cobra"
	"os"
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

var hdm = hdm2.Hdm{}

func RootCommand(Version string, BuildTime string) *cobra.Command {
	var logLevel string
	var version bool
	var hdmConfig string

	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println("HDM")
				fmt.Println("Version :", Version)
				fmt.Println("Build Time :", BuildTime)
				os.Exit(0)
			}

			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					logs.WithField("value", logLevel).Fatal("Unknown log level")
				}
				logs.SetLevel(level)
			} else {
				//logs.SetLevel(logs.WARN)
			}

			if err := hdm.InitFromFile(hdmConfig); err != nil {
				logs.WithE(err).Fatal("Cannot start, failed to load configuration")
			}
		},
	}

	cmd.AddCommand(
		command("agent", []string{"ls"}, hdm.Agent,
			"Run an agent that self handle disks"),
		command("list", []string{"ls"}, hdm.List,
			"List known disks (even unplugged)"),
		commandWithDiskSelector("index", []string{}, hdm.Index,
			"Index files from disks"),
		commandWithDiskSelector("add", []string{}, hdm.Add,
			"Add disks as usable (mdadm,crypt,mount,restart)"),
		commandWithRequiredDiskSelector("backupable", []string{}, hdm.Backupable,
			"Find backup configs and run backups"),
		commandWithRequiredDiskSelector("backup", []string{}, hdm.Backup,
			"Find backup configs and run backups"),
		commandWithDiskSelector("location", []string{}, hdm.Location,
			"Get disk location"),
		commandWithRequiredDiskSelector("remove", []string{}, hdm.Remove,
			"Remove or cleanup removed disk (kill,umount,restart,mdadm,crypt)"),
		commandWithRequiredServerDiskAndLabel("prepare", []string{}, hdm.Prepare,
			"Make disks usable(luksOpen,mount,restart)"),
		)

	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")
	cmd.PersistentFlags().BoolVarP(&version, "version", "V", false, "Display version")
	cmd.PersistentFlags().StringVarP(&hdmConfig, "config", "C", homeDotConfigPath()+"/hdm/config.yaml", "configFile")

	return cmd
}

func homeDotConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		logs.WithError(err).Fatal("Failed to find user home folder")
	}
	return home + "/.config"
}
