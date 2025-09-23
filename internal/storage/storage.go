package storage

import (
	"apiservice/internal/config"
	"apiservice/internal/storage/collection"
	"apiservice/internal/storage/research"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type Storage struct {
	db                *sql.DB
	cache             *redis.Client
	CollectionStorage *collection.CollectionStorage
	ResearchStorage   *research.ResearchStorage
}

func New(logger *slog.Logger, cfg *config.Config) Storage {
	db := connectPostgres(cfg)
	logger.Info("db connected succesfully")

	cache := connectCache(cfg)
	logger.Info("cache connected succesfully")

	createFileFolder(cfg.ResearchSavePath)
	logger.Info("research save folder checked", slog.String("path", cfg.ResearchSavePath))

	colStorage := collection.New(db)
	resStorage := research.New(cfg.ResearchSavePath)

	return Storage{
		db:                db,
		cache:             cache,
		CollectionStorage: colStorage,
		ResearchStorage:   resStorage,
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

func createFileFolder(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("can't create folder for research store: %s", err.Error())
	}
}
