package repositories

import (
	"api/internal/models"
	"database/sql"
	"fmt"
)

// Users represents an users repositorie
type Users struct {
	db *sql.DB
}

// NewUsersRepositorie a new repositorie of users
func NewUsersRepository(db *sql.DB) *Users {
	return &Users{db}
}

// Create insert a user in database
func (repository Users) Create(user models.User) (ID int, err error) {
	stmt, err := repository.db.Prepare(
		"INSERT INTO users (name, userName, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
	)

	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Name, user.UserName, user.Email, user.Password).Scan(&ID)
	if err != nil {
		return
	}

	return
}

// FilterByUserName fetch all the users that match the userName
func (repository Users) FilterByUserName(userName string) (users []models.User, err error) {
	userName = fmt.Sprintf("%%%s%%", userName)

	rows, err := repository.db.Query(
		"SELECT id, name, userName, email FROM users WHERE userName LIKE $1;", userName,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.UserName,
			&user.Email,
		); err != nil {
			return
		}
		users = append(users, user)
	}

	return
}

// GetById fetch an user by id
func (repository Users) GetById(id int) (user models.User, err error) {
	row, err := repository.db.Query(
		"SELECT id, name, userName, email FROM users WHERE id=$1", id,
	)
	if err != nil {
		return
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(
			&user.ID,
			&user.Name,
			&user.UserName,
			&user.Email,
		)
		if err != nil {
			return
		}
	}

	return
}

// Update updates an user in database
func (repository Users) Update(user models.User) (err error) {
	stmt, err := repository.db.Prepare(
		"UPDATE users SET name=$1, userName=$2, email=$3 WHERE id=$4",
	)
	if err != nil {
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Name, user.UserName, user.Email, user.ID); err != nil {
		return
	}

	return
}

// Delete deletes an user in database
func (repository Users) Delete(id int) (err error) {
	stmt, err := repository.db.Prepare("DELETE FROM users WHERE id=$1")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return
	}

	return
}