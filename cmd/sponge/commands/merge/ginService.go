package merge

import (
	"github.com/spf13/cobra"
)

// GinServiceCode merge the gin service code
func GinServiceCode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gin-service",
		Short: "Merge the gin service code",
		Long: `merge the gin service code.

Examples:
  # merge gin service code
  sponge merge gin-service
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return cmd
}
