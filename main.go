package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {

	e := echo.New()

	initDb()
	e.GET("/messages", GetHandler)
	e.POST("/messages", PostHandler)
	e.PATCH("/messages/:id", PatchHandler)
	e.DELETE("/messages/:id", DeleteHandler)
	e.Start(":8081")
}

var db *gorm.DB

func initDb() {
	dsn := "host=localhost user=postgres password=... dbname=echo_project port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Message{})
}

func GetHandler(e echo.Context) error {
	var message []Message
	if err := db.Find(&message).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not get data",
		})
	}
	return e.JSON(http.StatusInternalServerError, &message)
}

func PostHandler(e echo.Context) error {
	var message Message
	if err := e.Bind(&message); err != nil {
		return e.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "message not added",
		})
	}

	if err := db.Create(&message).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not create data",
		})
	}

	return e.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "message added",
	})
}

func PatchHandler(e echo.Context) error {
	idParam := e.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not convert id to int",
		})
	}
	var update Message
	if err := e.Bind(&update); err != nil {
		return e.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not update message",
		})
	}

	if err := db.Model(&Message{}).Where("id = ?", id).Update("text", update.Text).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not update message",
		})
	}

	return e.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "message updated",
	})

}

func DeleteHandler(e echo.Context) error {
	idParam := e.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Could not convert id to int",
		})
	}

	if err := db.Where("id = ?", id).Delete(&Message{}).Error; err != nil {
		return e.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Could not delete message",
		})
	}

	return e.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message deleted",
	})
}
