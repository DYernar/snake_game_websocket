package main

import (
	"net/http"
	"text/template"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	var homeTempl = template.Must(template.ParseFiles("templates/home.html"))
	data := struct {
		Host       string
		RoomsCount int
	}{r.Host, roomsCount}
	homeTempl.Execute(w, data)
}
