// Copyright © 2017 Louis Person <lait.kureuil@gmail.com>
//
// Licensed under the MIT License (the "License");
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

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Module struct {
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

type StudentGrades struct {
	Modules []Module `json:"modules"`
}

type Student struct {
	Login       string  `json:"login"`
	Fullname    string  `json:"title"`
	Credits     int     `json:"credits"`
	GPABachelor float64 `json:"gpa-bachelor"`
	GPAMaster   float64 `json:"gpa-master"`
}

type studentJSON struct {
	Login    string `json:"login"`
	Fullname string `json:"title"`
	Credits  int    `json:"credits"`
	GPA      []struct {
		GPA   string `json:"gpa"`
		Cycle string `json:"cycle"`
	} `json:"gpa"`
}

type IntranetClient struct {
	SessionID string
}

func (c IntranetClient) fetch(URL string, payload interface{}) error {
	authCookie := http.Cookie{
		Name:  "PHPSESSID",
		Value: c.SessionID,
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req.AddCookie(&authCookie)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}
	return nil
}

// Fetch all the students from a promotion, given a city a year.
func (c IntranetClient) FetchPromotion(city string, year int, promo string) ([]Student, error) {
	offset := 0
	total := 1
	students := make([]Student, 0, 128)
	for offset < total {
		URL := fmt.Sprintf("https://intra.epitech.eu/user/filter/user?format=json&location=FR/%s&year=%d&course=bachelor/classic|bachelor/tek2ed&active=true&promo=%s&offset=%d", city, year, promo, offset)
		studentList := struct {
			Items []Student `json:"items"`
			Total int       `json:"total"`
		}{}
		err := c.fetch(URL, &studentList)
		if err != nil {
			return []Student{}, err
		}
		students = append(students, studentList.Items...)
		offset += len(studentList.Items)
		total = studentList.Total
	}
	return students, nil
}

// Fetch a student, given its login (email address).
func (c IntranetClient) FetchStudent(login string) (Student, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/?format=json", login)
	stud := studentJSON{}
	err := c.fetch(URL, &stud)
	if err != nil {
		return Student{}, err
	}
	student := Student{
		Login:    stud.Login,
		Fullname: stud.Fullname,
		Credits:  stud.Credits,
	}
	for _, gpa := range stud.GPA {
		f, err := strconv.ParseFloat(gpa.GPA, 32)
		if err != nil {
			return Student{}, err
		}
		if gpa.Cycle == "bachelor" {
			student.GPABachelor = f
		} else if gpa.Cycle == "master" {
			student.GPAMaster = f
		}
	}
	return student, nil
}

// Fetch a student's grades, given the student login (email address)
func (c IntranetClient) FetchStudentGrades(login string) (StudentGrades, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/notes/?format=json", login)
	grades := StudentGrades{}
	err := c.fetch(URL, &grades)
	if err != nil {
		return StudentGrades{}, err
	}
	return grades, nil
}
