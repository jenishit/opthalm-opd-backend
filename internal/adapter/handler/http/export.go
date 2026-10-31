package http

import (
	"bytes"
	"encoding/csv"

	"github.com/gin-gonic/gin"
	"github.com/jenishit/opthalm-opd-backend/internal/core/domain"
	"github.com/jung-kurt/gofpdf"
)

// respondTable renders headers/rows as JSON (default), CSV, or PDF depending
// on the ?format= query param. jsonData is what's sent for the JSON case
// (kept separate from headers/rows so JSON responses can stay structured
// instead of stringly-typed).
func respondTable(ctx *gin.Context, filename, title string, headers []string, rows [][]string, jsonData any) {
	switch ctx.Query("format") {
	case "csv":
		writeCSV(ctx, filename, headers, rows)
	case "pdf":
		writePDFTable(ctx, title, headers, rows)
	default:
		handleSuccess(ctx, jsonData)
	}
}

func writeCSV(ctx *gin.Context, filename string, headers []string, rows [][]string) {
	ctx.Header("Content-Disposition", "attachment; filename=\""+filename+".csv\"")
	ctx.Header("Content-Type", "text/csv")

	w := csv.NewWriter(ctx.Writer)
	_ = w.Write(headers)
	for _, row := range rows {
		_ = w.Write(row)
	}
	w.Flush()
}

func writePDFTable(ctx *gin.Context, title string, headers []string, rows [][]string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, title)
	pdf.Ln(12)

	colWidth := 190.0 / float64(len(headers))

	pdf.SetFont("Arial", "B", 9)
	for _, h := range headers {
		pdf.CellFormat(colWidth, 7, h, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		for _, cell := range row {
			pdf.CellFormat(colWidth, 7, cell, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		handleError(ctx, domain.ErrInternal)
		return
	}
	ctx.Data(200, "application/pdf", buf.Bytes())
}
