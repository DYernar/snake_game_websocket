package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Go websockets")
	mux := http.NewServeMux()

	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/ws", Websocket)

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
