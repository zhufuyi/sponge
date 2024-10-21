package patch

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

// DeleteJSONOmitemptyCommand delete json omitempty
func DeleteJSONOmitemptyCommand() *cobra.Command {
	var (
		dir        string
		suffixName string
	)

	cmd := &cobra.Command{
		Use:   "del-omitempty",
		Short: "Delete json tag omitempty",
		Long:  "Delete json tag omitempty.",
		Example: color.HiBlackString(`  # Delete all files that include the omitempty character
  sponge patch del-omitempty --dir=./api

  # Delete the specified suffix file including the omitempty character
  sponge patch del-omitempty --dir=./api --suffix-name=pb.go`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := replaceFiles(dir, suffixName)
			if err != nil {
				return err
			}

			fmt.Printf("delete the json tag omitempty was successful.\n")
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "input directory")
	_ = cmd.MarkFlagRequired("dir")
	cmd.Flags().StringVarP(&suffixName, "suffix-name", "s", "", "specified suffix file name, if empty it means all files")

	return cmd
}

func replaceFiles(dir string, suffixName string) error {
	opt := gofile.WithSuffix(suffixName)
	if suffixName == "" {
		opt = nil
	}
	files, err := gofile.ListFiles(dir, opt)
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		data = bytes.ReplaceAll(data, []byte(`,omitempty"`), []byte(`"`))
		err = os.WriteFile(file, data, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
