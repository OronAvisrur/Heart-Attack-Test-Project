package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const db_time_out = time.Second * 3

var db *sql.DB

// Models is the type for this package
type Models struct {
	User User
}

type User struct {
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New is the function used to create an instance of the data package
func New(db_pool *sql.DB) Models {
	db = db_pool

	return Models{
		User: User{},
	}
}

// This function returns a slice of all users, sorted by user name
func (user *User) GetAllTableRows() ([]*User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	query := `select user_name, email, password, created_at, updated_at 
	from users order by user_name`

	table_rows, possible_error := db.QueryContext(ctx, query)
	if possible_error != nil {
		return nil, possible_error
	}
	defer table_rows.Close()

	var table_users_list []*User

	for table_rows.Next() {
		var user_info User
		possible_error := table_rows.Scan(
			&user_info.UserName,
			&user_info.Email,
			&user_info.Password,
			&user_info.CreatedAt,
			&user_info.UpdatedAt,
		)

		if possible_error != nil {
			log.Panic("Error scanning", possible_error)
			return nil, possible_error
		}

		table_users_list = append(table_users_list, &user_info)
	}

	return table_users_list, nil

}

func (user *User) GetUserByName() (*User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	query := `select user_name, email, password, created_at, updated_at
	from users where id = $1`

	user_row := db.QueryRowContext(ctx, query, user.UserName)
	var user_info User

	possible_error := user_row.Scan(
		&user_info.UserName,
		&user_info.Email,
		&user_info.Password,
		&user_info.CreatedAt,
		&user_info.UpdatedAt,
	)

	if possible_error != nil {
		return nil, possible_error
	}

	return &user_info, nil
}

func (user *User) Insert() (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	hashed_password, possible_error := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if possible_error != nil {
		return "", possible_error
	}

	var new_user_name string
	query := `insert into users (email, password, created_at, updated_at)
	values ($1, $2, $3, $4) returning user_name`

	possible_error = db.QueryRowContext(ctx, query,
		user.Email,
		hashed_password,
		time.Now(),
		time.Now(),
	).Scan(&new_user_name)

	if possible_error != nil {
		return "", possible_error
	}

	return new_user_name, nil
}

func (user *User) Delete() error {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	query := `delete from users where user_name = $1`

	_, possible_error := db.ExecContext(ctx, query, user.UserName)
	if possible_error != nil {
		return possible_error
	}

	return nil
}

func (user *User) RestPassword(new_password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db_time_out)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where user_name = $2`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, user.UserName)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) IsPasswordMatches(plainText string) (bool, error) {
	possible_error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if possible_error != nil {
		switch {
		case errors.Is(possible_error, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, possible_error
		}
	}

	return true, nil
}
