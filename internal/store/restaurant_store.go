package store

import (
	"database/sql"
	"errors"
	"htrr-apis/internal/utils"
	"time"

	"github.com/google/uuid"
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

type RestaurantStore interface {
	Create(*Restaurant) error
	Search(SearchRestaurantParams) ([]Restaurant, int, error)
	Update(*Restaurant) error
	GetRestaurantById(string) (*Restaurant, error)
	Delete(string) error
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

