package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/zhxauda9/StayMate/microservice/models"
)

func GeneratePDFReceipt(tx models.Transaction, txID string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Company: MyLibrary Inc.")
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", txID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date & Time: %s", time.Now().Format(time.RFC1123)))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: %.2f", tx.Amount))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", tx.UserEmail))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Payment Method: %s (Card details encrypted)", tx.PaymentMethod))

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
