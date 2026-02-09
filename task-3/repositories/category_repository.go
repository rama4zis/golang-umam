package repositories

import (
	"database/sql"
	"errors"
	"task-3/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var category models.Category
	if err := row.Scan(&category.ID, &category.Name); err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1)`
	result, err := r.db.Exec(query, category.Name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not created")
	}

	return nil
}

func (r *CategoryRepository) Update(category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	result, err := r.db.Exec(query, category.Name, category.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not updated")
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("category not found")
	}

	return nil
}
