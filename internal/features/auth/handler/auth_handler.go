package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/zerodayz7/http-server/internal/errors"
	"github.com/zerodayz7/http-server/internal/features/auth/service"
	"github.com/zerodayz7/http-server/internal/shared"
	"github.com/zerodayz7/http-server/internal/shared/logger"
	"github.com/zerodayz7/http-server/internal/shared/security"
	"github.com/zerodayz7/http-server/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) InitSession(c *fiber.Ctx) error {
	// Pobierz sesję, Fiber automatycznie ją tworzy
	sess := c.Locals("session").(*session.Session)

	// Możesz tu opcjonalnie ustawić dane w sesji
	UserID := shared.GenerateUuid()
	sess.Set("UserID", UserID)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	// Zwróć JSON, nawet jeśli niczego nie potrzebujesz
	return c.JSON(fiber.Map{
		"sessionInitialized": true,
		"UserID":             UserID,
	})
}

func (h *AuthHandler) GetCSRFToken(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)

	csrfToken := sess.Get("csrfToken")
	if csrfToken == nil {
		csrfToken = shared.GenerateUuid()
		sess.Set("csrfToken", csrfToken)
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save session",
			})
		}
	}

	return c.JSON(fiber.Map{"csrf_token": csrfToken})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	log := logger.GetLogger()

	body := c.Locals("validatedBody").(validator.LoginRequest)
	log.InfoObj("Login attempt", map[string]any{
		"email": body.Email,
	})

	user, err := h.authService.GetUserByEmail(body.Email)
	if err != nil {
		log.WarnObj("Login failed: user not found", map[string]any{
			"email": body.Email,
		})
		return errors.SendAppError(c, errors.ErrInvalidCredentials)
	}

	log.InfoObj("User found", map[string]any{
		"id": user.ID,
	})

	valid, err := h.authService.VerifyPassword(user, body.Password)
	if err != nil {
		log.ErrorObj("Error verifying password", err)
		return errors.SendAppError(c, errors.ErrInternal)
	}
	if !valid {
		log.WarnObj("Invalid password attempt", map[string]any{
			"userID": user.ID,
		})
		return errors.SendAppError(c, errors.ErrInvalidCredentials)
	}

	if user.TwoFactorEnabled && user.TwoFactorSecret != "" {
		log.InfoObj("2FA required for user", map[string]any{
			"userID": user.ID,
		})
		return c.JSON(fiber.Map{"2fa_required": true})
	}

	token, err := security.GenerateToken(fmt.Sprint(user.ID))
	if err != nil {
		log.ErrorObj("Failed to generate JWT", err)
		return errors.SendAppError(c, errors.ErrInternal)
	}
	log.InfoObj("JWT generated for user", map[string]any{
		"userID": user.ID,
	})

	payload := fiber.Map{
		"2fa_required": false,
		"token":        token,
	}
	log.InfoObj("Login successful", payload)

	return c.JSON(payload)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(validator.RegisterRequest)

	user, err := h.authService.Register(body.Username, body.Email, body.Password)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			errors.AttachRequestMeta(c, appErr, "requestID")
			return appErr
		}
		return errors.ErrInternal
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

func (h *AuthHandler) Verify2FA(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(validator.TwoFARequest)

	sess := c.Locals("session").(*session.Session)
	userID := sess.Get("userID").(uint)

	ok, err := h.authService.Verify2FACodeByID(userID, body.Code)
	if err != nil {
		return errors.SendAppError(c, errors.ErrInvalidCredentials)
	}
	if !ok {
		return errors.SendAppError(c, errors.ErrInvalid2FACode)
	}

	return c.JSON(fiber.Map{
		"message": "2FA verified successfully",
	})
}
