package pg_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jonss/testcontainers-go-wrapper/pg"
)

func TestContainer_success(t *testing.T) {
	cfg := pg.PostgresCfg{
		ImageName: pg.DefaultImage,
		Password:  "a_secret_password",
		UserName:  "a_user",
		DbName:    "db_name",
	}
	pgInfo, err := pg.Container(context.Background(), cfg)
	defer pgInfo.TearDown()

	if err != nil {
		t.Fatalf("pg.Container(); unexpected error %v", err)
	}

	sqlDB, err := sql.Open("postgres", pgInfo.DbURL)
	if err != nil {
		t.Fatalf("sql.Open(); error %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("sql.DB.Ping(); error %v", err)
	}
}
