package shortener

import (
    "errors"
    "short-track/internal/domain"
    "short-track/pkg/base62"
)

type service struct {
    repo domain.URLRepository
}

func NewService(r domain.URLRepository) domain.ShortenerService {
    return &service{repo: r}
}

func (s *service) Shorten(originalURL string) (string, error) {
    // 1. Create URL object
    url := &domain.URL{OriginalURL: originalURL}
    
    // 2. Save to DB to get the auto-increment ID
    if err := s.repo.Save(url); err != nil {
        return "", err
    }
    
    // 3. Convert ID to Base62 Short Code
    code := base62.Encode(url.ID)
    return code, nil
}

func (s *service) Resolve(code string) (string, error) {
    url, err := s.repo.GetByCode(code)
    if err != nil {
        return "", errors.New("url not found")
    }
    return url.OriginalURL, nil
}