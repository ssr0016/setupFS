package main

// import (
// 	"ambassador/src/database"
// 	"ambassador/src/models"

// 	"math/rand"

// 	"github.com/bxcodec/faker/v3"
// )

// func main() {
// 	database.Connect()

// 	for i := 0; i < 30; i++ {
// 		product := models.Product{
// 			Title:       faker.Username(),
// 			Description: faker.Username(),
// 			Image:       faker.URL(),
// 			Price:       float64(rand.Intn(90) + 10),
// 		}

// 		database.DB.Create(&product)
// 	}
// }

// to run this faker docker compose exec backend sh
// go run src/commands/populateProducts.go
