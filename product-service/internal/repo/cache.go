package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(addr string) *CacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &CacheRepository{
		client: client,
	}
}

func (c *CacheRepository) GetProduct(ctx context.Context, id uint) (*model.Product, error) {
	key := fmt.Sprintf("product:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *CacheRepository) SetProduct(ctx context.Context, product *model.Product) error {
	key := fmt.Sprintf("product:%d", product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// Cache for 10 minutes
	return c.client.Set(ctx, key, data, 10*time.Minute).Err()
}

func (c *CacheRepository) InvalidateProduct(ctx context.Context, id uint) error {
	key := fmt.Sprintf("product:%d", id)
	return c.client.Del(ctx, key).Err()
}
