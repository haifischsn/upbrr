// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package unit3d

import (
	"fmt"
	"regexp"
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

var audioTagRegex = regexp.MustCompile(`(?i)-([^.-]+)\.(?:DUAL|MULTI)`)

func BuildCBRName(meta api.PreparedMetadata, customTag string) string {
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

	// If it is a TV or Anime, remove the year from the title.
	category := resolveUnit3DCategory(meta)
	if category == "TV" || meta.Anime {
		year := strconv.Itoa(meta.Release.Year)
		if meta.Release.Year > 0 && strings.Contains(name, year) {
			name = strings.ReplaceAll(name, "("+year+")", "")
			name = strings.ReplaceAll(name, year, "")
			name = strings.TrimSpace(name)
		}
	}

	origLang := strings.ToLower(resolveOriginalLanguage(meta))
	aka := ""
	if meta.ExternalMetadata.TMDB != nil {
		aka = meta.ExternalMetadata.TMDB.RetrievedAKA
	}

	ptbrVariations := []string{"português", "portuguese", "pt-br", "pt"}
	isPtBR := false
	for _, variation := range ptbrVariations {
		if origLang == variation {
			isPtBR = true
			break
		}
	}

	// Remove the AKA title, unless it is Brazilian
	if !isPtBR {
		if aka != "" {
			name = strings.ReplaceAll(name, aka, "")
		}
	}

	// If it is Brazilian, use only the AKA title, deleting the foreign title
	if isPtBR && aka != "" {
		akaClean := strings.TrimSpace(strings.ReplaceAll(aka, "AKA", ""))
		title := meta.Release.Title

		name = strings.ReplaceAll(name, aka, "")
		name = strings.ReplaceAll(name, strings.ReplaceAll(aka, " ", "."), "")
		name = strings.ReplaceAll(name, title, akaClean)
		name = strings.ReplaceAll(name, strings.ReplaceAll(title, " ", "."), akaClean)

		name = strings.TrimSpace(name)
	}

	cbrName := name

	if !isDiscType(meta.DiscType) {
		audioTag := ""
		hasPortuguese := false
		for _, l := range meta.AudioLanguages {
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

				cleanCustomTag := strings.TrimPrefix(customTag, "-")
				if cleanCustomTag != "" && strings.Contains(name, cleanCustomTag) {
					searchStr := meta.Filename
					if searchStr == "" {
						searchStr = meta.ReleaseName
					}

					if match := audioTagRegex.FindStringSubmatch(searchStr); len(match) > 1 {
						originalGroupTag := match[1]
						if !strings.EqualFold(originalGroupTag, meta.Release.Group) {
							cbrName = fmt.Sprintf("%s-%s%s-%s", parts[0], originalGroupTag, audioTag, parts[1])
						} else {
							cbrName = fmt.Sprintf("%s%s-%s", parts[0], audioTag, parts[1])
						}
					} else {
						cbrName = fmt.Sprintf("%s%s-%s", parts[0], audioTag, parts[1])
					}
				} else {
					cbrName = fmt.Sprintf("%s%s-%s", parts[0], audioTag, parts[1])
				}
			} else {
				cbrName += audioTag
			}
		}
	}

	return addNoGroupSuffix(cbrName, meta, "NoGroup")
}
