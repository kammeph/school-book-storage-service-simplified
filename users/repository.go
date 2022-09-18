package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kammeph/school-book-storage-service-simplified/common"
)

const (
	getAllQuery               = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE ORDER BY u.id"
	getByIdQuery              = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE AND u.id = ? ORDER BY u.id"
	getByNameQuery            = "SELECT u.id, u.school_id, u.username, u.locale, r.role FROM users u INNER JOIN roles r ON r.user_id = u.id WHERE u.active = TRUE AND u.username= ? ORDER BY u.id"
	getCredentialsByNameQuery = "SELECT password_hash FROM users WHERE active = TRUE AND username= ?"
	countByNameQuery          = "SELECT COUNT(id) FROM users WHERE active = TRUE AND username= ?"
	createQuery               = "INSERT INTO users (id, school_id, username, password_hash, active, locale) VALUES (?, ?, ?, ?, ?, ?)"
	updateQuery               = "UPDATE users SET school_id = ?, username = ?, locale = ? WHERE id = ?"
	addRolesQuery             = "INSERT INTO roles (user_id, role) VALUES (?, ?)"
	deleteRolesQuery          = "DELETE FROM roles WHERE user_id = ?"
)

type UsersRepository interface {
	GetAll(ctx context.Context) ([]UserDto, error)
	GetById(ctx context.Context, userId string) (UserDto, error)
	GetByName(ctx context.Context, username string) (UserDto, error)
	GetCredentialsByName(ctx context.Context, username string) (string, error)
	Create(ctx context.Context, user UserModel) error
}

type SqlUsersRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) *SqlUsersRepository {
	return &SqlUsersRepository{db}
}

func (r *SqlUsersRepository) GetAll(ctx context.Context) ([]UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getAllQuery)
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

func (r *SqlUsersRepository) GetById(ctx context.Context, userId string) (UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getByIdQuery)
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

func (r *SqlUsersRepository) GetByName(ctx context.Context, username string) (UserDto, error) {
	stmt, err := r.db.PrepareContext(ctx, getByNameQuery)
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

func (r *SqlUsersRepository) GetCredentialsByName(ctx context.Context, username string) (string, error) {
	stmt, err := r.db.PrepareContext(ctx, getCredentialsByNameQuery)
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

func (r *SqlUsersRepository) countByName(ctx context.Context, username string) (int, error) {
	stmt, err := r.db.PrepareContext(ctx, countByNameQuery)
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

func (r *SqlUsersRepository) Create(ctx context.Context, user UserModel) error {
	count, err := r.countByName(ctx, user.Username)
	if err != nil {
		return err
	}
	if count >= 1 {
		return fmt.Errorf("a user with the name %s already exists", user.Username)
	}
	stmt, err := r.db.PrepareContext(ctx, createQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, user.SchoolId, user.Username, user.PasswordHash, user.Active, user.Locale)
	if err != nil {
		return err
	}

	return r.addRoles(ctx, user.ID, user.Roles)
}

func (r *SqlUsersRepository) Update(ctx context.Context, user UserDto) error {
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.SchoolId, user.Username, user.Locale, user.ID)
	if err != nil {
		return err
	}

	return r.updateRoles(ctx, user.ID, user.Roles)
}

func (r *SqlUsersRepository) addRoles(ctx context.Context, userId string, roles []common.Role) error {
	stmt, err := r.db.PrepareContext(ctx, addRolesQuery)
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

func (r *SqlUsersRepository) deleteRoles(ctx context.Context, userId string) error {
	stmt, err := r.db.PrepareContext(ctx, deleteRolesQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, userId)
	return err
}

func (r *SqlUsersRepository) updateRoles(ctx context.Context, userId string, roles []common.Role) error {
	err := r.deleteRoles(ctx, userId)
	if err != nil {
		return err
	}
	return r.addRoles(ctx, userId, roles)
}
