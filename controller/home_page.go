package controller

import (
	"encoding/json"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(allRooms)

	w.Write(json)
}
