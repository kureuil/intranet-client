package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "USAGE: %s TOKEN CITY YEAR PROMOTION\n\n\tTOKEN\t\tSession ID from the intranet's session cookie\n\tCITY\t\tTargeted city (e.g: REN)\n\tYEAR\t\tTargeted year (e.g 2016)\n\tPROMOTION\tTargeted promotion (e.g: tek3)\n", os.Args[0])
		return
	}
	client := intranetClient{
		sessionID: os.Args[1],
	}
	year, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
	students, err := client.fetchPromotion(os.Args[2], year, os.Args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}
	fmt.Printf("login,actuels,engages\n")
	for _, student := range students {
		stud, err := client.fetchStudent(student.Login)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			return
		}
		grades, err := client.fetchStudentGrades(student.Login)
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
		fmt.Printf("%s,%d,%d\n", student.Login, stud.Credits, runningCredits)
	}
}
