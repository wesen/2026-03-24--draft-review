package articleassets

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateAssetRecord(ctx context.Context, ownerUserID, articleID string, input AssetRecordInput) (*AssetRecord, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	assetID, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid asset id")
	}

	var asset AssetRecord
	err = r.pool.QueryRow(ctx, `
insert into article_assets (
    id,
    article_id,
    storage_key,
    original_filename,
    content_type,
    byte_size
)
select
    $3,
    a.id,
    $4,
    $5,
    $6,
    $7
from articles a
where a.id = $1
  and a.owner_user_id = $2
returning
    id,
    article_id,
    storage_key,
    original_filename,
    content_type,
    byte_size,
    created_at
`, articleID, ownerUserID, assetID, input.StorageKey, input.OriginalFilename, input.ContentType, input.ByteSize).Scan(
		&asset.ID,
		&asset.ArticleID,
		&asset.StorageKey,
		&asset.OriginalFilename,
		&asset.ContentType,
		&asset.ByteSize,
		&asset.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to create article asset")
	}

	return &asset, nil
}

func (r *PostgresRepository) GetAssetByID(ctx context.Context, assetID string) (*AssetRecord, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	var asset AssetRecord
	err := r.pool.QueryRow(ctx, `
select
    id,
    article_id,
    storage_key,
    original_filename,
    content_type,
    byte_size,
    created_at
from article_assets
where id = $1
`, assetID).Scan(
		&asset.ID,
		&asset.ArticleID,
		&asset.StorageKey,
		&asset.OriginalFilename,
		&asset.ContentType,
		&asset.ByteSize,
		&asset.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to fetch article asset")
	}

	return &asset, nil
}
