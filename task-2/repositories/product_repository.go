package repositories

import (
	"database/sql"
	"errors"
	"task-2/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, c.name as category_name FROM products p INNER JOIN categories c ON p.category_id = c.id`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID.Name); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
	row := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID.ID)
	if err := row.Scan(&product.ID); err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, c.name as category_name FROM products p INNER JOIN categories c ON p.category_id = c.id WHERE p.id = $1`

	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5`
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID.ID, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}
