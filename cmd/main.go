package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kammeph/school-book-storage-service-simplified/auth"
	"github.com/kammeph/school-book-storage-service-simplified/db"
	"github.com/kammeph/school-book-storage-service-simplified/users"
)

func main() {
	db := db.NewSqlDB()
	defer db.Close()
	auth.AddAuthController(db)
	users.AddUsersController(db)
	repo := users.NewSqlUserRepository(db)
	http.HandleFunc("/api/users/update", func(w http.ResponseWriter, r *http.Request) {
		var user users.UserDto
		json.NewDecoder(r.Body).Decode(&user)
		err := repo.UpdateUser(context.Background(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	})
	port := os.Getenv("CONTAINER_PORT")
	log.Printf("App will be served on port: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
