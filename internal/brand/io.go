package brand

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	FileName  = "admin.yaml"
	AssetsDir = "assets"
	BrandDir  = "brand"
	MaxAsset  = 2 * 1024 * 1024
)

type Store struct {
	RootDir string
}

func (s Store) brandDir() string   { return filepath.Join(s.RootDir, BrandDir) }
func (s Store) filePath() string   { return filepath.Join(s.brandDir(), FileName) }
func (s Store) AssetsPath() string { return filepath.Join(s.brandDir(), AssetsDir) }

func (s Store) Load() (AdminBrand, error) {
	b := Defaults()
	data, err := os.ReadFile(s.filePath())
	if err != nil {
		if os.IsNotExist(err) {
			return b, nil
		}
		return b, err
	}
	if len(data) == 0 {
		return b, nil
	}
	if err := yaml.Unmarshal(data, &b); err != nil {
		return b, err
	}
	return b, nil
}

func (s Store) Save(b AdminBrand) error {
	if err := os.MkdirAll(s.brandDir(), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(b)
	if err != nil {
		return err
	}
	tmp := s.filePath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.filePath())
}

func (s Store) SaveAsset(slot, ext string, body []byte) (string, error) {
	if err := os.MkdirAll(s.AssetsPath(), 0755); err != nil {
		return "", err
	}
	matches, _ := filepath.Glob(filepath.Join(s.AssetsPath(), slot+".*"))
	for _, m := range matches {
		_ = os.Remove(m)
	}
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	if ext == "" {
		ext = "bin"
	}
	name := slot + "." + ext
	dst := filepath.Join(s.AssetsPath(), name)
	if err := os.WriteFile(dst+".tmp", body, 0644); err != nil {
		return "", err
	}
	if err := os.Rename(dst+".tmp", dst); err != nil {
		return "", err
	}
	return name, nil
}

func (s Store) DeleteAsset(slot string) error {
	matches, _ := filepath.Glob(filepath.Join(s.AssetsPath(), slot+".*"))
	for _, m := range matches {
		if err := os.Remove(m); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (s Store) AssetFilePath(name string) string {
	base := filepath.Base(name)
	if base != name || base == "." || base == ".." {
		return ""
	}
	p := filepath.Join(s.AssetsPath(), base)
	if info, err := os.Stat(p); err != nil || info.IsDir() {
		return ""
	}
	return p
}

func Compute(b AdminBrand, gatusHeader, gatusDescription string) Effective {
	name := b.Name
	tagline := b.Tagline
	if b.InheritFromGatus {
		if name == "" {
			name = gatusHeader
		}
		if tagline == "" {
			tagline = gatusDescription
		}
	}
	if name == "" {
		name = "gatus-admin"
	}
	logoFull := assetURL(b.LogoFull)
	logoSquare := assetURL(b.LogoSquare)
	favicon := assetURL(b.Favicon)
	return Effective{
		Name:            name,
		Tagline:         tagline,
		LogoFull:        logoFull,
		LogoSquare:      logoSquare,
		Favicon:         favicon,
		Colors:          b.Colors,
		Fonts:           b.Fonts,
		HasLogos:        logoFull != "" || logoSquare != "",
		CustomCSS:       b.CustomCSS,
		MirrorPublicCSS: b.MirrorPublicCSS,
		Light:           b.Light,
		Dark:            b.Dark,
		ColorMode:       b.ColorMode,
		Preset:          b.Preset,
		Radius:          b.Radius,
		Density:         b.Density,
		ActiveStyle:     b.ActiveStyle,
		HeadingCase:     b.HeadingCase,
		HeadingTrack:    b.HeadingTrack,
		LogoSize:        b.LogoSize,
	}
}

func assetURL(v string) string {
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "/") {
		return v
	}
	return "/admin/assets/" + v
}
