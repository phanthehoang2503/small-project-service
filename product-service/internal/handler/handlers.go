package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanthehoang2503/small-project/product-service/internal/model"
	"github.com/phanthehoang2503/small-project/product-service/internal/repo"
)

// ListProducts godoc
// @Summary List all products
// @Description Returns all products available in the store
// @Tags Products
// @Produce json
// @Success 200 {array} model.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func ListProducts(r *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) { // gin.Context have all the func inside, docs: https://pkg.go.dev/github.com/gin-gonic/gin#Context
		products, err := r.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // internal server error = 500
		}
		c.JSON(http.StatusOK, products) // status ok = 200
	}
}

// GetProducts godoc
// @Summary Get a product by ID
// @Description Returns product information based on the ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func GetProducts(r *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
			return
		}

		p, err := r.Get(id) // find the prod id
		if err != nil {     // in case not found
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, p)
	}
}

// CreateProducts godoc
// @Summary Create a new product
// @Description Add a new product to the store
// @Tags Products
// @Accept json
// @Produce json
// @Param payload body model.Product true "Product payload"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func CreateProducts(r *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in model.Product
		if err := c.ShouldBindJSON(&in); err != nil { //bind attempts to json body to in var
			c.JSON(400, gin.H{"error": err.Error()}) // invalid json or mismatched fields
			return
		}

		created, err := r.Create(in)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, created) // status created = 201
	}
}

// UpdateProducts godoc
// @Summary Update an existing product
// @Description Update product information by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param payload body model.Product true "Updated product data"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func UpdateProducts(r *repo.Database) gin.HandlerFunc {
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

		updated, err := r.Update(id, in) //product's struct, boolean
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // not found = 404
			return
		}

		c.JSON(200, updated)
	}
}

// DeleteProducts godoc
// @Summary Delete a product
// @Description Remove a product by ID
// @Tags Products
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func DeleteProducts(r *repo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqStr := c.Param("id")
		id, err := strconv.ParseInt(reqStr, 10, 64) // int, error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}) // status bad request = 400
		}

		if r.Delete(id) != nil { //Delete return true if deleted an id and vice versa
			c.JSON(404, r.Delete(id).Error())
			return
		}

		c.Status(http.StatusNoContent) // status no content = 204
	}
}
