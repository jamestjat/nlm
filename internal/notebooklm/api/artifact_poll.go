package api

import (
	"context"
	"fmt"
	"math"
	"time"
)

// ArtifactPollResult represents the outcome of polling an artifact.
type ArtifactPollResult struct {
	ArtifactID string
	Status     string // "completed", "failed", "not_found", "timeout"
	Error      string
}

// WaitForArtifact polls ListArtifacts until the given artifact reaches a
// terminal state (completed/failed) or the quota-detection heuristic fires.
//
// Google may silently remove quota-rejected artifacts from the list instead of
// marking them failed. When the artifact is absent for maxNotFound consecutive
// polls spanning at least minNotFoundWindow, this returns a descriptive error
// rather than spinning until timeout.
func (c *Client) WaitForArtifact(ctx context.Context, projectID, artifactID string) (*ArtifactPollResult, error) {
	const (
		initialInterval   = 2 * time.Second
		maxInterval       = 10 * time.Second
		timeout           = 5 * time.Minute
		maxNotFound       = 5
		minNotFoundWindow = 10 * time.Second
	)

	start := time.Now()
	interval := initialInterval
	consecutiveNotFound := 0
	totalNotFound := 0
	var firstNotFoundTime time.Time

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		status, err := c.pollArtifactStatus(projectID, artifactID)
		if err != nil {
			return nil, fmt.Errorf("poll artifact: %w", err)
		}

		switch status {
		case "completed", "failed":
			return &ArtifactPollResult{ArtifactID: artifactID, Status: status}, nil
		case "not_found":
			consecutiveNotFound++
			totalNotFound++
			now := time.Now()
			if firstNotFoundTime.IsZero() {
				firstNotFoundTime = now
			}
			notFoundElapsed := now.Sub(firstNotFoundTime)

			consecutiveTrigger := consecutiveNotFound >= maxNotFound && notFoundElapsed >= minNotFoundWindow
			totalTrigger := totalNotFound >= maxNotFound*2

			if consecutiveTrigger || totalTrigger {
				return &ArtifactPollResult{
					ArtifactID: artifactID,
					Status:     "failed",
					Error: "artifact was removed by the server — " +
						"this may indicate a daily quota/rate limit was exceeded, " +
						"an invalid notebook ID, or a transient API issue",
				}, nil
			}
		default:
			consecutiveNotFound = 0
		}

		if time.Since(start) > timeout {
			return &ArtifactPollResult{
				ArtifactID: artifactID,
				Status:     "timeout",
				Error:      fmt.Sprintf("timed out after %s (last status: %s)", timeout, status),
			}, nil
		}

		sleep := interval
		if remaining := timeout - time.Since(start); sleep > remaining {
			sleep = remaining
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleep):
		}

		interval = time.Duration(math.Min(float64(interval)*1.5, float64(maxInterval)))
	}
}

// pollArtifactStatus checks the current status of an artifact by listing all
// artifacts in the project and finding the matching one.
func (c *Client) pollArtifactStatus(projectID, artifactID string) (string, error) {
	artifacts, err := c.ListArtifacts(projectID)
	if err != nil {
		return "", err
	}

	for _, a := range artifacts {
		if a.ArtifactId == artifactID {
			switch int(a.State) {
			case int(ArtifactStatusCompleted):
				return "completed", nil
			case int(ArtifactStatusFailed):
				return "failed", nil
			case int(ArtifactStatusProcessing):
				return "processing", nil
			case int(ArtifactStatusPending):
				return "pending", nil
			default:
				return "unknown", nil
			}
		}
	}

	return "not_found", nil
}
