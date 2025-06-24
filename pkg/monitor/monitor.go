package monitor

import (
	"context"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	level := slog.LevelInfo
	cobra.CheckErr(level.UnmarshalText([]byte(strings.ToUpper(viper.GetString("log-level")))))
	slog.SetLogLoggerLevel(level)
	process(ctx)
}
