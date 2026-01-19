package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"msng/internal/app"
	"msng/internal/config"
)

func main() {
	//Init config
	config := config.InitConfiguration()

	log.SetFlags(log.LstdFlags | log.Llongfile)
	ctx := context.Background()

	handler, db := app.AppInit(config, ctx)
	defer db.Close()
	//Init server
	server := &http.Server{
		Addr:    ":8080",
		Handler: *handler,
	}

	fmt.Println("Starting server on 8080 port")

	log.Fatal(server.ListenAndServe())
}
