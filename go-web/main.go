package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID           int    `json:"id"`
	Name         string `json:"nome"`
	Color        string `json:"cor"`
	Price        int    `json:"preco"`
	Stock        int    `json:"estoque"`
	Code         string `json:"codigo"`
	Publication  bool   `json:"publicacao"`
	CreationDate string `json:"data_de_criacao"`
}

var AllProducts []Product

func InstantiateProducts() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Error opening file")

		}
	}()
	jsonFile, err := os.Open("products.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &AllProducts)

}

func getAllProducts(c *gin.Context) {
	filteredProducts := []Product{}

	for _, p := range AllProducts {
		if c.Query("color") != p.Color {
			fmt.Printf("%v,%v", p.Color, c.Query("color"))
			continue
		}
		filteredProducts = append(filteredProducts, p)
	}
	c.JSON(http.StatusOK, gin.H{
		"products": filteredProducts,
	})
}
func getProductById(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	for _, p := range AllProducts {
		if p.ID == id {
			c.JSON(http.StatusOK, gin.H{"product": p})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func main() {
	InstantiateProducts()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"products": "Ola Lucas",
		})
	})
	products := r.Group("/products")
	{
		products.GET("/", getAllProducts)
		products.GET("/:id", getProductById)

	}

	r.Run()

}
