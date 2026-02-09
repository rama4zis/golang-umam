package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"task-3/dto"
	"task-3/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("product id %d not enough stock", item.ProductID)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	log.Println("Total Amount: ", totalAmount)
	log.Println("Details: ", details)

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		log.Println("Error Insert Transaction: ", err)
		return nil, err
	}

	// Log id
	log.Println("Transaction ID: ", transactionID)

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			log.Println("Error Insert Transaction Detail: ", err)
			return nil, err
		}
	}

	// UPDATE STOCK
	for i := range details {
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", details[i].Quantity, details[i].ProductID)
		if err != nil {
			log.Println("Error Update Stock: ", err)
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) ReportToday(start_date string, end_date string) (*dto.ReportTodayResponse, error) {
	query := `
		SELECT
			COALESCE(SUM(tr.total_amount), 0) AS total_revenue,
			COALESCE(COUNT(*), 0) AS total_transaction,
			COALESCE(MAX(p.name), '') AS best_selling_product_name,
			COALESCE(MAX(td.quantity), 0) AS best_selling_product_qty
		FROM transactions tr
		LEFT JOIN transaction_details td ON tr.id = td.transaction_id
		LEFT JOIN products p ON td.product_id = p.id
	`

	args := []interface{}{}

	if start_date != "" {
		query += " WHERE DATE(tr.created_at) BETWEEN CAST($1 AS DATE) AND CAST($2 AS DATE)"
		args = append(args, start_date, end_date)
	} else {
		query += " WHERE DATE(tr.created_at) = DATE(CURRENT_DATE)"
	}

	row := repo.db.QueryRow(query, args...)

	var response dto.ReportTodayResponse
	if err := row.Scan(&response.TotalRevenue, &response.TotalTransaction, &response.BestSellingProduct.Name, &response.BestSellingProduct.Quantity); err != nil {
		return nil, err
	}

	return &response, nil
}
