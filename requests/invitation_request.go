package requests

import (
	"davet.link/pkg/flashmessages"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationRequest struct {
	InvitationKey     string   `form:"invitation_key" validate:"required,min=2"`
	UserID            uint     `form:"user_id" validate:"required,gt=0"`
	CategoryID        uint     `form:"category_id" validate:"required,gt=0"`
	Template          string   `form:"template"`
	Type              string   `form:"type"`
	Title             string   `form:"title" validate:"required,min=2"`
	Image             string   `form:"image"`
	Description       string   `form:"description"`
	Venue             string   `form:"venue"`
	Address           string   `form:"address"`
	Location          string   `form:"location"`
	Link              string   `form:"link"`
	Telephone         string   `form:"telephone"`
	Note              string   `form:"note"`
	Date              string   `form:"date"`
	Time              string   `form:"time"`
	IsConfirmed       string   `form:"is_confirmed"`
	IsParticipant     string   `form:"is_participant"`
	DetailTitle       string   `form:"detail_title"`
	DetailPerson      string   `form:"detail_person"`
}

func validateInvitationRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
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
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz davetiye bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}
	return nil
}

func ValidateInvitationRequest(c *fiber.Ctx) error {
	var req InvitationRequest
	errorMessages := map[string]string{
		"InvitationKey_required": "Davetiye anahtarı zorunludur",
		"InvitationKey_min":      "Davetiye anahtarı en az 2 karakter olmalıdır",
		"UserID_required":        "Kullanıcı seçimi zorunludur",
		"UserID_gt":              "Kullanıcı seçimi zorunludur",
		"CategoryID_required":    "Kategori seçimi zorunludur",
		"CategoryID_gt":          "Kategori seçimi zorunludur",
		"Title_required":         "Başlık zorunludur",
		"Title_min":              "Başlık en az 2 karakter olmalıdır",
	}
	if err := validateInvitationRequest(c, &req, errorMessages, "/dashboard/invitations/create"); err != nil {
		return err
	}
	c.Locals("invitationRequest", req)
	return nil
}
