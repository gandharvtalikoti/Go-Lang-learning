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

// dummy data
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
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
	db, err := connectDB()
	if err != nil {
		//c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"failed to connect to db"})
		return nil, errors.New("failed to connect to database")
	}
	defer db.Close()

	var b book // to store the book with id
	err = db.QueryRow("SELECT id, title, author, quantity FROM books_details WHERE id = ?", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err == sql.ErrNoRows {
		return nil, errors.New("book not found")
	} else if err != nil {
		return nil, errors.New("failed to fetch book")
	}
	return &b, nil
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
	db, err := connectDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}
	defer db.Close()

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

	res, err := db.Exec("UPDATE books_details SET quantity = ? WHERE id=?", book.Quantity-1, book.ID)
	if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to update book quantity"})
        return
    }

	log.Println(res.RowsAffected()) // this will print no. of rows affected in the db i suppose
	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)

}
func returnBook(c *gin.Context) {
	db, err := connectDB()
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error":"failed to connect to the db"})
	}
	defer db.Close()

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

	_, err = db.Exec("UPDATE books_details SET quantity = ? WHERE id=?", book.Quantity+1, book.ID)
	if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to update book quantity"})
        return
    }
	book.Quantity++
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
