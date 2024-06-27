package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Collection struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Item struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	CollectionID string `json:"collectionId"`
	IsCollected  bool   `json:"isCollected"`
	Rarity       int    `json:"rarity"`
}

const (
	host     = "postgres"
	port     = 5432
	user     = "user"
	password = "user"
	dbname   = "collection_db"
)

var (
	db *sql.DB
)

var secretKey = []byte("API_Secret_Key") //TODO set as env var

func main() {
	err := setDb()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/collections", validateTokenMiddleware(getCollections)).Methods("GET")
	r.HandleFunc("/collections", validateTokenMiddleware(createCollection)).Methods("POST")
	r.HandleFunc("/collections/{id}", validateTokenMiddleware(updateCollection)).Methods("PUT")
	r.HandleFunc("/collections/{id}", validateTokenMiddleware(deleteCollection)).Methods("DELETE")
	r.HandleFunc("/collections/{collectionId}/items", validateTokenMiddleware(getItems)).Methods("GET")
	r.HandleFunc("/collections/{collectionId}/items/{itemId}", validateTokenMiddleware(getItem)).Methods("GET")
	r.HandleFunc("/collections/{collectionId}/items", validateTokenMiddleware(createItem)).Methods("POST")
	r.HandleFunc("/collections/{collectionId}/items/{itemId}", validateTokenMiddleware(updateItem)).Methods("PUT")
	r.HandleFunc("/collections/{collectionId}/items/{itemId}", validateTokenMiddleware(deleteItem)).Methods("DELETE")

	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, HEAD, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, Accept, Origin, X-Custom-Header")
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Server is starting on port 8000")
	log.Fatal(http.ListenAndServe(":8000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "PATCH", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept", "Origin", "X-Custom-Header"}),
		)(r),
	))
}

func setDb() error {
	var err error
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	}

	err = createTables()
	if err != nil {
		db.Close()
		return err
	}

	err = insertExample()
	if err != nil {
		db.Close()
		log.Fatalf("Error inserting data: %v", err)
	}

	log.Println("Connected to PostgresQL database")
	return nil
}

func createTables() error {
	createUsersTable := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL
        );
    `

	createCollectionsTable := `
        CREATE TABLE IF NOT EXISTS collections (
            id UUID PRIMARY KEY,
            title VARCHAR(100) NOT NULL
        );
    `

	createItemsTable := `
        CREATE TABLE IF NOT EXISTS items (
            id UUID PRIMARY KEY,
            title VARCHAR(100) NOT NULL,
            collection_id UUID NOT NULL,
            is_collected BOOLEAN NOT NULL,
            rarity INTEGER NOT NULL,
            FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
        );
    `

	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createCollectionsTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createItemsTable)
	if err != nil {
		return err
	}

	return nil
}

func insertExample() error {
	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
    `
	_, err := db.Exec(query, "Karol", "pass")
	if err != nil {
		return err
	}
	_, err = db.Exec(query, "Rob", "pass")
	if err != nil {
		return err
	}

	insertCollectionQuery := `
        INSERT INTO collections (id, title)
        VALUES ($1, $2)
    `
	nbaCollectionID := uuid.New()
	_, err = db.Exec(insertCollectionQuery, nbaCollectionID, "NBA Cards")
	if err != nil {
		return err
	}

	nflCollectionID := uuid.New()
	_, err = db.Exec(insertCollectionQuery, nflCollectionID, "NFL Cards")
	if err != nil {
		return err
	}

	insertItemQuery := `
        INSERT INTO items (id, title, collection_id, is_collected, rarity)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = db.Exec(insertItemQuery, uuid.New(), "Platinum Michael Jordan '95", nbaCollectionID, true, 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertItemQuery, uuid.New(), "Gold Dennis Rodman '95", nbaCollectionID, false, 2)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertItemQuery, uuid.New(), "Silver Steph Curry '16", nbaCollectionID, false, 3)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertItemQuery, uuid.New(), "Orange Chad Ochocinco '09", nflCollectionID, true, 1)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertItemQuery, uuid.New(), "Silver Bo Jackson '85", nflCollectionID, true, 2)
	if err != nil {
		return err
	}
	_, err = db.Exec(insertItemQuery, uuid.New(), "Red O. J. Simpson '75", nflCollectionID, true, 3)
	if err != nil {
		return err
	}

	return nil
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("Creating user: %s\n", user.Username)
	log.Printf("Creating password: %s\n", user.Password)

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func login(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := requestBody["username"]
	password := requestBody["password"]

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if password != storedPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func getCollections(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title FROM collections")
	if err != nil {
		http.Error(w, "Failed to fetch collections", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var collection Collection
		err := rows.Scan(&collection.ID, &collection.Title)
		if err != nil {
			http.Error(w, "Failed to fetch collections", http.StatusInternalServerError)
			return
		}
		collections = append(collections, collection)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch collections", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(collections)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	var collection Collection
	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	collection.ID = uuid.New().String()

	_, err = db.Exec("INSERT INTO collections (id, title) VALUES ($1, $2)", collection.ID, collection.Title)
	if err != nil {
		http.Error(w, "Failed to create collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(collection)
}

func updateCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedCollection Collection
	err := json.NewDecoder(r.Body).Decode(&updatedCollection)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE collections SET title = $1 WHERE id = $2", updatedCollection.Title, params["id"])
	if err != nil {
		http.Error(w, "Failed to update collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedCollection)
}

func deleteCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM collections WHERE id = $1", params["id"]).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to check collection existence", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM collections WHERE id = $1", params["id"])
	if err != nil {
		http.Error(w, "Failed to delete collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	rows, err := db.Query("SELECT id, title, collection_id, is_collected, rarity FROM items WHERE collection_id = $1", params["collectionId"])
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Title, &item.CollectionID, &item.IsCollected, &item.Rarity)
		if err != nil {
			http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var item Item
	err := db.QueryRow("SELECT id, title, collection_id, is_collected, rarity FROM items WHERE id = $1 AND collection_id = $2", params["itemId"], params["collectionId"]).Scan(&item.ID, &item.Title, &item.CollectionID, &item.IsCollected, &item.Rarity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch item", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	item.ID = uuid.New().String()
	item.CollectionID = params["collectionId"]

	_, err = db.Exec("INSERT INTO items (id, title, collection_id, is_collected, rarity) VALUES ($1, $2, $3, $4, $5)",
		item.ID, item.Title, item.CollectionID, item.IsCollected, item.Rarity)
	if err != nil {
		http.Error(w, "Failed to create item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var updatedItem Item
	err := json.NewDecoder(r.Body).Decode(&updatedItem)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE items SET title = $1, is_collected = $2, rarity = $3 WHERE id = $4 AND collection_id = $5",
		updatedItem.Title, updatedItem.IsCollected, updatedItem.Rarity, params["itemId"], params["collectionId"])
	if err != nil {
		http.Error(w, "Failed to update item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM items WHERE id = $1 AND collection_id = $2", params["itemId"], params["collectionId"]).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to check item existence", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM items WHERE id = $1 AND collection_id = $2", params["itemId"], params["collectionId"])
	if err != nil {
		http.Error(w, "Failed to delete item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
