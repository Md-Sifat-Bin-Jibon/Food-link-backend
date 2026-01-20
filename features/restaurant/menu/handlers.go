package menu

import (
	"encoding/json"
	"foodlink_backend/errors"
	"foodlink_backend/features/auth"
	"foodlink_backend/utils"
	"net/http"
	"strings"

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

// GetAll handles GET /api/v1/restaurant/menu
// @Summary      List menu items
// @Description  Get all menu items for the authenticated restaurant
// @Tags         restaurant-menu
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   RestaurantMenuItem
// @Failure      401  {object}  errors.AppError
// @Router       /restaurant/menu [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	items, err := h.service.GetAllByUserID(userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve menu items", err.Error())
		return
	}
	utils.OKResponse(w, "Menu items retrieved successfully", items)
}

// GetByID handles GET /api/v1/restaurant/menu/:id
// @Summary      Get menu item by ID
// @Description  Get details of a specific menu item
// @Tags         restaurant-menu
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Menu Item ID"
// @Success      200  {object}  RestaurantMenuItem
// @Failure      401  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /restaurant/menu/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/restaurant/menu/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to retrieve menu item", err.Error())
		return
	}
	utils.OKResponse(w, "Menu item retrieved successfully", item)
}

// Create handles POST /api/v1/restaurant/menu
// @Summary      Add menu item
// @Description  Add a new item to restaurant menu
// @Tags         restaurant-menu
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateRestaurantMenuItemRequest  true  "Menu item data"
// @Success      201      {object}  RestaurantMenuItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Router       /restaurant/menu [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	var req CreateRestaurantMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Create(userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to create menu item", err.Error())
		return
	}
	utils.CreatedResponse(w, "Menu item created successfully", item)
}

// Update handles PUT /api/v1/restaurant/menu/:id
// @Summary      Update menu item
// @Description  Update an existing menu item
// @Tags         restaurant-menu
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                      true  "Menu Item ID"
// @Param        request  body      UpdateRestaurantMenuItemRequest true  "Menu item data"
// @Success      200      {object}  RestaurantMenuItem
// @Failure      400      {object}  errors.AppError
// @Failure      401      {object}  errors.AppError
// @Failure      403      {object}  errors.AppError
// @Failure      404      {object}  errors.AppError
// @Router       /restaurant/menu/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/restaurant/menu/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	var req UpdateRestaurantMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequestResponse(w, "Invalid request body", err.Error())
		return
	}
	item, err := h.service.Update(id, userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to update menu item", err.Error())
		return
	}
	utils.OKResponse(w, "Menu item updated successfully", item)
}

// Delete handles DELETE /api/v1/restaurant/menu/:id
// @Summary      Delete menu item
// @Description  Delete a menu item
// @Tags         restaurant-menu
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Menu Item ID"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  errors.AppError
// @Failure      403  {object}  errors.AppError
// @Failure      404  {object}  errors.AppError
// @Router       /restaurant/menu/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.BadRequestResponse(w, "Method not allowed", nil)
		return
	}
	userID, err := h.getUserID(r)
	if err != nil {
		utils.UnauthorizedResponse(w, "Authentication required")
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/restaurant/menu/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(w, "Invalid ID format", nil)
		return
	}
	if err := h.service.Delete(id, userID); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.ErrorResponse(w, appErr.Code, appErr.Message, nil)
			return
		}
		utils.InternalServerErrorResponse(w, "Failed to delete menu item", err.Error())
		return
	}
	utils.OKResponse(w, "Menu item deleted successfully", map[string]string{"message": "Deleted"})
}
