package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/pinger/pkg/utils"
)

type BackendStatusRepo struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     utils.LoggerInterface
}

func NewBackendStatusRepo(
	baseURL, apiKey string,
	logger utils.LoggerInterface,
) repositories.StatusRepository {
	return &BackendStatusRepo{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     logger,
	}
}

func (r *BackendStatusRepo) UpdateStatus(ctx context.Context, ip string, pingTime float64) error {
	url := fmt.Sprintf("%s/api/v1/container_status/%s", r.baseURL, ip)
	return r.sendRequest(ctx, "PATCH", url, map[string]interface{}{
		"ping_time":            pingTime,
		"last_successful_ping": time.Now().Format(time.RFC3339),
	})
}

func (r *BackendStatusRepo) CreateStatus(ctx context.Context, ip string, pingTime float64) error {
	url := fmt.Sprintf("%s/api/v1/container_status", r.baseURL)
	return r.sendRequest(ctx, "POST", url, map[string]interface{}{
		"ip_address":           ip,
		"ping_time":            pingTime,
		"last_successful_ping": time.Now().Format(time.RFC3339),
	})
}

func (r *BackendStatusRepo) sendRequest(
	ctx context.Context,
	method, url string,
	body interface{},
) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("X-Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("api returned error status: %s", resp.Status)
	}

	r.logger.Infof("Successfully processed request to %s", url)

	return nil
}
