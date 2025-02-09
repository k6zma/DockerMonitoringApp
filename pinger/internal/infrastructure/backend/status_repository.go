package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/k6zma/DockerMonitoringApp/pinger/internal/application/repositories"
	"github.com/k6zma/DockerMonitoringApp/pinger/internal/domain"
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

func (r *BackendStatusRepo) UpdateStatus(ctx context.Context, containerID string, pingTime int64, name, status string, success bool) error {
	url := fmt.Sprintf("%s/api/v1/container_status/%s", r.baseURL, containerID)
	r.logger.Debugf("Sending PATCH request to %s with data: name=%s, status=%s, ping_time=%d", url, name, status, pingTime)

	payload := map[string]interface{}{
		"ping_time": pingTime,
		"name":      name,
		"status":    status,
	}

	if success {
		payload["last_successful_ping"] = time.Now().Format(time.RFC3339)
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		r.logger.Errorf("JSON marshal failed: %v", err)
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		r.logger.Errorf("Request creation failed: %v", err)
		return fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("X-Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Errorf("Request execution failed: %v", err)
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		r.logger.Errorf("API returned error status: %s", resp.Status)
		return fmt.Errorf("api returned error status: %s", resp.Status)
	}

	r.logger.Debugf("Successfully updated status for container ID %s", containerID)
	return nil
}

func (r *BackendStatusRepo) CreateStatus(ctx context.Context, containerID, ip string, pingTime int64, name, status string) error {
	url := fmt.Sprintf("%s/api/v1/container_status", r.baseURL)
	r.logger.Debugf("Sending POST request to %s with data: container_id=%s, name=%s, status=%s, ping_time=%d", url, containerID, name, status, pingTime)

	payload := map[string]interface{}{
		"container_id":         containerID,
		"ip_address":           ip,
		"ping_time":            pingTime,
		"last_successful_ping": time.Now().Format(time.RFC3339),
		"name":                 name,
		"status":               status,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		r.logger.Errorf("JSON marshal failed: %v", err)
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		r.logger.Errorf("Request creation failed: %v", err)
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
		r.logger.Errorf("API returned error status: %s", resp.Status)
		return fmt.Errorf("api returned error status: %s", resp.Status)
	}

	r.logger.Infof("Successfully created status for container ID %s", containerID)
	return nil
}

func (r *BackendStatusRepo) GetStatuses(ctx context.Context) ([]domain.PingResult, error) {
	url := fmt.Sprintf("%s/api/v1/container_status", r.baseURL)
	r.logger.Debugf("Sending GET request to %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		r.logger.Errorf("Request creation failed: %v", err)
		return nil, fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("X-Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Errorf("Request execution failed: %v", err)
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		r.logger.Errorf("API returned error status: %s", resp.Status)
		return nil, fmt.Errorf("api returned error status: %s", resp.Status)
	}

	var statuses []domain.PingResult
	if err := json.NewDecoder(resp.Body).Decode(&statuses); err != nil {
		r.logger.Errorf("JSON decode failed: %v", err)
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	r.logger.Debugf("Received response: %+v", resp)
	r.logger.Debugf("Received statuses: %+v", statuses)
	return statuses, nil
}

func (r *BackendStatusRepo) DeleteStatus(ctx context.Context, containerID string) error {
	url := fmt.Sprintf("%s/api/v1/container_status/%s", r.baseURL, containerID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, http.NoBody)
	if err != nil {
		r.logger.Errorf("Request creation failed: %v", err)
		return fmt.Errorf("request creation failed: %w", err)
	}
	req.Header.Set("X-Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.logger.Errorf("Request execution failed: %v", err)
		return fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNoContent {
		r.logger.Errorf("API returned error status: %s", resp.Status)
		return fmt.Errorf("api returned error status: %s", resp.Status)
	}

	r.logger.Debugf("Successfully deleted status for container ID %s", containerID)
	return nil
}
