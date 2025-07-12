package supervisor

import (
	"context"
	"fmt"
	"log/slog"
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
	ctx, cancel := context.WithTimeout(context.Background(), 80 * time.Second)
	defer cancel() 

	su.initSQLite(ctx)
}

func (su *Supervisor) initSQLite(ctx context.Context) {
	var user entities.User
	stmt, err := data.CreateTable(ctx, su.SQLite.DB, user)

	slog.Info("Creating initial tables ", "stmt", stmt)
	if err != nil {
		panic(fmt.Errorf("[Supervisor.initSQLite] create initial tables failed :: %w", err))
	}
}