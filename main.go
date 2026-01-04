package main

import (
	"flag"
	"fmt"
	"htrr-apis/internal/app"
	"htrr-apis/internal/routes"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 5500, "BE served on port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	r := routes.SetupRoutes(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute, // TCP wait 1 minute utils close TCP connection
		ReadTimeout:  10 * time.Second,
		Handler:      r,
		WriteTimeout: 30 * time.Minute,
	}
	app.Logger.Printf("App start at port: %d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
