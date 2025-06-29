package combonamerepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	comboname "catalog/internal/entity/comboname"
	"catalog/internal/infra/cache/redis"
)

type postgresRepository struct {
	db    *sql.DB
	cache *redis.RedisCache
}

func NewPostgresRepository(db *sql.DB, cache *redis.RedisCache) comboname.Repository {
	return &postgresRepository{
		db:    db,
		cache: cache,
	}
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*comboname.ComboName, error) {
	cacheKey := redis.BuildKey("comboname", fmt.Sprint(id))

	if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
		var combo comboname.ComboName
		if err := json.Unmarshal([]byte(cached), &combo); err == nil {
			return &combo, nil
		}
	}

	query := `SELECT id, name, uuid, nickname, is_available FROM combo_names WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var combo comboname.ComboName
	err := row.Scan(&combo.ID, &combo.Name, &combo.ComboNameUuid, &combo.Nickname, &combo.IsAvailable)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(combo)
	_ = r.cache.Set(ctx, cacheKey, string(data))

	return &combo, nil
}

func (r *postgresRepository) GetAll(ctx context.Context, page, limit int) ([]comboname.ComboName, int64, error) {
	cacheKey := redis.BuildKey("comboname", "all", "page", fmt.Sprint(page), "limit", fmt.Sprint(limit))

	if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
		var result struct {
			Combos []comboname.ComboName
			Total  int64
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result.Combos, result.Total, nil
		}
	}

	offset := (page - 1) * limit
	query := `SELECT id, name, uuid, nickname, is_available FROM combo_names ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var combos []comboname.ComboName
	for rows.Next() {
		var combo comboname.ComboName
		if err := rows.Scan(&combo.ID, &combo.Name, &combo.ComboNameUuid, &combo.Nickname, &combo.IsAvailable); err != nil {
			return nil, 0, err
		}
		combos = append(combos, combo)
	}

	var total int64
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM combo_names`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	payload := struct {
		Combos []comboname.ComboName
		Total  int64
	}{
		Combos: combos,
		Total:  total,
	}
	data, _ := json.Marshal(payload)
	_ = r.cache.Set(ctx, cacheKey, string(data))

	return combos, total, nil
}

func (r *postgresRepository) Create(ctx context.Context, combo *comboname.ComboName) (int64, error) {
	query := `INSERT INTO combo_names (name, uuid, nickname, is_available) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, combo.Name, combo.ComboNameUuid, combo.Nickname, combo.IsAvailable).Scan(&combo.ID)
	if err != nil {
		return 0, err
	}

	_ = r.cache.DeleteByPrefix(ctx, redis.BuildKey("comboname"))
	return combo.ID, nil
}

func (r *postgresRepository) Update(ctx context.Context, combo *comboname.ComboName) error {
	query := `UPDATE combo_names SET name = $1, uuid = $2, nickname = $3, is_available = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, combo.Name, combo.ComboNameUuid, combo.Nickname, combo.IsAvailable, combo.ID)
	if err != nil {
		return err
	}

	_ = r.cache.DeleteByPrefix(ctx, redis.BuildKey("comboname"))
	return nil
}

func (r *postgresRepository) Disable(ctx context.Context, id int64) error {
	query := `UPDATE combo_names SET is_available = false WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	_ = r.cache.DeleteByPrefix(ctx, redis.BuildKey("comboname"))
	return nil
}
