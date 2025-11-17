package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/williamkoller/system-education/config"
	userRouter "github.com/williamkoller/system-education/internal/user/presentation/router"
	"github.com/williamkoller/system-education/shared/middleware"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database := config.NewDatabaseConnection()
	config.RunMigrations(database, "")

	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(middleware.GlobalErrorHandler())
	userRouter.UserRouter(g, database, cfg.Resend)

	address := ":" + strconv.Itoa(cfg.App.Port)
	srv := &http.Server{
		Addr:              address,
		Handler:           g,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server running at http://localhost:8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
