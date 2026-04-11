// Copyright (c) 2025-2026, Audionut and the autobrr contributors.
// SPDX-License-Identifier: GPL-2.0-or-later

package unit3d

import (
	"testing"

	"github.com/autobrr/upbrr/pkg/api"
)

func TestBuildCBRName(t *testing.T) {
	tests := []struct {
		name      string
		meta      api.PreparedMetadata
		customTag string
		want      string
	}{
		{
			name: "Basic movie",
			meta: api.PreparedMetadata{
				ReleaseName: "Movie.2023.1080p.WEB-DL.DDP5.1.H.264-GRP",
				Release:     api.ReleaseInfo{Title: "Movie", Year: 2023},
				Tag:         "GRP",
			},
			want: "Movie.2023.1080p.WEB-DL.DDP5.1.H.264-GRP",
		},
		{
			name: "Portuguese DUAL",
			meta: api.PreparedMetadata{
				ReleaseName:    "Movie.2023.1080p.WEB-DL.H.264-GRP",
				Release:        api.ReleaseInfo{Title: "Movie", Year: 2023},
				Tag:            "GRP",
				AudioLanguages: []string{"English", "Portuguese"},
			},
			want: "Movie.2023.1080p.WEB-DL.H.264 DUAL-GRP",
		},
		{
			name: "Custom tag with original group",
			meta: api.PreparedMetadata{
				ReleaseName:    "Movie.2023.1080p.WEB-DL.H.264-CBR",
				Filename:       "Movie.2023.1080p.WEB-DL.H.264-GRP.DUAL.mkv",
				Release:        api.ReleaseInfo{Title: "Movie", Year: 2023},
				Tag:            "CBR",
				AudioLanguages: []string{"English", "Portuguese"},
			},
			customTag: "CBR",
			want:      "Movie.2023.1080p.WEB-DL.H.264-GRP DUAL-CBR",
		},
		{
			name: "Brazilian PT override",
			meta: api.PreparedMetadata{
				ReleaseName: "A.Foreign.Movie.2023.1080p.WEB-DL.H.264-GRP",
				Release: api.ReleaseInfo{
					Title: "A Foreign Movie",
					Year:  2023,
				},
				ExternalIDs: api.ExternalIDs{
					Category: "MOVIE",
				},
				ExternalMetadata: api.ExternalMetadata{
					TMDB: &api.TMDBMetadata{
						OriginalLanguage: "pt",
						RetrievedAKA:     "Filme Brasileiro AKA",
					},
				},
				Tag: "GRP",
			},
			want: "Filme Brasileiro.2023.1080p.WEB-DL.H.264-GRP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildCBRName(tt.meta, tt.customTag); got != tt.want {
				t.Errorf("BuildCBRName() = %q, want %q", got, tt.want)
			}
		})
	}
}
