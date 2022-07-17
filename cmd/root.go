/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	rootCommandHeader = `

░█████╗░░██████╗░██╗░░░██╗██╗██╗░░░░░░█████╗░
██╔══██╗██╔═══██╗██║░░░██║██║██║░░░░░██╔══██╗
███████║██║██╗██║██║░░░██║██║██║░░░░░███████║
██╔══██║╚██████╔╝██║░░░██║██║██║░░░░░██╔══██║
██║░░██║░╚═██╔═╝░╚██████╔╝██║███████╗██║░░██║
╚═╝░░╚═╝░░░╚═╝░░░░╚═════╝░╚═╝╚══════╝╚═╝░░╚═╝

Manage docs seamlessly.
`
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aquila",
	Short: "A brief description of your application",
	Long:  rootCommandHeader,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aquila.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".aquila" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aquila")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
