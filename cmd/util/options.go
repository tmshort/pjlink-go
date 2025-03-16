// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tmshort/pjlink-go/pkg/pjlink"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFileBase = "pjlink"
	defaultConfigFileType = "yaml"
)

// root represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pjlink",
	Short: "A command line utility to control projectors via the PJLink Protocol",
	Long: `A command line utility to control projectors via the PJLink Protocol
The default configuration file is located in $HOME/pjlink.yaml and is in YAML format.
The values in the YAML configuration match the names of the arguments.
Values may also be in the environment with the PJLINK prefix.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", fmt.Sprintf("config file (default is $HOME/%s.%s)", defaultConfigFileBase, defaultConfigFileType))
	rootCmd.PersistentFlags().String("projectorip", "", "IP address of the Projector")
	rootCmd.PersistentFlags().String("projectorport", "", "Port of the Projector")
	rootCmd.PersistentFlags().String("password", "", "Password of the Projector")
}

func setDefaults() {
	viper.SetDefault("config", fmt.Sprintf("%s.%s", defaultConfigFileBase, defaultConfigFileType))
	viper.SetDefault("projectorip", "127.0.0.1")
	viper.SetDefault("projectorport", "4352")
	viper.SetDefault("password", "password")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	setDefaults()
	//viper.SetOptions(viper.WithLogger(slog.Default()))
	viper.BindPFlags(rootCmd.Flags())
	viper.SetEnvPrefix("PJLINK")
	viper.AutomaticEnv() // read in environment variables that match

	// Grab the config file from either the flags or the environment
	config := viper.GetString("config")
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType(defaultConfigFileType)
	viper.SetConfigName(config)
	fmt.Println("ConfigFile", viper.ConfigFileUsed())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}

func createProjector() *pjlink.Projector {
	return pjlink.NewProjector(
		viper.GetString("projectorip"),
		viper.GetString("projectorport"),
		viper.GetString("password"),
	)
}

func printResponse(resp *pjlink.Response) {
	blob, _ := json.Marshal(resp)
	fmt.Println(string(blob))
}
