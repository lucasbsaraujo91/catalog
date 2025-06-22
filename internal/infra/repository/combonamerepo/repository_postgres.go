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

func (r *postgresRepository) GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error) {
	var combos []comboname.ComboName
	var total int64

	offset := (page - 1) * limit

	// Conta total
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM combo_names").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Busca paginada
	query := `SELECT id, name, uuid, nickname, is_available FROM combo_names ORDER BY id ASC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var combo comboname.ComboName
		err := rows.Scan(&combo.ID, &combo.Name, &combo.ComboNameUuid, &combo.Nickname, &combo.IsAvailable)
		if err != nil {
			return nil, 0, err
		}
		combos = append(combos, combo)
	}

	return combos, total, nil
}
