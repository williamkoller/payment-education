package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	user_router "github.com/williamkoller/system-education/internal/user/router"
)

func main() {
	g := gin.Default()

	g.Use(gin.Recovery())

	user_router.UserRouter(g)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           g,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
