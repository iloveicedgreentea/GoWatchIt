package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
	"context"

	"github.com/gin-gonic/gin"

	"github.com/iloveicedgreentea/go-plex/models"
	"github.com/iloveicedgreentea/go-plex/internal/common"
	"github.com/stretchr/testify/assert"

	"bytes"
	"net/http"
	"time"
)

func TestCheckUUID(t *testing.T) {
	type uuidTest struct {
		uuid   string
		filter string
		pass   bool
	}

	tt := []uuidTest{
		{"fakeuuidtesting", "fakeuuidtesting", true},
		{"fakeuuidtestingfail", "fakeuuidtesting", false},
		{"fakeuuidtesting", "fakeuuidtesting, fakeuuidtesting2", true},
		{"fakeuuidtesting", "fakeuuidtesting,fakeuuidtesting2", true},
		{"fakeuuidtesting", "fakeuuidtestingfail, fakeuuidtesting2", false},
		{"fakeuuidtesting", "fakeuuidtesting, fakeuuidtesting2,fakeuuidtesting3", true},
		{"fakeuuidtesting", "", true},

	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.uuid, func(t *testing.T) {
			t.Parallel()
			pass := checkUUID(tc.uuid, tc.filter)
			if pass != tc.pass {
				t.Errorf("UUID %s should pass: %v", tc.uuid, tc.pass)
			}
		})
	}
}

func createMockMultipartRequest(rawBody string) *http.Request {
	req, _ := http.NewRequest("POST", "http://localhost:9999/plexwebhook", bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=------------------------9a74d1aa5a7ac807")
	return req
}

func TestProcessWebhook(t *testing.T) {
	// test multipart request body
	rawBody := "--------------------------9a74d1aa5a7ac807\r\nContent-Disposition: form-data; name=\"payload\"\r\nContent-Type: application/json\r\n\r\n{\"event\":\"media.play\",\"user\":false,\"owner\":true,\"Account\":{\"id\":1234,\"thumb\":\"https://plex.tv/users/1234/avatar?c=1234\",\"title\":\"o\"},\"Server\":{\"title\":\"123\",\"uuid\":\"fakeuuidtesting\"},\"Player\":{\"local\":false,\"publicAddress\":\"192.168.1.1\",\"title\":\"Player\",\"uuid\":\"player-id\"},\"Metadata\":{\"librarySectionType\":\"show\",\"ratingKey\":\"3019\",\"key\":\"/library/metadata/3019\",\"parentRatingKey\":\"3009\",\"grandparentRatingKey\":\"2958\",\"guid\":\"plex://episode/5d9c12a208fddd001f318b56\",\"parentGuid\":\"plex://season/602e680b9b7e9c002d71a5e0\",\"grandparentGuid\":\"plex://show/5d9c086d2192ba001f3101c6\",\"type\":\"episode\",\"title\":\"The One Where Rachel Quits\",\"titleSort\":\"One Where Rachel Quits\",\"grandparentKey\":\"/library/metadata/2958\",\"parentKey\":\"/library/metadata/3009\",\"librarySectionTitle\":\"TV Shows\",\"librarySectionID\":2,\"librarySectionKey\":\"/library/sections/2\",\"grandparentTitle\":\"Friends\",\"parentTitle\":\"Season 3\",\"contentRating\":\"TV-14\",\"summary\":\"Rachel makes a rash decision after Gunther tells her she needs to be retrained.  Phoebe tries to help Joey when he gets a job selling Christmas trees.  And Ross accidentally breaks a girl's leg and tries to make it up to her. [Christmas Episode]\",\"index\":10,\"parentIndex\":3,\"audienceRating\":8.1,\"viewCount\":1,\"lastViewedAt\":1696800036,\"year\":1996,\"thumb\":\"/library/metadata/3019/thumb/1687815756\",\"art\":\"/library/metadata/2958/art/1695823887\",\"parentThumb\":\"/library/metadata/3009/thumb/1687815755\",\"grandparentThumb\":\"/library/metadata/2958/thumb/1695823887\",\"grandparentArt\":\"/library/metadata/2958/art/1695823887\",\"grandparentTheme\":\"/library/metadata/2958/theme/1695823887\",\"duration\":1320000,\"originallyAvailableAt\":\"1996-12-12\",\"addedAt\":1669266057,\"updatedAt\":1687815756,\"audienceRatingImage\":\"themoviedb://image.rating\",\"Guid\":[{\"id\":\"imdb://tt0583474\"},{\"id\":\"tmdb://86334\"},{\"id\":\"tvdb://303878\"}],\"Rating\":[{\"image\":\"themoviedb://image.rating\",\"value\":8.1,\"type\":\"audience\"}],\"Director\":[{\"id\":17953,\"filter\":\"director=17953\",\"tag\":\"Terry Hughes\",\"tagKey\":\"5d7768384de0ee001fccc190\",\"thumb\":\"https://image.tmdb.org/t/p/original/ffU0D0Yn6RIjdufcviD3e5tn7Hu.jpg\"}],\"Writer\":[{\"id\":17812,\"filter\":\"writer=17812\",\"tag\":\"Michael Curtis\",\"tagKey\":\"5e1635494c78f7003e7f44ba\"},{\"id\":17813,\"filter\":\"writer=17813\",\"tag\":\"Greg Malins\",\"tagKey\":\"5d7768760ea56a001e2a5a4c\",\"thumb\":\"https://metadata-static.plex.tv/b/people/b9a7830f2754cca651abbefe7d64fdd1.jpg\"}],\"Role\":[{\"id\":15772,\"filter\":\"actor=15772\",\"tag\":\"Mae Whitman\",\"tagKey\":\"5d776831103a2d001f566b27\",\"role\":\"Sarah Tuttle\",\"thumb\":\"https://metadata-static.plex.tv/8/people/848114147b5a88bf0a6fab205d9524dc.jpg\"},{\"id\":17648,\"filter\":\"actor=17648\",\"tag\":\"Shelley Berman\",\"tagKey\":\"5d776827103a2d001f564674\",\"role\":\"Mr. Kaplan Jr.\",\"thumb\":\"https://metadata-static.plex.tv/f/people/fa5ceaa3e423b6ec48b116f19cd2a625.jpg\"},{\"id\":17695,\"filter\":\"actor=17695\",\"tag\":\"Kyla Pratt\",\"tagKey\":\"5d77682d8718ba001e31307a\",\"role\":\"Charla Nichols\",\"thumb\":\"https://metadata-static.plex.tv/8/people/801a2079ce5ddc1000a0373f6d353f2c.jpg\"},{\"id\":17698,\"filter\":\"actor=17698\",\"tag\":\"Romy Rosemont\",\"tagKey\":\"5d77682b999c64001ec2d66b\",\"role\":\"Troop Leader\",\"thumb\":\"https://metadata-static.plex.tv/people/5d77682b999c64001ec2d66b.jpg\"},{\"id\":17955,\"filter\":\"actor=17955\",\"tag\":\"Sandra Gould\",\"tagKey\":\"5d77683aeb5d26001f1e1db5\",\"role\":\"Old Woman (voice)\",\"thumb\":\"https://metadata-static.plex.tv/c/people/c361791218a21938dfa1bfa7e379afd5.jpg\"},{\"id\":17612,\"filter\":\"actor=17612\",\"tag\":\"James Michael Tyler\",\"tagKey\":\"5d776b0ffb0d55001f55a7fb\",\"role\":\"Gunther\",\"thumb\":\"https://metadata-static.plex.tv/3/people/3570d61e44686f5d15724609d9e5d059.jpg\"}]}}\r\n--------------------------9a74d1aa5a7ac807--\r\n"

	// Prepare the mock request and gin context
	req := createMockMultipartRequest(rawBody)
	// record response
	respRecorder := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	// get testing context with recorder
	c, _ := gin.CreateTestContext(respRecorder)
	c.Request = req

	// Prepare a channel to receive the decoded payload
	plexChan := make(chan models.PlexWebhookPayload, 10)
	ctx, cancel := context.WithCancel(context.Background())
	// Call the ProcessWebhook function
	ProcessWebhook(ctx,plexChan, c)
	defer cancel() // TODO: test cancelling

	// Test res
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	select {
	case payload := <-plexChan:
		assert.NotNil(t, payload)
		assert.Equal(t, "player-id", payload.Player.UUID)
	case <- ctx.Done():
		t.Error("ctx cancelled")
	case <-time.After(time.Second * 2): // Wait up to 1 second
		t.Error("Expected payload was not received on plexChan")
	}

	var jsonResponse map[string]interface{}
	err := json.Unmarshal(respRecorder.Body.Bytes(), &jsonResponse)
	assert.NoError(t, err)

	if message, exists := jsonResponse["message"]; exists {
		assert.Equal(t, "Payload processed", message)
	} else {
		t.Error("Expected message not found in response")
	}
}

// function to test handler with predetermined webhook
func TestDecodeWebhook(t *testing.T) {
	// open testing file
	jsonFile, err := os.ReadFile("testdata/media.pause.json")
	if err != nil {
		t.Fatal(err)
	}
	var jsonStr []string
	jsonStr = append(jsonStr, string(jsonFile))

	// mock request
	payload, code, err := common.DecodeWebhook(jsonStr)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if code != 0 {
		t.Fatalf("Code is not 0: %d", code)
	}

	// lazy check that parsing works
	expectedData := models.PlexWebhookPayload{Event: "media.pause", User: true, Owner: true, Account: models.Account{ID: models.IntOrString{IntValue: 123}, Thumb: "https:/", Title: "123"}, Server: models.Server{Title: "123", UUID: "123"}, Player: models.Player{Local: true, PublicAddress: "", Title: "SHIELD Android TV", UUID: "123"}, Metadata: models.Metadata{LibrarySectionType: "movie", RatingKey: "26", Key: "/library/metadata/26", GUID: "plex://movie/123", Studio: "Mikona Productions GmbH & Co. KG", Type: "movie", Title: "2 Fast 2 Furious", LibrarySectionTitle: "Movies", LibrarySectionID: 1, LibrarySectionKey: "/library/sections/1", ContentRating: "PG-13", Summary: "EX LAPD cop Brian O'Conner (Paul Walker) teams up with his ex-con friend Roman Pearce (Tyrese Gibson) and works with undercover U.S. Customs Service agent Monica Fuentes (Eva Mendes) to bring Miami-based drug lord Carter Verone (Cole Hauser) down.", Rating: 3.6, AudienceRating: 5, ViewCount: 1, SkipCount: 1, LastViewedAt: 123, Year: 2003, Tagline: "How fast do you like it?", Thumb: "/library/metadata/26/thumb/123", Art: "/library/metadata/26/art/123", Duration: 6420000, OriginallyAvailableAt: "2003-06-05", AddedAt: 123, UpdatedAt: 123, AudienceRatingImage: "rottentomatoes://image.rating.spilled", ChapterSource: "media", PrimaryExtraKey: "/library/metadata/2329", RatingImage: "rottentomatoes://image.rating.rotten", Genre: []models.Genre{{ID: 5, Filter: "genre=5", Tag: "Action", Count: 97}, {ID: 126, Filter: "genre=126", Tag: "Adventure", Count: 46}, {ID: 23, Filter: "genre=23", Tag: "Crime", Count: 43}, {ID: 25, Filter: "genre=25", Tag: "Thriller", Count: 82}}, Director: []models.Director{{ID: 9848, Filter: "director=9848", Tag: "John Singleton"}}, Writer: []models.Writer{{ID: 9849, Filter: "writer=9849", Tag: "Michael Brandt", Count: 0}, {ID: 9850, Filter: "writer=9850", Tag: "Derek Haas", Count: 0}, {ID: 9851, Filter: "writer=9851", Tag: "Gary Scott Thompson", Count: 2}}, Producer: []models.Producer{{ID: 9882, Filter: "producer=9882", Tag: "Neal H. Moritz", Count: 8}}, Country: []models.Country{{ID: 26, Filter: "country=26", Tag: "Germany", Count: 16}, {ID: 28, Filter: "country=28", Tag: "United States of America", Count: 143}}, GUID0: []models.GUID0{{ID: "imdb://tt0322259"}, {ID: "tmdb://584"}, {ID: "tvdb://20800"}}, Rating0: []models.Rating0{{Image: "imdb://image.rating", Value: 5.9, Type: "audience", Count: 157}, {Image: "rottentomatoes://image.rating.rotten", Value: 3.6, Type: "critic", Count: 44}, {Image: "rottentomatoes://image.rating.spilled", Value: 5, Type: "audience", Count: 29}, {Image: "themoviedb://image.rating", Value: 6.5, Type: "audience", Count: 159}}, Collection: []models.Collection{{ID: 8528, Filter: "collection=123", Tag: "Fast & Furious", Count: 8, GUID: "collection://123"}}, Role: []models.Role{{ID: 9852, Filter: "actor=9852", Tag: "Paul Walker", TagKey: "5d7768275af944001f1f6abf", Count: 6, Role: "Brian O'Conner", Thumb: "https://metadata-static.plex.tv/e/people/eb652fccb3a1455611aee35234f5fba7.jpg"}, {ID: 9853, Filter: "actor=9853", Tag: "Tyrese Gibson", TagKey: "5d7768275af944001f1f6ac0", Count: 5, Role: "Roman Pearce", Thumb: "https://metadata-static.plex.tv/f/people/f47582a4f219a318d48f2606b0d5d005.jpg"}, {ID: 9854, Filter: "actor=9854", Tag: "Eva Mendes", TagKey: "5d7768275af944001f1f6ac1", Count: 3, Role: "Monica Fuentes", Thumb: "https://metadata-static.plex.tv/b/people/badbc172b73bb8af7aa1bbdebae48c5b.jpg"}, {ID: 9855, Filter: "actor=9855", Tag: "Ludacris", TagKey: "5d7768275af944001f1f6ac2", Count: 6, Role: "Tej Parker", Thumb: "https://metadata-static.plex.tv/3/people/3cef46981718b30ce26b93b5598d2a00.jpg"}, {ID: 9856, Filter: "actor=9856", Tag: "Cole Hauser", TagKey: "5d776826880197001ec9070c", Count: 2, Role: "Carter Verone", Thumb: "https://metadata-static.plex.tv/f/people/f3f17528f7e9dfc6ac84cc3f265069dd.jpg"}, {ID: 9857, Filter: "actor=9857", Tag: "James Remar", TagKey: "5d776825999c64001ec2bf71", Count: 0, Role: "Agent Markham", Thumb: "https://metadata-static.plex.tv/5/people/54b2d87fe726a66e910a5d872986cf46.jpg"}, {ID: 9858, Filter: "actor=9858", Tag: "Devon Aoki", TagKey: "5d776826999c64001ec2c5cf", Count: 2, Role: "Suki", Thumb: "https://metadata-static.plex.tv/7/people/7932114625ec8e3adec57f367fa25a3a.jpg"}, {ID: 9859, Filter: "actor=9859", Tag: "Thom Barry", TagKey: "5d7768275af944001f1f6ac3", Count: 2, Role: "Agent Bilkins", Thumb: "https://metadata-static.plex.tv/3/people/36c10e9f1231d708277075378b04d169.jpg"}, {ID: 9860, Filter: "actor=9860", Tag: "Amaury Nolasco", TagKey: "5d7768275af944001f1f6ac4", Count: 0, Role: "Orange Julius", Thumb: "https://metadata-static.plex.tv/2/people/20a0b342299e20d1588e81856248a5ce.jpg"}, {ID: 9861, Filter: "actor=9861", Tag: "Michael Ealy", TagKey: "5d7768275af944001f1f6ac5", Count: 0, Role: "Slap Jack", Thumb: "https://metadata-static.plex.tv/9/people/9fdc65f280fc5285c1318b84d824ac99.jpg"}, {ID: 9862, Filter: "actor=9862", Tag: "Jin Au-Yeung", TagKey: "5d7768275af944001f1f6ac6", Count: 0, Role: "Jimmy", Thumb: "https://metadata-static.plex.tv/4/people/4406310587896220bcbaf1387072f363.jpg"}, {ID: 9863, Filter: "actor=9863", Tag: "Mark Boone Junior", TagKey: "5d776825880197001ec90395", Count: 0, Role: "Detective Whitworth", Thumb: "https://metadata-static.plex.tv/people/5d776825880197001ec90395.jpg"}, {ID: 9864, Filter: "actor=9864", Tag: "Mo Gallini", TagKey: "5d7768275af944001f1f6ac7", Count: 2, Role: "Enrique", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6ac7.jpg"}, {ID: 9865, Filter: "actor=9865", Tag: "Roberto 'Sanz' Sanchez", TagKey: "5d7768275af944001f1f6ac8", Count: 0, Role: "Roberto", Thumb: "https://metadata-static.plex.tv/0/people/091c935635355a5b5553ccffee805ba1.jpg"}, {ID: 9866, Filter: "actor=9866", Tag: "John Cenatiempo", TagKey: "5d776826151a60001f24a7c0", Count: 0, Role: "Korpi", Thumb: "https://metadata-static.plex.tv/6/people/665bb1f04186f2a6b64c0592da1c4ccc.jpg"}, {ID: 9867, Filter: "actor=9867", Tag: "Eric Etebari", TagKey: "5d7768275af944001f1f6ac9", Count: 0, Role: "Darden", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6ac9.jpg"}, {ID: 9868, Filter: "actor=9868", Tag: "Neal H. Moritz", TagKey: "5d7768275af944001f1f6abe", Count: 2, Role: "Swerving Cop", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6abe.jpg"}, {ID: 9869, Filter: "actor=9869", Tag: "Edward Finlay", TagKey: "5d7768275af944001f1f6aca", Count: 0, Role: "Agent Dunn", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6aca.jpg"}, {ID: 9870, Filter: "actor=9870", Tag: "Troy Brown", TagKey: "5d7768275af944001f1f6acb", Count: 0, Role: "Paul Hackett", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6acb.jpg"}, {ID: 9871, Filter: "actor=9871", Tag: "Corey Michael Eubanks", TagKey: "5d7768275af944001f1f6acc", Count: 2, Role: "Max Campisi", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6acc.jpg"}, {ID: 44986, Filter: "actor=44986", Tag: "Sammy Maloof", TagKey: "632158b268feb8052d6f5941", Count: 0, Role: "Joe Osborne", Thumb: "https://metadata-static.plex.tv/5/people/51fc80493f1f0c68257b39024af643b2.jpg"}, {ID: 9383, Filter: "actor=9383", Tag: "Troy Robinson", TagKey: "5d7768275af944001f1f6ace", Count: 0, Role: "Feliz Vispone", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6ace.jpg"}, {ID: 44987, Filter: "actor=44987", Tag: "Sincerely A. Ward", TagKey: "632491c83598608e855578d6", Count: 0, Role: "Slap Jack's Girlfriend", Thumb: "https://metadata-static.plex.tv/f/people/f7798e6b5c9bd0db1a677259ccdc4c78.jpg"}, {ID: 9874, Filter: "actor=9874", Tag: "Nievecita Dubuque", TagKey: "5d7768275af944001f1f6ad1", Count: 0, Role: "Suki's Girl", Thumb: "https://metadata-static.plex.tv/7/people/7b92bf9d6eee6cbaf6ca2662349cbb32.jpg"}, {ID: 44988, Filter: "actor=44988", Tag: "Mateo Herreros", TagKey: "6323ae1024eace258e6eb7fb", Count: 0, Role: "Detective", Thumb: ""}, {ID: 9876, Filter: "actor=9876", Tag: "Kerry Rossall", TagKey: "5d7768253c3c2a001fbca97b", Count: 0, Role: "Police Chase Cop", Thumb: ""}, {ID: 9877, Filter: "actor=9877", Tag: "Marc Macaulay", TagKey: "5d77682585719b001f3a0569", Count: 3, Role: "Agent", Thumb: "https://metadata-static.plex.tv/people/5d77682585719b001f3a0569.jpg"}, {ID: 9878, Filter: "actor=9878", Tag: "Tony Bolano", TagKey: "5d7768275af944001f1f6ada", Count: 0, Role: "Gardener", Thumb: "https://metadata-static.plex.tv/people/5d7768275af944001f1f6ada.jpg"}, {ID: 9879, Filter: "actor=9879", Tag: "Marianne M. Arreaga", TagKey: "5d7768275af944001f1f6adc", Count: 2, Role: "Police Chopper Pilot", Thumb: ""}, {ID: 9880, Filter: "actor=9880", Tag: "Tara Carroll", TagKey: "5d7768275af944001f1f6adb", Count: 0, Role: "Seductress", Thumb: "https://image.tmdb.org/t/p/original/8nd78lt723Q3HHja0JWdYAILJo5.jpg"}, {ID: 9881, Filter: "actor=9881", Tag: "Tamara Jones", TagKey: "5d7768275af944001f1f6add", Count: 0, Role: "Customs Technician", Thumb: ""}}}}

	// check our json was processed as expected. If this test fails, the function has broken
	assert.Equal(t, expectedData, payload)
}


// saving in case I need to test handlers directly

// func TestHandler(t *testing.T) {
// 	// open testing file
// 	jsonFile, err := os.Open("testdata/media.pause.json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer jsonFile.Close()

// 	// create a multipart form field
// 	var b bytes.Buffer
// 	mWriter := multipart.NewWriter(&b)
// 	fw, err := mWriter.CreateFormField("payload")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// copy file into writer
// 	_, err = io.Copy(fw, jsonFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// write the boundary
// 	mWriter.Close()
// 	req := httptest.NewRequest(http.MethodPost, "/plexwebhook", &b)
// 	req.Header.Set("Content-Type", mWriter.FormDataContentType())
// 	w := httptest.NewRecorder()

// 	// mock request
// 	ProcessWebhook(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()
// 	data, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Errorf("Error: %v", err)
// 	}

// 	t.Log(data)
// }
