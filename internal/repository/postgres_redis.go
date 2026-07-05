package repository

import (
    "context"
    "database/sql"
    "time"
    "short-track/internal/domain"
    "short-track/pkg/base62"
    "github.com/go-redis/redis/v8"
)

type repo struct {
    db    *sql.DB
    cache *redis.Client
}

func NewRepository(db *sql.DB, cache *redis.Client) domain.URLRepository {
    return &repo{db: db, cache: cache}
}

func (r *repo) Save(url *domain.URL) error {
    query := `INSERT INTO urls (original_url) VALUES ($1) RETURNING id, created_at`
    return r.db.QueryRow(query, url.OriginalURL).Scan(&url.ID, &url.CreatedAt)
}

func (r *repo) GetByCode(code string) (*domain.URL, error) {
    ctx := context.Background()
    
    // 1. Try Redis Cache
    val, err := r.cache.Get(ctx, code).Result()
    if err == nil {
        return &domain.URL{OriginalURL: val, ShortCode: code}, nil
    }

    // 2. Fallback to Postgres
    id, err := base62.Decode(code)
    if err != nil {
        return nil, err
    }

    var url domain.URL
    query := `SELECT original_url FROM urls WHERE id = $1`
    err = r.db.QueryRow(query, id).Scan(&url.OriginalURL)
    
    if err == nil {
        // 3. Prime the cache for next time (TTL of 24 hours)
        r.cache.Set(ctx, code, url.OriginalURL, 24*time.Hour)
    }
    
    return &url, err
}