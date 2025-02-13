package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zhxauda9/StayMate/microservice/models"
	"github.com/zhxauda9/StayMate/microservice/services"
)

type PaymentHandler struct {
	PGPool        *pgxpool.Pool
	EmailSettings services.EmailSettings
}

func (p *PaymentHandler) ProcessPayment(c *gin.Context) {
	var tx models.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request data"})
		return
	}
	tx.CreatedAt = time.Now()

	// Insert into Postgres
	// We'll return the new row's ID
	var insertedID int // your DB column is 'id serial primary key'
	insertSQL := `
        INSERT INTO transactions (user_email, amount, products, created_at, payment_method, card_details)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.PGPool.QueryRow(ctx, insertSQL,
		tx.UserEmail,
		tx.Amount,
		tx.Products, // TEXT[] in Postgres
		tx.CreatedAt,
		tx.PaymentMethod,
		tx.CardDetails,
	).Scan(&insertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save transaction"})
		return
	}

	// Convert insertedID to string
	insertedIDStr := /* e.g. */ (fmt.Sprintf("%d", insertedID))

	// Generate PDF
	pdfBytes, err := services.GeneratePDFReceipt(tx, insertedIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to generate PDF receipt"})
		return
	}

	// Email with PDF attachment
	emailBody := "Thank you for your purchase! Please find your receipt attached."
	err = services.SendReceiptEmail(p.EmailSettings, tx.UserEmail, "Your Receipt", emailBody, "receipt.pdf", pdfBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to send receipt email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"message":        "Payment processed",
		"transaction_id": insertedIDStr,
	})
}
