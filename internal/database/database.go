package database

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			os.Setenv(k, v)
		}
	}
}

// ensureDatabase kiểm tra kết nối tới database đích. Nếu database chưa tồn tại
// (PostgreSQL SQLSTATE 3D000) thì tự động kết nối vào database "postgres" để
// tạo mới. Trả về true nếu database vừa được tạo.
func ensureDatabase(ctx context.Context, config *pgxpool.Config) (bool, error) {
	dbName := config.ConnConfig.Database

	probe, err := pgxpool.NewWithConfig(ctx, config)
	if err == nil {
		err = probe.Ping(ctx)
		probe.Close()
	}
	if err == nil {
		return false, nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "3D000" {
		return false, err
	}

	maintConfig := config.Copy()
	maintConfig.ConnConfig.Database = "postgres"

	maint, mErr := pgxpool.NewWithConfig(ctx, maintConfig)
	if mErr != nil {
		return false, mErr
	}
	defer maint.Close()

	if mErr = maint.Ping(ctx); mErr != nil {
		return false, mErr
	}

	var exists bool
	if mErr = maint.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists); mErr != nil {
		return false, mErr
	}
	if exists {
		return false, err
	}

	if _, mErr = maint.Exec(ctx, "CREATE DATABASE "+pgx.Identifier{dbName}.Sanitize()); mErr != nil {
		return false, mErr
	}

	fmt.Printf("🆕 Đã tạo mới database %q.\n", dbName)
	return true, nil
}

// Connect thiết lập connection pool tới PostgreSQL local. Nếu database trong
// DATABASE_URL chưa tồn tại, hàm sẽ tự động tạo database đó trước khi kết nối.
// Trả về true nếu database vừa được tạo mới (khi đó main sẽ chạy seed dữ liệu).
func Connect() bool {
	loadEnvFile(".env")
	loadEnvFile("../backend/.env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/mo_da?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	created, err := ensureDatabase(context.Background(), config)
	if err != nil {
		fmt.Printf("\n⚠️ [PostgreSQL] Không thể kết nối/tạo database (%v).\n", err)
		fmt.Println("   Kiểm tra PostgreSQL local đã chạy chưa, ví dụ: make db-up")
		return false
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		fmt.Printf("\n⚠️ [PostgreSQL] Không thể kết nối tới database (%v).\n", err)
		fmt.Println("   Kiểm tra PostgreSQL local đã chạy chưa, ví dụ: make db-up")
		return false
	}

	if created {
		fmt.Printf("🎉 Đã tạo mới và kết nối PostgreSQL (%s/%s)!\n", config.ConnConfig.Host, config.ConnConfig.Database)
	} else {
		fmt.Printf("🎉 Connected to local PostgreSQL (%s/%s) successfully!\n", config.ConnConfig.Host, config.ConnConfig.Database)
	}
	return created
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
