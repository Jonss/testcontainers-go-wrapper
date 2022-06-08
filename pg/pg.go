// pg/wrapper is a PostgreSQL testcontainers-go wrapper
// the idea of this package is isolate a common call to testcontainers
package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

const (
	DefaultImage        = "postgres:12.3-alpine"
	dbPort              = "5432/tcp"
	postgres            = "postgres"
	postgresUserKey     = "POSTGRES_USER"
	postgresPasswordKey = "POSTGRES_PASSWORD"
	postgresDbKey       = "POSTGRES_DB"
)

// PostgresCfg contains the config data to set a postgres container
// such as docker image name, password, username and DB name
type PostgresCfg struct {
	ImageName string
	Password  string
	UserName  string
	DbName    string
}

// ContainerInfo contains the testcontainers-go container functions
// the URL generated and the tearDown function to be called on tests

type CointainerInfo struct {
	Container testcontainers.Container
	DbURL     string
	TearDown  func()
}

// Container provides a docker image using github.com/testcontainers/testcontainers-go
func Container(ctx context.Context, cfg PostgresCfg) (*CointainerInfo, error) {
	var env = map[string]string{
		postgresUserKey:     cfg.UserName,
		postgresPasswordKey: cfg.Password,
		postgresDbKey:       cfg.DbName,
	}

	var url string
	dbURL := func(port nat.Port) string {
		url = fmt.Sprintf("%s://%s:%s@localhost:%s/%s?sslmode=disable", postgres, cfg.UserName, cfg.Password, port.Port(), cfg.DbName)
		return url
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        cfg.ImageName,
			ExposedPorts: []string{dbPort},
			Cmd:          []string{postgres, "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(nat.Port(dbPort), postgres, dbURL).Timeout(time.Second * 10),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("pg.Container(); error: %v", err)
	}

	tearDown := func() {
		container.Terminate(ctx)
	}

	return &CointainerInfo{
		Container: container,
		DbURL:     url,
		TearDown:  tearDown,
	}, nil
}
