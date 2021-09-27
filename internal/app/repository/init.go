package repository

import (
	"fmt"
	"github.com/portnyagin/practicum_go/internal/app/database"
	"github.com/portnyagin/practicum_go/internal/app/repository/basedbhandler"
)

func InitDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("Database structure created successfully")
	return nil
}

func ClearDatabase(h basedbhandler.DBHandler) error {
	err := h.Execute(database.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
