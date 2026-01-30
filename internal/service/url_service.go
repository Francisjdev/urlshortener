package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/francisjdev/urlshortener/internal/model"
	"github.com/francisjdev/urlshortener/internal/repository"
	"github.com/google/uuid"
)

type URLService struct {
	urlRepo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{
		urlRepo: repo,
	}
}
func (s *URLService) CreateShortURL(ctx context.Context, url *model.URL) error {
	if url == nil {
		return errors.New("url is nil")
	}
	if url.LongURL == "" {
		return errors.New("long url is empty")
	}

	url.ID = uuid.New()
	url.CreatedAt = time.Now().UTC()
	url.HitCount = 0

	for i := 0; i < 5; i++ {
		code, err := s.createCharCode()
		if err != nil {
			return err
		}

		url.Code = code

		err = s.urlRepo.Create(ctx, url)
		if err != nil {
			if errors.Is(err, repository.ErrCodeAlreadyExists) {
				continue
			}
			return err
		}

		return nil
	}

	return errors.New("could not generate unique short code")
}

func (s *URLService) createCharCode() (string, error) {
	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	code := make([]byte, 6)
	max := big.NewInt(int64(len(base62Chars)))
	for i := 0; i < 6; i++ {
		number, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		code[i] = base62Chars[number.Int64()]
	}

	return string(code), nil
}

func (s *URLService) GetCode(ctx context.Context, code string) (*model.URL, error) {
	if code == "" {
		return nil, errors.New("Code is empty")
	}
	urlValue, err := s.urlRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return urlValue, nil
}
