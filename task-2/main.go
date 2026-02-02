package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"task-2/database"
	"task-2/handlers"
	"task-2/repositories"
	"task-2/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"DATABASE_PORT"`
	DBConn string `mapstructure:"DATABASE_CONN"`
}

func main() {
	// Viper Config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("DATABASE_PORT"),
		DBConn: viper.GetString("DATABASE_CONN"),
	}

	// Database Setup
	db, err := database.InitDb(config.DBConn)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on ", addr)

	// err = http.ListenAndServe(addr, nil)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	os.Exit(1)
	// }

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

}
