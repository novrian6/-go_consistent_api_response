package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
	Version string      `json:"version"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}

func successResponse(c *gin.Context, message string, data interface{}, meta *Meta) {
	response := APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
		Version: "v1",
	}
	c.JSON(http.StatusOK, response)
}

func errorResponse(c *gin.Context, statusCode int, message string, errors []string) {
	response := APIResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
		Version: "v1",
	}
	c.JSON(statusCode, response)
}

func getUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	start := (page - 1) * limit
	end := start + limit
	if start >= len(users) {
		successResponse(c, "Users retrieved successfully", []User{}, &Meta{
			Page:       page,
			Limit:      limit,
			Total:      len(users),
			TotalPages: (len(users) + limit - 1) / limit,
		})
		return
	}
	if end > len(users) {
		end = len(users)
	}
	paginatedUsers := users[start:end]
	meta := &Meta{
		Page:       page,
		Limit:      limit,
		Total:      len(users),
		TotalPages: (len(users) + limit - 1) / limit,
	}
	successResponse(c, "Users retrieved successfully", paginatedUsers, meta)
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid user ID", []string{"ID must be a valid number"})
		return
	}
	for _, user := range users {
		if user.ID == id {
			successResponse(c, "User found", user, nil)
			return
		}
	}
	errorResponse(c, http.StatusNotFound, "User not found", []string{"User with specified ID does not exist"})
}

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", getUsers)
		v1.GET("/users/:id", getUserByID)
	}
	r.Run(":8080")
}
