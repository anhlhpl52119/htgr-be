package store

import (
	"database/sql"
	"errors"
	"htrr-apis/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresRestaurantStore struct {
	db *sql.DB
}

func NewPostgresRestaurantStore(pgDB *sql.DB) *PostgresRestaurantStore {
	return &PostgresRestaurantStore{
		db: pgDB,
	}
}

type Restaurant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchRestaurantParams struct {
	Page     int
	PageSize int
	Name     string
}

type BulkDeleteResult struct {
	DeletedCount int
	FailedCount  int
	DeletedIDs   []string
	FailedIDs    []string
}

type RestaurantStore interface {
	Create(*Restaurant) error
	Search(SearchRestaurantParams) ([]Restaurant, int, error)
	Update(*Restaurant) error
	GetRestaurantById(string) (*Restaurant, error)
	Delete(string) error
	BulkDeleteAtomic([]string) (int, error)
	BulkDeletePartial([]string) (*BulkDeleteResult, error)
	BulkDeleteBestEffort([]string) (int, error)
}

func (pg *PostgresRestaurantStore) Create(restaurant *Restaurant) error {
	q := `
	INSERT INTO restaurants (name, address, phone, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, created_at
	`
	err := pg.db.QueryRow(q,
		restaurant.Name,
		restaurant.Address,
		restaurant.Phone,
		restaurant.IsActive).
		Scan(&restaurant.ID,
			&restaurant.Name,
			&restaurant.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresRestaurantStore) Search(params SearchRestaurantParams) ([]Restaurant, int, error) {
	q := `
	SELECT id, name, address, phone, is_active, created_at, updated_at,
			COUNT(*) OVER()
	FROM restaurants
	WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
	ORDER BY name
	LIMIT $2 OFFSET $3
	`
	limit, offset := utils.GetOffset(&params.Page, &params.PageSize)

	row, err := pg.db.Query(q, params.Name, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer row.Close()

	var total int
	var list []Restaurant
	for row.Next() {
		var rtr Restaurant
		err := row.Scan(
			&rtr.ID,
			&rtr.Name,
			&rtr.Address,
			&rtr.Phone,
			&rtr.IsActive,
			&rtr.CreatedAt,
			&rtr.UpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, rtr)
	}

	return list, total, nil
}

func (pg *PostgresRestaurantStore) Update(restaurant *Restaurant) error {
	_, err := uuid.Parse(restaurant.ID)
	if err != nil {
		return errors.New("Invalid id format")
	}

	q := `
	UPDATE restaurants
	SET name = $1, address = $2, phone = $3, is_active = $4
	WHERE id = $5
	`

	_, err = pg.db.Exec(q,
		restaurant.Name,
		restaurant.Address,
		restaurant.Phone,
		restaurant.IsActive,
		restaurant.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresRestaurantStore) GetRestaurantById(id string) (*Restaurant, error) {
	q := `
	SELECT id, name, address, phone, is_active, created_at, updated_at
	FROM restaurants
	WHERE id = $1
	`

	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("Invalid id format")
	}

	restaurant := &Restaurant{}
	err = pg.db.QueryRow(q, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Address,
		&restaurant.Phone,
		&restaurant.IsActive,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (pg *PostgresRestaurantStore) Delete(id string) error {
	q := `
	DELETE FROM restaurants WHERE id = $1
	`

	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("Invalid id format")
	}

	result, err := pg.db.Exec(q, id)
	if err != nil {
		return err
	}

	rowEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowEffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresRestaurantStore) BulkDeleteAtomic(ids []string) (int, error) {
	// Validate all IDs upfront
	for _, id := range ids {
		_, err := uuid.Parse(id)
		if err != nil {
			return 0, errors.New("Invalid id format")
		}
	}

	// Start transaction
	tx, err := pg.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Delete with ANY clause for all IDs
	q := `DELETE FROM restaurants WHERE id = ANY($1)`
	result, err := tx.Exec(q, pq.Array(ids))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// If no rows were deleted, return error
	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (pg *PostgresRestaurantStore) BulkDeletePartial(ids []string) (*BulkDeleteResult, error) {
	result := &BulkDeleteResult{
		DeletedIDs: []string{},
		FailedIDs:  []string{},
	}

	// Separate valid and invalid IDs
	validIDs := []string{}
	for _, id := range ids {
		_, err := uuid.Parse(id)
		if err != nil {
			result.FailedIDs = append(result.FailedIDs, id)
			result.FailedCount++
			continue
		}
		validIDs = append(validIDs, id)
	}

	// If no valid IDs, return early
	if len(validIDs) == 0 {
		return result, nil
	}

	// Start transaction
	tx, err := pg.db.Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	// First, query which valid IDs exist in the database
	existingQuery := `SELECT id FROM restaurants WHERE id = ANY($1)`
	rows, err := tx.Query(existingQuery, pq.Array(validIDs))
	if err != nil {
		return result, err
	}

	existingSet := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return result, err
		}
		existingSet[id] = true
	}
	rows.Close()

	// Identify which valid IDs don't exist
	for _, id := range validIDs {
		if !existingSet[id] {
			result.FailedIDs = append(result.FailedIDs, id)
			result.FailedCount++
		}
	}

	// Delete the valid IDs that exist
	if len(existingSet) > 0 {
		deleteQuery := `DELETE FROM restaurants WHERE id = ANY($1)`
		deleteResult, err := tx.Exec(deleteQuery, pq.Array(validIDs))
		if err != nil {
			return result, err
		}

		rowsAffected, err := deleteResult.RowsAffected()
		if err != nil {
			return result, err
		}

		result.DeletedCount = int(rowsAffected)

		// Track which IDs were deleted
		for id := range existingSet {
			result.DeletedIDs = append(result.DeletedIDs, id)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

func (pg *PostgresRestaurantStore) BulkDeleteBestEffort(ids []string) (int, error) {
	// Filter to only valid UUID IDs
	validIDs := []string{}
	for _, id := range ids {
		_, err := uuid.Parse(id)
		if err != nil {
			// Skip invalid IDs silently
			continue
		}
		validIDs = append(validIDs, id)
	}

	// If no valid IDs, return 0
	if len(validIDs) == 0 {
		return 0, nil
	}

	// Delete valid IDs
	q := `DELETE FROM restaurants WHERE id = ANY($1)`
	result, err := pg.db.Exec(q, pq.Array(validIDs))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
