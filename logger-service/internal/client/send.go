package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/phanthehoang2503/logger-service/internal/model"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

func SendLog(ctx context.Context, ingestURL string, le model.LogEntry) error {
	if le.Timestamp.IsZero() {
		le.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(le)
	if err != nil {
		return fmt.Errorf("logentry: %s", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ingestURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ingest returned status %d", resp.StatusCode)
	}
	return nil
}
