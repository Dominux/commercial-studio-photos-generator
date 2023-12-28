package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const (
	CheckProgressEndpoint = "/check_progress/"
	GetResultEndpoint     = "/get_result/"
)

var (
	pollingTempl = readTemplate("./static/polling.html")
	imageTempl   = readTemplate("./static/image.html")
)

type Handler struct {
	p                *Producer
	redisClient      *redis.Client
	minioClient      *minio.Client
	minioBucket      string
	minioResultsPath string
}

func NewHandler(p *Producer, redisClient *redis.Client, minioClient *minio.Client, minioBucket string, minioResultsPath string) *Handler {
	return &Handler{p, redisClient, minioClient, minioBucket, minioResultsPath}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		raise404(w)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w)
}

func (h *Handler) RunGenerating(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		raise404(w)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	product := r.FormValue("product")
	id, err := RunGenerating(product, h.p)
	if err != nil {
		log.Println(err)
		raise500(w)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	html := fmt.Sprintf(pollingTempl, id)
	fmt.Fprint(w, html)
}

func (h *Handler) CheckProgressHandler(w http.ResponseWriter, r *http.Request) {
	rawId := getIdFromPath(w, r, CheckProgressEndpoint)
	if rawId == "" {
		return
	}

	// trying to get id from redis
	if err := CheckProgress(r.Context(), rawId, h.redisClient); err != nil {
		raise404(w)
		return
	}

	html := fmt.Sprintf(imageTempl, rawId)
	fmt.Fprint(w, html)
}

func (h *Handler) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	rawId := getIdFromPath(w, r, GetResultEndpoint)
	if rawId == "" {
		return
	}

	// getting result
	result, err := GetResult(r.Context(), rawId, h.minioClient, h.minioBucket, h.minioResultsPath)
	if err != nil {
		log.Println(err)
		raise404(w)
		return
	}

	w.Header().Add("Content-Type", "image/png")
	w.Write(result)
}

func raise404(w http.ResponseWriter) {
	http.Error(w, "404 not Found", http.StatusNotFound)
}

func raise500(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func getIdFromPath(w http.ResponseWriter, r *http.Request, rel string) string {
	rawId := strings.TrimPrefix(r.URL.Path, rel)
	_, err := uuid.Parse(rawId)
	if err != nil {
		http.Error(w, "Path must be a valid UUID", http.StatusBadRequest)
		return ""
	}
	return rawId
}
