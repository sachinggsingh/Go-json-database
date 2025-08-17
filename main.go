package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/sachinggsingh/database/model"
)

var db *model.Driver

// init the database once
func initDB() {
	var err error
	db, err = model.New("./", nil)
	if err != nil {
		panic(err)
	}
}

// Process user and save to DB
func processUser(user model.User) model.User {
	if err := db.Write("users", user.Name, user); err != nil {
		fmt.Println("Error writing user:", err)
	}
	return user
}

// Handle POST request
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dataResponse model.User
	err := json.NewDecoder(r.Body).Decode(&dataResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process user (e.g. save to DB)
	result := processUser(dataResponse)

	// Return as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"users": result})
}

// Handle GET request - return all users
func handleGetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// read all user records
	records, err := db.ReadAll("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userList []model.User
	for _, rec := range records {
		var u model.User
		if err := json.Unmarshal([]byte(rec), &u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

func main() {
	// Init DB
	initDB()

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Welcome to my Go Server!")
	})
	http.HandleFunc("/process", handlePost)
	http.HandleFunc("/users", handleGetAll)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
