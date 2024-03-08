package main

import (
	"fmt"
	"log"

	"net/http"
	"uts/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rooms", controllers.GetAllRooms).Methods("GET")
	router.HandleFunc("/roomsDetail/{id}", controllers.GetDetailRooms).Methods("GET")
	router.HandleFunc("/rooms/{id}", controllers.InsertRoom).Methods("POST")
	router.HandleFunc("/rooms/{id}", controllers.LeaveRoom).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8980")
	log.Println("Connected to port 8980")
	log.Fatal(http.ListenAndServe(":8980", router))
}
