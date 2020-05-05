package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/didiyudha/sse-redis/domain/product/model"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

// ProductCache - product cache APIs.
type ProductCache interface {
	Store(product *model.Product) error
	GetByID(id uuid.UUID) (model.Product, error)
	Streams(ctx context.Context, id uuid.UUID, prodChan chan model.Product)
}

type productCache struct {
	Conn redis.Conn
}

// NewProductCache is a factory function of product cache.
func NewProductCache(conn redis.Conn) ProductCache {
	return &productCache{
		Conn: conn,
	}
}

func (p *productCache) Store(product *model.Product) error {
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("product-%s", product.ID)
	_, err = p.Conn.Do("SET", key, string(b))
	return err
}

func (p *productCache) GetByID(id uuid.UUID) (model.Product, error) {
	key := fmt.Sprintf("product-%s", id)
	b, err := redis.Bytes(p.Conn.Do("GET", key))
	if err != nil {
		return model.Product{}, err
	}
	var product model.Product
	if err := json.Unmarshal(b, &product); err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (p *productCache) Streams(ctx context.Context, id uuid.UUID, prodChan chan model.Product) {
	psc := redis.PubSubConn{Conn: p.Conn}
	key := fmt.Sprintf("product-%s", id)
	keyspace := fmt.Sprintf("__keyspace@*__:%s", key)

	psc.PSubscribe(keyspace, "set")

	recvChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	go func() {
		for {
			switch m := psc.Receive().(type) {
			case redis.Message:
				log.Printf("message: %v\n", m)
				recvChan <- m
			case error:
				errChan <- fmt.Errorf("error receive message %v", m)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// Client closed the connection.
			log.Printf("client closed connection: %v\n", ctx.Err())
			break
		case <-recvChan:
			product, err := p.GetByID(id)
			if err != nil {
				log.Println(err)
			}
			prodChan <- product
		case err := <-errChan:
			log.Println(err)
		}
	}
}
