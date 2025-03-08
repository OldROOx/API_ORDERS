package main

import (
	"app.initial/src/core/database"
	"app.initial/src/core/router"
	"app.initial/src/core/server"
	"app.initial/src/orders/domain/entities"
	orderRoutes "app.initial/src/orders/infrastructure/routes"
	"log"
)

func main() {
	// Database connection
	db := database.NewMySQLConnection()

	// Auto-migrate order entities
	err := db.AutoMigrate(&entities.Order{}, &entities.OrderItem{})
	if err != nil {
		log.Fatal("Error migrating order entities:", err)
	}

	// RabbitMQ URL
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"

	// Router setup
	r := router.NewRouter(db)

	// Setup order routes
	api := r.GetEngine().Group("/api")
	orderRoutes.SetupOrderRoutes(api, db, rabbitMQURL)

	// Create server
	srv := server.NewServer("8080", r.GetEngine())

	// Start server
	log.Fatal(srv.Start())
}
