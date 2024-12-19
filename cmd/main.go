package main

import (
	"log"
	"os"

	"ivf_calculator/api"
	"ivf_calculator/internal/repo"
	"ivf_calculator/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[ivf_calculator]: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	ivfRepo := repo.NewIVFFormula(&repo.Config{
		FilePath: "internal/repo/data/ivf_success_formulas.csv",
		Logger:   logger,
	})
	ivfService := server.NewSuccessCalculator(&server.Config{
		Logger: logger,
		Repo:   ivfRepo,
	})

	s := api.New(&api.Config{
		Port:       ":8080",
		Logger:     logger,
		IVFService: ivfService,
	})

	s.Start()
}
