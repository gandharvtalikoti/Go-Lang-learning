package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// we'll build a library api
// - store a bunch of books
// - check in a book
// - check out a book
// - view all books
// - get books by id

type book struct {
	ID       string `json:"id"` // json field initializer
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}


func getBooks(c *gin.Context) {
	db, err := connectDB()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"failed to connect to db"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, author, quantity FROM books_details")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch books"})
		return
	}
	defer rows.Close()

	var books []book // books is array which contains all books
	for rows.Next() {
		var b book // a single book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch books"})
			return
		}
		books = append(books, b)
	}

	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

// getting a book by id
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id") // "/books/2" from url path
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	// here we are going to use query paramenter
	// example ?id=2
	id, ok := c.GetQuery("id") // "/books/?id=2"
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "book not available"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id") // "/books/?id=2"
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing id query parameter"})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "book not available"})
		return
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func connectDB() (*sql.DB, error) {
	uri := "root:1910@tcp(localhost:3306)/library_system_go"
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	return db, nil

}
func main() {
	
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.GET("/books/:id", bookById)
	// router.Run("localhost:8080")


	log.Println("Starting server on http://localhost:8080")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
