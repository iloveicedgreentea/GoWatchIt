package models

import (
	"encoding/xml"
)

type MediaContainer struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Size                string   `xml:"size,attr"`
	AllowSync           string   `xml:"allowSync,attr"`
	Identifier          string   `xml:"identifier,attr"`
	LibrarySectionID    string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string   `xml:"mediaTagVersion,attr"`
	Video               struct {
		RatingKey             string `xml:"ratingKey,attr"`
		Key                   string `xml:"key,attr"`
		ParentRatingKey       string `xml:"parentRatingKey,attr"`
		GrandparentRatingKey  string `xml:"grandparentRatingKey,attr"`
		AttrGuid              string `xml:"guid,attr"`
		ParentGuid            string `xml:"parentGuid,attr"`
		GrandparentGuid       string `xml:"grandparentGuid,attr"`
		EditionTitle string `xml:"editionTitle,attr"`
		Type                  string `xml:"type,attr"`
		Title                 string `xml:"title,attr"`
		GrandparentKey        string `xml:"grandparentKey,attr"`
		ParentKey             string `xml:"parentKey,attr"`
		LibrarySectionTitle   string `xml:"librarySectionTitle,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		LibrarySectionKey     string `xml:"librarySectionKey,attr"`
		GrandparentTitle      string `xml:"grandparentTitle,attr"`
		ParentTitle           string `xml:"parentTitle,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Summary               string `xml:"summary,attr"`
		Index                 string `xml:"index,attr"`
		ParentIndex           string `xml:"parentIndex,attr"`
		AudienceRating        string `xml:"audienceRating,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Art                   string `xml:"art,attr"`
		ParentThumb           string `xml:"parentThumb,attr"`
		GrandparentThumb      string `xml:"grandparentThumb,attr"`
		GrandparentArt        string `xml:"grandparentArt,attr"`
		GrandparentTheme      string `xml:"grandparentTheme,attr"`
		Duration              string `xml:"duration,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		AddedAt               string `xml:"addedAt,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		AudienceRatingImage   string `xml:"audienceRatingImage,attr"`
		Media                 struct {
			ID              string `xml:"id,attr"`
			Duration        string `xml:"duration,attr"`
			Bitrate         string `xml:"bitrate,attr"`
			Width           string `xml:"width,attr"`
			Height          string `xml:"height,attr"`
			AspectRatio     string `xml:"aspectRatio,attr"`
			AudioChannels   string `xml:"audioChannels,attr"`
			AudioCodec      string `xml:"audioCodec,attr"`
			VideoCodec      string `xml:"videoCodec,attr"`
			VideoResolution string `xml:"videoResolution,attr"`
			Container       string `xml:"container,attr"`
			VideoFrameRate  string `xml:"videoFrameRate,attr"`
			AudioProfile    string `xml:"audioProfile,attr"`
			VideoProfile    string `xml:"videoProfile,attr"`
			Part            struct {
				ID           string `xml:"id,attr"`
				Key          string `xml:"key,attr"`
				Duration     string `xml:"duration,attr"`
				File         string `xml:"file,attr"`
				Size         string `xml:"size,attr"`
				AudioProfile string `xml:"audioProfile,attr"`
				Container    string `xml:"container,attr"`
				VideoProfile string `xml:"videoProfile,attr"`
				Stream       []struct {
					ID                   string `xml:"id,attr"`
					StreamType           string `xml:"streamType,attr"`
					Default              string `xml:"default,attr"`
					Codec                string `xml:"codec,attr"`
					Index                string `xml:"index,attr"`
					Bitrate              string `xml:"bitrate,attr"`
					BitDepth             string `xml:"bitDepth,attr"`
					ChromaLocation       string `xml:"chromaLocation,attr"`
					ChromaSubsampling    string `xml:"chromaSubsampling,attr"`
					CodedHeight          string `xml:"codedHeight,attr"`
					CodedWidth           string `xml:"codedWidth,attr"`
					ColorRange           string `xml:"colorRange,attr"`
					FrameRate            string `xml:"frameRate,attr"`
					Height               string `xml:"height,attr"`
					Level                string `xml:"level,attr"`
					Profile              string `xml:"profile,attr"`
					RefFrames            string `xml:"refFrames,attr"`
					Width                string `xml:"width,attr"`
					DisplayTitle         string `xml:"displayTitle,attr"`
					ExtendedDisplayTitle string `xml:"extendedDisplayTitle,attr"`
					Selected             string `xml:"selected,attr"`
					Channels             string `xml:"channels,attr"`
					Language             string `xml:"language,attr"`
					LanguageTag          string `xml:"languageTag,attr"`
					LanguageCode         string `xml:"languageCode,attr"`
					AudioChannelLayout   string `xml:"audioChannelLayout,attr"`
					SamplingRate         string `xml:"samplingRate,attr"`
					Title         string `xml:"title,attr"`
				} `xml:"Stream"`
			} `xml:"Part"`
		} `xml:"Media"`
		Director struct {
			ID     string `xml:"id,attr"`
			Filter string `xml:"filter,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Director"`
		Writer struct {
			ID     string `xml:"id,attr"`
			Filter string `xml:"filter,attr"`
			Tag    string `xml:"tag,attr"`
		} `xml:"Writer"`
		Guid []struct {
			ID   string `xml:"id,attr"`
		} `xml:"Guid"`
		Rating struct {
			Image string `xml:"image,attr"`
			Value string `xml:"value,attr"`
			Type  string `xml:"type,attr"`
		} `xml:"Rating"`
		Role []struct {
			ID     string `xml:"id,attr"`
			Filter string `xml:"filter,attr"`
			Tag    string `xml:"tag,attr"`
			TagKey string `xml:"tagKey,attr"`
			Role   string `xml:"role,attr"`
			Thumb  string `xml:"thumb,attr"`
		} `xml:"Role"`
	} `xml:"Video"`
} 

type AllMediaContainer struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Text                string   `xml:",chardata"`
	Size                string   `xml:"size,attr"`
	AllowSync           string   `xml:"allowSync,attr"`
	Art                 string   `xml:"art,attr"`
	Identifier          string   `xml:"identifier,attr"`
	LibrarySectionID    string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string   `xml:"mediaTagVersion,attr"`
	Thumb               string   `xml:"thumb,attr"`
	Title1              string   `xml:"title1,attr"`
	Title2              string   `xml:"title2,attr"`
	ViewGroup           string   `xml:"viewGroup,attr"`
	ViewMode            string   `xml:"viewMode,attr"`
	Video               []struct {
		Text                  string `xml:",chardata"`
		RatingKey             string `xml:"ratingKey,attr"`
		Key                   string `xml:"key,attr"`
		Guid                  string `xml:"guid,attr"`
		Studio                string `xml:"studio,attr"`
		Type                  string `xml:"type,attr"`
		Title                 string `xml:"title,attr"`
		ContentRating         string `xml:"contentRating,attr"`
		Summary               string `xml:"summary,attr"`
		Rating                string `xml:"rating,attr"`
		AudienceRating        string `xml:"audienceRating,attr"`
		ViewCount             string `xml:"viewCount,attr"`
		SkipCount             string `xml:"skipCount,attr"`
		LastViewedAt          string `xml:"lastViewedAt,attr"`
		Year                  string `xml:"year,attr"`
		Tagline               string `xml:"tagline,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Art                   string `xml:"art,attr"`
		Duration              string `xml:"duration,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		AddedAt               string `xml:"addedAt,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		AudienceRatingImage   string `xml:"audienceRatingImage,attr"`
		ChapterSource         string `xml:"chapterSource,attr"`
		PrimaryExtraKey       string `xml:"primaryExtraKey,attr"`
		RatingImage           string `xml:"ratingImage,attr"`
		ViewOffset            string `xml:"viewOffset,attr"`
		TitleSort             string `xml:"titleSort,attr"`
		OriginalTitle         string `xml:"originalTitle,attr"`
		UserRating            string `xml:"userRating,attr"`
		LastRatedAt           string `xml:"lastRatedAt,attr"`
		Media                 []struct {
			Text                  string `xml:",chardata"`
			ID                    string `xml:"id,attr"`
			Duration              string `xml:"duration,attr"`
			Bitrate               string `xml:"bitrate,attr"`
			Width                 string `xml:"width,attr"`
			Height                string `xml:"height,attr"`
			AspectRatio           string `xml:"aspectRatio,attr"`
			AudioChannels         string `xml:"audioChannels,attr"`
			AudioCodec            string `xml:"audioCodec,attr"`
			VideoCodec            string `xml:"videoCodec,attr"`
			VideoResolution       string `xml:"videoResolution,attr"`
			Container             string `xml:"container,attr"`
			VideoFrameRate        string `xml:"videoFrameRate,attr"`
			AudioProfile          string `xml:"audioProfile,attr"`
			VideoProfile          string `xml:"videoProfile,attr"`
			OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
			Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
			Part                  struct {
				Text                  string `xml:",chardata"`
				ID                    string `xml:"id,attr"`
				Key                   string `xml:"key,attr"`
				Duration              string `xml:"duration,attr"`
				File                  string `xml:"file,attr"`
				Size                  string `xml:"size,attr"`
				AudioProfile          string `xml:"audioProfile,attr"`
				Container             string `xml:"container,attr"`
				VideoProfile          string `xml:"videoProfile,attr"`
				HasThumbnail          string `xml:"hasThumbnail,attr"`
				Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
				OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
			} `xml:"Part"`
		} `xml:"Media"`
		Genre []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Genre"`
		Director []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Director"`
		Writer []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Writer"`
		Country []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Country"`
		Collection struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Collection"`
		Role []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Role"`
	} `xml:"Video"`
} 

type SessionMediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Text    string   `xml:",chardata"`
	Size    string   `xml:"size,attr"`
	Video   []struct {
		Text                  string `xml:",chardata"`
		AddedAt               string `xml:"addedAt,attr"`
		Art                   string `xml:"art,attr"`
		CreatedAtAccuracy     string `xml:"createdAtAccuracy,attr"`
		CreatedAtTZOffset     string `xml:"createdAtTZOffset,attr"`
		Duration              string `xml:"duration,attr"`
		Guid                  string `xml:"guid,attr"`
		Key                   string `xml:"key,attr"`
		LibrarySectionID      string `xml:"librarySectionID,attr"`
		LibrarySectionKey     string `xml:"librarySectionKey,attr"`
		LibrarySectionTitle   string `xml:"librarySectionTitle,attr"`
		OriginallyAvailableAt string `xml:"originallyAvailableAt,attr"`
		RatingKey             string `xml:"ratingKey,attr"`
		SessionKey            string `xml:"sessionKey,attr"`
		SkipCount             string `xml:"skipCount,attr"`
		Subtype               string `xml:"subtype,attr"`
		Thumb                 string `xml:"thumb,attr"`
		Title                 string `xml:"title,attr"`
		Type                  string `xml:"type,attr"`
		UpdatedAt             string `xml:"updatedAt,attr"`
		ViewOffset            string `xml:"viewOffset,attr"`
		Year                  string `xml:"year,attr"`
		Media                 struct {
			Text                  string `xml:",chardata"`
			ID                    string `xml:"id,attr"`
			VideoProfile          string `xml:"videoProfile,attr"`
			AudioChannels         string `xml:"audioChannels,attr"`
			AudioCodec            string `xml:"audioCodec,attr"`
			Container             string `xml:"container,attr"`
			Duration              string `xml:"duration,attr"`
			Height                string `xml:"height,attr"`
			OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
			Protocol              string `xml:"protocol,attr"`
			VideoCodec            string `xml:"videoCodec,attr"`
			VideoFrameRate        string `xml:"videoFrameRate,attr"`
			VideoResolution       string `xml:"videoResolution,attr"`
			Width                 string `xml:"width,attr"`
			Selected              string `xml:"selected,attr"`
			Part                  struct {
				Text                  string `xml:",chardata"`
				ID                    string `xml:"id,attr"`
				VideoProfile          string `xml:"videoProfile,attr"`
				Container             string `xml:"container,attr"`
				Duration              string `xml:"duration,attr"`
				Height                string `xml:"height,attr"`
				OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
				Protocol              string `xml:"protocol,attr"`
				Width                 string `xml:"width,attr"`
				Decision              string `xml:"decision,attr"`
				Selected              string `xml:"selected,attr"`
				Stream                []struct {
					Text                 string `xml:",chardata"`
					Bitrate              string `xml:"bitrate,attr"`
					Codec                string `xml:"codec,attr"`
					Default              string `xml:"default,attr"`
					DisplayTitle         string `xml:"displayTitle,attr"`
					ExtendedDisplayTitle string `xml:"extendedDisplayTitle,attr"`
					FrameRate            string `xml:"frameRate,attr"`
					Height               string `xml:"height,attr"`
					ID                   string `xml:"id,attr"`
					Language             string `xml:"language,attr"`
					LanguageCode         string `xml:"languageCode,attr"`
					LanguageTag          string `xml:"languageTag,attr"`
					StreamType           string `xml:"streamType,attr"`
					Width                string `xml:"width,attr"`
					Decision             string `xml:"decision,attr"`
					Location             string `xml:"location,attr"`
					BitrateMode          string `xml:"bitrateMode,attr"`
					Channels             string `xml:"channels,attr"`
					Selected             string `xml:"selected,attr"`
				} `xml:"Stream"`
			} `xml:"Part"`
		} `xml:"Media"`
		User struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Thumb string `xml:"thumb,attr"`
			Title string `xml:"title,attr"`
		} `xml:"User"`
		Player struct {
			Text              string `xml:",chardata"`
			Address           string `xml:"address,attr"`
			Device            string `xml:"device,attr"`
			MachineIdentifier string `xml:"machineIdentifier,attr"`
			Model             string `xml:"model,attr"`
			Platform          string `xml:"platform,attr"`
			PlatformVersion   string `xml:"platformVersion,attr"`
			Product           string `xml:"product,attr"`
			Profile           string `xml:"profile,attr"`
			State             string `xml:"state,attr"`
			Title             string `xml:"title,attr"`
			Version           string `xml:"version,attr"`
			Local             string `xml:"local,attr"`
			Relayed           string `xml:"relayed,attr"`
			Secure            string `xml:"secure,attr"`
			UserID            string `xml:"userID,attr"`
		} `xml:"Player"`
		Session struct {
			Text      string `xml:",chardata"`
			ID        string `xml:"id,attr"`
			Bandwidth string `xml:"bandwidth,attr"`
			Location  string `xml:"location,attr"`
		} `xml:"Session"`
		TranscodeSession struct {
			Text                     string `xml:",chardata"`
			Key                      string `xml:"key,attr"`
			Throttled                string `xml:"throttled,attr"`
			Complete                 string `xml:"complete,attr"`
			Progress                 string `xml:"progress,attr"`
			Size                     string `xml:"size,attr"`
			Speed                    string `xml:"speed,attr"`
			Error                    string `xml:"error,attr"`
			Duration                 string `xml:"duration,attr"`
			Remaining                string `xml:"remaining,attr"`
			Context                  string `xml:"context,attr"`
			SourceVideoCodec         string `xml:"sourceVideoCodec,attr"`
			SourceAudioCodec         string `xml:"sourceAudioCodec,attr"`
			VideoDecision            string `xml:"videoDecision,attr"`
			AudioDecision            string `xml:"audioDecision,attr"`
			Protocol                 string `xml:"protocol,attr"`
			Container                string `xml:"container,attr"`
			VideoCodec               string `xml:"videoCodec,attr"`
			AudioCodec               string `xml:"audioCodec,attr"`
			AudioChannels            string `xml:"audioChannels,attr"`
			TranscodeHwRequested     string `xml:"transcodeHwRequested,attr"`
			TranscodeHwEncoding      string `xml:"transcodeHwEncoding,attr"`
			TranscodeHwEncodingTitle string `xml:"transcodeHwEncodingTitle,attr"`
			TimeStamp                string `xml:"timeStamp,attr"`
			MaxOffsetAvailable       string `xml:"maxOffsetAvailable,attr"`
			MinOffsetAvailable       string `xml:"minOffsetAvailable,attr"`
		} `xml:"TranscodeSession"`
	} `xml:"Video"`
} 