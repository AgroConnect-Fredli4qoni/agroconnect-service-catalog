package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/agroconnect/service-catalog/models"
	"github.com/agroconnect/service-catalog/repository"
	"github.com/gorilla/mux"
)

type FarmerHandler struct {
	repo repository.FarmerRepository
}

func NewFarmerHandler(repo repository.FarmerRepository) *FarmerHandler {
	return &FarmerHandler{repo: repo}
}

func (h *FarmerHandler) GetAllFarmers(w http.ResponseWriter, r *http.Request) {
	farmers, err := h.repo.FindAll(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(farmers)
}

func (h *FarmerHandler) GetFarmerBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	farmer, err := h.repo.FindBySlug(r.Context(), slug)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(farmer)
}

func (h *FarmerHandler) CreateFarmer(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateFarmerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload format"})
		return
	}

	if dto.Name == "" || dto.Slug == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Validation failed: Name and Slug are required"})
		return
	}

	farmer := models.Farmer{
		Slug:           dto.Slug,
		Name:           dto.Name,
		OriginRegion:   dto.OriginRegion,
		AvatarURL:      dto.AvatarURL,
		BannerURL:      dto.BannerURL,
		Description:    dto.Description,
		Phone:          dto.Phone,
		Address:        dto.Address,
		OperatingHours: dto.OperatingHours,
		LandArea:       dto.LandArea,
		IsVerified:     dto.IsVerified,
		FarmingMethods: dto.FarmingMethods,
		Certifications: dto.Certifications,
	}

	if err := h.repo.Create(r.Context(), &farmer); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(farmer)
}
