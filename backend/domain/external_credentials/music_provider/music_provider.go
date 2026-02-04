package music_provider

type MusicProvider int

const (
	YouTube MusicProvider = iota
	Spotify MusicProvider = iota
)

var musicProviderString = map[MusicProvider]string{
	YouTube: "youtube",
	Spotify: "spotify",
}

var musicProviderValue = map[string]MusicProvider{
	"youtube": YouTube,
	"spotify": Spotify,
}

func (mp MusicProvider) String() string {
	return musicProviderString[mp]
}

func ParseMusicProvider(s string) *MusicProvider {
	msicProvider, ok := musicProviderValue[s]
	if !ok {
		return nil
	}
	return &msicProvider
}
