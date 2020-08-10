package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	print(freeRooms)
	json, _ := json.Marshal(freeRooms)
	fmt.Println(json)
	w.Write(json)
}
