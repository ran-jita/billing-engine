package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ran-jita/billing-engine/internal/domain"
	"github.com/ran-jita/billing-engine/internal/handler"
	"github.com/ran-jita/billing-engine/internal/repository"
	"github.com/ran-jita/billing-engine/internal/usecase"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/ran-jita/billing-engine/pkg/database"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Initialize database connection
	dbConfig := database.GetConfig()
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	pingHandler := handler.NewPingHandler()

	loanRepository := repository.NewLoanRepository(db)
	loanDomain := domain.NewLoanDomain(loanRepository)
	loanUsecase := usecase.NewLoanUsecase(loanDomain)
	loanHandler := handler.NewLoanHandler(loanUsecase)

	borrowerRepository := repository.NewBorrowerRepository(db)
	borrowerDomain := domain.NewBorrowerDomain(borrowerRepository)
	borrowerUsecase := usecase.NewBorrowerUsecase(borrowerDomain)
	borrowerHandler := handler.NewBorrowerHandler(borrowerUsecase)

	paymentRepository := repository.NewPaymentRepository(db)
	paymentDomain := domain.NewPaymentDomain(paymentRepository)
	paymentUsecase := usecase.NewPaymentUsecase(paymentDomain)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	// Set Gin mode (release/debug)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/ping", pingHandler.Ping)

	group := router.Group("/api/v1")
	{
		// Loan routes
		loans := group.Group("/loans")
		{
			loans.POST("", loanHandler.GetAll)
			loans.GET("/:id", loanHandler.GetById)
		}

		// Borrowers routes
		borrowers := group.Group("/borrowers")
		{
			borrowers.GET("/:id", borrowerHandler.GetById)
		}

		// Payments routes
		payments := group.Group("/payments")
		{
			payments.POST("", paymentHandler.Create)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

// healthCheck handler untuk cek koneksi database
func healthCheck(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ping database
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error","message":"database connection failed"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"service is healthy"}`))
	}
}
