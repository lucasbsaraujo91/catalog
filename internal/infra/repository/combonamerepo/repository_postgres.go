package combonamerepo

import (
	"context"
	"database/sql"
	"log"

	comboname "catalog/internal/entity/comboname"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) comboname.Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*comboname.ComboName, error) {
	log.Printf("Repository → Searching combo with ID: %d", id)

	query := `SELECT id, name, uuid, nickname, is_available FROM combo_names WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var combo comboname.ComboName
	err := row.Scan(&combo.ID, &combo.Name, &combo.ComboNameUuid, &combo.Nickname, &combo.IsAvailable)
	if err != nil {
		log.Printf("Repository → Query error: %v", err)
		return nil, err
	}

	log.Printf("Repository → Found combo: %+v", combo)

	return &combo, nil
}
