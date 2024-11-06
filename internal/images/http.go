package images

import (
	"encoding/json"
	"fmt"
	"image-processing-service/internal/services/server/util"
	"image-processing-service/internal/users"
	"io"
	"net/http"
)

type Config struct {
	repo    Repository
	storage StorageService
}

func NewConfig(repo Repository, storage StorageService) *Config {
	return &Config{repo: repo, storage: storage}
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

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "image uploaded successfully"})
}

func (cfg *Config) Download(user *users.User, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	type response struct {
		Image []byte `json:"image"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	blobName := fmt.Sprintf("%s-%s", user.ID, p.Name)
	imageBytes, err := cfg.storage.DownloadObject(r.Context(), blobName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to download image")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, response{Image: imageBytes})
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

	blobName := fmt.Sprintf("%s-%s", user.ID, p.Name)
	err = cfg.storage.DeleteObject(r.Context(), blobName)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "failed to delete image")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "image deleted successfully"})
}
