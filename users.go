package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// model
type User struct {
	IDUser     int    `json:"id_user"`
	FirstName  string `json:"first_name"`
	MidName    string `json:"mid_name"`
	SurName    string `json:"sur_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	PasswordUs string `json:"password_us"`
	Phone      string `json:"phone"`
	Col        string `json:"col"`
	Street     string `json:"street"`
	StNum      string `json:"st_num"`
	ZipCode    string `json:"zip_code"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

// In-memory data store
var userDatabase = make(map[string]User)

// GET handler
func GetUser(writer http.ResponseWriter, reques *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	//validate HTTP method
	if reques.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Extract the email from the URL query
	email := reques.URL.Query().Get("email")
	if email == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Email query parameter is required"})
		return
	}

	// Look up the user in the map
	user, exists := userDatabase[email]
	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(map[string]string{"error": "User not found"})
		return
	}

	// Return the found user
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(user)
}

func GetAllUsers(writer http.ResponseWriter, reques *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Validate HTTP Method
	if reques.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Convert map to slice for standard JSON array serialization
	var users []User
	for _, user := range userDatabase {
		users = append(users, user)
	}

	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(users)
}

// POST handler
func CreateUser(writer http.ResponseWriter, reques *http.Request) {
	// JSON Header
	writer.Header().Set("Content-Type", "application/json")

	// Validate HTTP method
	if reques.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Decode incoming JSON body directly into your Go struct model
	var newUser User
	err := json.NewDecoder(reques.Body).Decode(&newUser)

	// Handle decode error
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Save to DB using Email as the map key
	userDatabase[newUser.Email] = newUser

	// Return success response
	writer.WriteHeader(http.StatusCreated) // 201 created
	json.NewEncoder(writer).Encode(newUser)
}

// PUT handler
func UpdateUser(writer http.ResponseWriter, reques *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Validate HTTP Method
	if reques.Method != http.MethodPut {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	// Extract the email to identify the user
	email := reques.URL.Query().Get("email")
	if email == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Email query parameter is required"})
		return
	}

	// Check if the user actually exists before updating
	_, exists := userDatabase[email]
	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(writer).Encode(map[string]string{"error": "User not found"})
		return
	}

	// Decode the incoming JSON into a new User struct
	var updatedUser User
	err := json.NewDecoder(reques.Body).Decode(&updatedUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Enforce that the email field matches the URL parameter so the map key does not break
	updatedUser.Email = email

	// Overwrite the existing user in the map
	userDatabase[email] = updatedUser

	// Return the updated user data
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(updatedUser)

}

func main() {
	fmt.Println("Starting server on port 8080...")

	// Register POST routes
	http.HandleFunc("/user/create", CreateUser)

	http.HandleFunc("/user/update", UpdateUser)

	// Register GET routes
	http.HandleFunc("/user", GetUser)
	http.HandleFunc("/users", GetAllUsers)

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
