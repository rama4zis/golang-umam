package repositories

import (
	"database/sql"
	"errors"
	"task-3/dto"
	"task-3/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll(name string) ([]dto.ProductResponse, error) {
	query := `SELECT
		p.id,
		p.name,
		p.price,
		p.stock,
		p.category_id,
		c.name
	FROM products p
	JOIN categories c ON p.category_id = c.id`

	args := []interface{}{}

	if name != "" {
		query += " WHERE p.name LIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]dto.ProductResponse, 0)
	for rows.Next() {
		var product dto.ProductResponse
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (dto.ProductResponse, error) {
	query := `SELECT
		p.id,
		p.name,
		p.price,
		p.stock,
		p.category_id,
		c.name
	FROM products p
	JOIN categories c ON p.category_id = c.id
	WHERE p.id = $1`

	row := r.DB.QueryRow(query, id)

	var product dto.ProductResponse
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID, &product.CategoryName); err != nil {
		return dto.ProductResponse{}, err
	}
	return product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)`
	result, err := r.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Product not created")
	}

	return nil
}

func (r *ProductRepository) Update(product *models.Product) error {
	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5`
	result, err := r.DB.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("error, product not updated")
	}

	return nil
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("error, product not found")
	}

	return nil
}
