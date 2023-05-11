package db

import (
	"database/sql"
)

func GetUserByTgId(tgId int) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT id, username, tg_id FROM users WHERE tg_id = $1", tgId).Scan(&user.Id, &user.Username, &user.TgId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func insertUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (username, tg_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err := DB.QueryRow(query, user.Username, user.TgId).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetOrCreateUser(username string, tgId int) (*User, error) {
	user, err := GetUserByTgId(tgId)
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

	return insertUser(newUser)
}