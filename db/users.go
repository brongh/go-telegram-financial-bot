package db

import (
	"database/sql"
	"fmt"
	"log"
)

func GetUserByTgId(tgId int) (*User, error) {
	log.Print("user telegram id: ", tgId)
	var user User
	err := GetDB().QueryRow("SELECT id, username, tg_id FROM users WHERE tg_id = $1", tgId).Scan(&user.Id, &user.Username, &user.TgId)
	if err != nil {
		log.Print("[-] user not found: ", err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	log.Print("[+] user found: ", user.Username)
	return &user, nil
}

func insertUser(user *User) (*User, error) {
	log.Print("creating userId: ", user.TgId)
	query := `
		INSERT INTO users (username, tg_id)
		VALUES ($1, $2)
		RETURNING id;
	`

	err := GetDB().QueryRow(query, user.Username, user.TgId).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetOrCreateUser(username string, tgId int) (*User, error) {
	user, err := GetUserByTgId(tgId)
	if err != nil {
		if err.Error() != "user not found" {
			return nil, err
		}

		newUser := &User{
			Username: username,
			TgId: tgId,
		}
		user, err = insertUser(newUser)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}