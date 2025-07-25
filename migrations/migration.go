package migrations

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

// Ошибка для обработки в вызывающем коде (main)
var ErrNoChange = migrate.ErrNoChange

func New(sourceURL, databaseURL string) (*migrate.Migrate, error) {
	return migrate.New(sourceURL, databaseURL)
}

func init() {
	source.Register("custom-file-sprintf", &fmtSprintfSource{})
}

//  Кастомный драйвер поддерживающий форматирование. Подробнее в документации
type fmtSprintfSource struct {
	fileDriver *file.File
	schemaName string
	tmplData   map[string]interface{} // подстановки
}

func (s *fmtSprintfSource) processTemplate(content string) (string, error) {
	tmpl, err := template.New("migration").Parse(content)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s.tmplData); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

func (s *fmtSprintfSource) Open(rawURL string) (source.Driver, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	rawPath := u.Host + u.Path

	filePath, err := filepath.Abs(rawPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	schemaName := u.Query().Get("schema")
	if schemaName == "" {
		return nil, fmt.Errorf("schema name not provided in query")
	}

	// Подключаем file-драйвер
	underlyingURL := "file:" + filePath 
	fileDriver, err := (&file.File{}).Open(underlyingURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open underlying file source: %w", err)
	}

	return &fmtSprintfSource{
		fileDriver: fileDriver.(*file.File),
		schemaName: schemaName,
		tmplData: map[string]interface{}{
			"Schema": schemaName,
		},
	}, nil
}

func (s *fmtSprintfSource) Close() error {
	return s.fileDriver.Close()
}

func (s *fmtSprintfSource) First() (version uint, err error) {
	return s.fileDriver.First()
}

func (s *fmtSprintfSource) Next(version uint) (nextVersion uint, err error) {
	return s.fileDriver.Next(version)
}

func (s *fmtSprintfSource) Prev(version uint) (prevVersion uint, err error) {
	return s.fileDriver.Prev(version)
}

func (s *fmtSprintfSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = s.fileDriver.ReadUp(version)
	if err != nil {
		return nil, "", err
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read UP migration content for version %d: %w", version, err)
	}

	substituted, err := s.processTemplate(string(content))
	if err != nil {
		return nil, "", fmt.Errorf("template processing failed: %w", err)
	}

	return io.NopCloser(bytes.NewReader([]byte(substituted))), identifier, nil
}

func (s *fmtSprintfSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	r, identifier, err = s.fileDriver.ReadDown(version)
	if err != nil {
		return nil, "", err
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read DOWN migration content for version %d: %w", version, err)
	}

	substituted, err := s.processTemplate(string(content))
	if err != nil {
		return nil, "", fmt.Errorf("template processing failed: %w", err)
	}

	return io.NopCloser(bytes.NewReader([]byte(substituted))), identifier, nil
}

