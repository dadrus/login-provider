package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"login-provider/cmd/server"
	"login-provider/internal/config"
	"login-provider/internal/logging"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "login-provider",
	Short: "Hydra login provider",
	Long:  "Hydra login provider offering UI controls for OIDC login, consent and logout flows",
	Run:   server.Serve,
}

func init() {
	cobra.OnInitialize(config.Load(cfgFile), logging.ConfigureLogging)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $PWD/config.yaml)")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
