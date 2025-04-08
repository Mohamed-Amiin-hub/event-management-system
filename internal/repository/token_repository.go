package repository

import "example.com/EVENT-MANAGEMENT-SYSTEM/internal/entity"

type TokenRepository interface {
	FindByToken(token string) (*entity.Token, error)
	Create(token *entity.Token) error
}
