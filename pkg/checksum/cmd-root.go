package checksum

import (
	_ "github.com/n0rad/go-erlog/register"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/integrity"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	var config = &Config{}

	cmd := &cobra.Command{
		Use:   "checksum",
		Short: "Handle files checksum",
	}

	cmd.AddCommand(
		removeCommand(config),
		checkCommand(config),
		listCommand(config),
		setCommand(config),
		sumCommand(),
	)

	return cmd
}

func runCmdForPath(config *Config, path string, f func(d integrity.Directory) func(path string) error) error {
	directory := integrity.Directory{
		Regex:     config.regex,
		Inclusive: config.PatternIsInclusive,
		Strategy:  integrity.NewStrategy(config.Strategy, config.Hash),
	}

	return f(directory)(path)
}
