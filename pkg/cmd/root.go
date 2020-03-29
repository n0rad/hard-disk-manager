package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/n0rad/go-erlog/errs"
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
//hdm disks sync  			// sync to other disks
//hdm load
//hdm unload

func RootCommand(version string, buildTime string) (*cobra.Command, error) {
	var logLevel string
	var hdmHome string

	cmd := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if logLevel != "" {
				level, err := logs.ParseLevel(logLevel)
				if err != nil {
					return err
				}
				logs.SetLevel(level)
			}

			if err := hdm.HDM.Init(hdmHome); err != nil {
				return errs.WithE(err, "Failed to init hdm")
			}
			return nil
		},
	}

	path, err := homeDotConfigPath()
	if err != nil {
		return cmd, err
	}

	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "", "Set log level")
	cmd.PersistentFlags().StringVarP(&hdmHome, "home", "H", path+"/hdm", "configFile")

	cmd.AddCommand(prepareCommand())
	cmd.AddCommand(versionCommand(version, buildTime))
	cmd.AddCommand(passwordCommand())
	cmd.AddCommand(agentCommand())
	cmd.AddCommand(listCommand())
	cmd.AddCommand(addCommand())
	cmd.AddCommand(removeCommand())
	return cmd, nil
}

func homeDotConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", errs.WithE(err, "Failed to find user home folder")
	}
	return home + "/.config", nil
}
