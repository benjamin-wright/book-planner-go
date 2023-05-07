package cockroach

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/postgres"
)

type MigrationsClient struct {
	conn       *pgx.Conn
	deployment string
	database   string
}

func NewMigrations(deployment string, namespace string, database string) (*MigrationsClient, error) {
	cfg := postgres.ConnectConfig{
		Host:     fmt.Sprintf("%s.%s.svc.cluster.local", deployment, namespace),
		Port:     26257,
		Username: "root",
		Database: database,
	}

	conn, err := postgres.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cockroach db at %s: %+v", database, err)
	}

	return &MigrationsClient{
		conn:       conn,
		deployment: deployment,
		database:   database,
	}, nil
}

func (c *MigrationsClient) Stop() {
	zap.S().Infof("Closing connection to DB %s[%s]", c.deployment, c.database)
	c.conn.Close(context.Background())
}

func (c *MigrationsClient) HasMigrationsTable() (bool, error) {
	rows, err := c.conn.Query(context.TODO(), "SELECT DISTINCT(tablename) FROM pg_catalog.pg_tables WHERE tablename = $1", "migrations")
	if err != nil {
		zap.S().Errorf("failed to check migrations table: %+v", err)
		return false, fmt.Errorf("failed to check for migrations: %+v", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (d *MigrationsClient) CreateMigrationsTable() error {
	_, err := d.conn.Exec(
		context.TODO(),
		`
			BEGIN;

			SAVEPOINT migration_restart;

			CREATE TABLE migrations (
				id INT PRIMARY KEY NOT NULL UNIQUE
			);

			RELEASE SAVEPOINT migration_restart;

			COMMIT;
		`,
	)

	return err
}

func (c *MigrationsClient) LatestMigration() int64 {
	var found int64
	err := c.conn.QueryRow(context.Background(), "SELECT MAX(id) FROM migrations").Scan(&found)
	if err != nil {
		return 0
	}
	return found
}

func (c *MigrationsClient) AddMigration(id int64) error {
	_, err := c.conn.Exec(context.Background(), "INSERT INTO migrations (id) VALUES ($1)", id)
	return err
}

func (c *MigrationsClient) RunMigration(query string) error {
	_, err := c.conn.Exec(context.TODO(), query)

	return err
}

func (c *MigrationsClient) GetTables() ([]string, error) {
	rows, err := c.conn.Query(context.TODO(), "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := []string{}

	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, nil
}

func (c *MigrationsClient) GetTableSchema(tableName string) (map[string]string, error) {
	rows, err := c.conn.Query(context.TODO(), "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := map[string]string{}

	for rows.Next() {
		var column string
		var dataType string

		if err = rows.Scan(&column, &dataType); err != nil {
			return nil, err
		}

		columns[column] = dataType
	}

	return columns, err
}

func (c *MigrationsClient) GetContents(tableName string) ([][]interface{}, error) {
	rows, err := c.conn.Query(context.TODO(), fmt.Sprintf("SELECT * FROM %s", pgx.Identifier{tableName}.Sanitize()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := [][]interface{}{}

	for rows.Next() {
		row, err := rows.Values()
		if err != nil {
			return nil, err
		}

		contents = append(contents, row)
	}

	return contents, err
}
