package interfaces

import (
	"encoding/base64"
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
	ImagesService *application.ImagesService
}

func NewAPI(imagesService *application.ImagesService) *ImageAPI {
	return &ImageAPI{ImagesService: imagesService}
}

func (a *ImageAPI) Upload(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")

	err := r.ParseMultipartForm(domain.MaxImageSize)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput(fmt.Sprintf("image size exceeds %d bytes", domain.MaxImageSize)))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("image file not found"))
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid image file"))
		return
	}

	err = a.ImagesService.Upload(userID, name, description, bytes)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusCreated)
}

func (a *ImageAPI) Get(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	type response struct {
		Metadata domain.ImageMetadata `json:"metadata"`
		Image    string               `json:"image"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	metadata, bytes, err := a.ImagesService.Get(userID, p.Name)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	encodedImage := base64.StdEncoding.EncodeToString(bytes)

	respond.WithJSON(w, http.StatusOK, response{
		Metadata: *metadata,
		Image:    encodedImage,
	})
}

func (a *ImageAPI) GetAll(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type responseImageMetadata struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		UpdatedAt   string `json:"updated_at"`
		CreatedAt   string `json:"created_at"`
	}

	type responseImage struct {
		Metadata     responseImageMetadata `json:"metadata"`
		ImagePreview string                `json:"image_preview"`
	}

	type response struct {
		Images     []responseImage `json:"images"`
		TotalCount int             `json:"total_count"`
		Page       int             `json:"page"`
		Limit      int             `json:"limit"`
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid page"))
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid limit"))
		return
	}

	metadata, previews, totalCount, err := a.ImagesService.GetAll(userID, page, limit)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	var respImages []responseImage
	for i, m := range metadata {
		encodedPreview := base64.StdEncoding.EncodeToString(previews[i])
		respImages = append(respImages, responseImage{
			Metadata: responseImageMetadata{
				Name:        m.Name,
				Description: m.Description,
				UpdatedAt:   m.UpdatedAt.String(),
				CreatedAt:   m.CreatedAt.String(),
			},
			ImagePreview: encodedPreview,
		})
	}

	respond.WithJSON(w, http.StatusOK, response{Images: respImages, TotalCount: totalCount, Page: page, Limit: limit})
}

func (a *ImageAPI) UpdateDetails(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		OldName        string `json:"old_name"`
		NewName        string `json:"new_name"`
		NewDescription string `json:"new_description"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.ImagesService.UpdateDetails(userID, p.OldName, p.NewName, p.NewDescription)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *ImageAPI) Transform(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
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

	err = a.ImagesService.Transform(userID, p.Name, p.Transformations)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *ImageAPI) Delete(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
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

	err = a.ImagesService.Delete(userID, p.Name)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *ImageAPI) AdminListAllImages(w http.ResponseWriter, r *http.Request) {
	type responseImage struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	type response struct {
		Images     []responseImage `json:"images"`
		Page       int             `json:"page"`
		Limit      int             `json:"limit"`
		TotalCount int             `json:"total_count"`
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid page"))
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid limit"))
		return
	}

	images, total, err := a.ImagesService.AdminListAllImages(page, limit)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	var responseImages []responseImage
	for _, img := range images {
		responseImages = append(responseImages, responseImage{
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

func (a *ImageAPI) AdminDeleteImage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid user ID"))
		return
	}

	err = a.ImagesService.AdminDeleteImage(id)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
