// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"os"

	"github.com/kureuil/intranet-client/client"
	"github.com/spf13/cobra"
)

var gpaCity string
var gpaYear int
var gpaPromo string

// creditsCmd represents the credits command
var gpaCmd = &cobra.Command{
	Use:   "gpa",
	Short: "Fetch GPAs associated to students of a promotion",
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := client.IntranetClient{
			SessionID: cmd.Flag("sessionid").Value.String(),
		}
		students, err := apiClient.FetchPromotion(gpaCity, gpaYear, gpaPromo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return
		}

		channel := make(chan client.Student)
		for _, student := range students {
			go (func(channel chan client.Student, client client.IntranetClient, login string, year int) {
				stud, err := apiClient.FetchStudent(login)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
					return
				}
				channel <- stud
			})(channel, apiClient, student.Login, year)
		}

		fmt.Printf("login;gpa\n")
		for _ = range students {
			student := <-channel
			fmt.Printf("%s;%.2f\n", student.Login, student.GPABachelor)
		}
	},
}

func init() {
	RootCmd.AddCommand(gpaCmd)

	// Here you will define your flags and configuration settings.

	gpaCmd.Flags().StringVar(&gpaCity, "city", "", "Targeted city (e.g: REN)")
	gpaCmd.Flags().IntVar(&gpaYear, "year", 0, "Target year (e.g: 2016)")
	gpaCmd.Flags().StringVar(&gpaPromo, "promo", "", "Targeted promotion (e.g: tek3)")
}
