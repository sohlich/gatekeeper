package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/sohlich/gatekeeper/handlers"
	"github.com/sohlich/gatekeeper/mail"
	"github.com/sohlich/gatekeeper/model"
)

func main() {
	log.Println("Application starting...")
	err := model.InitDB()
	if err != nil {
		log.Printf("Database init fails: %+v \n", errors.Cause(err))
		return
	}

	appMux := http.NewServeMux()
	appMux.HandleFunc("/register", handlers.Register)
	appMux.HandleFunc("/login", handlers.Login)
	appMux.HandleFunc("/activate", handlers.ActivateUser)
	server := manners.NewServer()
	server.Handler = handlers.LoggingHandler(appMux)

	errChan := make(chan (error))
	go func() {
		server.Addr = ":8080"
		errChan <- server.ListenAndServe()
	}()

	mail.InitFakeMailer()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			server.BlockingClose()
			model.CloseDB()
			os.Exit(0)
		}
	}
}
