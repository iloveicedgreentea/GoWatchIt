
## How BEQ Support Works
On play and resume, it will load the profile. On pause and stop, it will unload it (so you don't forget to). It has some logic to cache the profile so if you pause and unpause, the profile will get loaded much faster as it skips searching the DB and stuff. 

If enabled, it will also send a notification to Home Assistant via Notify so you can send an alert to your phone for example. 

For safety, the application tries to unload the profile when it loads up each time in case it crashed or was killed previously, and will unload before playing anything so it doesn't start playing something with the wrong profile. 

### Matching
The application will search the catalog and match based on codec (Atmos, DTS-X, etc), title, year, TMDB, and edition. I have tested with multiple titles and everything matched as expected.

> ⚠️ *If you get an incorrect match, please open a github issue with the full log output and expected codec and title*

Jellyfin may have some issues matching as I have found it will sometimes just not return a TMDB. This has nothing to do with me. Jellyfin is generally just quite buggy. There is a configuration option that you should probably enable in the Jellyfin section which lets you skip TMDB matching. It will instead use the title name which could be prone to false negatives. 

### Editions

This application will do its best to match editions. It will look for one of the following:
1) Plex edition metadata. Set this from your server in the Plex UI
2) Looking at the file name if it contains `Unrated, Ultimate, Theatrical, Extended, Director, Criterion` (case insensitive)

There is no other reliable way to get the edition. If an edition is not matched, BEQ will fail to load for safety reasons (different editions have different masterings, etc). If a BEQCatalog entry has a blank edition, then edition will not matter and it will match based on the usual criteria.

If you want to match an edition when none is found locally, you can enable Loose Edition Matching. It will match if:

1) We sent a blank edition in the request
2) BEQ catalogue returns an edition

This is only useful if you have issues matching editions.


### Audio stuff
Here are some examples of what kind of codec tags Plex will have based on file metadata

 TrueHD 7.1
Unknown (TRUEHD 7.1) --- Surround 7.1 (TRUEHD) 
English (TRUEHD 7.1) --- TrueHD 7.1 (English) 
English (TRUEHD 7.1) --- Surround 7.1 (English TRUEHD) 
English (TRUEHD 7.1) --- English (TRUEHD 7.1) 

 Atmos
English (TRUEHD 7.1) --- TrueHD Atmos 7.1 Remixed (English) 
English (TRUEHD 7.1) --- TrueHD Atmos 7.1 (English) 
English (TRUEHD 7.1) --- TrueHD 7.1. Atmos (English 7.1) 
English (TRUEHD 7.1) --- TrueHD 7.1 Atmos (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1+13 objects / 48 kHz / 4691 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1+11 objects / 48 kHz / 4309 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 5026 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 4396 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 4353 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 3258 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby TrueHD/Atmos Audio / 7.1 / 48 kHz / 3041 kbps / 16-bit (English) 
English (TRUEHD 7.1) --- Dolby Atmos/TrueHD Audio / 7.1-Atmos / 48 kHz / 5826 kbps / 24-bit (English) 
English (TRUEHD 7.1) --- Dolby Atmos/TrueHD Audio / 7.1-Atmos / 48 kHz / 4535 kbps / 24-bit (English) 

 TrueHD 5.1
English (TRUEHD 5.1) --- English Dolby TrueHD (5.1) 
English (TRUEHD 5.1) --- English (TRUEHD 5.1) 
English (TRUEHD 5.1) --- Dolby TrueHD Audio / 5.1 / 48 kHz / 4130 kbps / 24-bit (English) 
English (TRUEHD 5.1) --- Dolby TrueHD Audio / 5.1 / 48 kHz / 1522 kbps / 16-bit (English) 
English (TRUEHD 5.1) --- TrueHD 5.1 (English) 
English (TRUEHD 5.1) --- Dolby TrueHD 5.1 (English) 

 DTS:X 
English (DTS-HD MA 7.1) --- DTS:X/DTS-HD Master Audio / 7.1-X / 48 kHz / 4458 kbps / 24-bit (English) 
English (DTS-HD MA 7.1) --- DTS:X/DTS-HD Master Audio / 7.1-X / 48 kHz / 4255 kbps / 24-bit (English) 
English (DTS-HD MA 7.1) --- DTS:X 7.1 (English DTS-HD MA) 

 DTS-HD MA 7.1
English (DTS-HD MA 7.1) --- Surround 7.1 (English DTS-HD MA) 
English (DTS-HD MA 7.1) --- DTS-HD MA 7.1 (English) 

 DTS-MA 6.1
中文 (DTS-HD MA 6.1) --- DTS-HD Master Audio / 6.1 / 48 kHz / 4667 kbps / 24-bit (中文) 

 DTS-HD MA 5.1
中文 (DTS-HD MA 5.1) --- Mandarin DTS-HD MA 5.1 (中文) 
中文 (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 2360 kbps / 16-bit (中文) 
Indonesia (DTS-HD MA 5.1) --- DTS-HD Master Audio / Indonesian / 5.1 / 48 kHz / 3531 kbps / 24-bit (Indonesia) 
Indonesia (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3576 kbps / 24-bit (Indonesia) 
English (DTS-HD MA 5.1) --- Surround 5.1 (English DTS-HD MA) 
English (DTS-HD MA 5.1) --- English / DTS-HD Master Audio / 5.1 / 48 kHz / 4104 kbps / 24-bit (DTS Core: 5.1 / 48 kHz / 1509 kbps / 24-bit) 
English (DTS-HD MA 5.1) --- English / DTS-HD Master Audio / 5.1 / 48 kHz / 2688 kbps / 24-bit 
English (DTS-HD MA 5.1) --- English (DTS-HD MA 5.1) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio English 2877 kbps 5.1 / 48 kHz / 2877 kbps / 24-bit 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / English / 3336 kbps / 5.1 Channels / 48 kHz / 24-bit (DTS Core: 5.1 Channels / 48 kHz / 1509 kbps / 24-bit) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48kHz / 5128 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 4189 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 4107 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3900 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3746 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3600 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3596 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD Master Audio / 5.1 / 48 kHz / 3233 kbps / 24-bit (English) 
English (DTS-HD MA 5.1) --- DTS-HD MA 5.1 (English) 
English (DTS-HD MA 5.1) --- DTS-HD MA @ 1509 kbps (English 5.1) 
English (DTS-HD MA 5.1) --- DTS HD MA 5.1 (English) 
English (DTS-HD MA 5.1) --- DTS (MA) / 2181 kbps / 48 KHz / 24-Bit / 5.1 (English) 

 DTS HD HRA 7.1
English (DTS-HD HRA 7.1) --- Surround 7.1 (English DTS-HD HRA) 

 DTS 5.1
English (DTS 5.1) --- English (DTS 5.1) 
English (DTS 5.1) --- DTS 5.1 (English) 

Unknown (AAC Stereo) --- Unknown (AAC Stereo) 

 AC3 stereo
English (AC3 Stereo) --- English (AC3 Stereo) 

 AC3
English (AC3 5.1) --- Surround (English AC3 5.1) 
English (AC3 5.1) --- English (AC3 5.1) 
English (AC3 5.1) --- AC3 5.1 @ 640 Kbps (English) 

 Misc
English (PCM Mono) --- Mono (English PCM) 
English (FLAC Stereo) --- Original Dolby Stereo (Laserdisc USA LD68993) (English FLAC) 
English (FLAC Stereo) --- FLAC Audio / 1266 kbps / 2.0 / 48 kHz / 24-bit (English) 

English (FLAC 5.1) --- Main Audio (English FLAC 5.1) 
English (FLAC 5.1) --- FLAC 5.1 @ 2954 kbps / 24-bit (English) 
English (FLAC 5.1) --- English (FLAC 5.1) 

English (EAC3 Stereo) --- English (EAC3 Stereo) 
English (EAC3 5.1) --- Main Audio (English EAC3 5.1) 
English (EAC3 5.1) --- English (EAC3 5.1) 