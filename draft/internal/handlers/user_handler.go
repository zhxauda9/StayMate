package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/zhxauda9/StayMate/draft/internal/db"
	"github.com/zhxauda9/StayMate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Создание нового пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("users")
	user.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Получение всех пользователей
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	collection := db.GetCollection("users")
	cursor, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
