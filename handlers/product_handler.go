package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/agroconnect/service-catalog/models"
	"github.com/agroconnect/service-catalog/repository"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) isUserAuthorized(r *http.Request, product *models.Product) bool {
	userRole := r.Header.Get("X-User-Role")
	if userRole == "admin" {
		return true
	}

	userIDStr := r.Header.Get("X-User-ID")
	userName := r.Header.Get("X-User-Name")

	if product.FarmerID > 0 && userIDStr != "" {
		if strconv.Itoa(product.FarmerID) == userIDStr {
			return true
		}
	}

	if userName != "" && product.FarmerName != "" {
		if strings.EqualFold(strings.TrimSpace(product.FarmerName), strings.TrimSpace(userName)) {
			return true
		}
	}

	return false
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	category := r.URL.Query().Get("category")

	products, err := h.repo.FindAll(r.Context(), search, category)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload format"})
		return
	}

	if dto.Name == "" || dto.PricePerKg <= 0 || dto.StockKg < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed: Name, positive Price, and non-negative Stock are required"})
		return
	}

	unit := dto.Unit
	if unit == "" {
		unit = "Kg"
	}

	farmerName := dto.FarmerName
	if headerName := r.Header.Get("X-User-Name"); headerName != "" && farmerName == "" {
		farmerName = headerName
	}

	farmerID := dto.FarmerID
	if headerID := r.Header.Get("X-User-ID"); headerID != "" && farmerID == 0 {
		if idNum, err := strconv.Atoi(headerID); err == nil {
			farmerID = idNum
		}
	}

	product := models.Product{
		Name:            dto.Name,
		Category:        dto.Category,
		PricePerKg:      dto.PricePerKg,
		StockKg:         dto.StockKg,
		Unit:            unit,
		OriginRegion:    dto.OriginRegion,
		FarmerID:        farmerID,
		FarmerName:      farmerName,
		FarmerAvatarURL: dto.FarmerAvatarURL,
		IsOrganic:       dto.IsOrganic,
		Description:     dto.Description,
		ImageURL:        dto.ImageURL,
	}

	if err := h.repo.Create(r.Context(), &product); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	existing, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}

	if !h.isUserAuthorized(r, existing) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Anda tidak memiliki hak akses untuk mengubah komoditas ini"})
		return
	}

	var dto models.UpdateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload format"})
		return
	}

	if dto.Name == "" || dto.PricePerKg <= 0 || dto.StockKg < 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed: Name, positive Price, and non-negative Stock are required"})
		return
	}

	unit := dto.Unit
	if unit == "" {
		unit = existing.Unit
		if unit == "" {
			unit = "Kg"
		}
	}

	existing.Name = dto.Name
	existing.Category = dto.Category
	existing.PricePerKg = dto.PricePerKg
	existing.StockKg = dto.StockKg
	existing.Unit = unit
	existing.OriginRegion = dto.OriginRegion
	existing.IsOrganic = dto.IsOrganic
	existing.Description = dto.Description
	if dto.ImageURL != "" {
		existing.ImageURL = dto.ImageURL
	}

	if err := h.repo.Update(r.Context(), id, existing); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existing)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	existing, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if !h.isUserAuthorized(r, existing) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Anda tidak memiliki hak akses untuk menghapus komoditas ini"})
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product successfully deleted"})
}

func (h *ProductHandler) DeductProductStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dto models.StockUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	if dto.Quantity <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Quantity must be positive"})
		return
	}

	if err := h.repo.DeductStock(r.Context(), id, dto.Quantity); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Stock successfully updated",
		"id":      id,
	})
}
