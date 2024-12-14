package models

type JellyfinExternalLookup struct {
	Name            string `json:"Name"`
	Key             string `json:"Key"`
	Type            string `json:"Type"`
	URLFormatString string `json:"UrlFormatString"`
}

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
	Name             string         `json:"Name"`
	OriginalTitle    string         `json:"OriginalTitle"`
	ServerID         string         `json:"ServerId"`
	ID               string         `json:"Id"`
	SourceType       string         `json:"SourceType"`
	Container        string         `json:"Container"`
	SortName         string         `json:"SortName"`
	ForcedSortName   string         `json:"ForcedSortName"`
	Path             string         `json:"Path"`
	AspectRatio      string         `json:"AspectRatio"`
	ProductionYear   int            `json:"ProductionYear"`
	IsFolder         bool           `json:"IsFolder"`
	ParentID         string         `json:"ParentId"`
	Type             string         `json:"Type"`
	UserData         UserData       `json:"UserData"`
	SeriesName       string         `json:"SeriesName"`
	SeriesID         string         `json:"SeriesId"`
	SeasonID         string         `json:"SeasonId"`
	MediaStreams     []MediaStreams `json:"MediaStreams"`
	VideoType        string         `json:"VideoType"`
	PartCount        int            `json:"PartCount"`
	MediaSourceCount int            `json:"MediaSourceCount"`
	Audio            string         `json:"Audio"`
	ExternalUrls     []ExternalUrls `json:"ExternalUrls"`
	IsMovie          bool           `json:"IsMovie"`
	IsSports         bool           `json:"IsSports"`
	IsSeries         bool           `json:"IsSeries"`
	IsLive           bool           `json:"IsLive"`
	IsNews           bool           `json:"IsNews"`
	IsKids           bool           `json:"IsKids"`
	IsPremiere       bool           `json:"IsPremiere"`
}
type ExternalUrls struct {
	Name string `json:"Name"`
	URL  string `json:"Url"`
}
type MediaStreams struct {
	Codec         string `json:"Codec"`
	CodecTag      string `json:"CodecTag"`
	Language      string `json:"Language"`
	Title         string `json:"Title"`
	DisplayTitle  string `json:"DisplayTitle"`
	ChannelLayout string `json:"ChannelLayout"`
	Profile       string `json:"Profile"`
	Type          string `json:"Type"`
	AspectRatio   string `json:"AspectRatio"`
	Path          string `json:"Path"`
}
type MediaAttachments struct {
	Codec    string `json:"Codec"`
	CodecTag string `json:"CodecTag"`
	FileName string `json:"FileName"`
}
type MediaSources struct {
	Protocol         string             `json:"Protocol"`
	ID               string             `json:"Id"`
	Path             string             `json:"Path"`
	Type             string             `json:"Type"`
	Container        string             `json:"Container"`
	Name             string             `json:"Name"`
	Video3DFormat    string             `json:"Video3DFormat"`
	MediaStreams     []MediaStreams     `json:"MediaStreams"`
	MediaAttachments []MediaAttachments `json:"MediaAttachments"`
	Formats          []string           `json:"Formats"`
	Bitrate          int                `json:"Bitrate"`
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
	Key    string `json:"Key"`
	ItemID string `json:"ItemId"`
}

type JellyfinSession struct {
	PlayState             PlayState         `json:"PlayState"`
	AdditionalUsers       []AdditionalUsers `json:"AdditionalUsers"`
	Capabilities          Capabilities      `json:"Capabilities"`
	RemoteEndPoint        string            `json:"RemoteEndPoint"`
	PlayableMediaTypes    []string          `json:"PlayableMediaTypes"`
	ID                    string            `json:"Id"`
	UserID                string            `json:"UserId"`
	UserName              string            `json:"UserName"`
	Client                string            `json:"Client"`
	DeviceName            string            `json:"DeviceName"`
	DeviceType            string            `json:"DeviceType"`
	DeviceID              string            `json:"DeviceId"`
	ApplicationVersion    string            `json:"ApplicationVersion"`
	IsActive              bool              `json:"IsActive"`
	SupportsMediaControl  bool              `json:"SupportsMediaControl"`
	SupportsRemoteControl bool              `json:"SupportsRemoteControl"`
	HasCustomDeviceName   bool              `json:"HasCustomDeviceName"`
	PlaylistItemID        string            `json:"PlaylistItemId"`
	ServerID              string            `json:"ServerId"`
	SupportedCommands     []string          `json:"SupportedCommands"`
}
type PlayState struct {
	PositionTicks       int    `json:"PositionTicks"`
	CanSeek             bool   `json:"CanSeek"`
	IsPaused            bool   `json:"IsPaused"`
	IsMuted             bool   `json:"IsMuted"`
	VolumeLevel         int    `json:"VolumeLevel"`
	AudioStreamIndex    int    `json:"AudioStreamIndex"`
	SubtitleStreamIndex int    `json:"SubtitleStreamIndex"`
	MediaSourceID       string `json:"MediaSourceId"`
	PlayMethod          string `json:"PlayMethod"`
	RepeatMode          string `json:"RepeatMode"`
	LiveStreamID        string `json:"LiveStreamId"`
}
type AdditionalUsers struct {
	UserID   string `json:"UserId"`
	UserName string `json:"UserName"`
}
type Identification struct {
	FriendlyName     string `json:"FriendlyName"`
	ModelNumber      string `json:"ModelNumber"`
	SerialNumber     string `json:"SerialNumber"`
	ModelName        string `json:"ModelName"`
	ModelDescription string `json:"ModelDescription"`
	ModelURL         string `json:"ModelUrl"`
	Manufacturer     string `json:"Manufacturer"`
	ManufacturerURL  string `json:"ManufacturerUrl"`
}
type XMLRootAttributes struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}
type DirectPlayProfiles struct {
	Container  string `json:"Container"`
	AudioCodec string `json:"AudioCodec"`
	VideoCodec string `json:"VideoCodec"`
	Type       string `json:"Type"`
}
type Conditions struct {
	Condition  string `json:"Condition"`
	Property   string `json:"Property"`
	Value      string `json:"Value"`
	IsRequired bool   `json:"IsRequired"`
}
type ContainerProfiles struct {
	Type       string       `json:"Type"`
	Conditions []Conditions `json:"Conditions"`
	Container  string       `json:"Container"`
}
type ApplyConditions struct {
	Condition  string `json:"Condition"`
	Property   string `json:"Property"`
	Value      string `json:"Value"`
	IsRequired bool   `json:"IsRequired"`
}
type CodecProfiles struct {
	Type            string            `json:"Type"`
	Conditions      []Conditions      `json:"Conditions"`
	ApplyConditions []ApplyConditions `json:"ApplyConditions"`
	Codec           string            `json:"Codec"`
	Container       string            `json:"Container"`
}
type ResponseProfiles struct {
	Container  string       `json:"Container"`
	AudioCodec string       `json:"AudioCodec"`
	VideoCodec string       `json:"VideoCodec"`
	Type       string       `json:"Type"`
	OrgPn      string       `json:"OrgPn"`
	MimeType   string       `json:"MimeType"`
	Conditions []Conditions `json:"Conditions"`
}
type SubtitleProfiles struct {
	Format    string `json:"Format"`
	Method    string `json:"Method"`
	DidlMode  string `json:"DidlMode"`
	Language  string `json:"Language"`
	Container string `json:"Container"`
}
type DeviceProfile struct {
	Name                             string               `json:"Name"`
	ID                               string               `json:"Id"`
	Identification                   Identification       `json:"Identification"`
	FriendlyName                     string               `json:"FriendlyName"`
	Manufacturer                     string               `json:"Manufacturer"`
	ManufacturerURL                  string               `json:"ManufacturerUrl"`
	ModelName                        string               `json:"ModelName"`
	ModelDescription                 string               `json:"ModelDescription"`
	ModelNumber                      string               `json:"ModelNumber"`
	ModelURL                         string               `json:"ModelUrl"`
	SerialNumber                     string               `json:"SerialNumber"`
	EnableAlbumArtInDidl             bool                 `json:"EnableAlbumArtInDidl"`
	EnableSingleAlbumArtLimit        bool                 `json:"EnableSingleAlbumArtLimit"`
	EnableSingleSubtitleLimit        bool                 `json:"EnableSingleSubtitleLimit"`
	SupportedMediaTypes              string               `json:"SupportedMediaTypes"`
	UserID                           string               `json:"UserId"`
	AlbumArtPn                       string               `json:"AlbumArtPn"`
	MaxAlbumArtWidth                 int                  `json:"MaxAlbumArtWidth"`
	MaxAlbumArtHeight                int                  `json:"MaxAlbumArtHeight"`
	MaxIconWidth                     int                  `json:"MaxIconWidth"`
	MaxIconHeight                    int                  `json:"MaxIconHeight"`
	MaxStreamingBitrate              int                  `json:"MaxStreamingBitrate"`
	MaxStaticBitrate                 int                  `json:"MaxStaticBitrate"`
	MusicStreamingTranscodingBitrate int                  `json:"MusicStreamingTranscodingBitrate"`
	MaxStaticMusicBitrate            int                  `json:"MaxStaticMusicBitrate"`
	SonyAggregationFlags             string               `json:"SonyAggregationFlags"`
	ProtocolInfo                     string               `json:"ProtocolInfo"`
	TimelineOffsetSeconds            int                  `json:"TimelineOffsetSeconds"`
	RequiresPlainVideoItems          bool                 `json:"RequiresPlainVideoItems"`
	RequiresPlainFolders             bool                 `json:"RequiresPlainFolders"`
	EnableMSMediaReceiverRegistrar   bool                 `json:"EnableMSMediaReceiverRegistrar"`
	IgnoreTranscodeByteRangeRequests bool                 `json:"IgnoreTranscodeByteRangeRequests"`
	XMLRootAttributes                []XMLRootAttributes  `json:"XmlRootAttributes"`
	DirectPlayProfiles               []DirectPlayProfiles `json:"DirectPlayProfiles"`
	ContainerProfiles                []ContainerProfiles  `json:"ContainerProfiles"`
	CodecProfiles                    []CodecProfiles      `json:"CodecProfiles"`
	ResponseProfiles                 []ResponseProfiles   `json:"ResponseProfiles"`
	SubtitleProfiles                 []SubtitleProfiles   `json:"SubtitleProfiles"`
}
type Capabilities struct {
	PlayableMediaTypes           []string      `json:"PlayableMediaTypes"`
	SupportedCommands            []string      `json:"SupportedCommands"`
	SupportsMediaControl         bool          `json:"SupportsMediaControl"`
	SupportsContentUploading     bool          `json:"SupportsContentUploading"`
	MessageCallbackURL           string        `json:"MessageCallbackUrl"`
	SupportsPersistentIdentifier bool          `json:"SupportsPersistentIdentifier"`
	SupportsSync                 bool          `json:"SupportsSync"`
	DeviceProfile                DeviceProfile `json:"DeviceProfile"`
	AppStoreURL                  string        `json:"AppStoreUrl"`
	IconURL                      string        `json:"IconUrl"`
}
