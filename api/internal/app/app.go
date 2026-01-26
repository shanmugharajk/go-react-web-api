package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shanmugharajk/go-react-web-api/api/internal/config"
	"github.com/shanmugharajk/go-react-web-api/api/internal/db"
	httpserver "github.com/shanmugharajk/go-react-web-api/api/internal/http"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/auth"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/customer"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/inventory"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/payment"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/product"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/purchase"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/receiving"
	"github.com/shanmugharajk/go-react-web-api/api/internal/modules/vendor"
	"github.com/shanmugharajk/go-react-web-api/api/internal/pkg/logger"
)

// App represents the application.
type App struct {
	cfg    *config.Config
	db     *db.DB
	server *httpserver.Server
}

// New creates a new application instance.
func New() (*App, error) {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database
	database, err := db.New(cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Auto-migrate models (order matters for FK constraints)
	if err := database.AutoMigrate(
		&auth.User{},
		&product.ProductCategory{},
		&customer.Customer{},
		&product.Product{},
		&inventory.ProductBatch{},
		// Receivables module models
		&vendor.Vendor{},
		&purchase.PurchaseOrder{},
		&purchase.PurchaseOrderItem{},
		&receiving.StockReceipt{},
		&receiving.StockReceiptItem{},
		&payment.VendorPayment{},
	); err != nil {
		return nil, fmt.Errorf("failed to run auto-migrations: %w", err)
	}

	// Create HTTP server
	server := httpserver.New(cfg, database)

	return &App{
		cfg:    cfg,
		db:     database,
		server: server,
	}, nil
}

// Run starts the application and handles graceful shutdown.
func (a *App) Run() error {
	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		logger.Info("Application starting")
		serverErrors <- a.server.Start()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive an error or interrupt signal
	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}

	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig.String())

		// Create context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Graceful shutdown failed", "error", err)
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		// Close database connection
		if err := a.db.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
			return fmt.Errorf("failed to close database: %w", err)
		}

		logger.Info("Application stopped gracefully")
	}

	return nil
}
