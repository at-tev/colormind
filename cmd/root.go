// Package cmd consists commands for COLORMIND.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	schemeURL = "http://colormind.io/api/"
	modelsURL = "http://colormind.io/list/"
)

var rootCmd = &cobra.Command{
	Use:   "cs",
	Short: "This is cli app for COLORMIND web application",
	Long:  `This is cli app that will help you to create color schemes by COLORMIND web application`,
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	models.lastUpdate = viper.GetTime("models.last_update")
	models.List = viper.GetStringSlice("models.list")
	if !models.updated() {
		Models()
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
