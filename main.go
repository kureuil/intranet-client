package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type module struct {
	ScolarYear    int    `json:"scolaryear"`
	IDUserHistory string `json:"id_user_history"`
	CodeModule    string `json:"codemodule"`
	CodeInstance  string `json:"codeinstance"`
	Title         string `json:"title"`
	DateIns       string `json:"date_ins"`
	Cycle         string `json:"cycle"`
	Grade         string `json:"grade"`
	Credits       int    `json:"credits"`
	Barrage       int    `json:"barrage"`
}

type studentGrades struct {
	Modules []module `json:"modules"`
}

type student struct {
	Login    string `json:"login"`
	Fullname string `json:"title"`
	Credits  int    `json:"credits"`
}

type intranetClient struct {
	sessionID string
}

func (c intranetClient) fetch(URL string) (*http.Response, error) {
	authCookie := http.Cookie{
		Name:  "PHPSESSID",
		Value: c.sessionID,
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&authCookie)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c intranetClient) fetchPromotion(city string, year int, promo string) ([]student, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/filter/user?format=json&location=FR/%s&year=%d&course=bachelor/classic&active=true&promo=%s&offset=0", city, year, promo)
	resp, err := c.fetch(URL)
	if err != nil {
		return []student{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []student{}, err
	}
	studentList := struct {
		Items []student `json:"items"`
		Total int       `json:"total"`
	}{}
	if err := json.Unmarshal(body, &studentList); err != nil {
		return []student{}, err
	}
	return studentList.Items, nil
}

func (c intranetClient) fetchStudent(login string) (student, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/?format=json", login)
	resp, err := c.fetch(URL)
	if err != nil {
		return student{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return student{}, err
	}
	stud := student{}
	if err := json.Unmarshal(body, &stud); err != nil {
		return student{}, err
	}
	return stud, nil
}

func (c intranetClient) fetchStudentGrades(login string) (studentGrades, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/notes/?format=json", login)
	resp, err := c.fetch(URL)
	if err != nil {
		return studentGrades{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return studentGrades{}, err
	}
	grades := studentGrades{}
	if err := json.Unmarshal(body, &grades); err != nil {
		return studentGrades{}, err
	}
	return grades, nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "USAGE: %s TOKEN CITY YEAR PROMOTION\n", os.Args[0])
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
