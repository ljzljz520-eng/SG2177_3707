package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"equipmentlending/internal/config"
	"equipmentlending/internal/persistence"
	"equipmentlending/internal/service"
	"equipmentlending/internal/transport"
)

func main() {
	settings, err := config.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := settings.EnsureParent(); err != nil {
		log.Fatal(err)
	}
	store, err := persistence.Open(settings.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	business := service.New(store)
	if err := business.Load(context.Background()); err != nil {
		log.Fatal(err)
	}
	if settings.ServeHTTP {
		log.Printf("equipment service listening on %s", settings.Listen)
		if err := transport.NewServer(business).ListenAndServe(context.Background(), settings.Listen); err != nil {
			log.Fatal(err)
		}
		return
	}
	fmt.Println("equipment lending service ready")
	fmt.Println(settings.Summary())
}
