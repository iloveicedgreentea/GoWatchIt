package plex

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/imdb"
	"github.com/anaskhan96/soup"
	"github.com/iloveicedgreentea/go-plex/logger"
	"github.com/iloveicedgreentea/go-plex/models"
)

var log = logger.GetLogger()

// Stuff to interface directly with Plex
// of course their api is undocumented and worst of all, in xml. I had to piece it together reading various unofficial API implementations

type PlexClient struct {
	ServerURL  string
	Port       string
	HTTPClient http.Client
	ImdbClient *http.Client
}

// return a new instance of a plex client
func NewClient(url, port string) *PlexClient {
	return &PlexClient{
		ServerURL: url,
		Port:      port,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
		ImdbClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &customTransport{http.DefaultTransport},
		},
	}
}

// only used for get all movies
func parseAllMediaContainer(payload []byte) (models.AllMediaContainer, error) {
	var data models.AllMediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// unmarshal xml into a struct
func parseMediaContainer(payload []byte) (models.MediaContainer, error) {
	var data models.MediaContainer
	err := xml.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// pass the path (/library/123) to the plex server
func (c *PlexClient) getPlexReq(path string) ([]byte, error) {
	res, err := c.HTTPClient.Get(fmt.Sprintf("%s:%s%s", c.ServerURL, c.Port, path))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return data, err
}

// send a request to Plex to get data about something
func (c *PlexClient) GetMediaData(libraryKey string) (models.MediaContainer, error) {
	res, err := c.getPlexReq(libraryKey)
	if err != nil {
		return models.MediaContainer{}, err
	}

	data, err := parseMediaContainer(res)
	if err != nil {
		return models.MediaContainer{}, err
	}

	return data, nil
}

func insensitiveContains(s string, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

// map a plex to a beq codec name
func mapPlexToBeqAudioCodec(codecTitle, codecExtendTitle string) string {
	log.Debugf("Full codec from plex received: %v", codecExtendTitle)
	switch {
	// DTS:X
	case insensitiveContains(codecExtendTitle, "DTS:X"):
		return "DTS-X"
	case insensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-X"
	// DTS MA 7.1 containers but not DTS:X codecs
	case insensitiveContains(codecTitle, "DTS-HD MA 7.1") && !insensitiveContains(codecExtendTitle, "DTS:X") && !insensitiveContains(codecExtendTitle, "DTS-X"):
		return "DTS-HD MA 7.1"
	// DTS HA MA 5.1
	case insensitiveContains(codecTitle, "DTS-HD MA 5.1"):
		return "DTS-HD MA 5.1"
	// DTS 5.1
	case insensitiveContains(codecTitle, "DTS 5.1"):
		return "DTS 5.1"
	case insensitiveContains(codecTitle, "TRUEHD 5.1"):
		return "TrueHD 5.1"
	// TrueHD 6.1
	case insensitiveContains(codecTitle, "TRUEHD 6.1"):
		return "TrueHD 6.1"
	// some Atmos titles return True HD 7.1 annoyingly, so lets just assume atmos for now
	case insensitiveContains(codecExtendTitle, "Atmos") || insensitiveContains(codecExtendTitle, "TRUEHD"):
		return "Atmos"
		// TrueHD 5.1
	// DTS HRA
	case insensitiveContains(codecTitle, "DTS-HD HRA 7.1"):
		return "DTS-HD HR 7.1"
	case insensitiveContains(codecTitle, "DTS-HD HRA 5.1"):
		return "DTS-HD HR 5.1"
	// LPCM
	case insensitiveContains(codecTitle, "LPCM 5.1"):
		return "LPCM 5.1"
	case insensitiveContains(codecTitle, "LPCM 7.1"):
		return "LPCM 7.1"
	case insensitiveContains(codecTitle, "LPCM 2.0"):
		return "LPCM 2.0"
	//DD+
	//English (EAC3 5.1) -> dd+ atmos?
	// Assuming EAC3 5.1 is DD+ Atmos, thats how plex seems to call it
	// may not always be the case but easier to assume so
	case insensitiveContains(codecTitle, "EAC3 5.1"):
		return "DD+ Atmos"
	// probably not accurate, but what can you do
	case insensitiveContains(codecTitle, "EAC3 Stereo"):
		return "DD+"
	// disabled because most movies report real atmos as tHD 7.1
	// case strings.Contains(codecExtendTitle, "Surround 7.1") && strings.Contains(codecExtendTitle, "TRUEHD"):
	// 	return "TrueHD 7.1"
	default:
		return "Empty"
	}

}

// get the type of audio codec for BEQ purpose like atmos, dts-x, etc
func (c *PlexClient) GetAudioCodec(data models.MediaContainer) (string, error) {
	var plexAudioCodec string
	// loop over streams, find the FIRST stream with ID = 2 (this is primary audio track) and read that val
	// loop instead of index because of edge case with two or more video streams
	for _, val := range data.Video.Media.Part.Stream {
		if val.StreamType == "2" {
			return mapPlexToBeqAudioCodec(val.DisplayTitle, val.ExtendedDisplayTitle), nil
		}
	}

	return plexAudioCodec, nil
}

// remove garbage from imdb string and convert to float64
func imdbStoFloat64(s string) (r float64) {
	// remove spaces
	// remove stuff after " : 1" and return the first part e.g 2.39
	// return it as float64 so we can do math
	log.Debugf("Plex: converting imdb string: %s", s)
	splitStr := strings.Split(s, ":")
	// first index of split string without spaces
	big := strings.TrimSpace(splitStr[0])
	// convert it to float
	firstVal, err := strconv.ParseFloat(big, 64)
	if err != nil {
		return r
	}
	// if the first number is larger than 3, its almost definitely something like 16:9
	// basically if its "16:9"
	if firstVal > 3 {
		// get the second value, "9" as float
		little, err := strconv.ParseFloat(strings.TrimSpace(splitStr[1]), 64)
		if err != nil {
			return r
		}
		// get the ratio, so 1.78 for 16:9
		r = firstVal / little
	} else {
		// set the comparison to the first value, like 2.39
		r = firstVal
	}

	log.Debugf("Plex: Converted val: %v", r)
	if err != nil {
		log.Error(err)
		log.Debugf("Plex: Length of input is %v", len(s))
		return r
	}
	// there is an edge case where annoyingly imdb will list 16:9 instead of 1.78:1
	if r > 3 && r <= 16 {
		r = 1.78
	}
	// for those crazy people who film in 4:3
	if r > 3 && r <= 4 {
		r = 1.33
	}
	// just in case something has 17:9
	if r > 16 && r <= 17 {
		r = 1.85
	}

	return r
}

// http request to the tech info page
func getImdbTechInfo(titleID string, client *http.Client) ([]soup.Root, error) {
	// use our slow client
	resp, err := soup.GetWithClient(fmt.Sprintf("https://www.imdb.com/title/%s/technical", titleID), client)

	if err != nil {
		return []soup.Root{}, err
	}
	if len(resp) == 0 {
		return []soup.Root{}, errors.New("soup response was empty")
	}
	log.Debug("Done getting soup response")
	docs := soup.HTMLParse(resp)

	// the page uses tr/td to display the info
	techSoup := docs.Find("div", "id", "technical_content")

	// catch nil pointer dereference
	if techSoup.Pointer == nil {
		return []soup.Root{}, techSoup.Error
	}
	res := techSoup.Find("table").FindAll("tr")

	return res, nil
}

// return the table name given a soup.Root schema
func parseImdbTableSchema(input soup.Root) string {
	//table schema:
	// <tr class="odd">
	// 	<td class="label"> Runtime </td>
	// 	<td>
	// 			2 hr 30 min (150 min)
	// 	</td>
	// </tr>
	return strings.TrimSpace(input.Find("td", "class", "label").Text())
}

// return the ratios as float64 given a schema of ratios
func parseImdbAspectSchema(input soup.Root) []float64 {
	var aspects []float64
	// get the items as text
	text := input.FullText()
	// split text by newline
	htmlLines := strings.Split(text, " \n")
	// get only the number
	for _, s := range htmlLines {
		// ignore empty strings
		if len(s) > 8 {
			aspects = append(aspects, imdbStoFloat64(s))
		}
	}
	log.Debugf("discovered aspects: %v", aspects)
	sort.Float64s(aspects)

	return aspects
}

// determine the aspect ratio(s) from a given title
func parseImdbTechnicalInfo(titleID string, client *http.Client) (float64, error) {
	log.Debugf("parsing info for %s", titleID)
	var res []soup.Root
	var err error
	// try up to n times, it seems to sometimes return nil from imdb scraping
	for i := 0; i < 5; i++ {
		res, err = getImdbTechInfo(titleID, client)
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second * 1)
			continue
		}
	}

	// schema
	// 	<td class="label"> Aspect Ratio </td>
	//     <td>
	//         1.43 : 1
	//         (IMAX 70 mm and Laser: some scenes)
	//          <br/>
	//         1.78 : 1
	//         (some scenes: IMAX Blu-ray)
	//          <br/>
	//         1.90 : 1
	//         (Digital IMAX: some scenes)
	//          <br/>
	//         2.20 : 1
	//         (70 mm and Digital)
	//          <br/>
	//         2.39 : 1
	//         (35 mm)

	//     </td>
	//   </tr>
	for _, val := range res {
		if val.Pointer == nil {
			return 0, val.Error
		}

		tableName := parseImdbTableSchema(val)
		// loop through and search until we get to "camera", its after AR so we can exit faster
		if tableName == "Camera" {
			break
		}
		// if its not camera or AR, keep searching
		if tableName != "Aspect Ratio" {
			continue
		}
		if tableName == "Aspect Ratio" {
			// loop through and find all aspects
			// the second element will be the data
			// find the max ratio and return it - max because I would rather zoom to scope and have 16:9 shots cropped
			aspects := parseImdbAspectSchema(val.FindAll("td")[1])
			log.Debug("finished searching")
			// return the maximum value in slice
			return aspects[len(aspects)-1], nil
		}
	}

	return 0, nil
}

// Prevent IMDB rate limiting
type customTransport struct {
	http.RoundTripper
}

func (e *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	defer time.Sleep(time.Second) // don't go too fast or risk being blocked by AWS WAF
	// headers to get around anti-scraping
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	r.Header.Set("Authority", "m.media-amazon.com")
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	r.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	r.Header.Set("Sec-Fetch-Dest", "image")
	r.Header.Set("Sec-Fetch-Mode", "no-cors")
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	r.Header.Set("Sec-Fetch-User", "?1")
	r.Header.Set("Referer", "https://www.imdb.com/")
	r.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	r.Header.Set("Origin", "https://www.imdb.com")
	r.Header.Set("X-Imdb-Client-Name", "imdb-ics-tpanswerscta")

	return e.RoundTripper.RoundTrip(r)
}

// get the aspect ratio like 1.78 (16:9) 1.85 ~17:9 from IMDB
func (c *PlexClient) GetAspectRatio(title string, year int, imdbID string) (float64, error) {
	// Plex directly not useful since almost everything is in a 1.78:1 container

	// poll IMDB to get title id if its blank
	if imdbID == "" {
		results, err := imdb.SearchTitle(c.ImdbClient, title)
		if err != nil {
			return 0, err
		}
		if len(results) == 0 {
			return 0, errors.New("not found")
		}

		// get the title based on name and year match
		for _, result := range results {
			if result.Year == year && strings.Contains(strings.ToLower(result.Name), title) {
				// get technical info
				log.Debugf("found year match: ID %s, Name %s", result.ID, result.Name)
				return parseImdbTechnicalInfo(result.ID, c.ImdbClient)
			}
		}
	}

	// return aspect ratio
	return parseImdbTechnicalInfo(imdbID, c.ImdbClient)

}
