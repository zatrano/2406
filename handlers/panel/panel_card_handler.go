package handlers

import (
	"net/http"

	"davet.link/models"
	"davet.link/pkg/queryparams"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

type PanelCardHandler struct {
	cardService        services.ICardService
	userService        services.IUserService
	bankService        services.IBankService
	socialMediaService services.ISocialMediaService
}

func NewPanelCardHandler() *PanelCardHandler {
	return &PanelCardHandler{
		cardService:        services.NewCardService(),
		userService:        services.NewUserService(),
		bankService:        services.NewBankService(),
		socialMediaService: services.NewSocialMediaService(),
	}
}

func (h *PanelCardHandler) ListCards(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	params := queryparams.ListParams{
		Page:    1,
		PerPage: 1,
	}
	cards, totalCount, err := h.cardService.GetCardsByUserID(userID)
	paginatedResult := &queryparams.PaginatedResult{
		Data: cards,
		Meta: queryparams.PaginationMeta{
			CurrentPage: 1,
			PerPage:     1,
			TotalItems:  totalCount,
			TotalPages:  1,
		},
	}
	renderData := fiber.Map{
		"Title":  "Kartım",
		"Result": paginatedResult,
		"Params": params,
	}
	if err != nil {
		renderData[renderer.FlashErrorKeyView] = "Kart getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.Card{},
			Meta: queryparams.PaginationMeta{CurrentPage: 1, PerPage: 1},
		}
	}
	return renderer.Render(c, "panel/cards/list", "layouts/panel", renderData, http.StatusOK)
}

func (h *PanelCardHandler) ShowCreateCard(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	// Check if user already has a card
	cards, _, _ := h.cardService.GetCardsByUserID(userID)
	if len(cards) > 0 {
		return c.Redirect("/panel/cards", http.StatusFound)
	}
	banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
		"Title":        "Yeni Kart Oluştur",
		"Banks":        banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	}, http.StatusOK)
}

func (h *PanelCardHandler) CreateCard(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	// Check if user already has a card
	cards, _, _ := h.cardService.GetCardsByUserID(userID)
	if len(cards) > 0 {
		return c.Redirect("/panel/cards", http.StatusFound)
	}
	if err := requests.ValidateCardRequest(c); err != nil {
		banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
		return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
			"Title":        "Yeni Kart Oluştur",
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
			renderer.FlashErrorKeyView: err.Error(),
			"FormData": c.Locals("cardRequest"),
		}, http.StatusBadRequest)
	}
	req := c.Locals("cardRequest").(requests.CardRequest)
	card := &models.Card{
		Name:      req.Name,
		Slug:      req.Slug,
		UserID:    userID,
		Photo:     req.Photo,
		Telephone: req.Telephone,
		Email:     req.Email,
		Location:  req.Location,
		WebsiteUrl: req.WebsiteUrl,
		StoreUrl: req.StoreUrl,
		IsActive:  req.IsActive == "true",
	}
	for _, bank := range req.CardBanks {
		card.CardBanks = append(card.CardBanks, models.CardBank{CardID: card.ID, BankID: bank.BankID, IBAN: bank.IBAN})
	}
	for _, sm := range req.CardSocialMedia {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{CardID: card.ID, SocialMediaID: sm.SocialMediaID, URL: sm.URL})
	}
	if err := h.cardService.CreateCard(c.UserContext(), card); err != nil {
		banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
		return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
			"Title":        "Yeni Kart Oluştur",
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
			renderer.FlashErrorKeyView: "Kart oluşturulamadı",
			"FormData": req,
		}, http.StatusInternalServerError)
	}
	return c.Redirect("/panel/cards", http.StatusFound)
}

func (h *PanelCardHandler) ShowUpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	userID, _ := c.Locals("userID").(uint)
	card, err := h.cardService.GetCardByID(c.UserContext(), uint(id))
	if err != nil || card.UserID != userID {
		return c.Status(http.StatusNotFound).SendString("Kart bulunamadı")
	}
	banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
	socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
	return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
		"Title":        "Kartı Düzenle",
		"Card":         card,
		"Banks":        banksResult.Data,
		"SocialMedias": socialMediasResult.Data,
	}, http.StatusOK)
}

func (h *PanelCardHandler) UpdateCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	userID, _ := c.Locals("userID").(uint)
	card, err := h.cardService.GetCardByID(c.UserContext(), uint(id))
	if err != nil || card.UserID != userID {
		return c.Status(http.StatusNotFound).SendString("Kart bulunamadı")
	}
	if err := requests.ValidateCardRequest(c); err != nil {
		banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
		return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
			"Title":        "Kartı Düzenle",
			"Card":         card,
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
			renderer.FlashErrorKeyView: err.Error(),
			"FormData": c.Locals("cardRequest"),
		}, http.StatusBadRequest)
	}
	req := c.Locals("cardRequest").(requests.CardRequest)
	card.Name = req.Name
	card.Slug = req.Slug
	card.Photo = req.Photo
	card.Telephone = req.Telephone
	card.Email = req.Email
	card.Location = req.Location
	card.WebsiteUrl = req.WebsiteUrl
	card.StoreUrl = req.StoreUrl
	card.IsActive = req.IsActive == "true"
	card.CardBanks = nil
	for _, bank := range req.CardBanks {
		card.CardBanks = append(card.CardBanks, models.CardBank{CardID: uint(id), BankID: bank.BankID, IBAN: bank.IBAN})
	}
	card.CardSocialMedia = nil
	for _, sm := range req.CardSocialMedia {
		card.CardSocialMedia = append(card.CardSocialMedia, models.CardSocialMedia{CardID: uint(id), SocialMediaID: sm.SocialMediaID, URL: sm.URL})
	}
	if err := h.cardService.UpdateCard(c.UserContext(), uint(id), card); err != nil {
		banksResult, _ := h.bankService.GetAllBanks(queryparams.ListParams{PerPage: 1000})
		socialMediasResult, _ := h.socialMediaService.GetAllSocialMedias(queryparams.ListParams{PerPage: 1000})
		return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
			"Title":        "Kartı Düzenle",
			"Card":         card,
			"Banks":        banksResult.Data,
			"SocialMedias": socialMediasResult.Data,
			renderer.FlashErrorKeyView: "Kart güncellenemedi",
			"FormData": req,
		}, http.StatusInternalServerError)
	}
	return c.Redirect("/panel/cards", http.StatusFound)
}

func (h *PanelCardHandler) DeleteCard(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	userID, _ := c.Locals("userID").(uint)
	card, err := h.cardService.GetCardByID(c.UserContext(), uint(id))
	if err != nil || card.UserID != userID {
		return c.Status(http.StatusNotFound).SendString("Kart bulunamadı")
	}
	if err := h.cardService.DeleteCard(c.UserContext(), uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Kart silinemedi")
	}
	return c.Redirect("/panel/cards", http.StatusFound)
}
