package controllers

import (
	"Mongo/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type UserController struct {
	Collection *mongo.Collection
}

func NewUserController(client *mongo.Client) *UserController {
	// Access the "users" collection
	collection := client.Database("mongo-golang").Collection("users")
	return &UserController{Collection: collection}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Convert the id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var u models.User
	// Query the "users" collection for the user by id
	err = uc.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Handle other errors
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the user object to JSON
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u models.User

	// Decode the request body into the user struct
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a new ObjectID and assign it to the user
	u.Id = primitive.NewObjectID()
	_, err = uc.Collection.InsertOne(context.Background(), u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert the user to JSON for the response
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Convert the id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the user from the collection
	_, err = uc.Collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Handle other errors
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s\n", id)
}
