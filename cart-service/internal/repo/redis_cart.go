package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/phanthehoang2503/small-project/cart-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisCartRepo struct {
	client *redis.Client
}

func NewRedisCartRepo(addr string) *RedisCartRepo {
	return &RedisCartRepo{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisCartRepo) AddNewItems(i *model.Cart) (model.Cart, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", i.UserID)
	field := strconv.Itoa(int(i.ProductID))

	// Get existing item if any
	var current model.Cart
	val, err := r.client.HGet(ctx, key, field).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(val), &current)
		// Update existing
		current.Quantity += i.Quantity
		current.Subtotal = current.Price * int64(current.Quantity)
		current.UpdatedAt = time.Now()
	} else {
		// Create new
		i.ID = i.ProductID // Use ProductID as ID in Redis context
		i.Subtotal = i.Price * int64(i.Quantity)
		i.CreatedAt = time.Now()
		i.UpdatedAt = time.Now()
		current = *i
	}

	data, err := json.Marshal(current)
	if err != nil {
		return model.Cart{}, err
	}

	// Save to hash
	if err := r.client.HSet(ctx, key, field, data).Err(); err != nil {
		return model.Cart{}, err
	}

	// Set TTL for cart (e.g. 24 hours)
	r.client.Expire(ctx, key, 24*time.Hour)

	return current, nil
}

func (r *RedisCartRepo) List(UserID uint) ([]model.Cart, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", UserID)

	val, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	items := make([]model.Cart, 0, len(val))
	for _, v := range val {
		var item model.Cart
		if err := json.Unmarshal([]byte(v), &item); err == nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *RedisCartRepo) UpdateQuantity(userID, id uint, quantity int) (model.Cart, error) {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", userID)
	field := strconv.Itoa(int(id))

	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		return model.Cart{}, fmt.Errorf("item not found")
	}

	var item model.Cart
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return model.Cart{}, err
	}

	if quantity <= 0 {
		r.client.HDel(ctx, key, field)
		return model.Cart{}, nil
	}

	item.Quantity = quantity
	item.Subtotal = item.Price * int64(quantity)
	item.UpdatedAt = time.Now()

	data, _ := json.Marshal(item)
	if err := r.client.HSet(ctx, key, field, data).Err(); err != nil {
		return model.Cart{}, err
	}

	r.client.Expire(ctx, key, 24*time.Hour)
	return item, nil
}

func (r *RedisCartRepo) Remove(UserID, id uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", UserID)
	// 'id' here is interpreted as ProductID
	field := strconv.Itoa(int(id))
	return r.client.HDel(ctx, key, field).Err()
}

func (r *RedisCartRepo) ClearCart(userID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("cart:%d", userID)
	return r.client.Del(ctx, key).Err()
}
