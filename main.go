package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/project", project).Methods("GET")
	r.HandleFunc("/add-project", addProject).Methods("POST")
	r.HandleFunc("/detail/{index}", detail).Methods("GET")
	r.HandleFunc("/delete/{index}", delete).Methods("GET")
	r.HandleFunc("/edit/{index}", update).Methods("GET")
	r.HandleFunc("/edit-project/{index}", editProject).Methods("POST")

	fmt.Println("server on in port 8000")
	http.ListenAndServe("localhost:8000", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	card := map[string]interface{}{
		"Add": data,
	}

	tmpl.Execute(w, card)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/addProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

type Project struct {
	Name         string
	Start_date   string
	End_date     string
	Duration     string
	Desc         string
	Technologies []string
}

var data = []Project{}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")
	var technologies []string
	technologies = r.Form["technologies"]

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date )

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := daysInHours / 30
	yearInMonth := monthInDay / 12

	var duration string
	var month, _ float64 = math.Modf(monthInDay)
	var year, _ float64 = math.Modf(yearInMonth)

	if year > 0 {
		duration = strconv.FormatFloat(year, 'f', 0, 64) + " Years"
	} else if month > 0 {
		duration = strconv.FormatFloat(month, 'f', 0, 64) + " Months"
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
	} else {
		duration = "0 Days"

	}

	var newData = Project{
		Name:         name,
		Start_date:   start_date,
		End_date:     end_date,
		Duration:     duration,
		Desc:         desc,
		Technologies: technologies,
	}

	data = append(data, newData)
	// fmt.Println(data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Detail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)

	for i, data := range data {
		if index == i {
			Detail = Project{
				Name:         data.Name,
				Start_date:   data.Start_date,
				End_date:     data.End_date,
				Desc:         data.Desc,
				Duration:     data.Duration,
				Technologies: data.Technologies,
			}
		}
	}

	data := map[string]interface{}{
		"Details": Detail,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	data = append(data[:index], data[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")
	var technologies []string
	technologies = r.Form["technologies"]

	// fmt.Println(technologies)
	// fmt.Println(start_date)

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date)

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := daysInHours / 30
	yearInMonth := monthInDay / 12

	var duration string
	var month, _ float64 = math.Modf(monthInDay)
	var year, _ float64 = math.Modf(yearInMonth)

	if year > 0 {
		duration = strconv.FormatFloat(year, 'f', 0, 64) + " Years"
	} else if month > 0 {
		duration = strconv.FormatFloat(month, 'f', 0, 64) + " Months"
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
	} else {
		duration = "0 Days"
	}

	var newData = Project{
		Name:         name,
		Start_date:   start_date,
		End_date:     end_date,
		Duration:     duration,
		Desc:         desc,
		Technologies: technologies,
	}

	data[index] = newData
	// fmt.Println(data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Edit = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	fmt.Println(index)

	for i, data := range data {
		if index == i {
			Edit = Project{
				Name:         data.Name,
				Start_date:   data.Start_date,
				End_date:     data.End_date,
				Desc:         data.Desc,
				Technologies: data.Technologies,
			}
		}
	}

	data := map[string]interface{}{
		"Index": index,
		"Edit": Edit,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
