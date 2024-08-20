package models

type EZBEQConfig struct {
	ID                            int64  `json:"-" db:"id"`
	AdjustMasterVolumeWithProfile bool   `json:"adjustmastervolumewithprofile" db:"adjust_master_volume_with_profile"`
	DenonIP                       string `json:"denonip" db:"denon_ip"`
	DenonPort                     string `json:"denonport" db:"denon_port"`
	DryRun                        bool   `json:"dryrun" db:"dry_run"`
	Enabled                       bool   `json:"enabled" db:"enabled"`
	EnableTVBEQ                   bool   `json:"enabletvbeq" db:"enable_tv_beq"`
	NotifyEndpointName            string `json:"notifyendpointname" db:"notify_endpoint_name"`
	NotifyOnLoad                  bool   `json:"notifyonload" db:"notify_on_load"`
	Port                          string `json:"port" db:"port"`
	PreferredAuthor               string `json:"preferredauthor" db:"preferred_author"`
	Slots                         []int `json:"slots" db:"slots"` // Store as JSON string in DB
	StopPlexIfMismatch            bool   `json:"stopplexifmismatch" db:"stop_plex_if_mismatch"`
	URL                           string `json:"url" db:"url"`
	UseAVRCodecSearch             bool   `json:"useavrcodecsearch" db:"use_avr_codec_search"`
	AVRBrand                      string `json:"avrbrand" db:"avr_brand"`
	AVRURL                        string `json:"avrurl" db:"avr_url"`
}

type HomeAssistantConfig struct {
	ID                                  int64  `json:"-" db:"id"`
	Enabled                             bool   `json:"enabled" db:"enabled"`
	PauseScriptName                     string `json:"pausescriptname" db:"pause_script_name"`
	PlayScriptName                      string `json:"playscriptname" db:"play_script_name"`
	Port                                string `json:"port" db:"port"`
	RemoteEntityName                    string `json:"remoteentityname" db:"remote_entity_name"`
	StopScriptName                      string `json:"stopscriptname" db:"stop_script_name"`
	Token                               string `json:"token" db:"token"`
	TriggerAspectRatioChangeOnEvent     bool   `json:"triggeraspectratiochangeonevent" db:"trigger_aspect_ratio_change_on_event"`
	TriggerAVRMasterVolumeChangeOnEvent bool   `json:"triggeravrmastervolumechangeonevent" db:"trigger_avr_master_volume_change_on_event"`
	TriggerLightsOnEvent                bool   `json:"triggerlightsonevent" db:"trigger_lights_on_event"`
	URL                                 string `json:"url" db:"url"`
}

type JellyfinConfig struct {
	ID               int64  `json:"-" db:"id"`
	APIToken         string `json:"apitoken" db:"api_token"`
	DeviceUUIDFilter string `json:"deviceuuidfilter" db:"device_uuid_filter"`
	Enabled          bool   `json:"enabled" db:"enabled"`
	OwnerNameFilter  string `json:"ownernamefilter" db:"owner_name_filter"`
	UserID           string `json:"userID" db:"user_id"`
	Port             string `json:"port" db:"port"`
	URL              string `json:"url" db:"url"`
	SkipTMDB         bool   `json:"skiptmdb" db:"skip_tmdb"`
}

type MainConfig struct {
	ID         int64  `json:"-" db:"id"`
	ListenPort string `json:"listenport" db:"listen_port"`
}

type MQTTConfig struct {
	ID                        int64  `json:"-" db:"id"`
	Enabled                   bool   `json:"enabled" db:"enabled"`
	Password                  string `json:"password" db:"password"`
	TopicAspectRatio          string `json:"topicaspectratio" db:"topic_aspect_ratio"`
	TopicAspectRatioMADVROnly string `json:"topicaspectratiomadvronly" db:"topic_aspect_ratio_madvr_only"`
	TopicBEQCurrentProfile    string `json:"topicbeqcurrentprofile" db:"topic_beq_current_profile"`
	TopicLights               string `json:"topiclights" db:"topic_lights"`
	TopicMiniDSPMuteStatus    string `json:"topicminidspmutestatus" db:"topic_minidsp_mute_status"`
	TopicPlayingStatus        string `json:"topicplayingstatus" db:"topic_playing_status"`
	TopicVolume               string `json:"topicvolume" db:"topic_volume"`
	URL                       string `json:"url" db:"url"`
	Username                  string `json:"username" db:"username"`
}

type PlexConfig struct {
	ID                   int64  `json:"-" db:"id"`
	DeviceUUIDFilter     string `json:"deviceuuidfilter" db:"device_uuid_filter"`
	Enabled              bool   `json:"enabled" db:"enabled"`
	EnableTrailerSupport bool   `json:"enabletrailersupport" db:"enable_trailer_support"`
	OwnerNameFilter      string `json:"ownernamefilter" db:"owner_name_filter"`
	Port                 string `json:"port" db:"port"`
	Token                string `json:"token" db:"token"`
	URL                  string `json:"url" db:"url"`
}

type HDMISyncConfig struct {
	ID                      int64  `json:"-" db:"id"`
	Enabled                 bool   `json:"enabled" db:"enabled"`
	Source                  string `json:"source" db:"source"`
	Time                    string `json:"time" db:"time"`
	Envy                    string `json:"envy" db:"envy"`
	PlayerIP                string `json:"playerip" db:"player_ip"`
	PlayerMachineIdentifier string `json:"playermachineidentifier" db:"player_machine_identifier"`
}
