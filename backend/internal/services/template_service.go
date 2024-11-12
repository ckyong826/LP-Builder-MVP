package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TemplateService struct {
    repo *repositories.Repository[models.Template]
}

func NewTemplateService(db *gorm.DB) *TemplateService {
    return &TemplateService{
        repo: repositories.NewRepository[models.Template](db),
    }
}

func (s *TemplateService) FindAll() ([]models.Template, error) {
    return s.repo.FindAll()
}

func (s *TemplateService) FindOneById(id uint) (models.Template, error) {
    template, err := s.repo.FindOneByID(id)
    if err != nil {
        return models.Template{}, err
    }
    return *template, nil
}

func (s *TemplateService) Create(template *models.Template) error {
	return s.repo.Create(template)
}



func (s *TemplateService) ConvertUrlToFile(template *models.Template, request models.ConvertUrlToFile) error {
	// Step 1: Update the template status to "in_progress" and save to DB
	template.OriginalURL = request.URL
	template.Status = "in_progress"
	template.CreatedAt = time.Now()

	template.FilePaths="{}"
	
	if err := s.repo.Create(template); err != nil {
		return fmt.Errorf("failed to initialize template record: %w", err)
	}

	// Step 2: Download HTML
	html, err := s.getHTML(request.URL)
	if err != nil {
		template.Status = "failed"
		template.ErrorMessage = fmt.Sprintf("failed to download HTML: %v", err)
		s.repo.Update(template)
		return err
	}

	// Step 3: Save HTML and extract assets
	baseDir := fmt.Sprintf("output/%d", template.ID)
	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		template.Status = "failed"
		template.ErrorMessage = fmt.Sprintf("failed to create output directory: %v", err)
		s.repo.Update(template)
		return err
	}

	htmlPath := filepath.Join(baseDir, "index.html")
	err = os.WriteFile(htmlPath, []byte(html), os.ModePerm)
	if err != nil {
		template.Status = "failed"
		template.ErrorMessage = fmt.Sprintf("failed to save HTML: %v", err)
		s.repo.Update(template)
		return err
	}

	// Step 4: Download assets and update file paths
	assets := s.extractAssets(html)
	filePaths, err := s.downloadAssets(baseDir, request.URL, assets)
	if err != nil {
		template.Status = "failed"
		template.ErrorMessage = fmt.Sprintf("failed to download assets: %v", err)
		s.repo.Update(template)
		return err
	}

	// Step 5: Update template with paths and mark as "completed"
	filePathsJson, _ := json.Marshal(filePaths)
	template.Status = "completed"
	template.HTMLPath = htmlPath
	template.FilePaths = string(filePathsJson)
	s.repo.Update(template)

	return nil
}

func (s *TemplateService) getHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func (s *TemplateService) extractAssets(html string) []string {
	var assets []string
	patterns := []string{
		`<link.*?href="(.*?\.css)"`,
		`<script.*?src="(.*?\.js)"`,
		`<img.*?src="(.*?\.(jpg|jpeg|png|gif|svg))"`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) > 1 {
				assets = append(assets, match[1])
			}
		}
	}
	return assets
}

func (s *TemplateService) downloadAssets(baseDir string, baseURL string, assets []string) (map[string][]string, error) {
	filePaths := map[string][]string{"css": {}, "js": {}, "images": {}}

	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	for _, asset := range assets {
		// Parse the asset URL relative to the base URL
		assetURL, err := url.Parse(asset)
		if err != nil {
			return nil, fmt.Errorf("failed to parse asset URL %s: %w", asset, err)
		}

		// Resolve the asset URL against the base URL
		fullURL := base.ResolveReference(assetURL)

		// Determine the folder based on asset type
		var folder string
		var assetType string
		if strings.HasSuffix(asset, ".css") {
			folder = filepath.Join(baseDir, "css")
			assetType = "css"
		} else if strings.HasSuffix(asset, ".js") {
			folder = filepath.Join(baseDir, "js")
			assetType = "js"
		} else {
			folder = filepath.Join(baseDir, "images")
			assetType = "images"
		}

		// Create the folder if it doesn't exist
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create folder %s: %w", folder, err)
		}

		// Determine the filename for saving the asset
		filename := filepath.Join(folder, filepath.Base(assetURL.Path))

		// Download the asset file and save it locally
		err = s.downloadFile(fullURL.String(), filename)
		if err != nil {
			return nil, fmt.Errorf("failed to download asset %s: %w", fullURL.String(), err)
		}

		filePaths[assetType] = append(filePaths[assetType], filename)
	}

	return filePaths, nil
}

func (s *TemplateService) downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (s *TemplateService) Update(template models.Template) error {
    return s.repo.Update(&template)
}

func (s *TemplateService) Delete(id uint) error {
    return s.repo.Delete(id)
}
