package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/csrf"

	"github.com/onatakduman/gatus-admin/internal/brand"
)

// daOf reads a DefaultAlert-pointer field off any provider struct pointer.
// Returns nil if the provider is nil or has no DefaultAlert field/value.
func daOf(provider any) any {
	if provider == nil {
		return nil
	}
	v := reflect.ValueOf(provider)
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	f := v.FieldByName("DefaultAlert")
	if !f.IsValid() {
		return nil
	}
	if f.Kind() == reflect.Pointer && f.IsNil() {
		return nil
	}
	return f.Interface()
}

//go:embed templates/*.html
var TemplatesFS embed.FS

//go:embed static
var StaticFS embed.FS

type Templates struct {
	dev bool
	t   *template.Template
}

func NewTemplates(dev bool) *Templates {
	return &Templates{dev: dev, t: parseAll()}
}

func parseAll() *template.Template {
	funcs := template.FuncMap{
		"csrfField": func() template.HTML { return "" },
		"dict": func(values ...any) map[string]any {
			if len(values)%2 != 0 {
				return nil
			}
			m := map[string]any{}
			for i := 0; i < len(values); i += 2 {
				key, _ := values[i].(string)
				m[key] = values[i+1]
			}
			return m
		},
		"deref": func(b *bool) bool {
			return b != nil && *b
		},
		"join": strings.Join,
		"list": func(items ...any) []any { return items },
		"daOf":          daOf,
		"daOfGeneric":   daOf,
		"docsURL":       docsURL,
		"docsLabel":     docsLabel,
		"docsContent":   docsContent,
		"presetNames":   func() []string { return brand.PresetNames() },
		"presetPalette": func(name string) any { p := brand.Presets[name]; return p },
		"firstHex": func(v string) string {
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "#") && (len(v) == 4 || len(v) == 7 || len(v) == 9) {
				return v[:7]
			}
			return ""
		},
		"tooltip": func(key string) template.HTML {
			txt := TooltipFor(key)
			if txt == "" {
				return ""
			}
			esc := template.HTMLEscapeString(txt)
			return template.HTML(`<span class="tt" title="` + esc + `" aria-label="` + esc + `"><i data-lucide="info"></i></span>`)
		},
	}
	tpl := template.New("").Funcs(funcs)
	tpl = template.Must(tpl.ParseFS(TemplatesFS, "templates/*.html"))
	return tpl
}

func (t *Templates) Render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) error {
	if t.dev {
		t.t = parseAll()
	}
	if data == nil {
		data = map[string]any{}
	}
	data["CSRFField"] = csrf.TemplateField(r)
	data["CSRFToken"] = csrf.Token(r)
	active := activeFromPath(r.URL.Path)
	data["Active"] = active
	if _, ok := data["PageTitle"]; !ok {
		data["PageTitle"] = titleFor(active)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.t.ExecuteTemplate(w, name, data)
}

func activeFromPath(p string) string {
	trim := strings.TrimPrefix(p, "/admin")
	trim = strings.Trim(trim, "/")
	if trim == "" {
		return "endpoints"
	}
	parts := strings.SplitN(trim, "/", 2)
	return parts[0]
}

func titleFor(active string) string {
	switch active {
	case "endpoints":
		return "Endpoints"
	case "alerting":
		return "Alerting"
	case "ui":
		return "UI / Public Page"
	case "brand":
		return "Admin Brand"
	case "announcements":
		return "Announcements"
	case "maintenance":
		return "Maintenance"
	case "core":
		return "Core"
	case "advanced":
		return "Advanced YAML"
	case "security":
		return "Security"
	case "tunneling":
		return "Tunneling"
	case "external-endpoints":
		return "External Endpoints"
	case "suites":
		return "Suites"
	case "login":
		return "Login"
	}
	return "Admin"
}

func StaticHandler() http.HandlerFunc {
	sub, _ := fs.Sub(StaticFS, "static")
	fileServer := http.FileServer(http.FS(sub))
	return func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, "/admin/static")
		fileServer.ServeHTTP(w, r2)
	}
}
