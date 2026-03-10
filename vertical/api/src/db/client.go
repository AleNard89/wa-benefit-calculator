package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBClient struct {
	C       *pgxpool.Pool
	Builder sq.StatementBuilderType
}

var (
	database *DBClient
	once     sync.Once
)

func DB() *DBClient {
	if database == nil {
		DBHost := os.Getenv("POSTGRES_HOST")
		DBPort := os.Getenv("POSTGRES_PORT")
		DBName := os.Getenv("POSTGRES_DB")
		DBUser := os.Getenv("POSTGRES_USER")
		DBPassword := os.Getenv("POSTGRES_PASSWORD")

		once.Do(func() {
			connectionURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName)
			connPool, err := pgxpool.New(context.Background(), connectionURI)
			if err != nil {
				zap.S().Fatal("Cannot connect to postgres db on ", DBHost)
			} else {
				zap.S().Info("Successfully connected to postgres on ", DBHost)
			}

			database = &DBClient{connPool, sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
		})
	}

	return database
}

func Close() {
	if database != nil {
		database.C.Close()
	}
}
