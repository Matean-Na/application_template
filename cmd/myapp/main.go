package main

import (
	"bitbucket.org/microret/oxus/internal/server"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

// @title        Microfin API
// @version      0.0.1
// @description  Microfin backend api
// @BasePath     /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := server.Server{}
	r, e := s.Init()
	if e != nil {
		return
	}

	s.Run(r)

	<-ctx.Done()

	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")
	s.CloseAll()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
