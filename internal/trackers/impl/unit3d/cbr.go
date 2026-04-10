// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package unit3d

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/autobrr/upbrr/pkg/api"
)

func siteCBRProfile() unit3DSiteProfile {
	return unit3DSiteProfile{resolveCategoryID: resolveUnit3DCBRCategoryID}
}

func resolveUnit3DCBRCategoryID(meta api.PreparedMetadata) string {
	if strings.EqualFold(resolveUnit3DCategory(meta), "TV") && meta.Anime {
		return "4"
	}
	return resolveUnit3DCategoryID(meta)
}

func BuildCBRName(meta api.PreparedMetadata) string {
	name := baseReleaseName(meta)
	if name == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"DD+ ", "DDP",
		"DD ", "DD",
		"AAC ", "AAC",
		"FLAC ", "FLAC",
		"Dubbed", "",
		"Dual-Audio", "",
	)
	name = replacer.Replace(name)

	// If it is a Series or Anime, remove the year from the title.
	category := resolveUnit3DCategory(meta)
	if category == "TV" || meta.Anime {
		year := strconv.Itoa(meta.Release.Year)
		if meta.Release.Year > 0 && strings.Contains(name, year) {
			name = strings.ReplaceAll(name, "("+year+")", "")
			name = strings.ReplaceAll(name, year, "")
			name = strings.TrimSpace(name)
		}
	}

	origLang := resolveOriginalLanguage(meta)
	aka := ""
	if meta.ExternalMetadata.TMDB != nil {
		aka = meta.ExternalMetadata.TMDB.RetrievedAKA
	}

	// Remove the AKA title, unless it is Brazilian
	if origLang != "pt" {
		if aka != "" {
			name = strings.ReplaceAll(name, aka, "")
		}
	}

	// If it is Brazilian, use only the AKA title, deleting the foreign title
	if origLang == "pt" && aka != "" {
		akaClean := strings.TrimSpace(strings.ReplaceAll(aka, "AKA", ""))
		title := meta.Release.Title
		name = strings.ReplaceAll(name, aka, "")
		name = strings.ReplaceAll(name, title, akaClean)
		name = strings.TrimSpace(name)
	}

	cbrName := name

	if !isDiscType(meta.DiscType) {
		audioTag := ""
		hasPortuguese := false
		for _, l := range meta.AudioLanguages {
			// log audio languages
			log.Println("Audio language:", l)
			lang := strings.ToLower(l)
			if lang == "portuguese" || lang == "português" {
				hasPortuguese = true
				break
			}
		}

		if hasPortuguese {
			if len(meta.AudioLanguages) >= 3 {
				audioTag = " MULTI"
			} else if len(meta.AudioLanguages) == 2 {
				audioTag = " DUAL"
			}
		}

		if audioTag != "" {
			if strings.Contains(cbrName, "-") {
				idx := strings.LastIndex(cbrName, "-")
				parts := []string{cbrName[:idx], cbrName[idx+1:]}

				cbrName = fmt.Sprintf("%s%s-%s", parts[0], audioTag, parts[1])
			} else {
				cbrName += audioTag
			}
		}
	}

	return addNoGroupSuffix(cbrName, meta, "NoGroup")
}
