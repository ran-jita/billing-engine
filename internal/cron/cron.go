package cron

import (
	"context"
	"github.com/ran-jita/billing-engine/internal/domain"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"

	"github.com/ran-jita/billing-engine/internal/repository"
	"github.com/ran-jita/billing-engine/internal/usecase"
)

type CronJobs struct {
	cron            *cron.Cron
	db              *sqlx.DB
	borrowerUsecase *usecase.BorrowerUsecase
}

func NewCronJobs(db *sqlx.DB) *CronJobs {
	billRepository := repository.NewBillRepository(db)
	loanRepository := repository.NewLoanRepository(db)
	loanDomain := domain.NewLoanDomain(loanRepository, billRepository)

	borrowerRepository := repository.NewBorrowerRepository(db)
	borrowerDomain := domain.NewBorrowerDomain(borrowerRepository)
	borrowerUsecase := usecase.NewBorrowerUsecase(borrowerDomain, loanDomain)

	return &CronJobs{
		cron:            cron.New(),
		db:              db,
		borrowerUsecase: borrowerUsecase,
	}
}

// Start memulai semua cron jobs
func (c *CronJobs) Start() {
	log.Println("Starting cron jobs...")

	// Job 1: Check overdue billings - Setiap hari jam 00:01
	c.cron.AddFunc("1 0 * * *", func() {
		log.Println("Running check overdue billings...")
		if err := c.CheckOverdueBillings(); err != nil {
			log.Printf("Error checking overdue billings: %v", err)
		}
	})

	c.cron.Start()
	log.Println("Cron jobs started successfully")
}

// Stop menghentikan semua cron jobs
func (c *CronJobs) Stop() {
	log.Println("Stopping cron jobs...")
	ctx := c.cron.Stop()
	<-ctx.Done()
	log.Println("Cron jobs stopped")
}

func (c *CronJobs) CheckOverdueBillings() error {
	ctx := context.Background()
	processDate := time.Now()
	return c.borrowerUsecase.UpdateDelinquent(ctx, processDate)
}
