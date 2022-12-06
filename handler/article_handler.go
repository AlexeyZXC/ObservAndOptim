package handler

import (
	"context"
	"elastic/m"
	"elastic/store"
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type ArticleHandler struct {
	S store.ArticleStore
}

func NewArticleHandler(s store.ArticleStore) ArticleHandler {
	return ArticleHandler{S: s}
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
func (h ArticleHandler) Err(r render.Render, err error) {
	r.JSON(http.StatusInternalServerError, err)
}
func (h ArticleHandler) BadRequest(r render.Render, err error) {
	r.JSON(http.StatusBadRequest, err)
}
