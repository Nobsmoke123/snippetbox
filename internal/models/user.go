package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Created  time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

type UserModelInterface interface{
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

func (m *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, password, created) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

	_, err = m.DB.Exec(context.Background(), stmt, name, email, string(hash))

	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return ErrDuplicateEmail
			}
		}
	}

	return nil
}

func (m *UserModel) Get(id int) (User, error) {
	var user User

	stmt:= 	`SELECT id, name, email, created FROM users WHERE id =$1`

	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return User{}, ErrNoRecord
		}else{
			return User{}, err
		}
	}

	return user, nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, password FROM users WHERE email = $1`

	err := u.DB.QueryRow(context.Background(), stmt, email).Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	var exists bool
	
	stmt := `SELECT EXISTS(SELECT true FROM users WHERE id=$1)`

	err := u.DB.QueryRow(context.Background(), stmt, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}


func (u *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var hashedPassword []byte

	stmt := `SELECT password FROM users WHERE id=$1`

	err := u.DB.QueryRow(context.Background(), stmt, id).Scan(&hashedPassword)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword))
	if err != nil {
		return ErrInvalidCredentials
	}

	hash, err:= bcrypt.GenerateFromPassword([]byte(newPassword), 12)
    if err != nil {
		return err
	}

	stmt = `UPDATE users SET password=$2 WHERE id=$1 `

	_, err = u.DB.Exec(context.Background(), stmt, id, hash)
	if err != nil {
		return err
	}

	return nil
}