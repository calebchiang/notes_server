package services

import (
	"encoding/json"
	"os"
	"strings"
)

type RawTranscript struct {
	Events []struct {
		StartMs    int `json:"tStartMs"`
		DurationMs int `json:"dDurationMs"`
		Segs       []struct {
			UTF8 string `json:"utf8"`
		} `json:"segs"`
	} `json:"events"`
}

type CleanSegment struct {
	StartMs int    `json:"start_ms"`
	EndMs   int    `json:"end_ms"`
	Text    string `json:"text"`
}

type CleanTranscript struct {
	VideoID  string         `json:"video_id"`
	Segments []CleanSegment `json:"segments"`
}

func ExtractStructuredTranscript(jsonFilePath string, videoID string) (*CleanTranscript, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var raw RawTranscript
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, err
	}

	var segments []CleanSegment

	for _, event := range raw.Events {
		if len(event.Segs) == 0 {
			continue
		}

		var builder strings.Builder
		for _, seg := range event.Segs {
			text := seg.UTF8
			if strings.TrimSpace(text) == "" || text == "\n" {
				continue
			}
			builder.WriteString(text)
		}

		cleanText := strings.TrimSpace(builder.String())
		if cleanText == "" {
			continue
		}

		segments = append(segments, CleanSegment{
			StartMs: event.StartMs,
			EndMs:   event.StartMs + event.DurationMs,
			Text:    cleanText,
		})
	}

	return &CleanTranscript{
		VideoID:  videoID,
		Segments: segments,
	}, nil
}
