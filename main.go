package main

import (
	"log"

	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/db/mysql"
	"nearbyassist/internal/routes"
	"nearbyassist/internal/routing_engine"
	"nearbyassist/internal/server"
	"nearbyassist/internal/storage"
	"nearbyassist/internal/websocket"
)

func main() {
	// Load configuration file
	config := config.LoadConfig()

	// Load file store
	store := storage.NewStorage(config)
	store.Initialize()

	// Load database configuration
	db := mysql.NewMysqlDatabase(config)

	// Load authenticator configuration
	auth := authenticator.NewJWTAuthenticator(config.JwtSecret)

	// Load websocket configuration
	ws := websocket.NewWebsocket(db)

	// Load Routing Engine configuration
	engine := routing_engine.NewOSRM(config)

	// Create and start the server
	server := server.NewServer(config, ws, db, store, auth, engine)
	routes.RegisterRoutes(server)

	go server.Websocket.SaveMessages()
	go server.Websocket.ForwardMessages()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
