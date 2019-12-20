package cmd

import (
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/hdm"
	"github.com/spf13/cobra"
)

func errorLoggerWrap(f func(cmd *cobra.Command, args []string) error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			logs.WithE(err).Fatal("Command failed")
		}
	}
}

//func command(use string, aliases []string, hdmCommand func() error, short string) *cobra.Command {
//	return &cobra.Command{
//		Use:     use,
//		Aliases: aliases,
//		Short:   short,
//		Run: func(cmd *cobra.Command, args []string) {
//			if err := hdmCommand(); err != nil {
//				logs.WithE(err).Fatal("Command failed")
//			}
//		},
//	}
//}


//func withDiskSelector(use string, aliases []string, hdmCommand func(selector hdm.DisksSelector) error, short string) (*cobra.Command, *hdm.DisksSelector) {
//	selector := hdm.DisksSelector{}
//	cmd := &cobra.Command{
//		Use:     use,
//		Aliases: aliases,
//		Short:   short,
//		Run: func(cmd *cobra.Command, args []string) {
//			if err := hdmCommand(selector); err != nil {
//				logs.WithE(err).Fatal("Command failed")
//			}
//		},
//	}
//	cmd.Flags().StringVarP(&selector.Server, "server", "s", "", "Server")
//	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")
//	cmd.Flags().StringVarP(&selector.Label, "label", "l", "", "Label")
//	return cmd, &selector
//}


func withDiskSelector(selector *hdm.DisksSelector, cmd *cobra.Command) {
	cmd.Flags().StringVarP(&selector.Server, "server", "s", "", "Server")
	cmd.Flags().StringVarP(&selector.Disk, "disk", "d", "", "Disk")
	cmd.Flags().StringVarP(&selector.Label, "label", "l", "", "Label")
}

//func commandWithDiskSelector(use string, aliases []string, hdmCommand func(selector hdm.DisksSelector) error, short string) *cobra.Command {
//	cmd, _ := withDiskSelector(use, aliases, hdmCommand, short)
//	return cmd
//}
//
//func commandWithRequiredServerDiskAndLabel(use string, aliases []string, hdmCommand func(selector hdm.DisksSelector) error, short string) *cobra.Command {
//	cmd, _ := withDiskSelector(use, aliases, hdmCommand, short)
//	_ = cmd.MarkFlagRequired("server")
//	_ = cmd.MarkFlagRequired("disk")
//	_ = cmd.MarkFlagRequired("label")
//	return cmd
//}
//
//func commandWithRequiredDiskSelector(use string, aliases []string, hdmCommand func(selector hdm.DisksSelector) error, short string) *cobra.Command {
//	cmd, selector := withDiskSelector(use, aliases, hdmCommand, short)
//	realRun := cmd.Run
//	cmd.Run = func(cmd *cobra.Command, args []string) {
//		if err := selector.IsValid(); err != nil {
//			_, _ = fmt.Fprintln(cmd.OutOrStderr(), err.Error())
//			_ = cmd.Usage()
//			os.Exit(1)
//		}
//		realRun(cmd, args)
//	}
//	_ = cmd.MarkFlagRequired("server")
//	return cmd
//}
