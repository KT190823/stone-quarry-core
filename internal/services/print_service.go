package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/repositories"
)

// printTemplatesDir is where printable template layout JSON files are stored on
// disk. The database keeps only template metadata (name, doc type, size, ...).
const printTemplatesDir = "./data/print_templates"

type PrintTemplateService struct {
	*BaseService
}

func NewPrintTemplateService() *PrintTemplateService {
	repo := repositories.NewPrintTemplateRepo()
	return &PrintTemplateService{BaseService: NewBaseService(repo.BaseRepo)}
}

func printLayoutPath(id string) string {
	return filepath.Join(printTemplatesDir, fmt.Sprintf("%s.json", id))
}

// writeLayoutFile persists the layout JSON to disk.
func writeLayoutFile(id string, layout interface{}) error {
	if err := os.MkdirAll(printTemplatesDir, 0o755); err != nil {
		return err
	}
	var content []byte
	switch v := layout.(type) {
	case string:
		content = []byte(v)
	case nil:
		content = []byte("{}")
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		content = b
	}
	return os.WriteFile(printLayoutPath(id), content, 0o644)
}

// readLayoutFile loads the layout JSON from disk. Returns an empty string when
// the file does not exist (e.g. a legacy row, or no custom layout yet).
func readLayoutFile(id string) string {
	b, err := os.ReadFile(printLayoutPath(id))
	if err != nil {
		return ""
	}
	return string(b)
}

func deleteLayoutFile(id string) {
	_ = os.Remove(printLayoutPath(id))
}

// Create persists template metadata to DB and the layout JSON to disk.
func (s *PrintTemplateService) Create(data map[string]interface{}) (map[string]interface{}, error) {
	layout, hasLayout := data["layout"]
	delete(data, "layout")

	result, err := s.BaseService.Create(data)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%v", result["id"])
	if hasLayout {
		if werr := writeLayoutFile(id, layout); werr != nil {
			return nil, werr
		}
	}
	result["layoutFile"] = printLayoutPath(id)
	return result, nil
}

// Update persists metadata to DB and, when present, writes the layout to disk.
func (s *PrintTemplateService) Update(id string, data map[string]interface{}) (map[string]interface{}, error) {
	layout, hasLayout := data["layout"]
	delete(data, "layout")

	result, err := s.BaseService.Update(id, data)
	if err != nil {
		return nil, err
	}
	if hasLayout {
		if werr := writeLayoutFile(id, layout); werr != nil {
			return nil, werr
		}
	}
	result["layoutFile"] = printLayoutPath(id)
	return result, nil
}

// GetByID returns metadata plus the layout read from disk (falling back to any
// legacy value previously stored in the DB layout column).
func (s *PrintTemplateService) GetByID(id string) (map[string]interface{}, error) {
	result, err := s.BaseService.GetByID(id)
	if err != nil {
		return nil, err
	}
	layout := readLayoutFile(id)
	if layout == "" {
		if legacy, ok := result["layout"].(string); ok && legacy != "" {
			layout = legacy
		}
	}
	if layout != "" {
		result["layout"] = layout
	}
	result["layoutFile"] = printLayoutPath(id)
	return result, nil
}

// Delete removes the DB row and any layout file on disk.
func (s *PrintTemplateService) Delete(id string) error {
	err := s.BaseService.Delete(id)
	if err != nil {
		return err
	}
	deleteLayoutFile(id)
	return nil
}

// ClearDefault unsets the default flag for all templates of the given doc type.
func (s *PrintTemplateService) ClearDefault(ctx context.Context, docType string) error {
	if docType == "" {
		return nil
	}
	_, err := database.Pool.Exec(ctx,
		`UPDATE print_templates SET is_default = FALSE WHERE doc_type = $1`, docType)
	return err
}
