package supervisor

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/entities"

	_ "github.com/mattn/go-sqlite3"
)

// Manage explicitly defined dependencies.
type Supervisor struct {
	SQLite *data.SQLite
}

// Initialize supervisor dependencies explicitly.
func New(ctx context.Context) (*Supervisor, error) {
	sqlite, err := data.NewSQLite(ctx)
	if err != nil {
		return nil, fmt.Errorf("[pkg supervisor : func New] sqlite init failed :: %w", err)
	}

	return &Supervisor{
		SQLite: sqlite,
	}, nil
}

func (su *Supervisor) Init() {
	su.initSQLite()
}

func (su *Supervisor) initSQLite() {
	initialEntities := []data.SQLEntity{
		entities.User{},
		entities.Mock{},
		entities.MockQuestion{},
		entities.MockOption{},
		entities.Session{},
	}

	for _, entity := range initialEntities {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 80 * time.Second)
			defer cancel()

			typ := reflect.TypeOf(entity)
			slog.Info("Entity", "type", typ)

			stmt, err := data.CreateTable(ctx, su.SQLite.DB, entity)
			slog.Info("Creating table...", "stmt", stmt)

			if err != nil {
				slog.Error("Failed to create table", "error", err)
			}
		}()
	}
}
