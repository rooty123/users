package actions

import (
	"net/http"
	"strconv"
	"users/db"
	"users/metrics"

	"github.com/labstack/echo/v4"
)

// Placeholder for User struct
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Mock in-memory "database"
var users = []User{
	{ID: "1", Name: "John Doe", Email: "john@example.com"},
}

// GetUsers returns a list of users
func GetUsers(c echo.Context) error {
	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	users, err := dbh.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get users")
	}

	// Update users gauge
	metrics.UsersGauge.Set(float64(len(users)))
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	// Extract the chat ID from the URL path
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	// Create a new database handler
	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	// Get the user from the database
	user, err := dbh.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Return the user as JSON
	return c.JSON(http.StatusOK, user)
}

// CreateUser adds a new user
func CreateUser(c echo.Context) error {
	var newUser struct {
		ChatID       int64  `json:"chatID"`
		TelegramID   int64  `json:"telegramID"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		LanguageCode string `json:"languageCode"`
		Username     string `json:"username"`
	}

	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	if err := dbh.CreateUser(newUser.ChatID, newUser.TelegramID, newUser.FirstName, newUser.LastName, newUser.LanguageCode, newUser.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create user")
	}

	// Update users gauge
	users, _ := dbh.ListUsers()
	metrics.UsersGauge.Set(float64(len(users)))

	return c.JSON(http.StatusCreated, newUser)
}

// UpdateUser updates an existing user
func UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	var updatedUser struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		LanguageCode string `json:"languageCode"`
		Username     string `json:"username"`
	}

	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	if err := dbh.UpdateUser(id, updatedUser.FirstName, updatedUser.LastName, updatedUser.LanguageCode, updatedUser.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	dbh := db.UsersDBHandler{}
	dbh.ConnectPg()
	defer dbh.Conn.Close()

	if err := dbh.DeleteUser(id); err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// Update users gauge
	users, _ := dbh.ListUsers()
	metrics.UsersGauge.Set(float64(len(users)))

	return c.NoContent(http.StatusNoContent)
}
