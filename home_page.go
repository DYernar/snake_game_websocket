package main

import (
	"encoding/json"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(freeRooms)
	w.Write(json)
}
