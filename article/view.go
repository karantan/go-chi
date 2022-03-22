package article

import (
	"context"
	"errors"
	"gochi/logger"
	"gochi/user"
	"gochi/utils"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var log = logger.New("article")
var ErrNotFound = &utils.ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func ListArticles(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewArticleListResponse(articles)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *Article
		var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = DBGetArticle(articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			article, err = DBGetArticleBySlug(articleSlug)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			log.Error(err)
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SearchArticles searches the Articles data for a matching article.
// It's just a stub, but you get the idea.
func SearchArticles(w http.ResponseWriter, r *http.Request) {
	render.RenderList(w, r, NewArticleListResponse(articles))
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	data := &ArticlePayload{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	article := data.Article
	DBNewArticle(article)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewArticleResponse(article))
}

// GetArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value("article").(*Article)

	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, utils.ErrRender(err))
		return
	}
}

// UpdateArticle updates an existing Article in our persistent store.
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	data := &ArticlePayload{Article: article}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}
	article = data.Article
	DBUpdateArticle(article.ID, article)

	render.Render(w, r, NewArticleResponse(article))
}

// DeleteArticle removes an existing Article from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value("article").(*Article)

	article, err = DBRemoveArticle(article.ID)
	if err != nil {
		render.Render(w, r, utils.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewArticleResponse(article))
}

// ArticlePayload is the request payload for Article data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type ArticlePayload struct {
	*Article

	User *user.UserPayload `json:"user,omitempty"`

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *ArticlePayload) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Article == nil {
		return errors.New("missing required Article fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	a.ProtectedID = ""                                 // unset the protected ID
	a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

func NewArticleResponse(article *Article) *ArticlePayload {
	resp := &ArticlePayload{Article: article}

	if resp.User == nil {
		if u, _ := user.DBGetUser(resp.UserID); u != nil {
			resp.User = user.NewUserPayloadResponse(u)
		}
	}

	return resp
}

func (rd *ArticlePayload) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewArticleListResponse(articles []*Article) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}
	return list
}
