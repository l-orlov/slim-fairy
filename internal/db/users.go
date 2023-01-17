package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/internal/model"
)

func (db *Database) CreateUser(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (name, age) VALUES ($1, $2);"

	_, err := db.pool.Exec(ctx, query, user.Name, user.Age)
	if err != nil {
		return dbError(err)
	}

	return nil
}

func (db *Database) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := "SELECT * FROM users WHERE id=$1;"

	user := &model.User{}
	err := pgxscan.Get(ctx, db.pool, user, query, id)
	if err != nil {
		return nil, dbError(err)
	}

	return user, nil
}

func (db *Database) GetUsers(ctx context.Context) ([]*model.User, error) {
	query := "SELECT * FROM users;"

	var users []*model.User
	err := pgxscan.Select(ctx, db.pool, &users, query)
	if err != nil {
		return nil, dbError(err)
	}

	return users, nil
}
