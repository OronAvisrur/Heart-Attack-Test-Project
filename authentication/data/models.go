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

	// Using the context we wait for the response from the DB to our query
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

// This function returns a user by the name
func (user *User) GetUserByName() (*User, error) {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	query := `select user_name, email, password, created_at, updated_at
	from users where id = $1`

	// Using the context we wait for the response from the DB to our query
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

// This function insert new entry to the DB
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

	// Execute the query by adding the given user to the DB
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

// This function delete user from the DB
func (user *User) Delete() error {
	ctx, cancle := context.WithTimeout(context.Background(), db_time_out)
	defer cancle()

	query := `delete from users where user_name = $1`

	// Execute the query by deleting the given user from the DB and ignore the sql result
	_, possible_error := db.ExecContext(ctx, query, user.UserName)
	if possible_error != nil {
		return possible_error
	}

	return nil
}

// This function reset user password
func (user *User) RestPassword(new_password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db_time_out)
	defer cancel()

	// Generating new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_password), 12)
	if err != nil {
		return err
	}

	query := `update users set password = $1 where user_name = $2`
	// Execute the query

	_, err = db.ExecContext(ctx, query, hashedPassword, user.UserName)
	if err != nil {
		return err
	}

	return nil
}

// This function check if the given password is the user password
func (user *User) IsPasswordMatches(plainText string) (bool, error) {

	// Comparing user.Password and plainText
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
