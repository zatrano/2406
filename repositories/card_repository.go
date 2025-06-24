package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICardRepository interface {
	CreateCardWithRelations(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, card *models.Card, updatedBy uint) error
	GetCardByID(id uint) (*models.Card, error)
	DeleteCard(ctx context.Context, id uint) error
	GetAllCards() ([]models.Card, error)
}

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository() ICardRepository {
	return &CardRepository{db: databaseconfig.GetDB()}
}

func (r *CardRepository) CreateCardWithRelations(ctx context.Context, card *models.Card) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(card).Error
}

func (r *CardRepository) UpdateCard(ctx context.Context, id uint, card *models.Card, updatedBy uint) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("card_id = ?", id).Delete(&models.CardBank{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("card_id = ?", id).Delete(&models.CardSocialMedia{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	card.ID = id
	card.UpdatedBy = updatedBy
	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(card).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.db.Preload("CardBanks").Preload("CardSocialMedia").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) DeleteCard(ctx context.Context, id uint) error {
	var card models.Card
	card.ID = id
	return r.db.WithContext(ctx).Select(clause.Associations).Delete(&card).Error
}

func (r *CardRepository) GetAllCards() ([]models.Card, error) {
	var cards []models.Card
	err := r.db.Preload("CardBanks").Preload("CardSocialMedia").Find(&cards).Error
	return cards, err
}
