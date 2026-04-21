package brand

// AdminBrand controls the admin panel's own visual identity.
// Separate from Gatus UI branding (ui.yaml). When Name/Tagline are empty and
// InheritFromGatus is true, effective values fall back to gatus ui.header/description.
type AdminBrand struct {
	Name             string `yaml:"name,omitempty"`
	Tagline          string `yaml:"tagline,omitempty"`
	LogoFull         string `yaml:"logo_full,omitempty"`
	LogoSquare       string `yaml:"logo_square,omitempty"`
	Favicon          string `yaml:"favicon,omitempty"`
	InheritFromGatus bool   `yaml:"inherit_from_gatus"`

	// Legacy short-hand colour fields (older admin.yaml files).
	Colors Colors `yaml:"colors,omitempty"`
	Fonts  Fonts  `yaml:"fonts,omitempty"`

	// Full shadcn-style theme tokens — light + dark modes.
	Light ThemePalette `yaml:"light,omitempty"`
	Dark  ThemePalette `yaml:"dark,omitempty"`

	// Default colour mode the admin opens in.
	ColorMode string `yaml:"color_mode,omitempty"` // light | dark | system (default: dark)

	// Layout & shape — no-code controls in the Brand panel.
	Preset       string `yaml:"preset,omitempty"`        // empty | doon | linear | vercel | rounded | shadcn-neutral | shadcn-slate | shadcn-zinc | shadcn-red | shadcn-blue | shadcn-violet
	Radius       string `yaml:"radius,omitempty"`        // 0 | 2 | 4 | 6 | 8 | 10 | 12 | 16 (px) OR rem like 0.5rem
	Density      string `yaml:"density,omitempty"`       // compact | normal | comfortable
	ActiveStyle  string `yaml:"active_style,omitempty"`  // stripe | pill | glow | bg
	HeadingCase  string `yaml:"heading_case,omitempty"`  // normal | uppercase
	HeadingTrack string `yaml:"heading_track,omitempty"` // tight | normal | wide
	LogoSize     string `yaml:"logo_size,omitempty"`     // small | medium | large

	// CustomCSS targets the admin panel only — escape hatch for the few
	// tweaks the form-based controls don't expose.
	CustomCSS string `yaml:"custom_css,omitempty"`
	// MirrorPublicCSS — when true, the Gatus public-page ui.custom-css is
	// also injected into the admin head so both surfaces share styling.
	MirrorPublicCSS bool `yaml:"mirror_public_css,omitempty"`
}

type Colors struct {
	// Legacy short-hand fields kept for back-compat with older admin.yaml.
	Accent      string `yaml:"accent,omitempty"`
	AccentHover string `yaml:"accent_hover,omitempty"`
	Background  string `yaml:"background,omitempty"`
	Surface     string `yaml:"surface,omitempty"`
	Border      string `yaml:"border,omitempty"`
	Text        string `yaml:"text,omitempty"`
	TextMuted   string `yaml:"text_muted,omitempty"`
}

type Fonts struct {
	Heading string `yaml:"heading,omitempty"`
	Body    string `yaml:"body,omitempty"`
	Sans    string `yaml:"sans,omitempty"`
	Serif   string `yaml:"serif,omitempty"`
	Mono    string `yaml:"mono,omitempty"`
}

// ThemePalette holds the full shadcn-style token set for a single colour mode.
// All values are CSS colour strings (#hex, rgb(), oklch(), etc.).
type ThemePalette struct {
	Background           string `yaml:"background,omitempty"`
	Foreground           string `yaml:"foreground,omitempty"`
	Card                 string `yaml:"card,omitempty"`
	CardForeground       string `yaml:"card_foreground,omitempty"`
	Popover              string `yaml:"popover,omitempty"`
	PopoverForeground    string `yaml:"popover_foreground,omitempty"`
	Primary              string `yaml:"primary,omitempty"`
	PrimaryForeground    string `yaml:"primary_foreground,omitempty"`
	Secondary            string `yaml:"secondary,omitempty"`
	SecondaryForeground  string `yaml:"secondary_foreground,omitempty"`
	Muted                string `yaml:"muted,omitempty"`
	MutedForeground      string `yaml:"muted_foreground,omitempty"`
	Accent               string `yaml:"accent,omitempty"`
	AccentForeground     string `yaml:"accent_foreground,omitempty"`
	Destructive          string `yaml:"destructive,omitempty"`
	DestructiveForeground string `yaml:"destructive_foreground,omitempty"`
	Border               string `yaml:"border,omitempty"`
	Input                string `yaml:"input,omitempty"`
	Ring                 string `yaml:"ring,omitempty"`

	// Sidebar — separate slot so the sidebar can have a distinct vibe.
	Sidebar                  string `yaml:"sidebar,omitempty"`
	SidebarForeground        string `yaml:"sidebar_foreground,omitempty"`
	SidebarPrimary           string `yaml:"sidebar_primary,omitempty"`
	SidebarPrimaryForeground string `yaml:"sidebar_primary_foreground,omitempty"`
	SidebarAccent            string `yaml:"sidebar_accent,omitempty"`
	SidebarAccentForeground  string `yaml:"sidebar_accent_foreground,omitempty"`
	SidebarBorder            string `yaml:"sidebar_border,omitempty"`
	SidebarRing              string `yaml:"sidebar_ring,omitempty"`
}

// Defaults returns the out-of-the-box brand. Colors are intentionally empty so
// the admin uses the same shadcn dark theme as Gatus by default — only when the
// operator picks colours via the Brand panel does the inline override kick in.
func Defaults() AdminBrand {
	return AdminBrand{
		Name:             "gatus-admin",
		Tagline:          "Status page control panel",
		InheritFromGatus: true,
	}
}

// Effective merges saved brand with Gatus UI values and sane defaults so
// templates never see a blank field where something must render.
type Effective struct {
	Name            string
	Tagline         string
	LogoFull        string
	LogoSquare      string
	Favicon         string
	Colors          Colors
	Fonts           Fonts
	HasLogos        bool
	CustomCSS       string
	MirrorPublicCSS bool
	Light           ThemePalette
	Dark            ThemePalette
	ColorMode       string
	Preset          string
	Radius          string
	Density         string
	ActiveStyle     string
	HeadingCase     string
	HeadingTrack    string
	LogoSize        string
}
