package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/istvzsig/eve-gate-watch/external/sde"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) SearchSystems(
	ctx context.Context,
	query string,
) ([]System, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name,
			constellation_id,
			region_id,
			security_status,
			wormhole_class_id
		FROM systems
		WHERE name ILIKE $1
		ORDER BY name
		LIMIT 25
	`, query+"%")
	if err != nil {
		return nil, fmt.Errorf("search systems: %w", err)
	}
	defer rows.Close()

	systems := make([]System, 0, 25)

	for rows.Next() {
		var system System

		if err := rows.Scan(
			&system.ID,
			&system.Name,
			&system.ConstellationID,
			&system.RegionID,
			&system.SecurityStatus,
			&system.WormholeClassID,
		); err != nil {
			return nil, fmt.Errorf("scan system: %w", err)
		}

		systems = append(systems, system)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate systems: %w", err)
	}

	return systems, nil
}

func (r *PostgresRepository) UpsertRegion(
	ctx context.Context,
	region sde.Region,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO regions (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id)
		DO UPDATE SET name = EXCLUDED.name
	`,
		region.ID,
		region.Name,
	)

	if err != nil {
		return fmt.Errorf("upsert region %d: %w", region.ID, err)
	}

	return nil
}

func (r *PostgresRepository) UpsertConstellation(
	ctx context.Context,
	constellation sde.Constellation,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO constellations (
			id,
			name,
			region_id
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			region_id = EXCLUDED.region_id
	`,
		constellation.ID,
		constellation.Name,
		constellation.RegionID,
	)

	if err != nil {
		return fmt.Errorf(
			"upsert constellation %d: %w",
			constellation.ID,
			err,
		)
	}

	return nil
}

func (r *PostgresRepository) UpsertSystem(
	ctx context.Context,
	system sde.SolarSystem,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO systems (
			id,
			name,
			constellation_id,
			region_id,
			security_status,
			wormhole_class_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			constellation_id = EXCLUDED.constellation_id,
			region_id = EXCLUDED.region_id,
			security_status = EXCLUDED.security_status,
			wormhole_class_id = EXCLUDED.wormhole_class_id
	`,
		system.ID,
		system.Name,
		system.ConstellationID,
		system.RegionID,
		system.SecurityStatus,
		system.WormholeClassID,
	)

	if err != nil {
		return fmt.Errorf(
			"upsert system %d (%s): %w",
			system.ID,
			system.Name,
			err,
		)
	}

	return nil
}
