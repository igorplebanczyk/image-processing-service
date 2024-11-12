package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"image-processing-service/internal/common/server/respond"
	"image-processing-service/internal/images/application"
	"image-processing-service/internal/images/domain"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ImageServer struct {
	ImagesService *application.ImageService
}

func NewServer(imagesService *application.ImageService) *ImageServer {
	return &ImageServer{ImagesService: imagesService}
}

func (s *ImageServer) Upload(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type response struct {
		Name      string `json:"name"`
		Size      int64  `json:"size"`
		CreatedAt string `json:"created_at"`
	}

	name := r.FormValue("name")

	err := r.ParseMultipartForm(domain.MaxImageSize)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "image too large")
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "missing image")
		return
	}
	defer imageFile.Close()

	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading image: %v", err))
		return
	}

	slog.Info(fmt.Sprintf("format at upload: %v", http.DetectContentType(imageBytes)))

	image, err := s.ImagesService.UploadImage(userID, name, imageBytes)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("error uploading image: %v", err))
		return
	}

	respond.WithJSON(w, http.StatusCreated, response{
		Name:      image.Name,
		Size:      int64(len(imageBytes)),
		CreatedAt: image.CreatedAt.String(),
	})
}

func (s *ImageServer) List(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching images: %v", err))
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

func (s *ImageServer) Info(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	imageData, err := s.ImagesService.GetImageData(userID, p.Name)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get image: %v", err))
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		Name:      imageData.Name,
		CreatedAt: imageData.CreatedAt,
		UpdatedAt: imageData.UpdatedAt,
	})
}

func (s *ImageServer) Download(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	imageBytes, err := s.ImagesService.DownloadImage(userID, p.Name)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to download image: %v", err))
		return
	}

	slog.Info(fmt.Sprintf("format at download: %v", http.DetectContentType(imageBytes)))

	respond.WithImage(w, http.StatusOK, imageBytes, p.Name)
}

func (s *ImageServer) Delete(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.ImagesService.DeleteImage(userID, p.Name)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete image: %v", err))
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *ImageServer) Transform(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name            string                  `json:"name"`
		Transformations []domain.Transformation `json:"transformations"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	slog.Info(fmt.Sprintf("Transformations: %v", p.Transformations))

	err = s.ImagesService.ApplyTransformations(userID, p.Name, p.Transformations)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to apply transformations: %v", err))
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
