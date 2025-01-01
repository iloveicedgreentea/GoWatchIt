package models

type EZBEQConfig struct {
	ID                            int64  `json:"-" db:"id"`
	AdjustMasterVolumeWithProfile bool   `json:"adjustmastervolumewithprofile" db:"adjust_master_volume_with_profile"`
	DenonIP                       string `json:"denonip" db:"denon_ip"`
	DenonPort                     string `json:"denonport" db:"denon_port"`
	DryRun                        bool   `json:"dryrun" db:"dry_run"`
	Enabled                       bool   `json:"enabled" db:"enabled"`
	EnableTVBEQ                   bool   `json:"enabletvbeq" db:"enable_tv_beq"`
	NotifyOnLoad                  bool   `json:"notifyonload" db:"notify_on_load"`
	NotifyOnUnLoad                bool   `json:"notifyonunload" db:"notify_on_unload"`
	PreferredAuthor               string `json:"preferredauthor" db:"preferred_author"`
	Slots                         []int  `json:"slots" db:"slots"` // Store as JSON string in DB
	StopPlexIfMismatch            bool   `json:"stopplexifmismatch" db:"stop_plex_if_mismatch"`
	Port                          string `json:"port" db:"port"`
	URL                           string `json:"url" db:"url"`
	Scheme                        string `json:"scheme" db:"scheme"`
	UseAVRCodecSearch             bool   `json:"useavrcodecsearch" db:"use_avr_codec_search"`
	AVRBrand                      string `json:"avrbrand" db:"avr_brand"`
	AVRURL                        string `json:"avrurl" db:"avr_url"`
	LooseEditionMatching          bool   `json:"looseeditionmatching" db:"loose_edition_matching"`
	SkipEditionMatching           bool   `json:"skipeditionmatching" db:"skip_edition_matching"`
}

type HomeAssistantConfig struct {
	ID                              int64  `json:"-" db:"id"`
	Enabled                         bool   `json:"enabled" db:"enabled"`
	RemoteEntityName                string `json:"remoteentityname" db:"remote_entity_name"`
	MediaPlayerEntityName           string `json:"mediaplayerentitynmae" db:"mediaplayer_entity_name"`
	Token                           string `json:"token" db:"token"`
	TriggerAspectRatioChangeOnEvent bool   `json:"triggeraspectratiochangeonevent" db:"trigger_aspect_ratio_change_on_event"`
	NotifyEndpointName              string `json:"notifyendpointname" db:"notify_endpoint_name"`
	URL                             string `json:"url" db:"url"`
	Port                            string `json:"port" db:"port"`
	Scheme                          string `json:"scheme" db:"scheme"`
}

type JellyfinConfig struct {
	ID               int64  `json:"-" db:"id"`
	APIToken         string `json:"apitoken" db:"api_token"`
	DeviceUUIDFilter string `json:"deviceuuidfilter" db:"device_uuid_filter"`
	Enabled          bool   `json:"enabled" db:"enabled"`
	OwnerNameFilter  string `json:"ownernamefilter" db:"owner_name_filter"`
	UserID           string `json:"userID" db:"user_id"`
	URL              string `json:"url" db:"url"`
	Port             string `json:"port" db:"port"`
	Scheme           string `json:"scheme" db:"scheme"`
	SkipTMDB         bool   `json:"skiptmdb" db:"skip_tmdb"`
}

type MainConfig struct {
	ID         int64  `json:"-" db:"id"`
	ListenPort string `json:"listenport" db:"listen_port"`
}

type PlexConfig struct {
	ID                   int64  `json:"-" db:"id"`
	DeviceUUIDFilter     string `json:"deviceuuidfilter" db:"device_uuid_filter"`
	Enabled              bool   `json:"enabled" db:"enabled"`
	EnableTrailerSupport bool   `json:"enabletrailersupport" db:"enable_trailer_support"`
	OwnerNameFilter      string `json:"ownernamefilter" db:"owner_name_filter"`
	Token                string `json:"token" db:"token"`
	URL                  string `json:"url" db:"url"`
	Port                 string `json:"port" db:"port"`
	Scheme               string `json:"scheme" db:"scheme"`
}

type HDMISyncConfig struct {
	ID                      int64  `json:"-" db:"id"`
	Enabled                 bool   `json:"enabled" db:"enabled"`
	Source                  string `json:"source" db:"source"`
	Time                    string `json:"time" db:"time"`
	Envy                    string `json:"envy" db:"envy"`
	PlayerIP                string `json:"playerip" db:"player_ip"`
	PlayerMachineIdentifier string `json:"playermachineidentifier" db:"player_machine_identifier"`
	Scheme                  string `json:"scheme" db:"scheme"`
}
