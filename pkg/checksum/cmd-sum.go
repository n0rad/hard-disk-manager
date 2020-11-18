package checksum

import (
	"fmt"
	"github.com/n0rad/go-checksum/pkg/checksum"
	"github.com/n0rad/hard-disk-manager/pkg/checksum/hashs"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func sumCommand() *cobra.Command {
	var hashString string

	cmd := &cobra.Command{
		Use:   "sum",
		Short: "Sum file(s) using hash",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h := hashs.MakeHashString(hashString)
			for _, arg := range args {
				fileSum, err := checksum.SumFilename(h, arg)
				if err != nil {
					fmt.Print(os.Args[0], ": ", err)
				}
				fmt.Println(fileSum)
				h.Reset()
			}
			return nil
		},
	}

	var hs strings.Builder
	for _, hash := range checksum.Hashs {
		hs.WriteString(string(hash))
		hs.WriteString(" ")
	}

	cmd.Flags().StringVar(&hashString, "hash", "sha1", "hash method among : "+hs.String())

	return cmd
}
