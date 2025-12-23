// Package prompter provides dynamic prompt loading and management tests
// Author: Done-0
// Created: 2025-08-31
package prompter

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/Done-0/gin-scaffold/configs"
)

func TestPrompt(t *testing.T) {
	testDir := filepath.Join(os.TempDir(), "prompt_test")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	configDir := filepath.Join(testDir, "configs")
	os.MkdirAll(configDir, 0755)

	configFile := filepath.Join(configDir, "config.local.yml")
	promptDir := filepath.Join(testDir, "prompts")
	os.MkdirAll(promptDir, 0755)

	configContent := `AI:
  PROMPT:
    DIR: ` + promptDir
	os.WriteFile(configFile, []byte(configContent), 0644)

	oldDir, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(oldDir)

	if err := configs.New(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	p := New()
	ctx := context.Background()

	t.Run("CreateTemplate", func(t *testing.T) {
		tmpl := &Template{
			Name:        "Test Template",
			Description: "test description",
			Messages: []Message{
				{Role: "system", Content: "You are {{.role}} assistant"},
				{Role: "user", Content: "{{.message}}"},
			},
		}
		if err := p.CreateTemplate(ctx, "test_template", tmpl); err != nil {
			t.Fatalf("CreateTemplate failed: %v", err)
		}
	})

	t.Run("GetTemplate_WithVariables", func(t *testing.T) {
		vars := map[string]any{"role": "AI", "message": "Hello"}
		result, err := p.GetTemplate(ctx, "test_template", &vars)
		if err != nil {
			t.Fatalf("GetTemplate failed: %v", err)
		}
		if result.Messages[0].Content != "You are AI assistant" {
			t.Errorf("Variable replacement failed, got: %s", result.Messages[0].Content)
		}
	})

	t.Run("GetTemplate_WithoutVariables", func(t *testing.T) {
		raw, err := p.GetTemplate(ctx, "test_template", nil)
		if err != nil {
			t.Fatalf("GetTemplate without vars failed: %v", err)
		}
		if raw.Messages[0].Content != "You are {{.role}} assistant" {
			t.Errorf("Raw template content incorrect, got: %s", raw.Messages[0].Content)
		}
	})

	t.Run("GetTemplate_NameIsMetadata", func(t *testing.T) {
		raw, err := p.GetTemplate(ctx, "test_template", nil)
		if err != nil {
			t.Fatalf("GetTemplate failed: %v", err)
		}
		if raw.Name != "Test Template" {
			t.Errorf("Name should be metadata, got: %s", raw.Name)
		}
	})

	t.Run("ListTemplates_All", func(t *testing.T) {
		names, err := p.ListTemplates(ctx, "")
		if err != nil {
			t.Fatalf("ListTemplates failed: %v", err)
		}
		if !slices.Contains(names, "test_template") {
			t.Errorf("Template 'test_template' not found in list: %v", names)
		}
	})

	t.Run("UpdateTemplate", func(t *testing.T) {
		updated := &Template{
			Name:        "Updated Name",
			Description: "updated description",
			Messages:    []Message{{Role: "system", Content: "Updated content"}},
		}
		if err := p.UpdateTemplate(ctx, "test_template", updated); err != nil {
			t.Fatalf("UpdateTemplate failed: %v", err)
		}
		result, _ := p.GetTemplate(ctx, "test_template", nil)
		if result.Messages[0].Content != "Updated content" {
			t.Errorf("Update failed, got: %s", result.Messages[0].Content)
		}
		if result.Name != "Updated Name" {
			t.Errorf("Name update failed, got: %s", result.Name)
		}
	})

	t.Run("DeleteTemplate_File", func(t *testing.T) {
		if err := p.DeleteTemplate(ctx, "test_template"); err != nil {
			t.Fatalf("DeleteTemplate failed: %v", err)
		}
		if _, err := p.GetTemplate(ctx, "test_template", nil); err == nil {
			t.Error("Deleted template should not exist")
		}
	})

	t.Run("Error_NonexistentTemplate", func(t *testing.T) {
		if _, err := p.GetTemplate(ctx, "nonexistent", nil); err == nil {
			t.Error("Should fail for nonexistent template")
		}
	})

	t.Run("Error_DuplicateTemplate", func(t *testing.T) {
		tmpl := &Template{Name: "Duplicate", Messages: []Message{{Role: "system", Content: "test"}}}
		p.CreateTemplate(ctx, "duplicate", tmpl)
		if err := p.CreateTemplate(ctx, "duplicate", tmpl); err == nil {
			t.Error("Should fail creating duplicate template")
		}
	})

	t.Run("Error_EmptyPath", func(t *testing.T) {
		tmpl := &Template{Name: "Test", Messages: []Message{{Role: "system", Content: "test"}}}
		if err := p.CreateTemplate(ctx, "", tmpl); err == nil {
			t.Error("Should fail for empty path")
		}
	})

	t.Run("Error_EmptyMessages", func(t *testing.T) {
		tmpl := &Template{Name: "Empty Messages", Messages: []Message{}}
		if err := p.CreateTemplate(ctx, "empty_messages", tmpl); err == nil {
			t.Error("Should fail for template with no messages")
		}
	})

	t.Run("Subdirectory_Create", func(t *testing.T) {
		tmpl := &Template{
			Name:     "午夜便利店",
			Messages: []Message{{Role: "system", Content: "Horror story"}},
		}
		if err := p.CreateTemplate(ctx, "stories/midnight_store", tmpl); err != nil {
			t.Fatalf("CreateTemplate in subdirectory failed: %v", err)
		}
		expectedPath := filepath.Join(promptDir, "stories", "midnight_store.json")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Template file not created at: %s", expectedPath)
		}
	})

	t.Run("Subdirectory_Get", func(t *testing.T) {
		result, err := p.GetTemplate(ctx, "stories/midnight_store", nil)
		if err != nil {
			t.Fatalf("GetTemplate from subdirectory failed: %v", err)
		}
		if result.Name != "午夜便利店" {
			t.Errorf("Template name mismatch, got: %s", result.Name)
		}
	})

	t.Run("Subdirectory_CreateNested", func(t *testing.T) {
		tmpl := &Template{
			Name:     "电梯游戏",
			Messages: []Message{{Role: "system", Content: "Elevator horror"}},
		}
		if err := p.CreateTemplate(ctx, "stories/horror/elevator_game", tmpl); err != nil {
			t.Fatalf("CreateTemplate in nested subdirectory failed: %v", err)
		}
	})

	t.Run("Subdirectory_ListAll", func(t *testing.T) {
		names, err := p.ListTemplates(ctx, "")
		if err != nil {
			t.Fatalf("ListTemplates failed: %v", err)
		}
		for _, exp := range []string{"stories/midnight_store", "stories/horror/elevator_game"} {
			if !slices.Contains(names, exp) {
				t.Errorf("Expected '%s' not found in: %v", exp, names)
			}
		}
	})

	t.Run("Subdirectory_ListWithPrefix", func(t *testing.T) {
		names, err := p.ListTemplates(ctx, "stories/horror")
		if err != nil {
			t.Fatalf("ListTemplates with prefix failed: %v", err)
		}
		if len(names) != 1 || names[0] != "stories/horror/elevator_game" {
			t.Errorf("Expected only 'stories/horror/elevator_game', got: %v", names)
		}
	})

	t.Run("Subdirectory_DeleteFile", func(t *testing.T) {
		if err := p.DeleteTemplate(ctx, "stories/horror/elevator_game"); err != nil {
			t.Fatalf("DeleteTemplate file failed: %v", err)
		}
		if _, err := p.GetTemplate(ctx, "stories/horror/elevator_game", nil); err == nil {
			t.Error("Deleted template should not exist")
		}
	})

	t.Run("Directory_Delete", func(t *testing.T) {
		for _, path := range []string{"game/level1", "game/level2", "game/bonus/secret"} {
			tmpl := &Template{Name: "Game Level", Messages: []Message{{Role: "system", Content: "content"}}}
			if err := p.CreateTemplate(ctx, path, tmpl); err != nil {
				t.Fatalf("Failed to create %s: %v", path, err)
			}
		}

		if err := p.DeleteTemplate(ctx, "game"); err != nil {
			t.Fatalf("DeleteTemplate directory failed: %v", err)
		}

		for _, path := range []string{"game/level1", "game/level2", "game/bonus/secret"} {
			if _, err := p.GetTemplate(ctx, path, nil); err == nil {
				t.Errorf("Template '%s' should not exist after directory delete", path)
			}
		}

		gamePath := filepath.Join(promptDir, "game")
		if _, err := os.Stat(gamePath); !os.IsNotExist(err) {
			t.Error("Directory 'game' should not exist after delete")
		}
	})

	t.Run("Directory_DeleteSubdirectory", func(t *testing.T) {
		for _, path := range []string{"category/sub/a", "category/sub/b", "category/other"} {
			tmpl := &Template{Name: "Category Item", Messages: []Message{{Role: "system", Content: "content"}}}
			p.CreateTemplate(ctx, path, tmpl)
		}

		if err := p.DeleteTemplate(ctx, "category/sub"); err != nil {
			t.Fatalf("DeleteTemplate subdirectory failed: %v", err)
		}

		if _, err := p.GetTemplate(ctx, "category/sub/a", nil); err == nil {
			t.Error("Template 'category/sub/a' should not exist")
		}

		if _, err := p.GetTemplate(ctx, "category/other", nil); err != nil {
			t.Error("Template 'category/other' should still exist")
		}
	})

	t.Run("FinalList", func(t *testing.T) {
		names, err := p.ListTemplates(ctx, "")
		if err != nil {
			t.Fatalf("ListTemplates failed: %v", err)
		}
		sort.Strings(names)
		t.Logf("Final template list: %v", names)
	})
}
