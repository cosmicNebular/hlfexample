package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Serve(c *Controller) {
	m := mux.NewRouter()
	m.HandleFunc("/create", c.create).Methods("POST")
	m.HandleFunc("/read/{passport}", c.read).Methods("GET")
	m.HandleFunc("/update", c.update).Methods("PUT")
	m.HandleFunc("/history/{passport}", c.history).Methods("GET")
	err := http.ListenAndServe(":80", m)
	if err != nil {
		log.Fatalln("There is problem during starting of rest endpoints")
	}
}
