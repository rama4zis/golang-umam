package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"task-3/database"
	"task-3/handlers"
	"task-3/repositories"
	"task-3/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"DATABASE_PORT"`
	DBConn string `mapstructure:"DATABASE_CONN"`
}

func main() {
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

	db, err := database.InitDb(config.DBConn)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running on ", addr)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	TransactionHandlerr := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Category
	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	// Transaction
	http.HandleFunc("/api/transaction", TransactionHandlerr.HandleCheckout)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
