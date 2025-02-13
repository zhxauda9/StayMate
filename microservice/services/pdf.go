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
	pdf.SetFont("Arial", "B", 20)
	pdf.ImageOptions("assets/logo.png", 10, 10, 30, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	pdf.Cell(190, 10, "Invoice")
	pdf.Ln(15)

	// Company and Customer Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(100, 10, "Staymate Inc.")
	pdf.Cell(90, 10, fmt.Sprintf("Date: %s", time.Now().Format("02 Jan 2006")))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(100, 10, "123 Staymate Kabanbay batyra, Astana")
	pdf.Cell(90, 10, fmt.Sprintf("Transaction ID: %s", txID))
	pdf.Ln(8)
	pdf.Cell(100, 10, "Customer:")
	pdf.Cell(90, 10, fmt.Sprintf("User Email: %s", tx.UserEmail))
	pdf.Ln(8)
	pdf.Cell(100, 10, "Payment Method:")
	pdf.Cell(90, 10, fmt.Sprintf("%s (Card Ending: **** **** **** %s)", tx.PaymentMethod, tx.CardDetails[len(tx.CardDetails)-4:]))

	pdf.Ln(15)
	// Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(100, 10, "Description", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 10, "Amount (USD)", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	// Table Content (Example data, you can replace it with dynamic data)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(100, 10, "Room Booking (Staymate)", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 10, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(45, 10, fmt.Sprintf("%.2f", tx.Amount), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	// Totals
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(145, 10, "Total Amount:")
	pdf.CellFormat(45, 10, fmt.Sprintf("%.2f USD", tx.Amount), "1", 0, "R", false, 0, "")
	pdf.Ln(15)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Thank you for your business. Please contact us for any inquiries.")

	// Convert to bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
