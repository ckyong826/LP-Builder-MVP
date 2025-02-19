package services

import (
	"backend/internal/models"
	"context"
	"database/sql"
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
)

type TemplateService struct {
    db *sql.DB
}

func NewTemplateService(db *sql.DB) *TemplateService {
    return &TemplateService{db: db}
}

func (s *TemplateService) FindAll(ctx context.Context, page, pageSize int, orderBy, sort string) ([]models.Template, int64, error) {
    var total int64
    err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM templates WHERE deleted_at IS NULL").Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("count error: %w", err)
    }

    query := `SELECT id, original_url, html_path, file_paths, status, error_message, 
              created_at, updated_at, deleted_at 
              FROM templates 
              WHERE deleted_at IS NULL`
    
    if orderBy != "" {
        direction := "ASC"
        if strings.ToLower(sort) == "desc" {
            direction = "DESC"
        }
        query += fmt.Sprintf(" ORDER BY %s %s", orderBy, direction)
    }
    
    if page > 0 && pageSize > 0 {
        offset := (page - 1) * pageSize
        query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
    }

    rows, err := s.db.QueryContext(ctx, query)
    if err != nil {
        return nil, 0, fmt.Errorf("query error: %w", err)
    }
    defer rows.Close()

    var templates []models.Template
    for rows.Next() {
        var t models.Template
        err := rows.Scan(
            &t.ID, &t.OriginalURL, &t.HTMLPath, &t.FilePaths,
            &t.Status, &t.ErrorMessage, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("scan error: %w", err)
        }
        templates = append(templates, t)
    }

    return templates, total, nil
}

func (s *TemplateService) FindOneById(ctx context.Context, id int64) (*models.Template, error) {
    t := &models.Template{}
    err := s.db.QueryRowContext(ctx, `
        SELECT id, original_url, html_path, file_paths, status, error_message, 
               created_at, updated_at, deleted_at 
        FROM templates 
        WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
        &t.ID, &t.OriginalURL, &t.HTMLPath, &t.FilePaths,
        &t.Status, &t.ErrorMessage, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("template not found")
    }
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    return t, nil
}

func (s *TemplateService) FindOneByUrl(ctx context.Context, url string) (*models.Template, error) {
    t := &models.Template{}
    err := s.db.QueryRowContext(ctx, `
        SELECT id, original_url, html_path, file_paths, status, error_message, 
               created_at, updated_at, deleted_at 
        FROM templates 
        WHERE original_url = $1 AND deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT 1`, url).Scan(
        &t.ID, &t.OriginalURL, &t.HTMLPath, &t.FilePaths,
        &t.Status, &t.ErrorMessage, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("query error: %w", err)
    }
    return t, nil
}

func (s *TemplateService) Create(ctx context.Context, template *models.Template) error {
    template.CreatedAt = time.Now()
    if template.Status == "" {
        template.Status = models.StatusPending
    }
    if template.FilePaths == "" {
        template.FilePaths = "{}"
    }

    err := s.db.QueryRowContext(ctx, `
        INSERT INTO templates (original_url, html_path, file_paths, status, error_message, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`,
        template.OriginalURL, template.HTMLPath, template.FilePaths,
        template.Status, template.ErrorMessage, template.CreatedAt,
    ).Scan(&template.ID)

    if err != nil {
        return fmt.Errorf("create error: %w", err)
    }
    return nil
}

func (s *TemplateService) Update(ctx context.Context, template *models.Template) error {
    template.UpdatedAt = time.Now()
    result, err := s.db.ExecContext(ctx, `
        UPDATE templates 
        SET original_url = $1, html_path = $2, file_paths = $3, 
            status = $4, error_message = $5, updated_at = $6 
        WHERE id = $7 AND deleted_at IS NULL`,
        template.OriginalURL, template.HTMLPath, template.FilePaths,
        template.Status, template.ErrorMessage, template.UpdatedAt, template.ID,
    )
    if err != nil {
        return fmt.Errorf("update error: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("rows affected error: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("template not found")
    }
    return nil
}

func (s *TemplateService) Delete(ctx context.Context, id int64) error {
    result, err := s.db.ExecContext(ctx, `
        UPDATE templates 
        SET deleted_at = $1 
        WHERE id = $2 AND deleted_at IS NULL`,
        time.Now(), id,
    )
    if err != nil {
        return fmt.Errorf("delete error: %w", err)
    }

    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("rows affected error: %w", err)
    }
    if rows == 0 {
        return fmt.Errorf("template not found")
    }
    return nil
}

// Your ConvertUrlToFile and its helper functions
func (s *TemplateService) ConvertUrlToFile(ctx context.Context, template *models.Template, request models.ConvertUrlToFile) error {
    existingTemplate, err := s.FindOneByUrl(ctx, request.URL)
    if err != nil {
        return fmt.Errorf("failed to check existing template: %w", err)
    }
    if existingTemplate != nil {
        template.ID = existingTemplate.ID
        return nil
    }
    
    // Initialize template
    template.OriginalURL = request.URL
    template.Status = models.StatusProgress
    template.CreatedAt = time.Now()
    template.FilePaths = "{}"

    if err := s.Create(ctx, template); err != nil {
        return fmt.Errorf("failed to initialize template record: %w", err)
    }

    // Download HTML
    html, err := s.getHTML(request.URL)
    if err != nil {
        template.Status = models.StatusFailed
        template.ErrorMessage = sql.NullString{String: fmt.Sprintf("failed to download HTML: %v", err), Valid: true}
        s.Update(ctx, template)
        return err
    }

    // Save HTML and extract assets
    baseDir := fmt.Sprintf("output/%d", template.ID)
    err = os.MkdirAll(baseDir, os.ModePerm)
    if err != nil {
        template.Status = models.StatusFailed
        template.ErrorMessage = sql.NullString{String: fmt.Sprintf("failed to create output directory: %v", err), Valid: true}
        s.Update(ctx, template)
        return err
    }

    htmlPath := filepath.Join(baseDir, "index.html")
    err = os.WriteFile(htmlPath, []byte(html), os.ModePerm)
    if err != nil {
        template.Status = models.StatusFailed
        template.ErrorMessage = sql.NullString{String: fmt.Sprintf("failed to save HTML: %v", err), Valid: true}
        s.Update(ctx, template)
        return err
    }

    // Download assets
    assets := s.extractAssets(html)
    filePaths, err := s.downloadAssets(baseDir, request.URL, assets)
    if err != nil {
        template.Status = models.StatusFailed
        template.ErrorMessage = sql.NullString{String: fmt.Sprintf("failed to download assets: %v", err), Valid: true}
        s.Update(ctx, template)
        return err
    }

    // Update template
    filePathsJson, _ := json.Marshal(filePaths)
    template.Status = models.StatusComplete
    template.HTMLPath = htmlPath
    template.FilePaths = string(filePathsJson)
    template.UpdatedAt = time.Now()

    return s.Update(ctx, template)
}

func (s *TemplateService) getHTML(urlStr string) (string, error) {
    resp, err := http.Get(urlStr)
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
        `<link[^>]*?href=["'](.*?\.css(?:\?[^"']*)?)"[^>]*>`,
        `<script.*?src="(.*?\.js(?:\?[^"]*)?)"`,
        `<img.*?src="(.*?\.(jpg|jpeg|png|gif|svg))"`,
    }
    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllStringSubmatch(html, -1)
        for _, match := range matches {
            if len(match) > 1 && match[1] != "" {
                // Remove query parameters for local file storage
                asset := strings.Split(match[1], "?")[0]
                assets = append(assets, asset)
            }
        }
    }
    return assets
}

func (s *TemplateService) downloadAssets(baseDir, baseURL string, assets []string) (map[string][]string, error) {
    filePaths := map[string][]string{"css": {}, "js": {}, "images": {}}
    base, err := url.Parse(baseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse base URL: %w", err)
    }

    client := &http.Client{
        Timeout: time.Second * 30,
    }

    for _, asset := range assets {
        assetURL, err := url.Parse(asset)
        if err != nil {
            continue // Skip invalid URLs instead of failing
        }

        // Handle both relative and absolute URLs
        var fullURL string
        if assetURL.IsAbs() {
            fullURL = asset
        } else {
            fullURL = base.ResolveReference(assetURL).String()
        }

        var folder string
        var assetType string

        switch {
        case strings.HasSuffix(strings.Split(asset, "?")[0], ".css"):
            folder = filepath.Join(baseDir, "assets/css")
            assetType = "css"
        case strings.HasSuffix(strings.Split(asset, "?")[0], ".js"):
            folder = filepath.Join(baseDir, "assets/js")
            assetType = "js"
        default:
            folder = filepath.Join(baseDir, "assets/images")
            assetType = "images"
        }

        if err := os.MkdirAll(folder, os.ModePerm); err != nil {
            return nil, fmt.Errorf("failed to create folder %s: %w", folder, err)
        }

        // Clean filename and remove query parameters
        cleanFilename := filepath.Base(strings.Split(assetURL.Path, "?")[0])
        filename := filepath.Join(folder, cleanFilename)

        // Download with proper HTTP client and headers
        req, err := http.NewRequest("GET", fullURL, nil)
        if err != nil {
            continue // Skip failed requests
        }

        // Add common headers
        req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
        req.Header.Set("Accept", "*/*")
        req.Header.Set("Accept-Language", "en-US,en;q=0.9")

        resp, err := client.Do(req)
        if err != nil {
            continue // Skip failed downloads
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            continue // Skip non-200 responses
        }

        out, err := os.Create(filename)
        if err != nil {
            continue
        }
        defer out.Close()

        if _, err := io.Copy(out, resp.Body); err != nil {
            continue
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

func (s *TemplateService) GetTemplateContent(ctx context.Context, templateID int64) (*models.FileContent, error) {
    // Get template record
    template, err := s.FindOneById(ctx, templateID)
    if err != nil {
        return nil, fmt.Errorf("failed to find template: %w", err)
    }

    content := &models.FileContent{
        CSS:    make(map[string]string),
        JS:     make(map[string]string),
        Images: make(map[string]string),
    }

    // Read HTML content
    htmlContent, err := os.ReadFile(template.HTMLPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read HTML file: %w", err)
    }
    content.HTML = string(htmlContent)

    // Parse FilePaths JSON
    var filePaths map[string][]string
    if err := json.Unmarshal([]byte(template.FilePaths), &filePaths); err != nil {
        return nil, fmt.Errorf("failed to parse file paths: %w", err)
    }

    // Read CSS files
    for _, path := range filePaths["css"] {
        cssContent, err := os.ReadFile(path)
        if err != nil {
            return nil, fmt.Errorf("failed to read CSS file %s: %w", path, err)
        }
        content.CSS[filepath.Base(path)] = string(cssContent)
    }

    // Read JS files
    for _, path := range filePaths["js"] {
        jsContent, err := os.ReadFile(path)
        if err != nil {
            return nil, fmt.Errorf("failed to read JS file %s: %w", path, err)
        }
        content.JS[filepath.Base(path)] = string(jsContent)
    }

    // Read image files
    for _, path := range filePaths["images"] {
        relPath := strings.TrimPrefix(path, "output")
        imageURL := "/static" + relPath
        
        content.Images[filepath.Base(path)] = imageURL
    }

    return content, nil
}