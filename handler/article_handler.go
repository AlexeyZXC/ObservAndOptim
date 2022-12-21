package handler

import (
	"context"
	"elastic/m"
	"elastic/store"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	render_chi "github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type ArticleHandler struct {
	S      store.ArticleStore
	Tracer opentracing.Tracer
}

func NewArticleHandler(s store.ArticleStore, t opentracing.Tracer) ArticleHandler {
	return ArticleHandler{S: s, Tracer: t}
}
func (h ArticleHandler) Id(r render.Render, params martini.Params) (interface{}, error) {
	id := params["id"]
	ctx := context.Background()
	article, err := h.S.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	r.JSON(http.StatusOK, article)
	return h.S.Get(ctx, id)
}

func (h ArticleHandler) Id_chi(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(r.Context(), h.Tracer, "Id path called")
	defer span.Finish()

	id := chi.URLParam(r, "id")
	article, err := h.S.Get(ctx, id)
	if err != nil {
		span.LogFields(log.Error(err), log.String("id", id))
		render_chi.Status(r, http.StatusInternalServerError)
		render_chi.JSON(w, r, err)
		return
	}
	render_chi.Status(r, http.StatusOK)
	render_chi.JSON(w, r, article)
}

func (h ArticleHandler) Add(r render.Render, req *http.Request) {
	ctx := context.Background()
	defer req.Body.Close()
	var article m.Article
	err := json.NewDecoder(req.Body).Decode(&article)
	if err != nil {
		h.Err(r, err)
		return
	}
	err = h.S.Add(ctx, article)
	if err != nil {
		h.Err(r, err)
		return
	}
	r.JSON(http.StatusOK, article)
}

func (h ArticleHandler) Add_chi(w http.ResponseWriter, r *http.Request) {
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

type SearchRequest struct {
	Query string `json:"query"`
}

func (h ArticleHandler) Search(r render.Render, req *http.Request) {
	ctx := context.Background()
	defer req.Body.Close()
	var query SearchRequest
	err := json.NewDecoder(req.Body).Decode(&query)
	if err != nil {
		h.BadRequest(r, err)
		return
	}
	articles, err := h.S.Search(ctx, query.Query)
	if err != nil {
		h.Err(r, err)
		return
	}
	r.JSON(http.StatusOK, articles)
}

func (h ArticleHandler) Search_chi(w http.ResponseWriter, r *http.Request) {
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

func (h ArticleHandler) Err(r render.Render, err error) {
	r.JSON(http.StatusInternalServerError, err)
}
func (h ArticleHandler) BadRequest(r render.Render, err error) {
	r.JSON(http.StatusBadRequest, err)
}
