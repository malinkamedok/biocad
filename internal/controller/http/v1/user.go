package v1

import (
	"biocad/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type userRoutes struct {
	u usecase.UserContract
}

func NewUserRoutes(routes chi.Router, u usecase.UserContract) {
	uc := &userRoutes{u: u}

	routes.Get("/get-guid-list", uc.getGUIDList)
	routes.Get("/{guid}/{page}/{limit}", uc.getDataByGUID)
}

func (u *userRoutes) getGUIDList(w http.ResponseWriter, r *http.Request) {
	list, err := u.u.GetAllUniqueGuidList(r.Context())
	if err != nil {
		w.WriteHeader(500)
		var message = "error in getting guid list"
		render.JSON(w, r, message)
		return
	}
	w.WriteHeader(200)
	render.JSON(w, r, list)
	return
}

func (u *userRoutes) getDataByGUID(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "guid")
	limit, err := strconv.ParseUint(chi.URLParam(r, "limit"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		var message = "incorrect limit number"
		render.JSON(w, r, message)
		return
	}
	page, err := strconv.ParseUint(chi.URLParam(r, "page"), 10, 64)
	if err != nil {
		w.WriteHeader(400)
		var message = "incorrect page number"
		render.JSON(w, r, message)
		return
	}
	data, err := u.u.GetAllDataByGuid(r.Context(), guid, limit, page)
	if err != nil {
		w.WriteHeader(500)
		var message = "error in getting data by guid"
		render.JSON(w, r, message)
		return
	}
	w.WriteHeader(200)
	render.JSON(w, r, data)
	return
}
