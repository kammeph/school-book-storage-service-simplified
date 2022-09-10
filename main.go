package main

import (
	"context"
	"encoding/json"
	"net/http"

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
	http.ListenAndServe(":9090", nil)
}
