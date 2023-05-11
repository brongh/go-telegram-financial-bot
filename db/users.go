package db

import (
	"database/sql"
)

func GetUserByTgId(db sql.DB, tgId int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, tg_id FROM users WHERE tg_id = $1", tgId).Scan(&user.Id, &user.Username, &user.TgId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func insertUser(db sql.DB, user *User) (*User, error) {
	query := `
		INSERT INTO users (username, tg_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err := db.QueryRow(query, user.Username, user.TgId).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetOrCreateUser(db sql.DB, username string, tgId int) (*User, error) {
	user, err := GetUserByTgId(db, tgId)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	newUser := &User{
		Username: username,
		TgId: tgId,
	}

	return insertUser(db, newUser)
}