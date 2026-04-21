package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/onatakduman/gatus-admin/internal/brand"
	"github.com/onatakduman/gatus-admin/internal/configio"
)

var allowedLogoTypes = map[string]string{
	"image/png":                "png",
	"image/jpeg":               "jpg",
	"image/svg+xml":            "svg",
	"image/webp":               "webp",
	"image/x-icon":             "ico",
	"image/vnd.microsoft.icon": "ico",
}

var allowedLogoSlots = map[string]bool{
	"logo_full":   true,
	"logo_square": true,
	"favicon":     true,
}

func (h *Handlers) BrandGet(w http.ResponseWriter, r *http.Request) {
	b, _ := h.BrandStore.Load()
	effective := h.computeBrand(b)
	h.render(w, r, "brand.html", map[string]any{
		"BrandRaw": b,
		"Brand":    effective,
	})
}

func (h *Handlers) BrandPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	current, _ := h.BrandStore.Load()
	b := brand.AdminBrand{
		Name:             strings.TrimSpace(r.FormValue("name")),
		Tagline:          strings.TrimSpace(r.FormValue("tagline")),
		LogoFull:         current.LogoFull,
		LogoSquare:       current.LogoSquare,
		Favicon:          current.Favicon,
		InheritFromGatus: r.FormValue("inherit_from_gatus") == "on",
		Colors: brand.Colors{
			Accent:      sanitizeColor(r.FormValue("color_accent")),
			AccentHover: sanitizeColor(r.FormValue("color_accent_hover")),
			Background:  sanitizeColor(r.FormValue("color_background")),
			Surface:     sanitizeColor(r.FormValue("color_surface")),
			Border:      sanitizeColor(r.FormValue("color_border")),
			Text:        sanitizeColor(r.FormValue("color_text")),
			TextMuted:   sanitizeColor(r.FormValue("color_text_muted")),
		},
		Fonts: brand.Fonts{
			Heading: strings.TrimSpace(r.FormValue("font_heading")),
			Body:    strings.TrimSpace(r.FormValue("font_body")),
			Sans:    strings.TrimSpace(r.FormValue("font_sans")),
			Serif:   strings.TrimSpace(r.FormValue("font_serif")),
			Mono:    strings.TrimSpace(r.FormValue("font_mono")),
		},
		Light:           parsePalette(r, "light"),
		Dark:            parsePalette(r, "dark"),
		ColorMode:       strings.TrimSpace(r.FormValue("color_mode")),
		CustomCSS:       r.FormValue("custom_css"),
		MirrorPublicCSS: r.FormValue("mirror_public_css") == "on",
		Preset:          strings.TrimSpace(r.FormValue("preset")),
		Radius:          strings.TrimSpace(r.FormValue("radius")),
		Density:         strings.TrimSpace(r.FormValue("density")),
		ActiveStyle:     strings.TrimSpace(r.FormValue("active_style")),
		HeadingCase:     strings.TrimSpace(r.FormValue("heading_case")),
		HeadingTrack:    strings.TrimSpace(r.FormValue("heading_track")),
		LogoSize:        strings.TrimSpace(r.FormValue("logo_size")),
	}

	// "Apply preset" button: when its name was submitted, overwrite the
	// per-field colour values with the preset bundle. The colour pickers
	// in the form already show the preset values after the redirect.
	if applyPreset := strings.TrimSpace(r.FormValue("apply_preset")); applyPreset != "" {
		if p, ok := brand.Presets[applyPreset]; ok {
			b.Light = p.Light
			b.Dark = p.Dark
			b.Preset = applyPreset
		}
	}
	if err := h.BrandStore.Save(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, "brand saved")
	http.Redirect(w, r, "/admin/brand", http.StatusFound)
}

func (h *Handlers) BrandUpload(w http.ResponseWriter, r *http.Request) {
	slot := chi.URLParam(r, "slot")
	if !allowedLogoSlots[slot] {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseMultipartForm(brand.MaxAsset + 1024); err != nil {
		h.flashErr(r, w, "upload too large or invalid: "+err.Error())
		http.Redirect(w, r, "/admin/brand", http.StatusFound)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		h.flashErr(r, w, "no file: "+err.Error())
		http.Redirect(w, r, "/admin/brand", http.StatusFound)
		return
	}
	defer file.Close()
	if header.Size > brand.MaxAsset {
		h.flashErr(r, w, fmt.Sprintf("file too large (max %d bytes)", brand.MaxAsset))
		http.Redirect(w, r, "/admin/brand", http.StatusFound)
		return
	}
	body, err := io.ReadAll(io.LimitReader(file, brand.MaxAsset+1))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(body) > brand.MaxAsset {
		h.flashErr(r, w, "file too large")
		http.Redirect(w, r, "/admin/brand", http.StatusFound)
		return
	}

	ctype := http.DetectContentType(body)
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(header.Filename), "."))
	if ext == "svg" {
		ctype = "image/svg+xml"
	}
	resolved, ok := allowedLogoTypes[ctype]
	if !ok {
		for _, v := range allowedLogoTypes {
			if v == ext {
				resolved = ext
				ok = true
				break
			}
		}
	}
	if !ok {
		h.flashErr(r, w, "unsupported file type (png, jpg, svg, webp, ico only); got "+ctype)
		http.Redirect(w, r, "/admin/brand", http.StatusFound)
		return
	}

	name, err := h.BrandStore.SaveAsset(slot, resolved, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := h.BrandStore.Load()
	switch slot {
	case "logo_full":
		b.LogoFull = name
	case "logo_square":
		b.LogoSquare = name
	case "favicon":
		b.Favicon = name
	}
	if err := h.BrandStore.Save(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, slot+" uploaded")
	http.Redirect(w, r, "/admin/brand", http.StatusFound)
}

func (h *Handlers) BrandReset(w http.ResponseWriter, r *http.Request) {
	slot := chi.URLParam(r, "slot")
	if !allowedLogoSlots[slot] {
		http.NotFound(w, r)
		return
	}
	if err := h.BrandStore.DeleteAsset(slot); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := h.BrandStore.Load()
	switch slot {
	case "logo_full":
		b.LogoFull = ""
	case "logo_square":
		b.LogoSquare = ""
	case "favicon":
		b.Favicon = ""
	}
	if err := h.BrandStore.Save(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.flashOk(r, w, slot+" removed")
	http.Redirect(w, r, "/admin/brand", http.StatusFound)
}

func (h *Handlers) AssetGet(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	p := h.BrandStore.AssetFilePath(name)
	if p == "" {
		http.NotFound(w, r)
		return
	}
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(p)), ".")
	if ct := mime.TypeByExtension("." + ext); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
	w.Header().Set("Content-Disposition", "inline")
	w.Header().Set("Cache-Control", "public, max-age=300")
	http.ServeFile(w, r, p)
}

func (h *Handlers) computeBrand(b brand.AdminBrand) brand.Effective {
	var header, description string
	if b.InheritFromGatus {
		uiDoc := &configio.UIFileDoc{}
		if err := configio.Load(h.path(configio.UIFile), uiDoc); err == nil {
			header = uiDoc.UI.Header
			if header == "" {
				header = uiDoc.UI.Title
			}
			description = uiDoc.UI.Description
			if description == "" {
				description = uiDoc.UI.DashboardSubheading
			}
		}
	}
	return brand.Compute(b, header, description)
}

// parsePalette extracts a ThemePalette for "light" or "dark" prefixed fields.
func parsePalette(r *http.Request, mode string) brand.ThemePalette {
	g := func(name string) string { return sanitizeColor(r.FormValue(mode + "_" + name)) }
	return brand.ThemePalette{
		Background:               g("background"),
		Foreground:               g("foreground"),
		Card:                     g("card"),
		CardForeground:           g("card_foreground"),
		Popover:                  g("popover"),
		PopoverForeground:        g("popover_foreground"),
		Primary:                  g("primary"),
		PrimaryForeground:        g("primary_foreground"),
		Secondary:                g("secondary"),
		SecondaryForeground:      g("secondary_foreground"),
		Muted:                    g("muted"),
		MutedForeground:          g("muted_foreground"),
		Accent:                   g("accent"),
		AccentForeground:         g("accent_foreground"),
		Destructive:              g("destructive"),
		DestructiveForeground:    g("destructive_foreground"),
		Border:                   g("border"),
		Input:                    g("input"),
		Ring:                     g("ring"),
		Sidebar:                  g("sidebar"),
		SidebarForeground:        g("sidebar_foreground"),
		SidebarPrimary:           g("sidebar_primary"),
		SidebarPrimaryForeground: g("sidebar_primary_foreground"),
		SidebarAccent:            g("sidebar_accent"),
		SidebarAccentForeground:  g("sidebar_accent_foreground"),
		SidebarBorder:            g("sidebar_border"),
		SidebarRing:              g("sidebar_ring"),
	}
}

func sanitizeColor(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if len(v) > 64 {
		return ""
	}
	first := v[0]
	if first == '#' {
		for _, c := range v[1:] {
			if !isHex(c) {
				return ""
			}
		}
		return v
	}
	lower := strings.ToLower(v)
	if strings.HasPrefix(lower, "rgb(") || strings.HasPrefix(lower, "rgba(") ||
		strings.HasPrefix(lower, "hsl(") || strings.HasPrefix(lower, "hsla(") ||
		strings.HasPrefix(lower, "oklch(") || strings.HasPrefix(lower, "lch(") ||
		strings.HasPrefix(lower, "lab(") || strings.HasPrefix(lower, "oklab(") {
		if !strings.HasSuffix(v, ")") {
			return ""
		}
		for _, c := range v {
			if !(isAlnum(c) || c == '(' || c == ')' || c == ',' || c == '.' || c == ' ' || c == '%' || c == '-' || c == '/') {
				return ""
			}
		}
		return v
	}
	for _, c := range v {
		if !(c >= 'a' && c <= 'z') && c != '-' {
			return ""
		}
	}
	return v
}

func isHex(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
func isAlnum(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
