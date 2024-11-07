package images

import (
	"encoding/json"
	"fmt"
	"image-processing-service/internal/services/server/util"
	"image-processing-service/internal/users"
	"io"
	"net/http"
	"time"
)

type Config struct {
	repo           Repository
	storage        StorageService
	cache          CacheService
	transformation TransformationService
}

func NewConfig(repo Repository, storage StorageService, cache CacheService, transformation TransformationService) *Config {
	return &Config{repo: repo, storage: storage, cache: cache, transformation: transformation}
}

func (cfg *Config) Upload(user *users.User, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	ok, err := validate(cfg.repo, user.ID, name)
	if !ok {
		util.RespondWithError(w, http.StatusBadRequest, "invalid image name")
		return
	}
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating image: %v", err))
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "image too large")
		return
	}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "missing image")
		return
	}
	defer imageFile.Close()

	imageBytes, err := io.ReadAll(imageFile)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error reading image: %v", err))
		return
	}

	_, err = cfg.repo.CreateImage(user.ID, name)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating image: %v", err))
		return
	}

	objectName := fmt.Sprintf("%s-%s", user.ID, name)
	err = cfg.storage.Upload(r.Context(), objectName, imageBytes)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error uploading image: %v", err))
		return
	}

	err = cfg.cache.Set(objectName, imageBytes, 30*time.Minute)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error caching image: %v", err))
		return
	}

	util.RespondWithoutContent(w, http.StatusCreated)
}

func (cfg *Config) Download(user *users.User, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	objectName := fmt.Sprintf("%s-%s", user.ID, p.Name)

	imageBytes, err := cfg.cache.Get(objectName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get image from cache: %v", err))
		return
	}
	if imageBytes != nil {
		w.Header().Set("X-Cache", "HIT")
		util.RespondWithImage(w, http.StatusOK, imageBytes, p.Name)
		return
	}

	imageBytes, err = cfg.storage.Download(r.Context(), objectName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to download image: %v", err))
		return
	}

	err = cfg.cache.Set(objectName, imageBytes, 30*time.Minute)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to cache image: %v", err))
		return
	}

	w.Header().Set("X-Cache", "MISS")
	util.RespondWithImage(w, http.StatusOK, imageBytes, p.Name)
}

func (cfg *Config) Delete(user *users.User, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = cfg.repo.DeleteImage(user.ID, p.Name)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete image from db: %v", err))
		return
	}

	err = cfg.cache.Delete(fmt.Sprintf("%s-%s", user.ID, p.Name))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete image from cache: %v", err))
		return
	}

	objectName := fmt.Sprintf("%s-%s", user.ID, p.Name)
	err = cfg.storage.Delete(r.Context(), objectName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete image from storage: %v", err))
		return
	}

	util.RespondWithoutContent(w, http.StatusNoContent)
}

func (cfg *Config) Transform(user *users.User, w http.ResponseWriter, r *http.Request) {
	type transformations struct {
		Type    string         `json:"type"`
		Options map[string]any `json:"options"`
	}

	type parameters struct {
		Name            string            `json:"name"`
		Transformations []transformations `json:"transformations"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	objectName := fmt.Sprintf("%s-%s", user.ID, p.Name)

	imageBytes, err := cfg.cache.Get(objectName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get image from cache: %v", err))
		return
	}

	if imageBytes != nil {
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("X-Cache", "MISS")
		imageBytes, err = cfg.storage.Download(r.Context(), objectName)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to download image: %v", err))
			return
		}
	}

	for _, transformation := range p.Transformations {
		imageBytes, err = cfg.transformation.Transform(imageBytes, transformation.Type, transformation.Options)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to apply transformation %s: %v", transformation.Type, err))
			return
		}
	}

	err = cfg.storage.Upload(r.Context(), objectName, imageBytes)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to upload transformed image: %v", err))
		return
	}

	err = cfg.repo.UpdateImage(user.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update image: %v", err))
		return
	}

	err = cfg.cache.Set(objectName, imageBytes, 30*time.Minute)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to cache image: %v", err))
		return
	}

	util.RespondWithImage(w, http.StatusOK, imageBytes, p.Name)
}
