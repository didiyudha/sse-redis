package usecase

import (
	"context"

	"github.com/didiyudha/sse-redis/domain/product/model"
	"github.com/didiyudha/sse-redis/domain/product/repository"
	"github.com/google/uuid"
)

type ProductUseCase interface {
	Store(product *model.Product) error
	StreamProduct(ctx context.Context, id uuid.UUID, prodChan chan model.Product)
}

type productUseCase struct {
	ProductCache repository.ProductCache
}

func NewProductUseCase(productCache repository.ProductCache) ProductUseCase {
	return &productUseCase{
		ProductCache: productCache,
	}
}

func (p *productUseCase) Store(product *model.Product) error {
	return p.ProductCache.Store(product)
}

func (p *productUseCase) StreamProduct(ctx context.Context, id uuid.UUID, prodChan chan model.Product) {
	p.ProductCache.Streams(ctx, id, prodChan)
}
