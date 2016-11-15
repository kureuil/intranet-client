package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (c intranetClient) fetch(URL string, payload interface{}) error {
	authCookie := http.Cookie{
		Name:  "PHPSESSID",
		Value: c.sessionID,
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

func (c intranetClient) fetchPromotion(city string, year int, promo string) ([]student, error) {
	offset := 0
	total := 1
	students := make([]student, 0, 128)
	for offset < total {
		URL := fmt.Sprintf("https://intra.epitech.eu/user/filter/user?format=json&location=FR/%s&year=%d&course=bachelor/classic|bachelor/tek2ed&active=true&promo=%s&offset=%d", city, year, promo, offset)
		studentList := struct {
			Items []student `json:"items"`
			Total int       `json:"total"`
		}{}
		err := c.fetch(URL, &studentList)
		if err != nil {
			return []student{}, err
		}
		students = append(students, studentList.Items...)
		offset += len(studentList.Items)
		total = studentList.Total
	}
	return students, nil
}

func (c intranetClient) fetchStudent(login string) (student, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/?format=json", login)
	stud := student{}
	err := c.fetch(URL, &stud)
	if err != nil {
		return student{}, err
	}
	return stud, nil
}

func (c intranetClient) fetchStudentGrades(login string) (studentGrades, error) {
	URL := fmt.Sprintf("https://intra.epitech.eu/user/%s/notes/?format=json", login)
	grades := studentGrades{}
	err := c.fetch(URL, &grades)
	if err != nil {
		return studentGrades{}, err
	}
	return grades, nil
}
