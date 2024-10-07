package models

import (
	"try-oauth/db"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" validate:"required,min=5,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=50,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
}

func (u *User) CheckUserByUsername(username string) *User {
	db.DB.QueryRow("SELECT username, password FROM users WHERE username = $1", username).Scan(&u.Username, &u.Password)
	return u
}

func (u *User) GetUserById(id int) *User {
	db.DB.QueryRow("SELECT id, username FROM users WHERE id = $1", id).Scan(&u.Id, &u.Username)
	return u
}

func (u *User) GetUserPagination(page int, size int) ([]User, int, error) {
	var total int
	_ = db.DB.QueryRow("SELECT COUNT(id) FROM users").Scan(&total)
	rows, err := db.DB.Query("SELECT id, username FROM users LIMIT $1 OFFSET $2", size, page)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username)
		users = append(users, user)
	}

	return users, total, nil
}

func (u *User) CreateUser() error {
	_, err := db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password)
	return err
}

func (u *User) ChangePassword(newPass string) error {
	_, err := db.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", newPass, u.Id)
	return err
}

func (u *User) DeleteUser() error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", u.Id)
	return err
}
