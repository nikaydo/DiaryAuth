package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nikaydo/DiaryAuth/internal/config"
)

type Database struct {
	Pool *pgxpool.Pool
	Env  config.Env
}

func InitBD(e config.Env) (Database, error) {
	pool, err := pgxpool.New(context.Background(), e.Postgresql)
	if err != nil {
		return Database{}, err
	}
	return Database{Pool: pool}, nil
}

func (db *Database) Create(login, password string) (uuid.UUID, error) {
	var id uuid.UUID
	if err := db.Pool.QueryRow(context.Background(), "INSERT INTO users (login,password) VALUES ($1,$2) RETURNING uuid;", login, password).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (db *Database) CheckExist(login, password string) (uuid.UUID, error) {
	var id uuid.UUID
	if err := db.Pool.QueryRow(context.Background(), "SELECT uuid FROM users WHERE login = $1 AND password = $2", login, password).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (db *Database) Delete(id uuid.UUID) error {
	if _, err := db.Pool.Exec(context.Background(), "DELETE FROM users WHERE uuid = $1", id); err != nil {
		return err
	}
	return nil
}

func (db *Database) RefreshUpdate(id uuid.UUID, refresh string) error {
	if _, err := db.Pool.Exec(context.Background(), "UPDATE users SET refresh_token = $1 WHERE uuid = $2;", refresh, id); err != nil {
		return err
	}
	return nil
}

func (db *Database) GetRefresh(id uuid.UUID) (string, error) {
	var refresh string
	if err := db.Pool.QueryRow(context.Background(), "SELECT refresh_token FROM users WHERE uuid = $1;", id).Scan(&refresh); err != nil {
		return refresh, err
	}
	return refresh, nil
}
