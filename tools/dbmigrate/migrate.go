package main

import (
	"log"

	"github.com/iloveicedgreentea/gowatchit/pkg/config"
	"github.com/iloveicedgreentea/gowatchit/pkg/database"
)

func main() {
	// TODO: get path from arg
	db, err := database.GetDB("./")
	if err != nil {
		log.Fatal(err)
	}
	config.RunMigrations(db)
}
