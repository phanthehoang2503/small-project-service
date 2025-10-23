package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/product-service/internal/model"
	"github.com/phanthehoang2503/product-service/internal/store"
)

func RegisterRoutes(r *gin.Engine, s *store.Store) {
	api := r.Group("/api/v1")
	{
		api.GET("/products", listProducts(s))
		api.GET("/products/:id", getProducts(s))
		api.POST("/products", createProducts(s))
		api.PUT("/products/:id", updateProducts(s))
		api.DELETE("/products/:id", deleteProducts(s))
	}
}

func listProducts(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) { // gin.Context have all the func inside, docs: https://pkg.go.dev/github.com/gin-gonic/gin#Context
		c.JSON(200, s.List()) // status ok = 200
	}
}

func getProducts(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		p, ok := s.Get(id) // find the prod id
		if !ok {           // in case not found
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, p)
	}
}

func createProducts(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in model.Product
		if err := c.ShouldBindJSON(&in); err != nil { //bind attempts to json body to in var
			c.JSON(400, gin.H{"error": err.Error()}) // invalid json or mismatched fields
			return
		}

		created := s.Create(in)
		c.JSON(http.StatusCreated, created) // status created = 201
	}
}

func updateProducts(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		var in model.Product
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, ok := s.Update(id, in) //product's struct, boolean
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"}) // not found = 404
			return
		}

		c.JSON(200, updated)
	}
}

func deleteProducts(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		if !s.Delete(id) { //Delete return true if deleted an id and vice versa
			c.JSON(404, gin.H{"error": "invalid id"})
			return
		}

		c.Status(http.StatusNoContent) // status no content = 204
	}
}
