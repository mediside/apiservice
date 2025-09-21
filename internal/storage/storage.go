package storage

import (
	"apiservice/internal/config"
	"apiservice/internal/storage/research"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type Storage struct {
	db              *sql.DB
	cache           *redis.Client
	ResearchStorage *research.ResearchStorage
}

func New(logger *slog.Logger, cfg *config.Config) Storage {
	db := connectPostgres(cfg)
	logger.Info("db connected succesfully")

	cache := connectCache(cfg)
	logger.Info("cache connected succesfully")

	resStorage := research.New(db)

	return Storage{
		db:              db,
		cache:           cache,
		ResearchStorage: resStorage,
	}
}

func (s Storage) Close() {
	s.db.Close()
	s.cache.Close()
}

func connectPostgres(cfg *config.Config) *sql.DB {
	dbOptions := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name,
	)

	db, err := sql.Open("postgres", dbOptions)
	if err != nil {
		log.Fatalf("can't open postgres: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("can't ping postgres: %s", err.Error())
	}

	return db
}

func connectCache(cfg *config.Config) *redis.Client {
	rdbOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	}

	rdb := redis.NewClient(rdbOptions)

	if _, err := rdb.Ping(context.TODO()).Result(); err != nil {
		log.Fatalf("can't ping redis: %s", err.Error())
	}

	return rdb
}
