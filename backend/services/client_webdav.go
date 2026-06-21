package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"latestpack/models"
)

var (
	ErrWebDAVOpenListHeadNotSupported = errors.New("HEAD not supported for OpenList mode")
	ErrUnexpectedWebDAVResponse       = errors.New("unexpected response from WebDAV server")
	ErrUnsupportedChannelType         = errors.New("unsupported channel type for download")
)

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Minute,
	}
}

func newHTTPClientNoRedirect() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func buildWebDAVURL(cfg models.ChannelConfig, version, relPath string) string {
	base := strings.TrimRight(cfg.Endpoint, "/")
	parts := []string{base}
	if cfg.Path != "" {
		parts = append(parts, strings.Trim(cfg.Path, "/"))
	}
	parts = append(parts, version, relPath)
	return strings.Join(parts, "/")
}

func (s *ClientService) downloadViaProxy(ctx context.Context, cfg models.ChannelConfig, webdavURL string, w http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequestWithContext(ctx, r.Method, webdavURL, nil)
	if err != nil {
		return fmt.Errorf("creating WebDAV request: %w", err)
	}
	req.SetBasicAuth(cfg.AccessKey, cfg.SecretKey)
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	resp, err := s.webdavClient.Do(req)
	if err != nil {
		return fmt.Errorf("requesting WebDAV: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		for _, key := range []string{"Content-Type", "Content-Length", "Content-Range", "Accept-Ranges", "ETag"} {
			if v := resp.Header.Get(key); v != "" {
				w.Header().Set(key, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		return err
	case http.StatusNotFound:
		return ErrFileNotFound
	case http.StatusRequestedRangeNotSatisfiable:
		for _, key := range []string{"Content-Range", "Content-Length"} {
			if v := resp.Header.Get(key); v != "" {
				w.Header().Set(key, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		return nil
	default:
		return fmt.Errorf("%w: status %d", ErrUnexpectedWebDAVResponse, resp.StatusCode)
	}
}

func (s *ClientService) downloadViaOpenList(ctx context.Context, cfg models.ChannelConfig, webdavURL string, w http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, webdavURL, nil)
	if err != nil {
		return fmt.Errorf("creating OpenList request: %w", err)
	}
	req.SetBasicAuth(cfg.AccessKey, cfg.SecretKey)

	resp, err := s.webdavClientNoRedirect.Do(req)
	if err != nil {
		return fmt.Errorf("requesting OpenList: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusFound, http.StatusMovedPermanently, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		location := resp.Header.Get("Location")
		if location == "" {
			return fmt.Errorf("%w: empty Location in redirect", ErrUnexpectedWebDAVResponse)
		}
		http.Redirect(w, r, location, http.StatusFound)
		return nil
	case http.StatusNotFound:
		return ErrFileNotFound
	default:
		return fmt.Errorf("%w: status %d (expected redirect)", ErrUnexpectedWebDAVResponse, resp.StatusCode)
	}
}

func (s *ClientService) headWebDAVFileViaProxy(ctx context.Context, cfg models.ChannelConfig, webdavURL string, w http.ResponseWriter) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, webdavURL, nil)
	if err != nil {
		return fmt.Errorf("creating WebDAV HEAD request: %w", err)
	}
	req.SetBasicAuth(cfg.AccessKey, cfg.SecretKey)

	resp, err := s.webdavClient.Do(req)
	if err != nil {
		return fmt.Errorf("requesting WebDAV HEAD: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		for _, key := range []string{"Content-Type", "Content-Length", "Accept-Ranges", "ETag"} {
			if v := resp.Header.Get(key); v != "" {
				w.Header().Set(key, v)
			}
		}
		w.WriteHeader(resp.StatusCode)
		return nil
	case http.StatusNotFound:
		return ErrFileNotFound
	default:
		return fmt.Errorf("%w: status %d", ErrUnexpectedWebDAVResponse, resp.StatusCode)
	}
}
