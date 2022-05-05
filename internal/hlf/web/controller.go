package web

import (
	"awesomeProject/internal/hlf/chaincode"
	"awesomeProject/internal/hlf/setup"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Controller struct {
	setup setup.FabricSetup
}

func NewController(setup setup.FabricSetup) *Controller {
	return &Controller{setup: setup}
}

func (c *Controller) create(w http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	r := new(chaincode.Record)
	err := d.Decode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	id, err := c.setup.Create(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) read(w http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)["passport"]
	r, err := c.setup.Read(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	rj, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(rj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) update(w http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	r := new(struct {
		Phone string `json:"phone"`
		Field string `json:"field"`
		Value string `json:"value"`
	})
	err := d.Decode(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	id, err := c.setup.Update(r.Phone, r.Field, r.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) history(w http.ResponseWriter, req *http.Request) {
	p := mux.Vars(req)["passport"]
	h, err := c.setup.History(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(h))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
