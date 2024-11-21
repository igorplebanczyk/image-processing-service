package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/common/server/respond"
	"image-processing-service/src/internal/images/application"
	"image-processing-service/src/internal/images/domain"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ImageAPI struct {
	ImagesService *application.ImageService
}

func NewServer(imagesService *application.ImageService) *ImageAPI {
	return &ImageAPI{ImagesService: imagesService}
}

func (s *ImageAPI) Upload(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type response struct {
		Name      string `json:"name"`
		Size      int64  `json:"size"`
		CreatedAt string `json:"created_at"`
	}

	name := r.FormValue("name")

	err := r.ParseMultipartForm(domain.MaxImageSize)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput(fmt.Sprintf("image size exceeds %d bytes", domain.MaxImageSize)))
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("image file not found"))
		return
	}
	defer imageFile.Close()

	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid image file"))
		return
	}

	image, err := s.ImagesService.UploadImage(userID, name, imageBytes)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusCreated, response{
		Name:      image.Name,
		Size:      int64(len(imageBytes)),
		CreatedAt: image.CreatedAt.String(),
	})
}

func (s *ImageAPI) GetDataAll(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type responseItem struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type response struct {
		Images     []responseItem `json:"images"`
		Page       int            `json:"page"`
		Limit      int            `json:"limit"`
		TotalCount int            `json:"total_count"`
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	images, total, err := s.ImagesService.ListUserImages(userID, &page, &limit)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	var responseImages []responseItem
	for _, img := range images {
		responseImages = append(responseImages, responseItem{
			Name:      img.Name,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		})
	}

	respond.WithJSON(w, http.StatusOK, response{
		Images:     responseImages,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
	})
}

func (s *ImageAPI) GetData(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	type response struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	imageData, err := s.ImagesService.GetImageData(userID, p.Name)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		Name:      imageData.Name,
		CreatedAt: imageData.CreatedAt,
		UpdatedAt: imageData.UpdatedAt,
	})
}

func (s *ImageAPI) Download(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	imageBytes, err := s.ImagesService.DownloadImage(userID, p.Name)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithImage(w, http.StatusOK, imageBytes, p.Name)
}

func (s *ImageAPI) Delete(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = s.ImagesService.DeleteImage(userID, p.Name)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *ImageAPI) Transform(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name            string                  `json:"name"`
		Transformations []domain.Transformation `json:"transformations"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = s.ImagesService.ApplyTransformations(userID, p.Name, p.Transformations)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *ImageAPI) AdminListAllImages(w http.ResponseWriter, r *http.Request) {
	type responseItem struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type response struct {
		Images     []responseItem `json:"images"`
		Page       int            `json:"page"`
		Limit      int            `json:"limit"`
		TotalCount int            `json:"total_count"`
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	images, total, err := s.ImagesService.AdminListAllImages(&page, &limit)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	var responseImages []responseItem
	for _, img := range images {
		responseImages = append(responseImages, responseItem{
			Name:      img.Name,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		})
	}

	respond.WithJSON(w, http.StatusOK, response{
		Images:     responseImages,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
	})
}

func (s *ImageAPI) AdminDeleteImage(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = s.ImagesService.AdminDeleteImage(p.ID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
