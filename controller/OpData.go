package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OpData struct {
	ID        string `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string `json:"email,omitempty" bson:"email,omitempty"`
}
type CustomResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UserController struct {
	collection *mongo.Collection
}

func NewUserController(r *mux.Router, collection *mongo.Collection) *UserController {
	uc := &UserController{collection}
	userRoutes(r, uc)
	return uc
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var opData OpData
	err := json.NewDecoder(r.Body).Decode(&opData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := uc.collection.InsertOne(ctx, opData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func (uc *UserController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var opData OpData
	err = uc.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&opData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opData)
}

func (uc *UserController) All(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := uc.collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var opDataList []OpData
	for cursor.Next(ctx) {
		var opData OpData
		err := cursor.Decode(&opData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		opDataList = append(opDataList, opData)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opDataList)
}

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var opData OpData
	err := json.NewDecoder(r.Body).Decode(&opData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": opData}

	result, err := uc.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "mongo: no documents in result", http.StatusNotFound)
		return
	}

	response := CustomResponse{
		Message: "Document updated successfully",
		Data:    result.ModifiedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uc *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID}

	result, err := uc.collection.DeleteOne(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "mongo: no documents in result", http.StatusNotFound)
		return
	}

	response := CustomResponse{
		Message: "Document deleted successfully",
		Data:    result.DeletedCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func userRoutes(r *mux.Router, uc *UserController) {
	r.HandleFunc("/opdata", uc.Create).Methods("POST")
	r.HandleFunc("/opdata", uc.All).Methods("GET")
	r.HandleFunc("/opdata/{id}", uc.Get).Methods("GET")
	r.HandleFunc("/opdata/{id}", uc.Update).Methods("PUT")
	r.HandleFunc("/opdata/{id}", uc.Delete).Methods("DELETE")
}
