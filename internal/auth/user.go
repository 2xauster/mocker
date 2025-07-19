package auth

// import (
// 	"context"
// 	"database/sql"

// 	"github.com/ashtonx86/mocker/internal/data"
// 	"github.com/ashtonx86/mocker/internal/entities"
// 	"github.com/ashtonx86/mocker/internal/schemas"
// )

// func CreateUser(ctx context.Context, db *sql.DB, userData schemas.UserCreateRequest) (*entities.User, error) {
// 	data.Insert(ctx, "User", data.SQLInsertArgs{
// 		What: []string{"ID", "Name"},
// 	})
// 	return nil, nil 
// }