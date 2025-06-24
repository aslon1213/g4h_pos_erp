package main

import (
	"aslon1213/magazin_pos/pkg/app"
	models "aslon1213/magazin_pos/pkg/repository"
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {

	app := app.NewApp()
	// app.Run()

	db := app.DB

	finances := Read_finance()
	suppliers := ReadSuppliers()

	db.Database("magazin").Collection("finance").Drop(context.Background())
	db.Database("magazin").Collection("suppliers").Drop(context.Background())

	for _, finance := range finances {
		log.Info().Interface("finance", finance).Msg("Inserting finance")
		db.Database("magazin").Collection("finance").InsertOne(context.Background(), finance)
	}

	for _, supplier := range suppliers {
		log.Info().Interface("supplier", supplier).Msg("Inserting supplier")
		db.Database("magazin").Collection("suppliers").InsertOne(context.Background(), supplier)
	}

}

func Read_finance() []models.BranchFinance {

	finances := []models.BranchFinance{}
	file, err := os.ReadFile("finances.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read finances.json")
	}
	json.Unmarshal(file, &finances)
	return finances
}

func ReadSuppliers() []map[string]interface{} {
	suppliers := []map[string]interface{}{}
	file, err := os.ReadFile("suppliers.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read suppliers.json")
	}
	json.Unmarshal(file, &suppliers)
	return suppliers
}
