package models

import "time"

type JellyfinWebhook struct {
	DeviceID           string `json:"DeviceId"`
	DeviceName         string `json:"DeviceName"`
	ClientName         string `json:"ClientName"`
	UserID             string `json:"UserId"`
	ItemID             string `json:"ItemId"`
	ItemType           string `json:"ItemType"`
	NotificationType   string `json:"NotificationType"`
	Year               string `json:"Year"`
	PlayedToCompletion string `json:"PlayedToCompletion"`
	IsPaused           string `json:"IsPaused"`
}

type JellyfinMetadata struct {
	Name                         string           `json:"Name"`
	OriginalTitle                string           `json:"OriginalTitle"`
	ServerID                     string           `json:"ServerId"`
	ID                           string           `json:"Id"`
	Etag                         string           `json:"Etag"`
	SourceType                   string           `json:"SourceType"`
	PlaylistItemID               string           `json:"PlaylistItemId"`
	DateCreated                  time.Time        `json:"DateCreated"`
	DateLastMediaAdded           time.Time        `json:"DateLastMediaAdded"`
	ExtraType                    string           `json:"ExtraType"`
	// AirsBeforeSeasonNumber       int              `json:"AirsBeforeSeasonNumber"`
	// AirsAfterSeasonNumber        int              `json:"AirsAfterSeasonNumber"`
	// AirsBeforeEpisodeNumber      int              `json:"AirsBeforeEpisodeNumber"`
	// CanDelete                    bool             `json:"CanDelete"`
	// CanDownload                  bool             `json:"CanDownload"`
	// HasSubtitles                 bool             `json:"HasSubtitles"`
	// PreferredMetadataLanguage    string           `json:"PreferredMetadataLanguage"`
	// PreferredMetadataCountryCode string           `json:"PreferredMetadataCountryCode"`
	// SupportsSync                 bool             `json:"SupportsSync"`
	Container                    string           `json:"Container"`
	SortName                     string           `json:"SortName"`
	ForcedSortName               string           `json:"ForcedSortName"`
	// Video3DFormat                string           `json:"Video3DFormat"`
	// PremiereDate                 time.Time        `json:"PremiereDate"`
	ExternalUrls                 []ExternalUrls   `json:"ExternalUrls"`
	MediaSources                 []MediaSources   `json:"MediaSources"`
	// CriticRating                 int              `json:"CriticRating"`
	// ProductionLocations          []string         `json:"ProductionLocations"`
	Path                         string           `json:"Path"`
	// EnableMediaSourceDisplay     bool             `json:"EnableMediaSourceDisplay"`
	// OfficialRating               string           `json:"OfficialRating"`
	// CustomRating                 string           `json:"CustomRating"`
	// ChannelID                    string           `json:"ChannelId"`
	// ChannelName                  string           `json:"ChannelName"`
	// Overview                     string           `json:"Overview"`
	// Taglines                     []string         `json:"Taglines"`
	// Genres                       []string         `json:"Genres"`
	// CommunityRating              int              `json:"CommunityRating"`
	// CumulativeRunTimeTicks       int              `json:"CumulativeRunTimeTicks"`
	// RunTimeTicks                 int              `json:"RunTimeTicks"`
	// PlayAccess                   string           `json:"PlayAccess"`
	AspectRatio                  string           `json:"AspectRatio"`
	ProductionYear               int              `json:"ProductionYear"`
	// IsPlaceHolder                bool             `json:"IsPlaceHolder"`
	// Number                       string           `json:"Number"`
	// ChannelNumber                string           `json:"ChannelNumber"`
	// IndexNumber                  int              `json:"IndexNumber"`
	// IndexNumberEnd               int              `json:"IndexNumberEnd"`
	// ParentIndexNumber            int              `json:"ParentIndexNumber"`
	// RemoteTrailers               []RemoteTrailers `json:"RemoteTrailers"`
	// ProviderIds                  ProviderIds      `json:"ProviderIds"`
	// IsHD                         bool             `json:"IsHD"`
	IsFolder                     bool             `json:"IsFolder"`
	ParentID                     string           `json:"ParentId"`
	Type                         string           `json:"Type"`
	// People                       []People         `json:"People"`
	// Studios                      []Studios        `json:"Studios"`
	// GenreItems                   []GenreItems     `json:"GenreItems"`
	// ParentLogoItemID             string           `json:"ParentLogoItemId"`
	// ParentBackdropItemID         string           `json:"ParentBackdropItemId"`
	// ParentBackdropImageTags      []string         `json:"ParentBackdropImageTags"`
	// LocalTrailerCount            int              `json:"LocalTrailerCount"`
	UserData                     UserData         `json:"UserData"`
	// RecursiveItemCount           int              `json:"RecursiveItemCount"`
	// ChildCount                   int              `json:"ChildCount"`
	SeriesName                   string           `json:"SeriesName"`
	SeriesID                     string           `json:"SeriesId"`
	SeasonID                     string           `json:"SeasonId"`
	// SpecialFeatureCount          int              `json:"SpecialFeatureCount"`
	// DisplayPreferencesID         string           `json:"DisplayPreferencesId"`
	// Status                       string           `json:"Status"`
	// AirTime                      string           `json:"AirTime"`
	// AirDays                      []string         `json:"AirDays"`
	// Tags                         []string         `json:"Tags"`
	// PrimaryImageAspectRatio      int              `json:"PrimaryImageAspectRatio"`
	// Artists                      []string         `json:"Artists"`
	// ArtistItems                  []ArtistItems    `json:"ArtistItems"`
	// Album                        string           `json:"Album"`
	// CollectionType               string           `json:"CollectionType"`
	// DisplayOrder                 string           `json:"DisplayOrder"`
	// AlbumID                      string           `json:"AlbumId"`
	// AlbumPrimaryImageTag         string           `json:"AlbumPrimaryImageTag"`
	// SeriesPrimaryImageTag        string           `json:"SeriesPrimaryImageTag"`
	// AlbumArtist                  string           `json:"AlbumArtist"`
	// AlbumArtists                 []AlbumArtists   `json:"AlbumArtists"`
	// SeasonName                   string           `json:"SeasonName"`
	MediaStreams                 []MediaStreams   `json:"MediaStreams"`
	VideoType                    string           `json:"VideoType"`
	PartCount                    int              `json:"PartCount"`
	MediaSourceCount             int              `json:"MediaSourceCount"`
	// ImageTags                    ImageTags        `json:"ImageTags"`
	// BackdropImageTags            []string         `json:"BackdropImageTags"`
	// ScreenshotImageTags          []string         `json:"ScreenshotImageTags"`
	// ParentLogoImageTag           string           `json:"ParentLogoImageTag"`
	// ParentArtItemID              string           `json:"ParentArtItemId"`
	// ParentArtImageTag            string           `json:"ParentArtImageTag"`
	// SeriesThumbImageTag          string           `json:"SeriesThumbImageTag"`
	// ImageBlurHashes              ImageBlurHashes  `json:"ImageBlurHashes"`
	// SeriesStudio                 string           `json:"SeriesStudio"`
	// ParentThumbItemID            string           `json:"ParentThumbItemId"`
	// ParentThumbImageTag          string           `json:"ParentThumbImageTag"`
	// ParentPrimaryImageItemID     string           `json:"ParentPrimaryImageItemId"`
	// ParentPrimaryImageTag        string           `json:"ParentPrimaryImageTag"`
	// Chapters                     []Chapters       `json:"Chapters"`
	// LocationType                 string           `json:"LocationType"`
	// IsoType                      string           `json:"IsoType"`
	// MediaType                    string           `json:"MediaType"`
	// EndDate                      time.Time        `json:"EndDate"`
	// LockedFields                 []string         `json:"LockedFields"`
	// TrailerCount                 int              `json:"TrailerCount"`
	// MovieCount                   int              `json:"MovieCount"`
	// SeriesCount                  int              `json:"SeriesCount"`
	// ProgramCount                 int              `json:"ProgramCount"`
	// EpisodeCount                 int              `json:"EpisodeCount"`
	// SongCount                    int              `json:"SongCount"`
	// AlbumCount                   int              `json:"AlbumCount"`
	// ArtistCount                  int              `json:"ArtistCount"`
	// MusicVideoCount              int              `json:"MusicVideoCount"`
	// LockData                     bool             `json:"LockData"`
	// Width                        int              `json:"Width"`
	// Height                       int              `json:"Height"`
	// CameraMake                   string           `json:"CameraMake"`
	// CameraModel                  string           `json:"CameraModel"`
	// Software                     string           `json:"Software"`
	// ExposureTime                 int              `json:"ExposureTime"`
	// FocalLength                  int              `json:"FocalLength"`
	// ImageOrientation             string           `json:"ImageOrientation"`
	// Aperture                     int              `json:"Aperture"`
	// ShutterSpeed                 int              `json:"ShutterSpeed"`
	// Latitude                     int              `json:"Latitude"`
	// Longitude                    int              `json:"Longitude"`
	// Altitude                     int              `json:"Altitude"`
	// IsoSpeedRating               int              `json:"IsoSpeedRating"`
	// SeriesTimerID                string           `json:"SeriesTimerId"`
	// ProgramID                    string           `json:"ProgramId"`
	// ChannelPrimaryImageTag       string           `json:"ChannelPrimaryImageTag"`
	// StartDate                    time.Time        `json:"StartDate"`
	// CompletionPercentage         int              `json:"CompletionPercentage"`
	// IsRepeat                     bool             `json:"IsRepeat"`
	// EpisodeTitle                 string           `json:"EpisodeTitle"`
	// ChannelType                  string           `json:"ChannelType"`
	Audio                        string           `json:"Audio"`
	IsMovie                      bool             `json:"IsMovie"`
	IsSports                     bool             `json:"IsSports"`
	IsSeries                     bool             `json:"IsSeries"`
	IsLive                       bool             `json:"IsLive"`
	IsNews                       bool             `json:"IsNews"`
	IsKids                       bool             `json:"IsKids"`
	IsPremiere                   bool             `json:"IsPremiere"`
	// TimerID                      string           `json:"TimerId"`
	// CurrentProgram               CurrentProgram   `json:"CurrentProgram"`
}
type ExternalUrls struct {
	Name string `json:"Name"`
	URL  string `json:"Url"`
}
type MediaStreams struct {
	Codec                     string `json:"Codec"`
	CodecTag                  string `json:"CodecTag"`
	Language                  string `json:"Language"`
	ColorRange                string `json:"ColorRange"`
	ColorSpace                string `json:"ColorSpace"`
	ColorTransfer             string `json:"ColorTransfer"`
	ColorPrimaries            string `json:"ColorPrimaries"`
	DvVersionMajor            int    `json:"DvVersionMajor"`
	DvVersionMinor            int    `json:"DvVersionMinor"`
	DvProfile                 int    `json:"DvProfile"`
	DvLevel                   int    `json:"DvLevel"`
	RpuPresentFlag            int    `json:"RpuPresentFlag"`
	ElPresentFlag             int    `json:"ElPresentFlag"`
	BlPresentFlag             int    `json:"BlPresentFlag"`
	DvBlSignalCompatibilityID int    `json:"DvBlSignalCompatibilityId"`
	Comment                   string `json:"Comment"`
	TimeBase                  string `json:"TimeBase"`
	CodecTimeBase             string `json:"CodecTimeBase"`
	Title                     string `json:"Title"`
	VideoRange                string `json:"VideoRange"`
	VideoRangeType            string `json:"VideoRangeType"`
	VideoDoViTitle            string `json:"VideoDoViTitle"`
	LocalizedUndefined        string `json:"LocalizedUndefined"`
	LocalizedDefault          string `json:"LocalizedDefault"`
	LocalizedForced           string `json:"LocalizedForced"`
	LocalizedExternal         string `json:"LocalizedExternal"`
	DisplayTitle              string `json:"DisplayTitle"`
	NalLengthSize             string `json:"NalLengthSize"`
	IsInterlaced              bool   `json:"IsInterlaced"`
	IsAVC                     bool   `json:"IsAVC"`
	ChannelLayout             string `json:"ChannelLayout"`
	// BitRate                   int    `json:"BitRate"`
	// BitDepth                  int    `json:"BitDepth"`
	// RefFrames                 int    `json:"RefFrames"`
	// PacketLength              int    `json:"PacketLength"`
	// Channels                  int    `json:"Channels"`
	// SampleRate                int    `json:"SampleRate"`
	// IsDefault                 bool   `json:"IsDefault"`
	// IsForced                  bool   `json:"IsForced"`
	// Height                    int    `json:"Height"`
	// Width                     int    `json:"Width"`
	// AverageFrameRate          int    `json:"AverageFrameRate"`
	// RealFrameRate             int    `json:"RealFrameRate"`
	Profile                   string `json:"Profile"`
	Type                      string `json:"Type"`
	AspectRatio               string `json:"AspectRatio"`
	// Index                     int    `json:"Index"`
	// Score                     int    `json:"Score"`
	// IsExternal                bool   `json:"IsExternal"`
	// DeliveryMethod            string `json:"DeliveryMethod"`
	// DeliveryURL               string `json:"DeliveryUrl"`
	// IsExternalURL             bool   `json:"IsExternalUrl"`
	// IsTextSubtitleStream      bool   `json:"IsTextSubtitleStream"`
	// SupportsExternalStream    bool   `json:"SupportsExternalStream"`
	Path                      string `json:"Path"`
	// PixelFormat               string `json:"PixelFormat"`
	// Level                     int    `json:"Level"`
	// IsAnamorphic              bool   `json:"IsAnamorphic"`
}
type MediaAttachments struct {
	Codec       string `json:"Codec"`
	CodecTag    string `json:"CodecTag"`
	Comment     string `json:"Comment"`
	Index       int    `json:"Index"`
	FileName    string `json:"FileName"`
	MimeType    string `json:"MimeType"`
	DeliveryURL string `json:"DeliveryUrl"`
}
type RequiredHTTPHeaders struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type MediaSources struct {
	Protocol                   string              `json:"Protocol"`
	ID                         string              `json:"Id"`
	Path                       string              `json:"Path"`
	EncoderPath                string              `json:"EncoderPath"`
	EncoderProtocol            string              `json:"EncoderProtocol"`
	Type                       string              `json:"Type"`
	Container                  string              `json:"Container"`
	Size                       int                 `json:"Size"`
	Name                       string              `json:"Name"`
	IsRemote                   bool                `json:"IsRemote"`
	ETag                       string              `json:"ETag"`
	RunTimeTicks               int                 `json:"RunTimeTicks"`
	ReadAtNativeFramerate      bool                `json:"ReadAtNativeFramerate"`
	IgnoreDts                  bool                `json:"IgnoreDts"`
	IgnoreIndex                bool                `json:"IgnoreIndex"`
	GenPtsInput                bool                `json:"GenPtsInput"`
	SupportsTranscoding        bool                `json:"SupportsTranscoding"`
	SupportsDirectStream       bool                `json:"SupportsDirectStream"`
	SupportsDirectPlay         bool                `json:"SupportsDirectPlay"`
	IsInfiniteStream           bool                `json:"IsInfiniteStream"`
	RequiresOpening            bool                `json:"RequiresOpening"`
	OpenToken                  string              `json:"OpenToken"`
	RequiresClosing            bool                `json:"RequiresClosing"`
	LiveStreamID               string              `json:"LiveStreamId"`
	BufferMs                   int                 `json:"BufferMs"`
	RequiresLooping            bool                `json:"RequiresLooping"`
	SupportsProbing            bool                `json:"SupportsProbing"`
	VideoType                  string              `json:"VideoType"`
	IsoType                    string              `json:"IsoType"`
	Video3DFormat              string              `json:"Video3DFormat"`
	MediaStreams               []MediaStreams      `json:"MediaStreams"`
	MediaAttachments           []MediaAttachments  `json:"MediaAttachments"`
	Formats                    []string            `json:"Formats"`
	Bitrate                    int                 `json:"Bitrate"`
	Timestamp                  string              `json:"Timestamp"`
	RequiredHTTPHeaders        RequiredHTTPHeaders `json:"RequiredHttpHeaders"`
	TranscodingURL             string              `json:"TranscodingUrl"`
	TranscodingSubProtocol     string              `json:"TranscodingSubProtocol"`
	TranscodingContainer       string              `json:"TranscodingContainer"`
	AnalyzeDurationMs          int                 `json:"AnalyzeDurationMs"`
	DefaultAudioStreamIndex    int                 `json:"DefaultAudioStreamIndex"`
	DefaultSubtitleStreamIndex int                 `json:"DefaultSubtitleStreamIndex"`
}
type RemoteTrailers struct {
	URL  string `json:"Url"`
	Name string `json:"Name"`
}
type ProviderIds struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Primary struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Art struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Backdrop struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Banner struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Logo struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Thumb struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Disc struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Box struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Screenshot struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Menu struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Chapter struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type BoxRear struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Profile struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type ImageBlurHashes struct {
	Primary    Primary    `json:"Primary"`
	Art        Art        `json:"Art"`
	Backdrop   Backdrop   `json:"Backdrop"`
	Banner     Banner     `json:"Banner"`
	Logo       Logo       `json:"Logo"`
	Thumb      Thumb      `json:"Thumb"`
	Disc       Disc       `json:"Disc"`
	Box        Box        `json:"Box"`
	Screenshot Screenshot `json:"Screenshot"`
	Menu       Menu       `json:"Menu"`
	Chapter    Chapter    `json:"Chapter"`
	BoxRear    BoxRear    `json:"BoxRear"`
	Profile    Profile    `json:"Profile"`
}
type People struct {
	Name            string          `json:"Name"`
	ID              string          `json:"Id"`
	Role            string          `json:"Role"`
	Type            string          `json:"Type"`
	PrimaryImageTag string          `json:"PrimaryImageTag"`
	ImageBlurHashes ImageBlurHashes `json:"ImageBlurHashes"`
}
type Studios struct {
	Name string `json:"Name"`
	ID   string `json:"Id"`
}
type GenreItems struct {
	Name string `json:"Name"`
	ID   string `json:"Id"`
}
type UserData struct {
	// Rating                int       `json:"Rating"`
	// PlayedPercentage      int       `json:"PlayedPercentage"`
	// UnplayedItemCount     int       `json:"UnplayedItemCount"`
	// PlaybackPositionTicks int       `json:"PlaybackPositionTicks"`
	// PlayCount             int       `json:"PlayCount"`
	// IsFavorite            bool      `json:"IsFavorite"`
	// Likes                 bool      `json:"Likes"`
	// LastPlayedDate        time.Time `json:"LastPlayedDate"`
	// Played                bool      `json:"Played"`
	Key                   string    `json:"Key"`
	ItemID                string    `json:"ItemId"`
}
type ArtistItems struct {
	Name string `json:"Name"`
	ID   string `json:"Id"`
}
type AlbumArtists struct {
	Name string `json:"Name"`
	ID   string `json:"Id"`
}
type ImageTags struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}
type Chapters struct {
	StartPositionTicks int       `json:"StartPositionTicks"`
	Name               string    `json:"Name"`
	ImagePath          string    `json:"ImagePath"`
	ImageDateModified  time.Time `json:"ImageDateModified"`
	ImageTag           string    `json:"ImageTag"`
}
type CurrentProgram struct {
}
