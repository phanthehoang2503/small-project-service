package util

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserID extracts "user_id" from Gin context (set by JWT middleware).
// Supports int, uint, int64, float64, or string types.
// Returns error if missing or invalid.
func GetUserID(c *gin.Context) (uint, error) {
	uidI, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("missing user_id in context")
	}

	switch v := uidI.(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		u64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id string: %w", err)
		}
		return uint(u64), nil
	default:
		return 0, errors.New("invalid user_id type in context")
	}
}
