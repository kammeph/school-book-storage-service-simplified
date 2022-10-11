package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kammeph/school-book-storage-service-simplified/auth"
	"github.com/kammeph/school-book-storage-service-simplified/db"
	"github.com/kammeph/school-book-storage-service-simplified/schools"
	"github.com/kammeph/school-book-storage-service-simplified/users"
)

func main() {
	db := db.NewSqlDB()
	defer db.Close()
	auth.AddAuthController(db)
	users.AddUsersController(db)
	schools.AddSchoolsController(db)
	port := os.Getenv("CONTAINER_PORT")
	log.Printf("App will be served on port: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}