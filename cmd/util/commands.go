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
	"fmt"
	"os"

	"github.com/tmshort/pjlink-go/pkg/pjlink"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// statusCmd represents the status command
	var configCmd = &cobra.Command{
		Use:   "writeconfig",
		Short: "Display current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Writing config file")
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, "error in WriteConfig", err)
			}
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display status of Projector",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stderr, "Status of projector... \n")

			stat, err := createProjector().GetPowerStatus()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err.Error())
			} else {
				printResponse(stat)
			}
		},
	}

	var powerOnOff = &cobra.Command{
		Use:   "power <on/off>",
		Short: "Power projector on or off",
		Run: func(cmd *cobra.Command, args []string) {
			proj := createProjector()

			if len(args) == 0 || (args[0] != pjlink.ON && args[0] != pjlink.OFF) {
				fmt.Fprintf(os.Stderr, "must specify action: 'on' or 'off'")
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Powering %s projector...\n", args[0])

			if args[0] == pjlink.ON {
				err := proj.PowerOn()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err.Error())
				}
			}

			if args[0] == pjlink.OFF {
				err := proj.PowerOff()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s", err.Error())
				}
			}

			stat, err := proj.GetPowerStatus()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err.Error())
			} else {
				printResponse(stat)
			}
		},
	}

	rootCmd.AddCommand(configCmd, statusCmd, powerOnOff)
}
