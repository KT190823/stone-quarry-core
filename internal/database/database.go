package database

import (
	"bufio"
	"context"
	"fmt"
	"log"
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
		dsn = "postgres://postgres.dgnglulokxxifvfxckgi:ppXS30yQ3eP342vR@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres?sslmode=require"
	}

	// Tự động chuyển đổi direct IPv6 Supabase sang IPv4 Pooler cho Render/Cloud không hỗ trợ IPv6
	if strings.Contains(dsn, "db.dgnglulokxxifvfxckgi.supabase.co") {
		dsn = "postgres://postgres.dgnglulokxxifvfxckgi:ppXS30yQ3eP342vR@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres?sslmode=require"
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
