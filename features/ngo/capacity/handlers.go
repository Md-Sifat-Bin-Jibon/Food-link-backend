package capacity

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) getUserID(r *http.Request) (uuid.UUID, error) {
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok || user == nil {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return user.ID, nil
}

// Get handles GET /api/v1/ngo/capacity
// @Summary      Get capacity settings
// @Description  Get NGO capacity settings for the authenticated user
// @Tags         ngo-capacity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  NGOCapacitySettings
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /ngo/capacity [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	settings, err := h.service.GetByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve capacity settings", err.Error())
		return
	}
	utils.OKResponse(w, "Capacity settings retrieved successfully", settings)
}

// CreateOrUpdate handles POST /api/v1/ngo/capacity
// @Summary      Create or update capacity settings
// @Description  Create or update NGO capacity settings
// @Tags         ngo-capacity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateNGOCapacitySettingsRequest  true  "Capacity settings data"
// @Success      200      {object}  NGOCapacitySettings
// @Success      201      {object}  NGOCapacitySettings
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /ngo/capacity [post]
func (h *Handler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateNGOCapacitySettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	settings, err := h.service.CreateOrUpdate(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to save capacity settings", err.Error())
		return
	}
	utils.OKResponse(w, "Capacity settings saved successfully", settings)
}
