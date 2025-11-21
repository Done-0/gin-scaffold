// Package i18n provides internationalization utility functions
// Author: Done-0
// Created: 2025-08-24
package i18n

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/Done-0/gin-scaffold/internal/types/consts"
)

// T translates a message key with template parameters
func T(c *gin.Context, key string, params ...string) string {
	localizer, exists := c.Get(consts.LocalizerContextKey)
	if !exists {
		return key
	}

	loc, ok := localizer.(*i18n.Localizer)
	if !ok {
		return key
	}

	config := &i18n.LocalizeConfig{MessageID: key}

	if len(params) > 0 {
		templateData := make(map[string]any)
		for i := 0; i < len(params)-1; i += 2 {
			templateData[params[i]] = params[i+1]
		}
		config.TemplateData = templateData
	}

	if msg, err := loc.Localize(config); err == nil {
		return msg
	}
	return key
}
