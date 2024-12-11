package models

type CodecName string

const (
	CodecDDP      CodecName = "ddp"
	CodecDDPlus   CodecName = "DD+"
	CodecEAC3     CodecName = "EAC3"
	CodecEAC351   CodecName = "EAC3 5.1"
	CodecEAC3Alt  CodecName = "e-ac3"
	CodecAtmos    CodecName = "Atmos"
	CodecDDPAtmos CodecName = "DD+ Atmos"

	CodecTrueHD71   CodecName = "TrueHD 7.1"
	CodecTrueHD51   CodecName = "TrueHD 5.1"
	CodecTrueHD61   CodecName = "TrueHD 6.1"
	CodecSurround71 CodecName = "Surround 7.1"

	CodecDTSX      CodecName = "DTS-X"
	CodecDTSXAlt   CodecName = "DTS:X"
	CodecDTSHDMA71 CodecName = "DTS-HD MA 7.1"
	CodecDTSHDMA51 CodecName = "DTS-HD MA 5.1"
	CodecDTSHDHR51 CodecName = "DTS-HD HR 5.1"
	CodecDTSHDHR71 CodecName = "DTS-HD HR 7.1"

	CodecLPCM51 CodecName = "LPCM 5.1"
	CodecLPCM71 CodecName = "LPCM 7.1"
	CodecLPCM20 CodecName = "LPCM 2.0"

	CodecAAC20 CodecName = "AAC 2.0"
	CodecAC351 CodecName = "AC3 5.1"

	CodecDTS51 CodecName = "DTS 5.1"

	CodecAACStereo CodecName = "AAC Stereo"

	// Maybe flags
	CodecDDPAtmos5Maybe CodecName = "DD+Atmos5.1Maybe"
	CodecDDPAtmos7Maybe CodecName = "DD+Atmos7.1Maybe"
	CodecAtmosMaybe     CodecName = "AtmosMaybe"
)
