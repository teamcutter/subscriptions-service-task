package repo_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/teamcutter/subscriptions-service-task/internal/model"
	"github.com/teamcutter/subscriptions-service-task/internal/repo"
)

var (
	db       *sql.DB
	testRepo *repo.SubscriptionRepo
	mockUUID uuid.UUID = uuid.New()
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE subscriptions (
			id SERIAL PRIMARY KEY,
			service_name TEXT NOT NULL,
			price INT NOT NULL,
			user_id UUID NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE
		);
	`)
	if err != nil {
		panic(err)
	}

	testRepo = repo.NewSubscriptionRepo(db)

	code := m.Run()
	_ = container.Terminate(ctx)
	db.Close()
	
	if code != 0 {
		panic(fmt.Sprintf("tests failed with code %d", code))
	}
}

func TestCreateAndGetAll(t *testing.T) {
	sub := &model.Subscription{
		ServiceName: "Netflix",
		Price:       1000,
		UserID:      mockUUID,
		StartDate:   "01-2024",
		EndDate:     "12-2024",
	}

	err := testRepo.Create(sub)
	require.NoError(t, err)

	subs, err := testRepo.GetAll()
	require.NoError(t, err)
	assert.Len(t, subs, 1)
	assert.Equal(t, "Netflix", subs[0].ServiceName)
	assert.Equal(t, 1000, subs[0].Price)
}

func TestDelete(t *testing.T) {
	sub := &model.Subscription{
		ServiceName: "Spotify",
		Price:       500,
		UserID:      mockUUID,
		StartDate:   "02-2024",
		EndDate:     "03-2024",
	}
	err := testRepo.Create(sub)
	require.NoError(t, err)

	var id int
	err = db.QueryRow(`SELECT id FROM subscriptions WHERE service_name = 'Spotify'`).Scan(&id)
	require.NoError(t, err)

	err = testRepo.Delete(id)
	require.NoError(t, err)

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM subscriptions WHERE id = $1`, id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestTotalCost(t *testing.T) {
	sub := &model.Subscription{
		ServiceName: "HBO",
		Price:       100,
		UserID:      mockUUID,
		StartDate:   "01-2024",
		EndDate:     "03-2024",
	}
	err := testRepo.Create(sub)
	require.NoError(t, err)

	total, err := testRepo.TotalCost(mockUUID.String(), "HBO", "01-2024", "03-2024")
	require.NoError(t, err)

	assert.Equal(t, 300, total)
}