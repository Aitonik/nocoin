package main

import (
	"gitlab.com/aitalina/nocoin/internal/repository/psql"
	"gitlab.com/aitalina/nocoin/internal/service"
	"gitlab.com/aitalina/nocoin/internal/transport/rest"
	"gitlab.com/aitalina/nocoin/pkg/database"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// init db
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "qwerty123",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// init deps
	restaurantsRepo := psql.NewRestaurant(db)
	restaurantsService := service.NewRestaurants(restaurantsRepo)

	profilesRepo := psql.NewProfile(db)
	profileService := service.NewProfiles(profilesRepo)

	tipRepo := psql.NewTip(db)
	tipService := service.NewTips(tipRepo)

	handler := rest.NewHandler(restaurantsService, profileService, tipService)

	// init & run server
	srv := http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
