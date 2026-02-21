package services

import (
	"net/url"
	"strings"
)

func ExtractVideoID(videoURL string) string {
	u, err := url.Parse(videoURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	if id := query.Get("v"); id != "" {
		return id
	}

	// fallback for youtu.be links
	if strings.Contains(videoURL, "youtu.be/") {
		parts := strings.Split(videoURL, "/")
		return parts[len(parts)-1]
	}

	return ""
}
