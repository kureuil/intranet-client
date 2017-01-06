// Copyright © 2017 Louis Person <lait.kureuil@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mit-license.org/
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

var creditsCity string
var creditsYear int
var creditsPromo string

type studentCredits struct {
	login          string
	credits        int
	runningCredits int
}

func fetchStudentCredits(channel chan studentCredits, client client.IntranetClient, login string, year int) {
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
		if module.ScolarYear != year || module.Grade != "-" {
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
	Run:   creditsCmdRun,
}

func creditsCmdRun(cmd *cobra.Command, args []string) {
	apiClient := client.IntranetClient{
		SessionID: cmd.Flag("sessionid").Value.String(),
	}
	students, err := apiClient.FetchPromotion(creditsCity, creditsYear, creditsPromo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
	fmt.Println("login;actuels;engages")
	channel := make(chan studentCredits)
	for _, student := range students {
		student := student
		go fetchStudentCredits(channel, apiClient, student.Login, creditsYear)
	}
	for range students {
		credits := <-channel
		fmt.Printf("%s;%d;%d\n", credits.login, credits.credits, credits.runningCredits)
	}
}

func init() {
	RootCmd.AddCommand(creditsCmd)
	creditsCmd.Flags().StringVar(&creditsCity, "city", "", "Targeted city (e.g: REN)")
	creditsCmd.Flags().IntVar(&creditsYear, "year", 0, "Target year (e.g: 2016)")
	creditsCmd.Flags().StringVar(&creditsPromo, "promo", "", "Targeted promotion (e.g: tek3)")
}
