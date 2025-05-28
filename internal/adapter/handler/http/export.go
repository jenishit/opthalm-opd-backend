package http

import (
	"bytes"
	"encoding/csv"

	"github.com/gin-gonic/gin"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
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
