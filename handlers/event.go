package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spectrum-team/tachartas/commons"
	"github.com/spectrum-team/tachartas/interfaces"
	"github.com/spectrum-team/tachartas/models"
	"github.com/spectrum-team/tachartas/repos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventHandler struct {
	dbConfig  *models.DatabaseConfig
	eventRepo interfaces.EventRepository
	userRepo  interfaces.UserRepository
}

func NewEventHandler(dbConfig *models.DatabaseConfig) *EventHandler {
	return &EventHandler{
		dbConfig:  dbConfig,
		eventRepo: repos.NewEventRepository(dbConfig),
		userRepo:  repos.NewUserRepository(dbConfig),
	}
}

func (e *EventHandler) FindOne(w http.ResponseWriter, r *http.Request) {

	idParam := mux.Vars(r)["id"]

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("There was an error parsing id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userEmail := commons.GetAuthCtx(r.Context())

	event, err := e.eventRepo.FindByID(userEmail, id)
	if err != nil {
		log.Println("Error looking by ID: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&event)
}

func (e *EventHandler) Find(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inception", r.Header)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := bson.M{}
	err = json.Unmarshal(body, &query)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	events, err := e.eventRepo.Find(query)
	if err != nil {
		log.Println("Error looking events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("keloke", w.Header())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&events)
}

func (e *EventHandler) Insert(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &models.Event{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Println("klok wawawa => ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := e.eventRepo.Insert(event)
	if err != nil {
		log.Println("Error Inserting events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&id)
}

func (e *EventHandler) Update(w http.ResponseWriter, r *http.Request) {

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("There was an error parsing id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := &models.Event{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event.ID = id
	err = e.eventRepo.Update(id, event)
	if err != nil {
		log.Println("Error updating events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (e *EventHandler) Assist(w http.ResponseWriter, r *http.Request) {

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	willAssist, err := strconv.Atoi(mux.Vars(r)["assist"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("There was an error parsing id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add it to the user too
	event, _ := e.eventRepo.FindByID(id)

	userEmail := commons.GetAuthCtx(r.Context())
	changed, err := e.userRepo.AssistToEvent(userEmail, event, willAssist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !changed {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = e.eventRepo.Assist(id, willAssist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (e *EventHandler) AddImageToEvent(w http.ResponseWriter, r *http.Request) {

	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("filename")
	if err != nil {
		log.Println("There was an error reading file from multipart form: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("There was an error reading file into memory: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Println("There was an error parsing id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = e.eventRepo.AddEventImage(id, handler.Filename, data)
	if err != nil {
		log.Println("there was an error saving image: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (e *EventHandler) FindHotEvents(w http.ResponseWriter, r *http.Request) {

	events, err := e.eventRepo.FindHotEvents()
	if err != nil {
		log.Println("Error looking events: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&events)
}
