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
	repo    Repository
	storage StorageService
	cache   CacheService
}

func NewConfig(repo Repository, storage StorageService, cache CacheService) *Config {
	return &Config{repo: repo, storage: storage, cache: cache}
}

func (cfg *Config) Upload(user *users.User, w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	ok, err := validate(cfg.repo, user.ID, name)
	if !ok {
		util.RespondWithError(w, http.StatusBadRequest, "image with this name already exists")
		return
	}
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to validate image")
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "file too large")
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
		util.RespondWithError(w, http.StatusInternalServerError, "failed to read image")
		return
	}

	_, err = cfg.repo.CreateImage(user.ID, name)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to create image")
		return
	}

	blobName := fmt.Sprintf("%s-%s", user.ID, name)
	err = cfg.storage.UploadObject(r.Context(), blobName, imageBytes)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to upload image")
		return
	}

	err = cfg.cache.Set(blobName, imageBytes, 30*time.Minute)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to cache image")
		return
	}

	util.RespondWithText(w, http.StatusOK, "image uploaded successfully")
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

	blobName := fmt.Sprintf("%s-%s", user.ID, p.Name)

	imageBytes, err := cfg.cache.Get(blobName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve image from cache")
		return
	}
	if imageBytes != nil {
		w.Header().Set("X-Cache", "HIT")
		util.RespondWithImage(w, http.StatusOK, imageBytes, p.Name)
		return
	}

	imageBytes, err = cfg.storage.DownloadObject(r.Context(), blobName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to download image")
		return
	}

	err = cfg.cache.Set(blobName, imageBytes, 30*time.Minute)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to cache image")
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
		util.RespondWithError(w, http.StatusInternalServerError, "failed to delete image")
		return
	}

	err = cfg.cache.Delete(fmt.Sprintf("%s-%s", user.ID, p.Name))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to delete image from cache")
		return
	}

	blobName := fmt.Sprintf("%s-%s", user.ID, p.Name)
	err = cfg.storage.DeleteObject(r.Context(), blobName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to delete image")
		return
	}

	util.RespondWithText(w, http.StatusOK, "image deleted successfully")
}
