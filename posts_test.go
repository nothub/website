package main

import (
	"testing"
	"time"
)

func TestParseMeta(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		wantErr bool
		check   func(t *testing.T, m Meta)
	}{
		{
			name: "valid",
			input: map[string]any{
				"title":       "Hello",
				"description": "A post",
				"date":        "2024-01-15",
				"tags":        []any{"go", "test"},
				"draft":       false,
			},
			check: func(t *testing.T, m Meta) {
				if m.Title != "Hello" {
					t.Errorf("Title = %q, want %q", m.Title, "Hello")
				}
				if m.Date != time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) {
					t.Errorf("Date = %v, wrong", m.Date)
				}
				if len(m.Tags) != 2 || m.Tags[0] != "go" {
					t.Errorf("Tags = %v, want [go test]", m.Tags)
				}
			},
		},
		{
			name:    "missing title",
			input:   map[string]any{"description": "x", "date": "2024-01-01", "tags": []any{}, "draft": false},
			wantErr: true,
		},
		{
			name:    "wrong title type",
			input:   map[string]any{"title": 42, "description": "x", "date": "2024-01-01", "tags": []any{}, "draft": false},
			wantErr: true,
		},
		{
			name:    "bad date format",
			input:   map[string]any{"title": "X", "description": "x", "date": "15-01-2024", "tags": []any{}, "draft": false},
			wantErr: true,
		},
		{
			name:    "missing tags",
			input:   map[string]any{"title": "X", "description": "x", "date": "2024-01-01", "draft": false},
			wantErr: true,
		},
		{
			name:    "non-string tag",
			input:   map[string]any{"title": "X", "description": "x", "date": "2024-01-01", "tags": []any{123}, "draft": false},
			wantErr: true,
		},
		{
			name:    "missing draft",
			input:   map[string]any{"title": "X", "description": "x", "date": "2024-01-01", "tags": []any{}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parseMeta(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.check != nil {
				tt.check(t, m)
			}
		})
	}
}
