package http

import (
	"bytes"
	"image/png"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http/dto"
	"github.com/jenish-brainztechs/go-backend/internal/core/domain"
	"github.com/jenish-brainztechs/go-backend/internal/core/port"
)

// ─── Inventory Items ─────────────────────────────────────────────

type InventoryHandler struct {
	svc port.InventoryItemService
}

func NewInventoryHandler(svc port.InventoryItemService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// CreateItem godoc
//
//	@Summary		Create an inventory item
//	@Description	Requires ROLE_ADMIN or ROLE_INVENTORY.
//	@Tags			inventory
//	@Accept			json
