package main

import (
	"log"

	"github.com/ryu111/special-area-booking/internal/config"
	"github.com/ryu111/special-area-booking/internal/database"
	"github.com/ryu111/special-area-booking/internal/routes"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = database.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err = database.Seed(db, cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("API listening on :%s", cfg.Port)
	if err = routes.Setup(db, cfg).Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
