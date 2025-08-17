# Go JSON Database

A lightweight JSON-backed database implemented in Go. This project store using JSON files, ideal for prototyping, small projects, or learning.

## Project Overview

- **Language**: Go  
- **Database**: Local JSON files for storage  
- **Features**:
  - Initialize a JSON-based datastore
  - Create, read (single/all), and delete records
  - RESTful HTTP API with JSON interface

## File Structure

```
.
├── main.go         # Server setup and HTTP handlers
├── model/          # JSON-based database driver
│   └── driver.go   # Database initialization and CRUD methods
└── go.mod          # Module definitions
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.18+ installed

### Installation

```bash
git clone https://github.com/sachinggsingh/Go-json-database.git
cd Go-json-database
go mod tidy
```

### Running the Server

```bash
go run main.go
```

Server will start on `localhost:8080`.

---

## API Endpoints

| HTTP Method | Endpoint   | Description                            | Payload / Response Example             |
|-------------|------------|----------------------------------------|----------------------------------------|
| `GET`       | `/`        | Welcome message                        | `"🚀 Welcome to my Go Server!"`         |
| `POST`      | `/process` | Create or update a user                | JSON user → response with created user |
| `GET`       | `/users`   | Retrieve all users                     | JSON array of users                    |
| `DELETE`    | `/delete`* | Delete a user by name                  | JSON result or success message         |


---

## Code Highlights

### `main.go`
- Initializes the JSON database via `initDB()`.
- Sets up HTTP routes for handling CRUD operations.
- Handles JSON input/output using Go’s `json.NewEncoder/Decoder`.

### `model.Driver`
Includes methods like:
- `New(path string, options interface{})` — initialize the driver.
- `Write(collection, key string, data interface{})` — write data.
- `ReadAll(collection string) ([]string, error)` — fetch all records.
- `Delete(collection, key string) error` — remove a record.

---

## Usage Examples via `curl`

1. **Create a user**
   ```bash
   curl -X POST http://localhost:8080/process
   -H "Content-Type: application/json"
   -d '
   {
   "Name":"Sachin",
   "Age":23,
   "Email":"sachin@example.com",
   "Contact":"1234567890",
   "Address":{
   "City":"Delhi",
   "State":"Delhi",
   "Country":"India",
   "Pincode":"
   110001"
     }
   }'
   ```

2. **Get all users**
   ```bash
   curl http://localhost:8080/users
   ```

3. **Delete a user**
   ```bash
   curl -X DELETE http://localhost:8080/delete
    -H "Content-Type: application/json"
   -d '
   {
   "Name":"Sachin"
   }'
   ```

---

## Why Use This?

- Great for **local prototyping**, simple web services, or learning Go.
- Minimal dependencies and easy-to-understand file-based storage.
- Easy to extend with features like update operations, structured logging, or concurrency control.

---

## Contribute

Contributions and improvements are welcome! Suggestions:
- Add `PUT` or `PATCH` endpoint for updating users
- Improve error handling and response structure
- Add tests and CI configurations
- Implement structured logging or concurrency safeguards

---
