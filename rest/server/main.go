package main

import (
	"fmt"
	"net/http"
	"networks-course/rest/products"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create Gin router
	router := gin.Default()

	// Instantiate Product Handler and provide a data store
	store := products.NewMemStore()
	productsHandler := NewProductsHandler(store)

	// Register Routes
	router.GET("/", homePage)
	router.GET("/products", productsHandler.Listproducts)
	router.POST("/product", productsHandler.CreateProduct)
	router.POST("/product/:id/image", productsHandler.SetImage)
	router.GET("/product/:id", productsHandler.GetProduct)
	router.GET("/product/:id/image", productsHandler.GetProductImage)
	router.PUT("/product/:id", productsHandler.UpdateProduct)
	router.DELETE("/product/:id", productsHandler.DeleteProduct)

	// Start the server
	err := router.Run()
	if err != nil {
		panic(err.Error())
	}
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
}

type productsHandler struct {
	store productstore
}

func NewProductsHandler(s productstore) *productsHandler {
	return &productsHandler{
		store: s,
	}
}

type productstore interface {
	Add(Product products.Product) (products.ProductView, error)
	Get(id int) (products.ProductView, error)
	List() ([]products.ProductView, error)
	Update(id int, Product products.Product) (products.ProductView, error)
	Remove(id int) (products.ProductView, error)
}

func (h productsHandler) CreateProduct(c *gin.Context) {
	// Get request body and convert it to products.product
	var product products.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	view, err := h.store.Add(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// return success payload
	c.JSON(http.StatusOK, view)
}

func ImgPath(id int, fileName string) string {
	return fmt.Sprintf("img/%d/%s", id, fileName)
}

func (h productsHandler) SetImage(c *gin.Context) {
	id_string := c.Param("id")

	id, err_conv := strconv.Atoi(id_string)
	if err_conv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err_conv.Error()})
		return
	}

	img_dir := fmt.Sprintf("img/%s", id_string)
	os.Mkdir(img_dir, os.ModePerm)

	file, err := c.FormFile("icon")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := ImgPath(id, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, filename)
}

func (h productsHandler) Listproducts(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}

func (h productsHandler) GetProduct(c *gin.Context) {
	id_string := c.Param("id")

	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (h productsHandler) GetProductImage(c *gin.Context) {
	id_string := c.Param("id")

	id, err := strconv.Atoi(id_string)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	fileName := product.Icon
	content, err := os.ReadFile(ImgPath(id, product.Icon))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "image/png")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write([]byte(content))
	c.JSON(http.StatusOK, gin.H{
		"msg": "Download file successfully",
	})
}

func (h productsHandler) UpdateProduct(c *gin.Context) {
	// Get request body and convert it to products.product

	id_string := c.Param("id")

	id, err_conv := strconv.Atoi(id_string)
	if err_conv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err_conv.Error()})
		return
	}

	var product products.Product
	old_view, err := h.store.Get(id)
	if err != nil {
		if err == products.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product = products.ProductViewToProduct(&old_view)

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	view, err := h.store.Update(id, product)
	if err != nil {
		if err == products.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, view)
}

func (h productsHandler) DeleteProduct(c *gin.Context) {
	id_string := c.Param("id")

	id, err_conv := strconv.Atoi(id_string)
	if err_conv != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err_conv.Error()})
		return
	}

	view, err := h.store.Remove(id)
	if err != nil {
		if err == products.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, view)
}
