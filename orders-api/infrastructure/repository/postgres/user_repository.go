package postgres

/*
struct UserPostgresRepository
*/

import (
	"context"
	"errors"

	"orders-api/internal/domain"
	"orders-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

type UserPostgresRepository struct {
	db repository.DBTX
}

func NewUserRepository(db repository.DBTX) repository.UserRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

func (r *UserPostgresRepository) Save(user *domain.User) error {

	_, err := r.db.Exec(
		context.Background(),
		`INSERT INTO users
		(id,name,email,password_hash)
		VALUES($1,$2,$3,$4)`,

		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserPostgresRepository) Update(user *domain.User) error {

	_, err := r.db.Exec(
		context.Background(),
		`UPDATE users
		SET
			name=$2,
			email=$3,
			password_hash=$4,
			updated_at=NOW()
		WHERE id=$1`,

		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserPostgresRepository) Delete(id string) error {

	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM users WHERE id=$1`,
		id,
	)

	return err
}

func (r *UserPostgresRepository) FindByID(id string) (*domain.User, error) {

	user := &domain.User{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id,name,email,password_hash
		 FROM users
		 WHERE id=$1`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserPostgresRepository) FindAll() ([]*domain.User, error) {

	rows, err := r.db.Query(
		context.Background(),
		`SELECT id,name,email,password_hash
		 FROM users
		 ORDER BY name`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*domain.User

	for rows.Next() {

		user := &domain.User{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserPostgresRepository) FindByEmail(email string) (*domain.User, error) {

	user := &domain.User{}

	err := r.db.QueryRow(
		context.Background(),
		`
			SELECT
				id,
				name,
				email,
				password_hash
			FROM users
			WHERE email = $1
			`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
