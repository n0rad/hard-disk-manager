package checksum

import (
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/integrity"
	"github.com/n0rad/hard-disk-manager/pkg/config"
	"github.com/spf13/cobra"
)

func RootCmd(conf *config.GlobalConfig) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "checksum",
		Short: "Handle files checksum",
	}

	cmd.AddCommand(
		removeCommand(conf),
		checkCommand(conf),
		listCommand(conf),
		setCommand(conf),
		sumCommand(),
	)

	return cmd
}

func runCmdForPath(config *config.GlobalConfig, path string, f func(pathConf config.PathConfig, d integrity.Directory) func(path string) error) error {
	pathConfig, err := config.GetPathConfig(path)
	if err != nil {
		return err
	}

	directory := integrity.Directory{
		Regex:     pathConfig.Checksum.GetRegex(),
		Exclusive: pathConfig.Checksum.PatternIsExclusive,
		Strategy:  integrity.NewStrategy(pathConfig.Checksum.Strategy, pathConfig.Checksum.Hash),
	}

	return f(pathConfig, directory)(path)
}
