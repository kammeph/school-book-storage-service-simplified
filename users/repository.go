package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service-simplified/common"
)

const (
	getUsersQuery                 = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE ORDER BY u.id"
	getUserByIdQuery              = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE AND u.id = ? ORDER BY u.id"
	getUserByNameQuery            = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE AND u.username= ? ORDER BY u.id"
	getUserCredentialsByNameQuery = "SELECT password_hash FROM users WHERE active = TRUE AND username= ?"
	countUserNameQuery            = "SELECT COUNT(id) FROM users WHERE active = TRUE AND username= ?"
	addUserQuery                  = "INSERT INTO users (id, school_id, username, password_hash, active, locale) VALUES (?, ?, ?, ?, ?, ?)"
	updateUserQuery               = "UPDATE users SET school_id = ?, username = ?, locale = ? WHERE id = ?"
	addUserRolesQuery             = "INSERT INTO roles (user_id, role) VALUES (?, ?)"
	deleteUserRolesQuery          = "DELETE FROM roles WHERE user_id = ?"
)

type UsersRepository interface {
	GetUsers(ctx context.Context) ([]UserDto, error)
	GetUserById(ctx context.Context, userId string) (UserDto, error)
	GetUserByName(ctx context.Context, username string) (UserDto, error)
	GetUserCredentialsByName(ctx context.Context, username string) (string, error)
	AddUser(ctx context.Context, user UserModel) error
}

type SqlUsersRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) *SqlUsersRepository {
	return &SqlUsersRepository{db}
}

func (r *SqlUsersRepository) GetUsers(ctx context.Context) ([]UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanUsers(rows)
}

func (r *SqlUsersRepository) GetUserById(ctx context.Context, userId string) (UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getUserByIdQuery)
	if err != nil {
		return UserDto{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return UserDto{}, err
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return UserDto{}, err
	}

	if len(users) == 0 {
		return UserDto{}, nil
	}

	if len(users) > 1 {
		return UserDto{}, errors.New("more than one user found")
	}

	return users[0], nil
}

func (r *SqlUsersRepository) GetUserByName(ctx context.Context, username string) (UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getUserByNameQuery)
	if err != nil {
		return UserDto{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return UserDto{}, err
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		return UserDto{}, err
	}

	if len(users) == 0 {
		return UserDto{}, nil
	}

	if len(users) > 1 {
		return UserDto{}, errors.New("more than one user found")
	}

	return users[0], nil
}

func (r *SqlUsersRepository) GetUserCredentialsByName(ctx context.Context, username string) (string, error) {
	stmt, err := r.db.PrepareContext(ctx, getUserCredentialsByNameQuery)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, username)
	if row.Err() != nil {
		return "", row.Err()
	}

	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}

	return password, nil
}

func (r *SqlUsersRepository) countUsername(ctx context.Context, username string) (int, error) {
	stmt, err := r.db.PrepareContext(ctx, countUserNameQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, username)
	if row.Err() != nil {
		return 0, err
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func scanUsers(rows *sql.Rows) ([]UserDto, error) {
	var users []UserDto
	previousID := ""
	for rows.Next() {
		var user UserDto
		var role common.Role
		err := rows.Scan(&user.ID, &user.SchoolId, &user.Username, &user.Locale, &role)
		if err != nil {
			return nil, err
		}
		if user.ID != previousID {
			user.Roles = append(user.Roles, role)
			users = append(users, user)
			previousID = user.ID
		} else {
			users[len(users)-1].Roles = append(users[len(users)-1].Roles, role)
		}
	}
	return users, nil
}

func (r *SqlUsersRepository) AddUser(ctx context.Context, user UserModel) error {
	count, err := r.countUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if count >= 1 {
		return fmt.Errorf("a user with the name %s already exists", user.Username)
	}
	stmt, err := r.db.PrepareContext(ctx, addUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.SchoolId, user.Username, user.PasswordHash, user.Active, user.Locale)
	if err != nil {
		return err
	}

	return r.addUserRoles(ctx, user.ID, user.Roles)
}

func (r *SqlUsersRepository) UpdateUser(ctx context.Context, user UserDto) error {
	stmt, err := r.db.PrepareContext(ctx, updateUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.SchoolId, user.Username, user.Locale, user.ID)
	if err != nil {
		return err
	}

	return r.updateUserRoles(ctx, user.ID, user.Roles)
}

func (r *SqlUsersRepository) addUserRoles(ctx context.Context, userId string, roles []common.Role) error {
	stmt, err := r.db.PrepareContext(ctx, addUserRolesQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, role := range roles {
		_, err = stmt.ExecContext(ctx, userId, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SqlUsersRepository) deleteUserRoles(ctx context.Context, userId string) error {
	stmt, err := r.db.PrepareContext(ctx, deleteUserRolesQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userId)
	return err
}

func (r *SqlUsersRepository) updateUserRoles(ctx context.Context, userId string, roles []common.Role) error {
	err := r.deleteUserRoles(ctx, userId)
	if err != nil {
		return err
	}
	return r.addUserRoles(ctx, userId, roles)
}
