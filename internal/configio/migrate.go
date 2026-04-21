package configio

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Migrate(configDir string) error {
	monolithPath := filepath.Join(configDir, "config.yaml")
	data, err := os.ReadFile(monolithPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse monolith: %w", err)
	}

	writeSection := func(name string, keys []string) error {
		out := map[string]any{}
		for _, k := range keys {
			if v, ok := raw[k]; ok {
				out[k] = v
			}
		}
		if len(out) == 0 {
			return nil
		}
		b, err := yaml.Marshal(out)
		if err != nil {
			return err
		}
		dst := filepath.Join(configDir, name)
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
		tmp := dst + ".tmp"
		if err := os.WriteFile(tmp, b, 0644); err != nil {
			return err
		}
		return os.Rename(tmp, dst)
	}

	if err := writeSection(CoreFile, []string{"web", "storage", "security", "metrics", "concurrency"}); err != nil {
		return err
	}
	if err := writeSection(UIFile, []string{"ui"}); err != nil {
		return err
	}
	if err := writeSection(AlertingFile, []string{"alerting"}); err != nil {
		return err
	}
	if err := writeSection(EndpointsFile, []string{"endpoints"}); err != nil {
		return err
	}
	if err := writeSection(AnnouncementsFile, []string{"announcements"}); err != nil {
		return err
	}
	if err := writeSection(MaintenanceFile, []string{"maintenance"}); err != nil {
		return err
	}

	bak := monolithPath + ".bak"
	if _, err := os.Stat(bak); os.IsNotExist(err) {
		return os.Rename(monolithPath, bak)
	}
	return nil
}
