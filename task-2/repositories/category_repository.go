package repositories

import (
	"database/sql"
	"errors"
	"task-2/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := repo.db.Query(query)
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

func (repo *CategoryRepository) Create(category *models.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	row := repo.db.QueryRow(query, category.Name)
	if err := row.Scan(&category.ID); err != nil {
		return err
	}
	return nil
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = $1`
	var category models.Category
	err := repo.db.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *CategoryRepository) Update(category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	result, err := repo.db.Exec(query, category.Name, category.ID)
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

func (repo *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := repo.db.Exec(query, id)
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
