package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"login-provider/cmd/server"
	"login-provider/internal/config"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "login-provider",
	Short: "Hydra login provider",
	Long:  "Hydra login provider offering UI controls for OIDC login, consent and logout flows",
	Run:   server.Serve,
}

func init() {
	var cfgFile string

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $PWD/config.yaml)")

	cobra.OnInitialize(config.Load(&cfgFile))
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
