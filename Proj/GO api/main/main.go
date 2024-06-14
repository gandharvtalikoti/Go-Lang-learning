package main

import (
	"log"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.GET("/books/:id", bookById)
	// router.Run("localhost:8080")

	router.POST("/users", createUser)

	log.Println("Starting server on http://localhost:8080")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
