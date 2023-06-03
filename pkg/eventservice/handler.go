package eventservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "embed"

	"event/internal/utils"
	"event/pkg/eventservice/service"

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
		sendResponse(w, err)
		return
	}

	defer r.Body.Close()

	var incomingData []EventModel
	if err := json.Unmarshal(body, &incomingData); err != nil {
		logrus.Errorf("can't unmarshal body request: %v", err)
		sendResponse(w, err)
		return
	}

	ip, ok := utils.ExtractIpAddr(r)
	if !ok {
		logrus.Warnf("can't extract ip addr from response: %v", r.Header)
		// todo
	}

	data := make([]service.ServiceEventModel, 0, len(incomingData))
	for _, v := range incomingData {
		data = append(data, *v.toServiseDataEventModel(time.Now().UTC(), ip))
	}
	fmt.Println(data)

	err1 := h.Service.Insert(r.Context(), data)
	logrus.Infof("request processed, error: %v", err1)
	sendResponse(w, err1)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	err := h.Service.Ping(r.Context())
	sendResponse(w, err)
}

type HttpResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, err error) {
	var msg string
	var code int

	if err != nil {
		var syntax = &json.SyntaxError{}
		switch {
		case errors.As(err, &syntax):
			code = http.StatusBadRequest
		default:
			code = http.StatusInternalServerError
		}
		msg = err.Error()
	} else {
		msg = "ok"
		code = http.StatusOK
	}

	rsp := HttpResponse{
		Status:  http.StatusText(code),
		Message: msg,
	}

	body, _ := json.Marshal(rsp)

	w.WriteHeader(code)
	w.Write(body)
}
