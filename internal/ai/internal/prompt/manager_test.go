// Package prompt provides dynamic prompt loading and management tests
// Author: Done-0
// Created: 2025-08-31
package prompt

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Done-0/gin-scaffold/configs"
)

func TestManager(t *testing.T) {
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

	manager := New()
	ctx := context.Background()

	t.Run("CreateTemplate", func(t *testing.T) {
		template := &Template{
			Name:        "test_template",
			Description: "test description",
			Messages: []Message{
				{Role: "system", Content: "You are {{.role}} assistant"},
				{Role: "user", Content: "{{.message}}"},
			},
		}

		t.Logf("Creating template: %s", template.Name)
		if err := manager.CreateTemplate(ctx, template); err != nil {
			t.Fatalf("CreateTemplate failed: %v", err)
		}
		t.Logf("Template created successfully")
	})

	t.Run("GetTemplate_WithVariables", func(t *testing.T) {
		vars := map[string]any{"role": "AI", "message": "Hello"}
		t.Logf("Getting template with variables: %+v", vars)

		result, err := manager.GetTemplate(ctx, "test_template", &vars)
		if err != nil {
			t.Fatalf("GetTemplate failed: %v", err)
		}

		t.Logf("Original: %s", "You are {{role}} assistant")
		t.Logf("Result:   %s", result.Messages[0].Content)

		if result.Messages[0].Content != "You are AI assistant" {
			t.Errorf("Variable replacement failed, got: %s", result.Messages[0].Content)
		}
		t.Logf("Variable replacement working correctly")
	})

	t.Run("GetTemplate_WithoutVariables", func(t *testing.T) {
		t.Logf("Getting raw template without variables")

		raw, err := manager.GetTemplate(ctx, "test_template", nil)
		if err != nil {
			t.Fatalf("GetTemplate without vars failed: %v", err)
		}

		t.Logf("Raw content: %s", raw.Messages[0].Content)

		if raw.Messages[0].Content != "You are {{.role}} assistant" {
			t.Errorf("Raw template content incorrect, got: %s", raw.Messages[0].Content)
		}
		t.Logf("Raw template content preserved")
	})

	t.Run("ListTemplates", func(t *testing.T) {
		t.Logf("Listing all templates")

		names, err := manager.ListTemplates(ctx)
		if err != nil {
			t.Fatalf("ListTemplates failed: %v", err)
		}

		t.Logf("Found templates: %v", names)

		found := false
		for _, name := range names {
			if name == "test_template" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Template 'test_template' not found in list: %v", names)
		}
		t.Logf("Template found in list")
	})

	t.Run("UpdateTemplate", func(t *testing.T) {
		updated := &Template{
			Name:        "test_template",
			Description: "updated description",
			Messages:    []Message{{Role: "system", Content: "Updated content"}},
		}

		t.Logf("Updating template with new content: %s", updated.Messages[0].Content)

		if err := manager.UpdateTemplate(ctx, "test_template", updated); err != nil {
			t.Fatalf("UpdateTemplate failed: %v", err)
		}

		result, _ := manager.GetTemplate(ctx, "test_template", nil)
		t.Logf("Updated content: %s", result.Messages[0].Content)
		t.Logf("Template updated successfully")
	})

	t.Run("UpdateTemplate_Rename", func(t *testing.T) {
		original := &Template{
			Name:        "original_name",
			Description: "original description",
			Messages:    []Message{{Role: "system", Content: "Original content"}},
		}
		manager.CreateTemplate(ctx, original)

		renamed := &Template{
			Name:        "renamed_template",
			Description: "renamed description",
			Messages:    []Message{{Role: "system", Content: "Renamed content"}},
		}

		t.Logf("Renaming template: original_name -> %s", renamed.Name)

		if err := manager.UpdateTemplate(ctx, "original_name", renamed); err != nil {
			t.Fatalf("UpdateTemplate rename failed: %v", err)
		}

		if _, err := manager.GetTemplate(ctx, "original_name", nil); err == nil {
			t.Error("Old template name should not exist after rename")
		} else {
			t.Logf("Old template correctly removed")
		}

		result, err := manager.GetTemplate(ctx, "renamed_template", nil)
		if err != nil {
			t.Fatalf("New template name should exist: %v", err)
		}
		t.Logf("New template found: %s", result.Name)
		t.Logf("Content: %s", result.Messages[0].Content)
	})

	t.Run("DeleteTemplate", func(t *testing.T) {
		t.Logf("Deleting template: test_template")

		if err := manager.DeleteTemplate(ctx, "test_template"); err != nil {
			t.Fatalf("DeleteTemplate failed: %v", err)
		}
		t.Logf("Template deleted successfully")
	})

	t.Run("ErrorHandling_NonexistentTemplate", func(t *testing.T) {
		t.Logf("Testing error handling for nonexistent template")

		if _, err := manager.GetTemplate(ctx, "nonexistent", nil); err == nil {
			t.Error("Should fail for nonexistent template")
		} else {
			t.Logf("Correctly failed with error: %v", err)
		}
	})

	t.Run("ErrorHandling_DuplicateTemplate", func(t *testing.T) {
		testTemplate := &Template{
			Name:     "duplicate",
			Messages: []Message{{Role: "system", Content: "test"}},
		}

		t.Logf("Creating template: %s", testTemplate.Name)
		manager.CreateTemplate(ctx, testTemplate)

		t.Logf("Attempting to create duplicate template")
		if err := manager.CreateTemplate(ctx, testTemplate); err == nil {
			t.Error("Should fail creating duplicate template")
		} else {
			t.Logf("✅ Correctly failed with error: %v", err)
		}
	})

	t.Run("ErrorHandling_EmptyName", func(t *testing.T) {
		emptyTemplate := &Template{
			Name:     "",
			Messages: []Message{{Role: "system", Content: "test"}},
		}

		t.Logf("Testing empty template name")
		if err := manager.CreateTemplate(ctx, emptyTemplate); err == nil {
			t.Error("Should fail for empty template name")
		} else {
			t.Logf("✅ Correctly failed with error: %v", err)
		}
	})

	t.Run("ErrorHandling_EmptyMessages", func(t *testing.T) {
		emptyMsgTemplate := &Template{
			Name:     "empty_messages",
			Messages: []Message{},
		}

		t.Logf("Testing template with no messages")
		if err := manager.CreateTemplate(ctx, emptyMsgTemplate); err == nil {
			t.Error("Should fail for template with no messages")
		} else {
			t.Logf("✅ Correctly failed with error: %v", err)
		}
	})
}
