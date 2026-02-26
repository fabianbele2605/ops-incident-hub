package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/api/handler"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/config"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/infrastructure/postgres"
	"github.com/fabianbele2605/ops-incident-hub/backend/internal/usecase/incident"
)

func main() {

	// Verificar si es comando health
	if len(os.Args) > 1 && os.Args[1] == "health" {
		runHealthCheck()
		return
	}
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Ops Incident Hub API")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Server: %s:%s", cfg.Server.Host, cfg.Server.Port)

	// Conectar a la base de datos
	db, err := postgres.NewDatabase(postgres.DatabaseConfig{
		DSN:             cfg.DatabaseDSN(),
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Ejecutar migraciones
	if err := postgres.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")

	// Inicializar repositorios
	incidentRepo := postgres.NewIncidentRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Inicializar casos de uso
	createIncidentUC := incident.NewCreateIncidentUseCase(incidentRepo, userRepo)
	assignIncidentUC := incident.NewAssignIncidentUseCase(incidentRepo, userRepo)
	listIncidentsUC := incident.NewListIncidentUseCase(incidentRepo)

	// Inicializar handlers
	healthHandler := handler.NewHealthHandler(db)


	// Inicializar handlers
	incidentHandler := handler.NewIncidentHandler(
		createIncidentUC,
		assignIncidentUC,
		listIncidentsUC,
	)

	// Configurar router
	router := api.SetupRouter(incidentHandler, healthHandler)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Canal para errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Server listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para señales del sistema
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Esperar por error o señal de shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case sig := <-shutdown:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

		// Crear contexto con timeout para shutdown
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		// Intentar shutdown graceful
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Failed to close server: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}
}


func runHealthCheck() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	url := fmt.Sprintf("http://%s:%s/health/live", cfg.Server.Host, cfg.Server.Port)
	resp, err := http.Get(url)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	os.Exit(0)
}
