package controllers

import (
	"Api-mongo/models"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserController struct {
	Session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}

	oid := bson.ObjectIdHex(id)

	u := models.User{}

	if err := uc.Session.DB("mongo-golang").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{}

	json.NewDecoder(r.Body).Decode(&u)

	u.Id = bson.NewObjectId()

	uc.Session.DB("mongo-golang").C("users").Insert(u)

	uj, err := json.Marshal(u)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(uj)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	if err := uc.Session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", oid)
}
