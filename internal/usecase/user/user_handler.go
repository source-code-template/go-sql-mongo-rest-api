package user

import (
	"context"
	"net/http"
	"reflect"

	"github.com/core-go/search"
	sv "github.com/core-go/service"
)

type UserHandler interface {
	Search(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service UserService, status sv.StatusConfig, validate func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error), logError func(context.Context, string)) UserHandler {
	searchModelType := reflect.TypeOf(UserFilter{})
	modelType := reflect.TypeOf(User{})
	keys, indexes, _ := sv.BuildHandlerParams(modelType)
	searchHandler := search.NewJSONSearchHandler(find, modelType, searchModelType, logError, nil)
	return &userHandler{SearchHandler: searchHandler, service: service, keys: keys, indexes: indexes, modelType: modelType, Status: status, Validate: validate, Error: logError}
}

type userHandler struct {
	*search.SearchHandler
	service   UserService
	keys      []string
	indexes   map[string]int
	modelType reflect.Type
	Status    sv.StatusConfig
	Validate  func(ctx context.Context, model interface{}) ([]sv.ErrorMessage, error)
	Error     func(context.Context, string)
}

func (h *userHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Load(r.Context(), id)
		sv.RespondModel(w, r, result, err, h.Error, nil)
	}
}
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, nil) {
			result, er3 := h.service.Create(r.Context(), &user)
			sv.AfterCreated(w, r, &user, result, er3, h.Status, h.Error, nil)
		}
	}
}
func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.DecodeAndCheckId(w, r, &user, h.keys, h.indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, nil) {
			result, er3 := h.service.Update(r.Context(), &user)
			sv.HandleResult(w, r, &user, result, er3, h.Status, h.Error, nil)
		}
	}
}
func (h *userHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var user User
	r, json, er1 := sv.BuildMapAndCheckId(w, r, &user, h.keys, h.indexes)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, nil) {
			result, er3 := h.service.Patch(r.Context(), json)
			sv.HandleResult(w, r, json, result, er3, h.Status, h.Error, nil)
		}
	}
}
func (h *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		result, err := h.service.Delete(r.Context(), id)
		sv.HandleDelete(w, r, result, err, h.Error, nil)
	}
}
