package main

import (
	"database/sql"
	"log"

	"daily-app-go/db/sqlc"
	"daily-app-go/internal/config"
	"daily-app-go/internal/handler"
	"daily-app-go/internal/repository"
	"daily-app-go/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	dbConn, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// Initialize sqlc queries
	queries := db.New(dbConn)

	// Initialize repositories
	taskRepo := repository.NewTaskRepository(queries)
	subTaskRepo := repository.NewSubTaskRepository(queries)
	blockerRepo := repository.NewBlockerRepository(queries)
	noteRepo := repository.NewNoteRepository(queries)
	dailySessionRepo := repository.NewDailySessionRepository(queries)

	// Initialize internal services
	taskPolicyService := service.NewTaskPolicyService(blockerRepo, subTaskRepo)
	subTaskPolicyService := service.NewSubTaskPolicyService(blockerRepo)
	taskSynchronizerService := service.NewTaskSynchronizerService(taskRepo, blockerRepo)
	subTaskSynchronizerService := service.NewSubTaskSynchronizerService(subTaskRepo, blockerRepo)
	taskOwnershipGuardService := service.NewTaskOwnershipGuardService(taskRepo, subTaskRepo)

	// Initialize main services
	blockerService := service.NewBlockerService(blockerRepo, taskOwnershipGuardService, taskSynchronizerService, subTaskSynchronizerService)
	noteService := service.NewNoteService(noteRepo, taskOwnershipGuardService)
	subTaskService := service.NewSubTaskService(subTaskRepo, taskRepo, taskSynchronizerService, subTaskPolicyService)
	taskService := service.NewTaskService(taskRepo, blockerService, subTaskService, noteService, taskPolicyService)
	dailySessionService := service.NewDailySessionService(dailySessionRepo)

	// Initialize handlers
	taskHandler := handler.NewTaskHandler(taskService)
	subTaskHandler := handler.NewSubTaskHandler(subTaskService)
	noteHandler := handler.NewNoteHandler(noteService)
	blockerHandler := handler.NewBlockerHandler(blockerService)
	dailySessionHandler := handler.NewDailySessionHandler(dailySessionService)

	// Setup Gin server
	r := gin.Default()

	// Register routes (order matters for Gin to avoid parameter conflicts)
	subTaskHandler.Register(r)
	noteHandler.Register(r)
	blockerHandler.Register(r)
	taskHandler.Register(r)
	dailySessionHandler.Register(r)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
