package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"net/http"
	
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	appLogger.Info(ctx, "Shutting down server")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	appLogger.Info(ctx, "Server gracefull stopped")
}
