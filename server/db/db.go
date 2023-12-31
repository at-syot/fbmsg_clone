package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB *sql.DB

	// envs
	dbHost string
	dbPort string
	dbUser string
	dbPw   string
	dbName string
)

func Init() error {
	err := loadEnvs()
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPw,
		dbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("database connection err - %s\n", err.Error())
		return err
	}

	DB = db
	return nil
}

func loadEnvs() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPw = os.Getenv("DB_PW")
	dbName = os.Getenv("DB_NAME")

	if len(dbHost) == 0 || len(dbPort) == 0 || len(dbUser) == 0 || len(dbPw) == 0 || len(dbName) == 0 {
		return errors.New("can't load db related envs")
	}

	return nil
}

type Conn struct {
	Ctx context.Context
	Tx  *sql.Tx
}

type ExecutorFn func(conn Conn) error

func ExecWithTx(ctx context.Context, executeFn ExecutorFn) error {
	conn, err := DB.Conn(ctx)
	if err != nil {
		log.Printf("DB: getting connection err - %s\n", err.Error())
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		log.Printf("DB: starting transaction err - %s\n", err.Error())
		return err
	}

	if err := executeFn(Conn{Ctx: ctx, Tx: tx}); err != nil {
		log.Printf("DB: executeFn called with err - %s\n", err.Error())
		log.Println("DB: rolling back")
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("DB: commiting err - %s\n", err.Error())
		return err
	}

	log.Println("DB: Tx committed successfully")
	return nil
}

func (c Conn) Execute(query string, args ...any) error {
	res, err := c.Tx.ExecContext(c.Ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("DB: get rowsAffected err - %s\n", err.Error())
		return err
	}
	if rows != 1 {
		errMsg := fmt.Sprintf("DB: expected to affect 1 row, got %d\n", rows)
		log.Printf(errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (c Conn) QueryRow(query string, args []any, dest ...any) error {
	err := c.Tx.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func QueryRowContext(ctx context.Context, query string, args []any, dest ...any) error {
	conn, err := DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = conn.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		return err
	}
	return nil
}

type RowReaderFn func(rows *sql.Rows) error

func QueryContext(ctx context.Context, query string, args []any, reader RowReaderFn) error {
	conn, err := DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("DB: query context err - %s\n", err)
		return err
	}

	for rows.Next() {
		if err := reader(rows); err != nil {
			break
		}
	}

	if err = rows.Close(); err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
