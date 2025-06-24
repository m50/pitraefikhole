package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/gookit/slog"
	"github.com/m50/traefik-pihole/pkg/pihole"
	"github.com/m50/traefik-pihole/pkg/traefik"
	"github.com/spf13/viper"
)

func process(ctx context.Context) {
	client := &http.Client{Timeout: 5 * time.Second}
	traefikClient := traefik.NewClient(client)
	piholeClient := pihole.NewClient(client)
	sleepPeriod := viper.GetInt64("poll-frequency-seconds")
	log := slog.WithContext(ctx)
	for {
		log.Debug("Checking for new hosts...")
		hosts, err := traefikClient.ListHosts(ctx)
		if err != nil {
			log.WithError(err).Error("failed to list hosts")
			time.Sleep(time.Duration(sleepPeriod) * time.Second)
			continue
		}
		log.Debug("Hosts found:", hosts)
		err = piholeClient.MergeHosts(ctx, hosts)
		if err != nil {
			log.WithError(err).Error("failed to merge hosts")
		}
		time.Sleep(time.Duration(sleepPeriod) * time.Second)
	}
}
