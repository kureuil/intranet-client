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

var city string
var year int
var promo string

type studentCredits struct {
	login          string
	credits        int
	runningCredits int
}

func fetchStudentCredits(channel chan studentCredits, client client.IntranetClient, login string) {
	stud, err := client.FetchStudent(login)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
	grades, err := client.FetchStudentGrades(login)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
	runningCredits := 0
	for _, module := range grades.Modules {
		if module.Grade != "-" {
			continue
		}
		runningCredits += module.Credits
	}
	credits := studentCredits{
		login:          login,
		credits:        stud.Credits,
		runningCredits: runningCredits,
	}
	channel <- credits
}

// creditsCmd represents the credits command
var creditsCmd = &cobra.Command{
	Use:   "credits",
	Short: "Fetch credits associated to students of a promotion",
	Run: func(cmd *cobra.Command, args []string) {
		client := client.IntranetClient{
			SessionID: cmd.Flag("sessionid").Value.String(),
		}
		students, err := client.FetchPromotion(city, year, promo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return
		}
		fmt.Printf("login;actuels;engages\n")
		channel := make(chan studentCredits)
		for _, student := range students {
			go fetchStudentCredits(channel, client, student.Login)
		}
		for _ = range students {
			credits := <-channel
			fmt.Printf("%s;%d;%d\n", credits.login, credits.credits, credits.runningCredits)
		}
	},
}

func init() {
	RootCmd.AddCommand(creditsCmd)

	// Here you will define your flags and configuration settings.

	creditsCmd.Flags().StringVar(&city, "city", "", "Targeted city (e.g: REN)")
	creditsCmd.Flags().IntVar(&year, "year", 0, "Target year (e.g: 2016)")
	creditsCmd.Flags().StringVar(&promo, "promo", "", "Targeted promotion (e.g: tek3)")
}
