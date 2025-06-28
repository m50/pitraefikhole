package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/m50/traefik-pihole/pkg/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "sync",
	Short: "",
	Long:  "",
	Run:   sync.Run,
}

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	flags := rootCmd.PersistentFlags()
	flags.StringVar(&cfgFile, "config", "", "config file (default: './config.yml' or '/config.yml')")
	flags.StringP("cname-address", "c", "", "the DNS record to create the cname to (REQUIRED)")
	flags.StringP("pihole-password", "p", "", "the password of the pihole (REQUIRED)")
	flags.String("pihole-password-file", "", "the password of the pihole to be read from a file (REQUIRED)")
	flags.String("pihole-address", "http://pihole/", "the address of the pihole for the API (default: 'http://pihole/')")
	flags.String("traefik-address", "http://traefik:8080/", "the address of the traefik instance to watch for routes (default: 'http://traefik:8080/')")
	flags.Int64("poll-frequency-seconds", 30, "the number of seconds to check between each poll of traefik to look for new hosts (default: 30)")
	flags.String("log-level", "INFO", "log level (default: 'INFO')")
	cobra.CheckErr(viper.BindPFlags(flags))
}

func initConfig() {
	viper.SetEnvPrefix("PITRAEFIKHOLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		workingDir, err := os.Getwd()
		if err == nil {
			viper.AddConfigPath(workingDir)
		}
		viper.AddConfigPath("/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	ppf := viper.GetString("pihole-password-file")
	if ppf == "" {
		return
	}
	piholePassword, err := os.ReadFile(ppf)
	if err != nil {
		fmt.Println("Can't read pihole password:", err)
		os.Exit(1)
	}
	viper.Set("pihole-password", piholePassword)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
