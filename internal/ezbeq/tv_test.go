package ezbeq

import (
	"testing"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/stretchr/testify/assert"
)

func TestEpisodesContain(t *testing.T) {
	testCases := []struct {
		name     string
		episodes string
		episode  int
		expected bool
	}{
		{name: "single episode", episodes: "3", episode: 3, expected: true},
		{name: "single episode no match", episodes: "3", episode: 4, expected: false},
		{name: "range start", episodes: "1-3", episode: 1, expected: true},
		{name: "range middle", episodes: "1-3", episode: 2, expected: true},
		{name: "range end", episodes: "1-3", episode: 3, expected: true},
		{name: "outside range", episodes: "1-3", episode: 4, expected: false},
		{name: "list of episodes", episodes: "1, 6, 10", episode: 6, expected: true},
		{name: "list of episodes no match", episodes: "1, 6, 10", episode: 7, expected: false},
		{name: "mixed list and range", episodes: "4, 6-9", episode: 8, expected: true},
		{name: "mixed list and range no match", episodes: "4, 6-9", episode: 5, expected: false},
		{name: "zero padded", episodes: "01-06", episode: 3, expected: true},
		{name: "empty", episodes: "", episode: 1, expected: false},
		{name: "unparseable", episodes: "all of them", episode: 1, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, episodesContain(tc.episodes, tc.episode))
		})
	}
}

// every season and episode format seen on a single show in the catalog, since
// authors are not consistent about where they record the episodes
func TestMatchTVEntryAgainstRealCatalog(t *testing.T) {
	oneShow := []models.BeqCatalog{
		{Author: "mobe1969", Season: "1", Episodes: "1,2,3"},
		{Author: "halcyon888", Season: "1", Episodes: "1,2,4,5,6,7,8,10"},
		{Author: "halcyon888", Season: "1", Episodes: "3,9"},
		{Author: "t1g8rsfan", Season: "1E1-2, 5, 10", Episodes: ""},
		{Author: "t1g8rsfan", Season: "1E3", Episodes: ""},
		{Author: "t1g8rsfan", Season: "1E4, 6-9", Episodes: ""},
		{Author: "t1g8rsfan", Season: "2E1, 6, 10", Episodes: ""},
		{Author: "t1g8rsfan", Season: "2E2-5, 7-9", Episodes: ""},
		{Author: "kaelaria", Season: "01", Episodes: ""},
		{Author: "kaelaria", Season: "02", Episodes: ""},
		{Author: "remixmark", Season: "02E1-10", Episodes: ""},
		{Author: "mikejl", Season: "03E01-06", Episodes: ""},
	}

	testCases := []struct {
		name            string
		season, episode int
		// the authors whose entry should name this episode
		expectedAuthors []string
		// whether a whole season entry is available as a fallback
		expectSeasonFallback bool
	}{
		{name: "S1E5", season: 1, episode: 5, expectedAuthors: []string{"halcyon888", "t1g8rsfan"}, expectSeasonFallback: true},
		{name: "S1E3", season: 1, episode: 3, expectedAuthors: []string{"mobe1969", "halcyon888", "t1g8rsfan"}, expectSeasonFallback: true},
		{name: "S1E9", season: 1, episode: 9, expectedAuthors: []string{"halcyon888", "t1g8rsfan"}, expectSeasonFallback: true},
		{name: "S2E3", season: 2, episode: 3, expectedAuthors: []string{"t1g8rsfan", "remixmark"}, expectSeasonFallback: true},
		{name: "S3E1", season: 3, episode: 1, expectedAuthors: []string{"mikejl"}},
		{name: "S3E10 is not covered", season: 3, episode: 10, expectedAuthors: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			playing := &models.SearchRequest{MediaType: "Episode", Season: tc.season, Episode: tc.episode}
			var matched []string
			var seasonFallbacks int
			for _, entry := range oneShow {
				switch matchTVEntry(entry, playing) {
				case tvMatchEpisode:
					matched = append(matched, entry.Author)
				case tvMatchSeason:
					seasonFallbacks++
				}
			}
			assert.ElementsMatch(t, tc.expectedAuthors, matched)
			assert.Equal(t, tc.expectSeasonFallback, seasonFallbacks > 0)
		})
	}
}

func TestMatchTVEntry(t *testing.T) {
	playing := &models.SearchRequest{MediaType: "Episode", Season: 1, Episode: 5}

	testCases := []struct {
		name     string
		entry    models.BeqCatalog
		expected tvMatch
	}{
		{
			name:     "entry names other episodes of the season",
			entry:    models.BeqCatalog{Season: "1", Episodes: "4, 6-9"},
			expected: tvMatchNone,
		},
		{
			name:     "entry covers the episode in a range",
			entry:    models.BeqCatalog{Season: "1", Episodes: "1-6"},
			expected: tvMatchEpisode,
		},
		{
			name:     "entry covers the whole season",
			entry:    models.BeqCatalog{Season: "1", Episodes: ""},
			expected: tvMatchSeason,
		},
		{
			name:     "zero padded season",
			entry:    models.BeqCatalog{Season: "01", Episodes: ""},
			expected: tvMatchSeason,
		},
		{
			name:     "different season",
			entry:    models.BeqCatalog{Season: "2", Episodes: ""},
			expected: tvMatchNone,
		},
		{
			name:     "different season, episode would otherwise match",
			entry:    models.BeqCatalog{Season: "2", Episodes: "1-6"},
			expected: tvMatchNone,
		},
		{
			name:     "movie entry has no season",
			entry:    models.BeqCatalog{Season: "", Episodes: ""},
			expected: tvMatchNone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, matchTVEntry(tc.entry, playing))
		})
	}
}

func TestIsTVSearch(t *testing.T) {
	testCases := []struct {
		name     string
		request  models.SearchRequest
		expected bool
	}{
		{name: "jellyfin episode", request: models.SearchRequest{MediaType: "Episode", Season: 1, Episode: 5}, expected: true},
		{name: "plex episode", request: models.SearchRequest{MediaType: "episode", Season: 1, Episode: 5}, expected: true},
		{name: "movie", request: models.SearchRequest{MediaType: "Movie"}, expected: false},
		{name: "episode without numbering", request: models.SearchRequest{MediaType: "Episode"}, expected: false},
		{name: "episode missing episode number", request: models.SearchRequest{MediaType: "Episode", Season: 1}, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := tc.request
			assert.Equal(t, tc.expected, isTVSearch(&request))
		})
	}
}

func TestBuildSearchEndpoint(t *testing.T) {
	t.Run("movie is filtered by year", func(t *testing.T) {
		m := &models.SearchRequest{MediaType: "Movie", Codec: "DD+ Atmos", Year: 2016, TMDB: "222222"}
		assert.Equal(t, "/api/1/search?audiotypes=DD%2B+Atmos&years=2016&tmdbid=222222", buildSearchEndpoint(m))
	})

	t.Run("episode is not filtered by year", func(t *testing.T) {
		m := &models.SearchRequest{MediaType: "Episode", Codec: "DD+ Atmos", Year: 2026, TMDB: "111111", Season: 3, Episode: 1}
		assert.Equal(t, "/api/1/search?audiotypes=DD%2B+Atmos&tmdbid=111111", buildSearchEndpoint(m))
	})

	t.Run("episode without numbering keeps the year", func(t *testing.T) {
		m := &models.SearchRequest{MediaType: "Episode", Codec: "DD+ Atmos", Year: 2026, TMDB: "111111"}
		assert.Equal(t, "/api/1/search?audiotypes=DD%2B+Atmos&years=2026&tmdbid=111111", buildSearchEndpoint(m))
	})
}
