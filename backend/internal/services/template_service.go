package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
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

func (s *TemplateService) Create(template models.Template) error {
    return s.repo.Create(&template)
}

func (s *TemplateService) ConvertUrlToFile(template models.Template, request models.ConvertUrlToFile) error {
    // Fetch the webpage content
    resp, err := http.Get(request.URL)
    if err != nil {
        return fmt.Errorf("failed to fetch URL: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return errors.New("failed to fetch URL content")
    }

    // Read the body content
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read URL content: %v", err)
    }

    // Parse HTML content
    doc, err := html.Parse(strings.NewReader(string(body)))
    if err != nil {
        return fmt.Errorf("failed to parse HTML: %v", err)
    }

    // Process HTML to download images and CSS files
    err = s.processHTML(doc, request.URL)
    if err != nil {
        return fmt.Errorf("failed to process HTML: %v", err)
    }

    // Save the processed HTML content to a file
    outputPath := "output.html"
    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("failed to create output file: %v", err)
    }
    defer outFile.Close()
    html.Render(outFile, doc)

    // Update template metadata
    template.Name = path.Base(outputPath)
    template.Type = "Converted Page"
    
    // Use repo to save the template to the database
    return s.repo.Create(&template)
}

// processHTML downloads images and CSS files, updating the HTML accordingly
func (s *TemplateService) processHTML(n *html.Node, baseURL string) error {
    if n.Type == html.ElementNode {
        switch n.Data {
        case "img":
            for i := range n.Attr {
                if n.Attr[i].Key == "src" {
                    imgURL := n.Attr[i].Val
                    localPath, err := s.downloadResource(imgURL, baseURL, "images")
                    if err != nil {
                        return err
                    }
                    n.Attr[i].Val = localPath
                }
            }
        case "link":
            for i := range n.Attr {
                if n.Attr[i].Key == "href" && strings.Contains(n.Attr[i].Val, ".css") {
                    cssURL := n.Attr[i].Val
                    localPath, err := s.downloadResource(cssURL, baseURL, "css")
                    if err != nil {
                        return err
                    }
                    n.Attr[i].Val = localPath
                }
            }
        }
    }

    // Recursively process child nodes
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if err := s.processHTML(c, baseURL); err != nil {
            return err
        }
    }
    return nil
}

// downloadResource fetches a resource (image or CSS) and saves it locally, returning the local path
func (s *TemplateService) downloadResource(resourceURL, baseURL, folder string) (string, error) {
    // Handle relative URLs
    if !strings.HasPrefix(resourceURL, "http") {
        if strings.HasPrefix(resourceURL, "/") {
            resourceURL = baseURL + resourceURL
        } else {
            resourceURL = baseURL + "/" + resourceURL
        }
    }

    // Fetch the resource
    resp, err := http.Get(resourceURL)
    if err != nil {
        return "", fmt.Errorf("failed to download resource: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to download resource, status code: %d", resp.StatusCode)
    }

    // Create the folder if it doesn't exist
    os.MkdirAll(folder, os.ModePerm)

    // Save the resource
    resourceName := path.Base(resourceURL)
    resourcePath := path.Join(folder, resourceName)
    outFile, err := os.Create(resourcePath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %v", err)
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to save resource: %v", err)
    }

    return resourcePath, nil
}

func (s *TemplateService) Update(template models.Template) error {
    return s.repo.Update(&template)
}

func (s *TemplateService) Delete(id uint) error {
    return s.repo.Delete(id)
}
