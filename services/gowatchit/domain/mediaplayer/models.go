package mediaplayer

type Status struct{}

// Codec represents the audio codec used in the media player
// Codec is based on EzBEQ/BEQ strings
// All users need to normalize to these in their own GetCodec mappers
type Codec string

const (
	// Dolby
	CodecDDP      Codec = "ddp"
	CodecDDPlus   Codec = "DD+"
	CodecEAC3     Codec = "EAC3"
	CodecEAC351   Codec = "EAC3 5.1"
	CodecEAC3Alt  Codec = "e-ac3"
	CodecAtmos    Codec = "Atmos"
	CodecDDPAtmos Codec = "DD+ Atmos"

	// Truehd
	CodecTrueHD71   Codec = "TrueHD 7.1"
	CodecTrueHD51   Codec = "TrueHD 5.1"
	CodecTrueHD61   Codec = "TrueHD 6.1"
	CodecSurround71 Codec = "Surround 7.1"

	// DTS
	CodecDTSX      Codec = "DTS-X"
	CodecDTSXAlt   Codec = "DTS:X"
	CodecDTSHDMA71 Codec = "DTS-HD MA 7.1"
	CodecDTSHDMA51 Codec = "DTS-HD MA 5.1"
	CodecDTSHDHR51 Codec = "DTS-HD HR 5.1"
	CodecDTSHDHR71 Codec = "DTS-HD HR 7.1"
	CodecDTS51     Codec = "DTS 5.1"

	// PCM
	CodecLPCM51 Codec = "LPCM 5.1"
	CodecLPCM71 Codec = "LPCM 7.1"
	CodecLPCM20 Codec = "LPCM 2.0"

	// Other
	CodecAAC20     Codec = "AAC 2.0"
	CodecAACStereo Codec = "AAC Stereo"
	CodecAC351     Codec = "AC3 5.1"

	// Maybe flags
	CodecDDPlusAtmos51Maybe Codec = "DD+Atmos5.1Maybe"
	CodecDDPlusAtmos71Maybe Codec = "DD+Atmos7.1Maybe"
	CodecAtmosMaybe         Codec = "AtmosMaybe"
	CodecDDPlus51           Codec = "DD+ 5.1"
	CodecDDPlus71           Codec = "DD+ 7.1"
)

// Edition represents the edition of the media being played
// Edition is based on EzBEQ/BEQ strings
// All users need to normalize to these in their own GetEdition mappers
type Edition string

const (
	EditionExtended       Edition = "Extended"
	EditionUnrated        Edition = "Unrated"
	EditionTheatrical     Edition = "Theatrical"
	EditionSpecialEdition Edition = "Special"
	EditionUltimate       Edition = "Ultimate"
	EditionDirectorsCut   Edition = "Director"
	EditionCriterion      Edition = "Criterion"
	EditionUnknown        Edition = "Unknown"
	EditionNone           Edition = "None"
)

type MetaData struct{}
