package utils

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/iambpn/go-email/src/config"
	imapsession "github.com/iambpn/go-email/src/imapSession"
)

func GracefulShutdown(app *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// logout the imapwrapper client
		if imapsession.IW.Host != "" {
			fmt.Println("Imap Wrapper logged out.")
			imapsession.IW.Logout()
		}

		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Error...! Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Run server.
	addr := fmt.Sprintf("%s:%s", config.GetConfig("host", ""), config.GetConfig("port", "3000"))
	if err := app.Listen(addr); err != nil {
		log.Printf("Error while starting server! Reason: %v", err)
	}

	<-idleConnsClosed
}
