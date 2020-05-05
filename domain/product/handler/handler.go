package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/didiyudha/sse-redis/domain/product/model"
	"github.com/didiyudha/sse-redis/domain/product/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ProductHandler - HTTP product handler.
type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

// NewProductHandler - a factory function of product handler.
func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: productUseCase,
	}
}

// Store product handler.
func (p *ProductHandler) Store(c echo.Context) error {
	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Qty      int    `json:"qty"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return err
	}
	product := model.Product{
		ID:        uuid.New(),
		Name:      payload.Name,
		Category:  payload.Category,
		Qty:       payload.Qty,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	if err := p.ProductUseCase.Store(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Streams a product update.
func (p *ProductHandler) Streams(c echo.Context) error {
	productID := c.Param("productId")
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	prodChan := make(chan model.Product, 1)

	go p.ProductUseCase.StreamProduct(ctx, productUUID, prodChan)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(c.Response())

	select {
	case <-ctx.Done():
		return nil
	default:
		for p := range prodChan {
			log.Printf("product p: %+v\n", p)
			if err := enc.Encode(p); err != nil {
				return err
			}
			c.Response().Flush()
		}
	}
	return nil
}
