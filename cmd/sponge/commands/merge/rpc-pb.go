package merge

import (
	"github.com/spf13/cobra"
)

// GRPCServiceCode merge the grpc service code
func GRPCServiceCode() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "rpc-pb",
		Short: "Merge the generated grpc related code into the template file",
		Long: `merge the generated grpc related code into the template file.

Examples:
  sponge merge rpc-pb
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir = adaptDir(dir)
			mergeGRPCECode(dir)
			mergeGRPCServiceTmpl(dir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")

	return cmd
}
