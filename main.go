package main

import (
	"errors"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", Quantity: 5},
	{ID: "2", Title: "Nineteen Eighty-Four", Author: "George Orwell", Quantity: 3},
	{ID: "3", Title: "Brave New World", Author: "Aldous Huxley", Quantity: 4},
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (book, error) {
	for i, book := range books {
		if book.ID == id {
			return books[i], nil
		}
	}
	return book{}, errors.New("Book not found")
}
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newBook.ID = uuid.New().String()
	for _, b := range books {
		if b.ID == newBook.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID already exists"})
			return
		}
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, http.StatusText(http.StatusCreated))
}

func checkinBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity++
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is out of stock."})
		return
	}

	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	// defining the routes
	router := gin.Default()

	// Get All Books
	router.GET("/books", getBooks)
	// Get Book by ID
	router.GET("/books/:id", bookById)
	// Create a new Book
	router.POST("/books", createBook)
	// Checkin a Book
	router.PUT("/books/:id/checkin", checkinBook)
	// Checkout a Book
	router.PUT("/books/:id/checkout", checkoutBook)
	// Start the server
	router.Run("localhost:8080")
}
