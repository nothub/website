package main

import (
	"testing"
	"testing/fstest"
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

func makePost(title, date string, draft bool) string {
	return "---\ntitle: " + title + "\ndescription: desc\ndate: " + date + "\ntags: []\ndraft: " + boolStr(draft) + "\n---\n\n# " + title + "\n"
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func TestLoadPostsSlugs(t *testing.T) {
	fsys := fstest.MapFS{
		"posts/alpha/index.md": {Data: []byte(makePost("Alpha", "2024-03-01", false))},
		"posts/beta/index.md":  {Data: []byte(makePost("Beta", "2024-01-01", false))},
	}
	entries, err := loadPosts(fsys, newGoldmark(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	slugs := map[string]bool{}
	for _, e := range entries {
		slugs[e.Slug] = true
	}
	if !slugs["alpha"] || !slugs["beta"] {
		t.Errorf("slug extraction wrong, got %v", slugs)
	}
}

func TestLoadPostsDraftFiltering(t *testing.T) {
	fsys := fstest.MapFS{
		"posts/pub/index.md":   {Data: []byte(makePost("Public", "2024-03-01", false))},
		"posts/draft/index.md": {Data: []byte(makePost("Draft", "2024-02-01", true))},
	}

	entries, err := loadPosts(fsys, newGoldmark(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Slug != "pub" {
		t.Errorf("without --drafts: got %v, want only pub", entries)
	}

	entries, err = loadPosts(fsys, newGoldmark(), true)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Errorf("with --drafts: got %d entries, want 2", len(entries))
	}
}

func TestLoadPostsNewestFirst(t *testing.T) {
	fsys := fstest.MapFS{
		"posts/old/index.md": {Data: []byte(makePost("Old", "2023-01-01", false))},
		"posts/mid/index.md": {Data: []byte(makePost("Mid", "2024-01-01", false))},
		"posts/new/index.md": {Data: []byte(makePost("New", "2025-01-01", false))},
	}
	entries, err := loadPosts(fsys, newGoldmark(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Meta.Date.After(entries[i-1].Meta.Date) {
			t.Errorf("not sorted newest-first: %s (%v) after %s (%v)",
				entries[i].Slug, entries[i].Meta.Date,
				entries[i-1].Slug, entries[i-1].Meta.Date)
		}
	}
}
