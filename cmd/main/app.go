package main

import (
	"console-application-service-age/internal/config"
	"console-application-service-age/internal/test"
	"console-application-service-age/internal/test/db"
	"console-application-service-age/pkg/client/redis"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("to launch the application, you need 2 parameters (-host, -port)")
	}

	log.Println("create router")
	router := httprouter.New()

	newClient := redis.NewClient(os.Args[1], os.Args[2])
	defer newClient.Close()

	repository := db.NewRepository(newClient)
	handler := test.NewHandler(repository)
	handler.Register(router)

	cfg := config.GetConfig()
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	log.Println("start application")

	var listener net.Listener
	var listenError error

	log.Println("listen tcp")
	listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
	log.Println(fmt.Sprintf("server is listening %s:%s", cfg.Listen.BindIp, cfg.Listen.Port))

	if listenError != nil {
		log.Fatalln(listenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.Serve(listener))
}
