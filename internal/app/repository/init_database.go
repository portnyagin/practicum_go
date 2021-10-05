package repository

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_go/internal/app/database"
	"github.com/portnyagin/practicum_go/internal/app/repository/basedbhandler"
)

func InitDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, database.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("Database structure created successfully")
	return nil
}

func ClearDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, database.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
