package eventservice

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "embed"

	"github.com/MrSwartz/event/internal/utils"
	"github.com/MrSwartz/event/pkg/eventservice/service"

	"github.com/flowchartsman/swaggerui"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service *service.Service
}

func (h *Handler) InitRoutes(exposeSwagger bool) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/events", h.storeEvents).Methods(http.MethodPost)
	router.HandleFunc("/v1/health", h.health).Methods(http.MethodGet)

	if exposeSwagger && len(swaggerSpec) > 0 {
		router.Path("/swagger").Handler(http.RedirectHandler("/swagger/", http.StatusPermanentRedirect))
		router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", swaggerui.Handler(swaggerSpec)))
	}
	return router
}

//go:embed swagger.json
var swaggerSpec []byte

func (h *Handler) storeEvents(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("can't read body from request: %v", err)
		sendResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var incomingData []EventModel
	if err := json.Unmarshal(body, &incomingData); err != nil {
		logrus.Errorf("can't unmarshal body request: %v", err)
		sendResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ip, ok := utils.ExtractIpAddr(r)
	if !ok {
		// тут ip должен вытягивается с учётом балансировщика нагрузки
		// странный баг, если не удалось вытащить ip из запроса,
		// вероятно хэдеры были модифицированы
		// если проблема не в хэдерах, то что то совсем не то
		// и сервис не должен ничего сохранять

		// чтобы не хардкодить ip, можно раскоментить 2 строки ниже
		// logrus.Warnf("can't extract ip addr from response: %v", r.Header)
		// sendResponse(w, "can't extract ip addr", http.StatusInternalServerError)

		ip = "8.8.8.8"
	}

	data := make([]service.ServiceEventModel, 0, len(incomingData))
	for _, v := range incomingData {
		data = append(data, *v.toServiseDataEventModel(time.Now().UTC(), ip))
	}

	if err := h.Service.Insert(r.Context(), data); err != nil {
		logrus.Infof("request processed with error: %v", err)
		sendResponse(w, err.Error(), http.StatusInternalServerError)
	}

	sendResponse(w, "ok", http.StatusOK)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	err := h.Service.Ping(r.Context())

	var msg string
	var code int

	switch err {
	case nil:
		msg = "ok"
		code = http.StatusOK
	default:
		msg = err.Error()
		code = http.StatusInternalServerError
	}

	sendResponse(w, msg, code)
}

type HttpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, msg string, code int) {
	rsp := HttpResponse{
		Status:  http.StatusText(code),
		Message: msg,
	}

	body, _ := json.Marshal(rsp)

	w.WriteHeader(code)
	//nolint:errcheck
	w.Write(body)
}
