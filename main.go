package main

import (
	"flag"
	"fmt"
	"gochi/article"
	"gochi/auth"
	"gochi/database"
	"gochi/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/oauth"
)

//
// REST
// ====
// This example demonstrates a HTTP REST web service with some fixture data.
// Follow along the example and patterns.
//
// Also check routes.json for the generated docs from passing the -routes flag,
// to run yourself do: `go run . -routes`
//
// Boot the server:
// ----------------
// $ go run main.go
//
// Client requests:
// ----------------
// $ curl http://localhost:3333/
// root.
//
// $ curl http://localhost:3333/articles
// [{"id":"1","title":"Hi"},{"id":"2","title":"sup"}]
//
// $ curl http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl -X DELETE http://localhost:3333/articles/1
// {"id":"1","title":"Hi"}
//
// $ curl http://localhost:3333/articles/1
// "Not Found"
//
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles/97
// {"id":"97","title":"awesomeness"}
//
// $ curl http://localhost:3333/articles
// [{"id":"2","title":"sup"},{"id":"97","title":"awesomeness"}]
//

const DB = "boltdb.db"
const OAUTH_KEY = "secret"

var log = logger.New("main")
var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()
	s := CreateNewServer(DB)
	s.MountHandlers()

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		fmt.Println(docgen.MarkdownRoutesDoc(s.Router, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi/v5",
			Intro:       "Welcome to the gochi generated docs.",
		}))
		return
	}
	addr := "0.0.0.0:3000"
	log.Infof("Starting server gochi on %s", addr)
	http.ListenAndServe(addr, s.Router)
	http.ListenAndServe(":3000", s.Router)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

type Server struct {
	Router *chi.Mux
	DB     *database.Database
}

func CreateNewServer(dbPath string) *Server {
	db, err := database.GetDatabase(dbPath, false)
	if err != nil {
		log.Error(err)
	}
	s := &Server{DB: db}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() {
	// Mount all Middleware here
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	// Mount all handlers here
	s.Router.Get("/", Hello)

	// auth
	bearerServer := oauth.NewBearerServer(OAUTH_KEY, time.Second*120, &auth.TestUserVerifier{}, nil)
	s.Router.Post("/token", bearerServer.UserCredentials)
	s.Router.Post("/auth", bearerServer.ClientCredentials)

	// oauth protected API routes
	s.Router.Route("/api/v1", func(r chi.Router) {
		// use the Bearer Authentication middleware
		r.Use(oauth.Authorize(OAUTH_KEY, nil))
		// RESTy routes for "articles" resource
		r.Route("/articles", func(r chi.Router) {
			r.Get("/", article.ListArticles)
			r.Post("/", article.CreateArticle)       // POST /articles
			r.Get("/search", article.SearchArticles) // GET /articles/search

			r.Route("/{articleID}", func(r chi.Router) {
				r.Use(article.ArticleCtx)            // Load the *Article on the request context
				r.Get("/", article.GetArticle)       // GET /articles/123
				r.Put("/", article.UpdateArticle)    // PUT /articles/123
				r.Delete("/", article.DeleteArticle) // DELETE /articles/123
			})

			// GET /articles/whats-up
			r.With(article.ArticleCtx).Get("/{articleSlug:[a-z-]+}", article.GetArticle)
		})
	})

}
