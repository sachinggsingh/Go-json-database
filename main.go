package main

import (
	"encoding/json"
	"fmt"

	model "github.com/sachinggsingh/database/model"
)

func main() {
	dir := "./"

	db, err := model.New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []model.User{
		{
			Name:     "John",
			Age:      30,
			Email:    "j@j.com",
			Password: "password",
			Contact:  "1234567890",
			Address: model.Address{
				City:    "New York",
				State:   "NY",
				Country: "USA",
				Pincode: "12345",
			}},
		{
			Name:     "Jane",
			Age:      25,
			Email:    "j@j.com",
			Password: "password",
			Contact:  "1234567890",
			Address: model.Address{
				City:    "New York",
				State:   "NY",
				Country: "USA",
				Pincode: "12345",
			}},
		{
			Name:     "Sachin",
			Age:      25,
			Email:    "sachin.com",
			Password: "password",
			Contact:  "1234567890",
			Address: model.Address{
				City:    "New York",
				State:   "NY",
				Country: "USA",
				Pincode: "12345",
			}},
	}

	// Write to DB
	for _, value := range employees {
		db.Write("users", value.Name, value)
	}

	// Read all users in JSON format
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}

	// Print raw JSON strings
	fmt.Println("Raw JSON records:")
	for _, r := range records {
		fmt.Println(r)
	}

	// If you need them back as structs
	var allUsers []model.User
	for _, jsonStr := range records {
		var u model.User
		if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {
			fmt.Println("Error:", err)
		}
		//nolint:sa4010
		allUsers = append(allUsers, u)
	}

	// fmt.Println("\nAs structs:")
	// for _, u := range allUsers {
	// 	fmt.Printf("%+v\n", u)
	// }

	// Example: Delete one
	// if err := db.Delete("users", "John"); err != nil {
	// 	fmt.Println("Error", err)
	// }
}
