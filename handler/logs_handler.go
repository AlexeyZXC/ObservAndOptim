package handler

import (
	"context"
	"elastic/m"
	"elastic/store"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	render_chi "github.com/go-chi/render"
	"github.com/martini-contrib/render"
)

type LogsHandler struct {
	S store.ArticleStore
}

func NewLogsHandler(s store.ArticleStore) ArticleHandler {
	return ArticleHandler{S: s}
}

func (h LogsHandler) Id_chi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := context.Background()
	article, err := h.S.Get(ctx, id)
	if err != nil {
		render_chi.Status(r, http.StatusInternalServerError)
		render_chi.JSON(w, r, err)
		return
	}
	render_chi.Status(r, http.StatusOK)
	render_chi.JSON(w, r, article)
}

func (h LogsHandler) Add_chi(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	var article m.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		render_chi.Status(r, http.StatusInternalServerError)
		render_chi.JSON(w, r, err)
		return
	}
	err = h.S.Add(ctx, article)
	if err != nil {
		render_chi.Status(r, http.StatusInternalServerError)
		render_chi.JSON(w, r, err)
		return
	}
	render_chi.Status(r, http.StatusOK)
	render_chi.JSON(w, r, article)
}

//declared in article_handler.go
// type SearchRequest struct {
// 	Query string `json:"query"`
// }

func (h LogsHandler) Search_chi(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer r.Body.Close()
	var query SearchRequest
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		render_chi.Status(r, http.StatusBadRequest)
		render_chi.JSON(w, r, err)
		return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		render_chi.Status(r, http.StatusInternalServerError)
		render_chi.JSON(w, r, err)
		return
	}
	render_chi.Status(r, http.StatusOK)
	render_chi.JSON(w, r, articles)
}

func (h LogsHandler) Err(r render.Render, err error) {
	r.JSON(http.StatusInternalServerError, err)
}
func (h LogsHandler) BadRequest(r render.Render, err error) {
	r.JSON(http.StatusBadRequest, err)
}
