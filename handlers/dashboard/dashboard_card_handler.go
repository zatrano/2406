package handlers

import (
	"net/http"
	"strconv"

	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardCardHandler struct {
	cardService services.ICardService
}

func NewDashboardCardHandler() *DashboardCardHandler {
	return &DashboardCardHandler{
		cardService: services.NewCardService(),
	}
}

func (h *DashboardCardHandler) ListCards(c *fiber.Ctx) error {
	cards, err := h.cardService.GetAllCards()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kartlar yüklenemedi")
		return renderer.Render(c, "dashboard/cards/list", "layouts/dashboard", fiber.Map{
			"Title": "Kartlar",
			"Cards": []models.Card{},
		}, http.StatusOK)
	}
	return renderer.Render(c, "dashboard/cards/list", "layouts/dashboard", fiber.Map{
		"Title": "Kartlar",
		"Cards": cards,
	}, http.StatusOK)
}

func (h *DashboardCardHandler) ShowCreateCard(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/cards/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Kart Oluştur",
	})
}

func (h *DashboardCardHandler) CreateCard(c *fiber.Ctx) error {
	if err := requests.ValidateCardRequest(c); err != nil {
		return err
	}
	req := c.Locals("cardRequest").(requests.CardRequest)

	card := &models.Card{
		Name:            req.Name,
		Slug:            req.Slug,
		Title:           req.Title,
		Photo:           req.Photo,
		Telephone:       req.Telephone,
		Email:           req.Email,
		Location:        req.Location,
		WebsiteUrl:      req.WebsiteUrl,
		StoreUrl:        req.StoreUrl,
		IsActive:        req.IsActive == "true",
		UserID:          req.UserID,
	}
	for _, cb := range req.CardBanks {
		card.CardBanks = append(card.CardBanks, models.CardBank{
			BankID: cb.BankID,
			IBAN:   cb.IBAN,
		})
	}
	for _, cs := range req.CardSocialMedia {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{
			SocialMediaID: cs.SocialMediaID,
			URL:           cs.URL,
		})
	}

	if err := h.cardService.CreateCard(c.UserContext(), card); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart oluşturulamadı: "+err.Error())
		return c.Redirect("/dashboard/cards/create", http.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) ShowUpdateCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	card, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart bulunamadı")
		return c.Redirect("/dashboard/cards", http.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/cards/update", "layouts/dashboard", fiber.Map{
		"Title": "Kart Güncelle",
		"Card":  card,
	})
}

func (h *DashboardCardHandler) UpdateCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := requests.ValidateCardRequestWithPath(c, "/dashboard/cards/update/"+strconv.Itoa(id)); err != nil {
		return err
	}
	req := c.Locals("cardRequest").(requests.CardRequest)

	card := &models.Card{
		Name:            req.Name,
		Slug:            req.Slug,
		Title:           req.Title,
		Photo:           req.Photo,
		Telephone:       req.Telephone,
		Email:           req.Email,
		Location:        req.Location,
		WebsiteUrl:      req.WebsiteUrl,
		StoreUrl:        req.StoreUrl,
		IsActive:        req.IsActive == "true",
		UserID:          req.UserID,
	}
	for _, cb := range req.CardBanks {
		card.CardBanks = append(card.CardBanks, models.CardBank{
			BankID: cb.BankID,
			IBAN:   cb.IBAN,
			BaseModel: models.BaseModel{ID: cb.ID},
		})
	}
	for _, cs := range req.CardSocialMedia {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{
			SocialMediaID: cs.SocialMediaID,
			URL:           cs.URL,
			BaseModel: models.BaseModel{ID: cs.ID},
		})
	}

	userID, _ := c.Locals("userID").(uint)
	if err := h.cardService.UpdateCard(c.UserContext(), uint(id), card, userID); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart güncellenemedi: "+err.Error())
		return c.Redirect("/dashboard/cards/update/"+strconv.Itoa(id), http.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla güncellendi.")
	return c.Redirect("/dashboard/cards", http.StatusFound)
}

func (h *DashboardCardHandler) DeleteCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.cardService.DeleteCard(c.UserContext(), uint(id)); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart silinemedi")
	} else {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla silindi.")
	}
	return c.Redirect("/dashboard/cards", fiber.StatusSeeOther)
}
