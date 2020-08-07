package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"snake_game/controller"

	"github.com/rs/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Go websockets")
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.HomePage)
	mux.HandleFunc("/ws", controller.Websocket)

	handler := cors.Default().Handler(mux)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
