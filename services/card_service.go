package services

import (
	"context"
	"errors"

	"davet.link/models"
	"davet.link/repositories"
)

type ICardService interface {
	CreateCard(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, card *models.Card, updatedBy uint) error
	GetCardByID(id uint) (*models.Card, error)
	DeleteCard(ctx context.Context, id uint) error
	GetAllCards() ([]models.Card, error)
}

type CardService struct {
	repo repositories.ICardRepository
}

func NewCardService() ICardService {
	return &CardService{repo: repositories.NewCardRepository()}
}

func (s *CardService) CreateCard(ctx context.Context, card *models.Card) error {
	if card == nil {
		return errors.New("geçersiz kart verisi")
	}
	return s.repo.CreateCardWithRelations(ctx, card)
}

func (s *CardService) UpdateCard(ctx context.Context, id uint, card *models.Card, updatedBy uint) error {
	return s.repo.UpdateCard(ctx, id, card, updatedBy)
}

func (s *CardService) GetCardByID(id uint) (*models.Card, error) {
	return s.repo.GetCardByID(id)
}

func (s *CardService) DeleteCard(ctx context.Context, id uint) error {
	return s.repo.DeleteCard(ctx, id)
}

func (s *CardService) GetAllCards() ([]models.Card, error) {
	return s.repo.GetAllCards()
}
