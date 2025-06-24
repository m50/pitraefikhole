package sync

import (
	"context"

	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	slog.SetLevelByName(viper.GetString("log-level"))
	process(ctx)
}
