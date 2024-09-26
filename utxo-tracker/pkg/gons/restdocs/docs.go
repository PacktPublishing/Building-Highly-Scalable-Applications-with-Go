package restdocs

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed *.gohtml
var templates embed.FS

//go:embed assets/*
var assets embed.FS

// UI represents the documentation page
type UI struct {
	params templateParams
	assets http.Handler
}

// templateParams is input to the .gohtml files
type templateParams struct {
	HTMLTitle      string
	OpenAPISpecURL string
	BaseURL        string
}

func New(options ...Option) *UI {
	ui := &UI{
		params: templateParams{
			HTMLTitle:      "API Docs",
			OpenAPISpecURL: "http://petstore.swagger.io/v2/swagger.json",
			BaseURL:        "",
		},
	}

	for _, opt := range options {
		opt(ui)
	}

	subAssets, _ := fs.Sub(assets, "assets")
	ui.assets = http.StripPrefix(
		ui.params.BaseURL+"/assets/",
		http.FileServer(http.FS(subAssets)),
	)

	return ui
}

// Option instances can be given to the New function.
type Option func(*UI)

// WithBaseURL sets the URL path where the document page are mounted
func WithBaseURL(url string) Option {
	return func(ui *UI) {
		ui.params.BaseURL = strings.TrimSuffix(url, "/")
	}
}

func WithSpecURL(url string) Option {
	return func(ui *UI) {
		ui.params.OpenAPISpecURL = url
	}
}

func WithHtmlTitle(title string) Option {
	return func(ui *UI) {
		ui.params.HTMLTitle = title
	}
}

// Handler creates and returns an HTTP handler that renders the service landing page
func (u *UI) Handler() http.HandlerFunc {
	tmpl := template.Must(
		template.ParseFS(templates, "*.gohtml"))

	return func(w http.ResponseWriter, r *http.Request) {
		apfx := u.params.BaseURL + "/assets/"
		if strings.HasPrefix(r.URL.Path, apfx) {
			u.assets.ServeHTTP(w, r)
			return
		}

		// Serve the documentation page
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.ExecuteTemplate(w,
			"redoc.gohtml", u.params)
		if err != nil {
			http.Error(w,
				"Failed to render template",
				http.StatusInternalServerError)
		}
	}

}
