package main

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func main() {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	//appLogger.Info(ctx, "Shutting down server")
	
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// srv.Shutdown(ctx)
	// appLogger.Info(ctx, "Server gracefull stopped")
}
