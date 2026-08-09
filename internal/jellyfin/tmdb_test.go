package jellyfin

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

func TestGetJfTMDB(t *testing.T) {
	c := &JellyfinClient{}

	testCases := []struct {
		name     string
		urls     []models.ExternalUrls
		expected string
		wantErr  bool
	}{
		{
			name:     "movie on jellyfin 10.10 and older",
			urls:     []models.ExternalUrls{{Name: "TheMovieDb", URL: "https://www.themoviedb.org/movie/222222"}},
			expected: "222222",
		},
		{
			name:     "movie on jellyfin 10.11",
			urls:     []models.ExternalUrls{{Name: "TMDB", URL: "https://www.themoviedb.org/movie/222222"}},
			expected: "222222",
		},
		{
			name:     "episode url resolves to the series id, not the episode number",
			urls:     []models.ExternalUrls{{Name: "TMDB", URL: "https://www.themoviedb.org/tv/111111/season/1/episode/3"}},
			expected: "111111",
		},
		{
			name:     "series url",
			urls:     []models.ExternalUrls{{Name: "TMDB", URL: "https://www.themoviedb.org/tv/111111"}},
			expected: "111111",
		},
		{
			name: "tmdb entry found alongside other providers",
			urls: []models.ExternalUrls{
				{Name: "IMDb", URL: "https://www.imdb.com/title/tt1111111"},
				{Name: "TVDB", URL: "https://www.thetvdb.com/?tab=series&id=12345"},
				{Name: "TMDB", URL: "https://www.themoviedb.org/tv/111111/season/2/episode/7"},
			},
			expected: "111111",
		},
		{
			name:    "no tmdb entry",
			urls:    []models.ExternalUrls{{Name: "IMDb", URL: "https://www.imdb.com/title/tt1111111"}},
			wantErr: true,
		},
		{
			name:    "tmdb entry with an unparseable url",
			urls:    []models.ExternalUrls{{Name: "TMDB", URL: "https://www.themoviedb.org/person/12345"}},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmdb, err := c.GetJfTMDB(models.JellyfinMetadata{ExternalUrls: tc.urls})
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, tmdb)
		})
	}
}
