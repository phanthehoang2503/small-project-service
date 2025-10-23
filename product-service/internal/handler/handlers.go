package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
)

func ListProducts(s *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) { // gin.Context have all the func inside, docs: https://pkg.go.dev/github.com/gin-gonic/gin#Context
		products, err := s.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // internal server error = 500
		}
		c.JSON(http.StatusOK, products) // status ok = 200
	}
}

func GetProducts(s *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		p, err := s.Get(id) // find the prod id
		if err != nil {     // in case not found
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, p)
	}
}

func CreateProducts(s *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in model.Product
		if err := c.ShouldBindJSON(&in); err != nil { //bind attempts to json body to in var
			c.JSON(400, gin.H{"error": err.Error()}) // invalid json or mismatched fields
			return
		}

		created, err := s.Create(in)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, created) // status created = 201
	}
}

func UpdateProducts(s *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		var in model.Product
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		updated, err := s.Update(id, in) //product's struct, boolean
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // not found = 404
			return
		}

		c.JSON(200, updated)
	}
}

func DeleteProducts(s *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		if s.Delete(id) != nil { //Delete return true if deleted an id and vice versa
			c.JSON(404, s.Delete(id).Error())
			return
		}

		c.Status(http.StatusNoContent) // status no content = 204
	}
}
