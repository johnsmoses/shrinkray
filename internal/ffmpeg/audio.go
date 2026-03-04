package ffmpeg

import "strings"

// mp4CompatibleAudioCodecs lists audio codecs that can be stream-copied into MP4 containers
// without re-encoding. Codecs not in this list must be transcoded (e.g., to AAC) when
// remuxing to MP4, as the muxer will silently drop or error on incompatible streams.
var mp4CompatibleAudioCodecs = map[string]bool{
	"aac":  true, // Most common; native MP4/M4A audio
	"mp3":  true, // MPEG-1/2 audio layer III; supported in MP4
	"ac3":  true, // Dolby Digital; supported in MP4
	"eac3": true, // Dolby Digital Plus; supported in MP4
	"alac": true, // Apple Lossless; native in M4A/MP4
}

// IsMP4CompatibleAudio returns true if the audio codec can be stream-copied to an MP4 container.
func IsMP4CompatibleAudio(codecName string) bool {
	return mp4CompatibleAudioCodecs[strings.ToLower(strings.TrimSpace(codecName))]
}

// IncompatibleMP4AudioCodecs returns the unique codec names from the given audio streams
// that cannot be stream-copied to MP4 (e.g., DTS, TrueHD, PCM variants).
// Returns nil if all streams are compatible or the slice is empty.
func IncompatibleMP4AudioCodecs(streams []AudioStream) []string {
	if len(streams) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var incompatible []string
	for _, s := range streams {
		codec := strings.ToLower(strings.TrimSpace(s.CodecName))
		if !mp4CompatibleAudioCodecs[codec] && !seen[codec] {
			seen[codec] = true
			incompatible = append(incompatible, s.CodecName)
		}
	}
	return incompatible
}
