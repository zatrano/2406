package requests

import (
	"davet.link/pkg/flashmessages"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"regexp"
	"sort"
)

type CardRequest struct {
	Name            string                   `form:"name" validate:"required,min=2"`
	Slug            string                   `form:"slug" validate:"required"`
	Title           string                   `form:"title"`
	Photo           string                   `form:"photo"`
	Telephone       string                   `form:"telephone"`
	Email           string                   `form:"email"`
	Location        string                   `form:"location"`
	WebsiteUrl      string                   `form:"website_url"`
	StoreUrl        string                   `form:"store_url"`
	IsActive        string                   `form:"is_active" validate:"required"`
	UserID          uint                     `form:"user_id"`
	CardBanks       []CardBankRequest        `form:"card_banks"`
	CardSocialMedia []CardSocialMediaRequest `form:"card_social_media"`
}

type CardBankRequest struct {
	ID     uint   // yeni eklendi
	BankID uint
	IBAN   string
}

type CardSocialMediaRequest struct {
	ID            uint   // yeni eklendi
	SocialMediaID uint
	URL           string
}

func validateCardRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
	if err := c.BodyParser(req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		if msg, ok := errorMessages[err.Field()+"_"+err.Tag()]; ok {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, msg)
		} else {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kart bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func parseCardBanksFromForm(c *fiber.Ctx) []CardBankRequest {
	params := c.Request().PostArgs()
	bankIndexes := map[string]struct{}{}
	re := regexp.MustCompile(`card_banks\[(\d+)\]`)
	params.VisitAll(func(key, _ []byte) {
		matches := re.FindSubmatch(key)
		if len(matches) == 2 {
			bankIndexes[string(matches[1])] = struct{}{}
		}
	})
	var indexes []int
	for idx := range bankIndexes {
		var i int
		fmt.Sscanf(idx, "%d", &i)
		indexes = append(indexes, i)
	}
	sort.Ints(indexes)
	var cardBanks []CardBankRequest
	for _, i := range indexes {
		bankID := c.FormValue(fmt.Sprintf("card_banks[%d][bank_id]", i))
		iban := c.FormValue(fmt.Sprintf("card_banks[%d][iban]", i))
		idStr := c.FormValue(fmt.Sprintf("card_banks[%d][id]", i))
		var id uint
		if idStr != "" {
			fmt.Sscanf(idStr, "%d", &id)
		}
		if bankID != "" && iban != "" {
			var bid uint
			_, err2 := fmt.Sscanf(bankID, "%d", &bid)
			if err2 == nil {
				cardBanks = append(cardBanks, CardBankRequest{ID: id, BankID: bid, IBAN: iban})
			}
		}
	}
	return cardBanks
}

func parseCardSocialMediaFromForm(c *fiber.Ctx) []CardSocialMediaRequest {
	params := c.Request().PostArgs()
	socialIndexes := map[string]struct{}{}
	re := regexp.MustCompile(`card_social_media\[(\d+)\]`)
	params.VisitAll(func(key, _ []byte) {
		matches := re.FindSubmatch(key)
		if len(matches) == 2 {
			socialIndexes[string(matches[1])] = struct{}{}
		}
	})
	var indexes []int
	for idx := range socialIndexes {
		var i int
		fmt.Sscanf(idx, "%d", &i)
		indexes = append(indexes, i)
	}
	sort.Ints(indexes)
	var cardSocialMedia []CardSocialMediaRequest
	for _, i := range indexes {
		smid := c.FormValue(fmt.Sprintf("card_social_media[%d][social_media_id]", i))
		url := c.FormValue(fmt.Sprintf("card_social_media[%d][url]", i))
		idStr := c.FormValue(fmt.Sprintf("card_social_media[%d][id]", i))
		var id uint
		if idStr != "" {
			fmt.Sscanf(idStr, "%d", &id)
		}
		if smid != "" && url != "" {
			var sid uint
			_, err2 := fmt.Sscanf(smid, "%d", &sid)
			if err2 == nil {
				cardSocialMedia = append(cardSocialMedia, CardSocialMediaRequest{ID: id, SocialMediaID: sid, URL: url})
			}
		}
	}
	return cardSocialMedia
}

func ValidateCardRequest(c *fiber.Ctx) error {
	var req CardRequest
	errorMessages := map[string]string{
		"Name_required":   "Kart adı zorunludur",
		"Name_min":        "Kart adı en az 2 karakter olmalıdır",
		"Slug_required":   "Slug alanı zorunludur",
		"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir",
	}

	// slug_without_at varsa onu slug'a ata
	slugWithoutAt := c.FormValue("slug_without_at")
	if slugWithoutAt != "" {
		req.Slug = slugWithoutAt
	}

	if err := validateCardRequest(c, &req, errorMessages, "/dashboard/cards/create"); err != nil {
		return err
	}

	req.CardBanks = parseCardBanksFromForm(c)
	req.CardSocialMedia = parseCardSocialMediaFromForm(c)

	c.Locals("cardRequest", req)
	return nil
}

// Yeni: Update işlemleri için redirect path parametreli fonksiyon
func ValidateCardRequestWithPath(c *fiber.Ctx, redirectPath string) error {
	var req CardRequest
	errorMessages := map[string]string{
		"Name_required":   "Kart adı zorunludur",
		"Name_min":        "Kart adı en az 2 karakter olmalıdır",
		"Slug_required":   "Slug alanı zorunludur",
		"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir",
	}
	// slug_without_at varsa onu slug'a ata
	slugWithoutAt := c.FormValue("slug_without_at")
	if slugWithoutAt != "" {
		req.Slug = slugWithoutAt
	}
	if err := validateCardRequest(c, &req, errorMessages, redirectPath); err != nil {
		return err
	}
	req.CardBanks = parseCardBanksFromForm(c)
	req.CardSocialMedia = parseCardSocialMediaFromForm(c)
	c.Locals("cardRequest", req)
	return nil
}
