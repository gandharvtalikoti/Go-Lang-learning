package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// who took the book
// every book has a log data
// authentication 
// user reports
// user lists 
// over due books
// book history

type user struct {
    ID    *int   `json:"id,omitempty"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type book struct {
	ID       *int   `json:"id,omitempty"` // json field initializer
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}


// creating user

func createUser(c *gin.Context) {
	db, er := connectDB()
	if er != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}
	defer db.Close()

	var newUser user
	if err := c.Bind(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	// Check if user already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", newUser.Email).Scan(&count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to check user existence"})
		return
	}
	if count > 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	res, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", newUser.Name, newUser.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	_, err = res.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user ID"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newUser)

}

// deleting user
func deleteUser(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}
	defer db.Close()

	id := c.Param("id")

	// Check if user exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to check user existence"})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Get user details before deletion
	var deletedUser user
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&deletedUser.ID, &deletedUser.Name, &deletedUser.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user details"})
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted successfully", "user": deletedUser})
}


// creating a book
func createBook(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to db"})
		return
	}
	defer db.Close()

	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	// books = append(books, newBook)
	res, err := db.Exec("INSERT into books_details (title, author, quantity) VALUES (?, ?, ?)", newBook.Title, newBook.Author, newBook.Quantity)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}
	_, err = res.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to get book ID"})
		return
	}
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
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to the db"})
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
