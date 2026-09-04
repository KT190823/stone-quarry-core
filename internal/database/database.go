package database

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
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

func Connect() {
	loadEnvFile(".env")
	loadEnvFile("../backend/.env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		pass := os.Getenv("SUPABASE_PASSWORD_DATABASE")
		if pass == "" {
			pass = "Mobeo@123#!@"
		}
		encodedPass := url.QueryEscape(pass)
		dsn = fmt.Sprintf("postgres://postgres:%s@db.nsdkctqxgrqvvgdtqyyr.supabase.co:5432/postgres?sslmode=require", encodedPass)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	// Disable statement caching to avoid prepared statement collisions with Supabase pooler
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		fmt.Printf("\n⚠️ [Supabase DB Auth] Không thể xác thực mật khẩu Postgres cho project (%v).\n", err)
		return
	}

	fmt.Printf("🎉 Connected to Supabase PostgreSQL (%s) successfully!\n", config.ConnConfig.Host)
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
