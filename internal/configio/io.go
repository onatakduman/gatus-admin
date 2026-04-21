package configio

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	CoreFile          = "core.yaml"
	UIFile            = "ui.yaml"
	AlertingFile      = "alerting.yaml"
	EndpointsFile     = "endpoints.yaml"
	AnnouncementsFile = "announcements.yaml"
	MaintenanceFile   = "maintenance.yaml"
)

var AllFiles = []string{CoreFile, UIFile, AlertingFile, EndpointsFile, AnnouncementsFile, MaintenanceFile}

func Load(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return yaml.Unmarshal(data, v)
}

func Save(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func LoadRaw(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

func SaveRaw(path, content string) error {
	var probe any
	if err := yaml.Unmarshal([]byte(content), &probe); err != nil {
		return fmt.Errorf("invalid YAML: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
