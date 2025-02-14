package handlers

import (
	"context"
	"fmt"
	"net/http"
	"payment-service/models"
	"payment-service/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PaymentHandler struct {
	PGPool        *pgxpool.Pool
	EmailSettings services.EmailSettings
	Logger        *zerolog.Logger
}

func (p *PaymentHandler) ProcessPayment(c *gin.Context) {
	startTime := time.Now()
	p.Logger.Info().Msg("Starting ProcessPayment handler")

	var request struct {
		UserID         int     `json:"user_id"`
		Email          string  `json:"email"`
		Amount         float64 `json:"amount"`
		CardNumber     string  `json:"card_number"`
		ExpirationDate string  `json:"expiration_date"`
		CVV            string  `json:"cvv"`
	}

	// Bind and validate the incoming JSON payload
	if err := c.ShouldBindJSON(&request); err != nil {
		p.Logger.Error().Err(err).Msg("Failed to bind request JSON")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request data"})
		return
	}

	p.Logger.Info().
		Int("user_id", request.UserID).
		Str("email", request.Email).
		Float64("amount", request.Amount).
		Str("expiration_date", request.ExpirationDate).
		Msg("Received payment request")

	// Create a transaction record
	tx := models.Transaction{
		UserEmail:     request.Email,
		Amount:        request.Amount,
		CreatedAt:     time.Now(),
		PaymentMethod: "credit_card", // Assume credit card for now
		CardDetails:   maskCardNumber(request.CardNumber),
	}

	p.Logger.Info().
		Str("user_email", tx.UserEmail).
		Float64("amount", tx.Amount).
		Str("payment_method", tx.PaymentMethod).
		Msg("Preparing to insert transaction into database")

	// Insert transaction into the database
	var insertedID int
	insertSQL := `
		INSERT INTO transactions (user_email, amount, created_at, payment_method, card_details)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.PGPool.QueryRow(ctx, insertSQL,
		tx.UserEmail,
		tx.Amount,
		tx.CreatedAt,
		tx.PaymentMethod,
		tx.CardDetails,
	).Scan(&insertedID)
	if err != nil {
		p.Logger.Error().Err(err).Msg("Failed to insert transaction into database")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save transaction"})
		return
	}

	p.Logger.Info().Int("transaction_id", insertedID).Msg("Transaction successfully saved in database")

	// Generate PDF receipt
	insertedIDStr := fmt.Sprintf("%d", insertedID)
	p.Logger.Info().Msg("Generating PDF receipt")

	pdfBytes, err := services.GeneratePDFReceipt(tx, insertedIDStr)
	if err != nil {
		p.Logger.Error().Err(err).Msg("Failed to generate PDF receipt")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate PDF receipt"})
		return
	}

	p.Logger.Info().Msg("PDF receipt generated successfully")

	// Send the receipt by email
	emailBody := "Thank you for your purchase! Please find your receipt attached."
	p.Logger.Info().Str("recipient_email", request.Email).Msg("Sending receipt email")

	err = services.SendReceiptEmail(p.EmailSettings, request.Email, "Your Receipt", emailBody, "receipt.pdf", pdfBytes)
	if err != nil {
		p.Logger.Error().Err(err).Msg("Failed to send receipt email")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to send receipt email"})
		return
	}

	p.Logger.Info().Msg("Receipt email sent successfully")

	// Return success response
	duration := time.Since(startTime)
	p.Logger.Info().
		Int("transaction_id", insertedID).
		Str("status", "success").
		Dur("duration", duration).
		Msg("Payment processed successfully")

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"message":        "Payment processed successfully",
		"transaction_id": insertedIDStr,
	})
}

func (p *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	email := c.Query("email") // Get the email from query parameters

	if email == "" {
		p.Logger.Error().Msg("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Email is required"})
		return
	}

	p.Logger.Info().Str("email", email).Msg("Fetching payment history")

	// Query to fetch transactions based on email
	query := `
		SELECT id, user_email, amount, created_at, payment_method, card_details
		FROM transactions
		WHERE user_email = $1
		ORDER BY created_at DESC
	`

	// Execute the query
	rows, err := p.PGPool.Query(context.Background(), query, email)
	if err != nil {
		p.Logger.Error().Err(err).Msg("Failed to fetch transactions from database")
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to fetch transaction history"})
		return
	}
	defer rows.Close()

	// Collect transactions
	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(&tx.ID, &tx.UserEmail, &tx.Amount, &tx.CreatedAt, &tx.PaymentMethod, &tx.CardDetails)
		if err != nil {
			p.Logger.Error().Err(err).Msg("Failed to scan transaction row")
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error processing transaction history"})
			return
		}
		transactions = append(transactions, tx)
	}

	// Check if no transactions were found
	if len(transactions) == 0 {
		p.Logger.Info().Str("email", email).Msg("No transactions found for this email")
		c.JSON(http.StatusNotFound, gin.H{"status": "success", "message": "No transactions found", "transactions": []models.Transaction{}})
		return
	}

	// Return transactions as JSON
	p.Logger.Info().Int("transaction_count", len(transactions)).Str("email", email).Msg("Returning transaction history")
	c.JSON(http.StatusOK, gin.H{"status": "success", "transactions": transactions})
}

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	return fmt.Sprintf("**** **** **** %s", cardNumber[len(cardNumber)-4:])
}
