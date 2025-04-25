package gateway

import (
	"database/sql"
	"log"
	"time"

	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"
	"example.com/EVENT-MANAGEMENT-SYSTEM/internal/repository"
)

// tokenRepositoryImpl is the implementation of TokenRepository
type tokenRepositoryImpl struct {
	db *sql.DB
}

// NewTokenRepository creates a new instance of TokenRepository.
func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &tokenRepositoryImpl{db: db}
}

// FindByToken implements repository.TokenRepository.
func (t *tokenRepositoryImpl) FindByToken(token string) (*entity.Token, error) {
	r := &entity.Token{}
	query := `SELECT id, user_id, token, expires_at, created_at, updated_at FROM tokens WHERE token = $1`
	row := t.db.QueryRow(query, token)

	err := row.Scan(&r.ID, &r.UserID, &r.Token, &r.ExpiresAt, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("token not found")
			//return nil, errors.New("token not found")
			return nil, nil
		}
		return nil, err
	}
	return r, err
}

// Create implements repository.TokenRepository.
func (t *tokenRepositoryImpl) Create(token *entity.Token) error {
	// ❌ Fix: Column name was "tokens" — should be "token" (as per SELECT query)
	query := `INSERT INTO tokens (id, user_id, token, expires_at, created_at, updated_at)
	          VALUES($1, $2, $3, $4, $5, $6)`

	// ✅ Use time.Now().UTC() for consistency
	now := time.Now().UTC()

	result, err := t.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, now, now)
	if err != nil {
		log.Printf("Error inserting token: %v", err)
		return err // ❌ Missing return on insert error
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		return err
	}

	log.Printf("Rows affected: %d", rowsAffected)
	return nil
}
