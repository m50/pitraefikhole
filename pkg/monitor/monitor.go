package monitor

import (
	"context"

	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	process(ctx)
}
