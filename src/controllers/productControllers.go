package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Create(&product)

	// clear cache so that frontend can get created data
	go database.ClearCache("products_frontend, products_backend")

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var product models.Product

	product.Id = uint(id)

	database.DB.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	// clear cache so that frontend can get updated data
	go database.ClearCache("products_frontend, products_backend")

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}
	product.Id = uint(id)

	database.DB.Delete(&product)

	// clear cache so that frontend can get updated data
	go database.ClearCache("products_frontend, products_backend")

	return nil
}

// redis cache
func ProductsFrontEnd(c *fiber.Ctx) error {
	var products []models.Product

	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()
	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(products)
}

func ProductsBackEnd(c *fiber.Ctx) error {
	var products []models.Product

	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_backend").Result()
	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute)
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	// search
	if s := c.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	} else {
		searchedProducts = products
	}

	// sortParam
	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price > searchedProducts[j].Price
			})
		}
	}

	// // pagination
	// var total = len(searchedProducts)
	// page, _ := strconv.Atoi(c.Query("page", "1"))
	// perPage := 9

	// var data []models.Product

	// if total <= page*perPage && total >= (page-1)*perPage {
	// 	data = searchedProducts[(page-1)*perPage : total]
	// } else if total >= page*perPage {
	// 	data = searchedProducts[(page-1)*perPage : page+perPage]
	// } else {
	// 	data = []models.Product{}
	// }

	// return c.JSON(fiber.Map{
	// 	"data":      data,
	// 	"total":     total,
	// 	"page":      page,
	// 	"last_page": total/perPage + 1,
	// })

	// pagination
	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage := 9

	// Calculate the start and end indices for the slice
	start := (page - 1) * perPage
	end := start + perPage

	// Ensure start and end indices are within valid bounds
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	// Slice the products accordingly
	var data []models.Product
	if start < total {
		data = searchedProducts[start:end]
	} else {
		data = []models.Product{}
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": (total + perPage - 1) / perPage,
	})
}
