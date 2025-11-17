package routes

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
)

type handlerCollection struct {
	pingHandler     *handler.PingHandler
	loanHandler     *handler.LoanHandler
	borrowerHandler *handler.BorrowerHandler
	paymentHandler  *handler.PaymentHandler
}

func InitHttpRoutes(db *sqlx.DB) {
	handlers := initHandler(db)

	// Set Gin mode (release/debug)
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/ping", handlers.pingHandler.Ping)

	group := router.Group("/api/v1")
	{
		// Loan routes
		loans := group.Group("/loans")
		{
			loans.POST("", handlers.loanHandler.GetAll)
			loans.GET("/:id", handlers.loanHandler.GetById)
		}

		// Borrowers routes
		borrowers := group.Group("/borrowers")
		{
			borrowers.GET("/:id", handlers.borrowerHandler.GetById)
			borrowers.PUT("/delinquent", handlers.borrowerHandler.UpdateStatusDelinquent)
		}

		// Payments routes
		payments := group.Group("/payments")
		{
			payments.POST("", handlers.paymentHandler.Create)
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

func initHandler(db *sqlx.DB) handlerCollection {
	var collection handlerCollection
	collection.pingHandler = handler.NewPingHandler()

	billRepository := repository.NewBillRepository(db)
	loanRepository := repository.NewLoanRepository(db)
	loanDomain := domain.NewLoanDomain(loanRepository, billRepository)
	loanUsecase := usecase.NewLoanUsecase(loanDomain)
	collection.loanHandler = handler.NewLoanHandler(loanUsecase)

	borrowerRepository := repository.NewBorrowerRepository(db)
	borrowerDomain := domain.NewBorrowerDomain(borrowerRepository)
	borrowerUsecase := usecase.NewBorrowerUsecase(borrowerDomain, loanDomain)
	collection.borrowerHandler = handler.NewBorrowerHandler(borrowerUsecase)

	paymentRepository := repository.NewPaymentRepository(db)
	paymentBillRepository := repository.NewPaymentBillRepository(db)
	paymentDomain := domain.NewPaymentDomain(paymentRepository, paymentBillRepository)
	paymentUsecase := usecase.NewPaymentUsecase(paymentDomain, loanDomain, db)
	collection.paymentHandler = handler.NewPaymentHandler(paymentUsecase)

	return collection
}
