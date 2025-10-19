package storage

import (
	"apiservice/internal/config"
	"apiservice/internal/storage/collection"
	"apiservice/internal/storage/inference"
	"apiservice/internal/storage/research"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"

	infGRPC "apiservice/internal/gen/go/inference/inference.v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

type Storage struct {
	db                *sql.DB
	CollectionStorage *collection.Storage
	ResearchStorage   *research.Storage
	InferenceStorage  *inference.Storage
}

func New(logger *slog.Logger, cfg *config.Config) Storage {
	db := connectPostgres(cfg)
	logger.Info("db connected succesfully and migrations applied")

	checkResearchesFolder(cfg.ResearchSavePath)
	logger.Info("research save folder checked", slog.String("path", cfg.ResearchSavePath))

	grpcConn := connectGRPC(cfg)
	logger.Info("create gRPC connection succesfully")

	colStorage := collection.New(cfg.ResearchSavePath, db)
	resStorage := research.New(cfg.ResearchSavePath, db)
	infStorage := inference.New(grpcConn)

	return Storage{
		db:                db,
		CollectionStorage: colStorage,
		ResearchStorage:   resStorage,
		InferenceStorage:  infStorage,
	}
}

func (s Storage) Close() {
	s.db.Close()
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

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("can't set dialect for migrations: %s", err.Error())
	}

	if err := goose.Up(db, cfg.Postgres.MigrationsDir); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	return db
}

func connectGRPC(cfg *config.Config) infGRPC.InferenceClient {
	addr := fmt.Sprintf("%s:%d", cfg.GrpcConfig.Host, cfg.GrpcConfig.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("can't create gRPC for inference: %s", err.Error()))
	}

	return infGRPC.NewInferenceClient(conn)
}

func checkResearchesFolder(path string) {
	if ex, err := exists(path); err != nil {
		log.Fatalf("can't check folder %s: %s", path, err.Error()) // volume в контейнере должен быть смонтирован
	} else if !ex {
		log.Fatalf("folder %s not exists", path) // volume в контейнере должен быть смонтирован
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
