package fixture

type FixturesMap struct {
	AppleTracks       map[string][]byte
	AppleAlbums       map[string][]byte
	AppleSearchTracks map[string][]byte
	AppleSearchAlbums map[string][]byte

	SpotifyTracks       map[string][]byte
	SpotifyAlbums       map[string][]byte
	SpotifySearchTracks map[string][]byte
	SpotifySearchAlbums map[string][]byte

	YandexTracks       map[string][]byte
	YandexAlbums       map[string][]byte
	YandexSearchTracks map[string][]byte
	YandexSearchAlbums map[string][]byte

	YoutubeTracks       map[string][]byte
	YoutubeAlbums       map[string][]byte
	YoutubeSearchTracks map[string][]byte
	YoutubeSearchAlbums map[string][]byte
}

func (fm *FixturesMap) Merge(mergeFm *FixturesMap) {
	fm.AppleTracks = mergeFm.AppleTracks
	fm.AppleAlbums = mergeFm.AppleAlbums
	fm.AppleSearchTracks = mergeFm.AppleSearchTracks
	fm.AppleSearchAlbums = mergeFm.AppleSearchAlbums

	fm.SpotifyAlbums = mergeFm.SpotifyAlbums
	fm.SpotifyTracks = mergeFm.SpotifyTracks
	fm.SpotifySearchAlbums = mergeFm.SpotifySearchAlbums
	fm.SpotifySearchTracks = mergeFm.SpotifySearchTracks

	fm.YandexAlbums = mergeFm.YandexAlbums
	fm.YandexTracks = mergeFm.YandexTracks
	fm.YandexSearchAlbums = mergeFm.YandexSearchAlbums
	fm.YandexSearchTracks = mergeFm.YandexSearchTracks

	fm.YoutubeAlbums = mergeFm.YoutubeAlbums
	fm.YoutubeTracks = mergeFm.YoutubeTracks
	fm.YoutubeSearchAlbums = mergeFm.YoutubeSearchAlbums
	fm.YoutubeSearchTracks = mergeFm.YoutubeSearchTracks
}

func (fm *FixturesMap) Reset() {
	fm.AppleTracks = map[string][]byte{}
	fm.AppleAlbums = map[string][]byte{}
	fm.AppleSearchTracks = map[string][]byte{}
	fm.AppleSearchAlbums = map[string][]byte{}

	fm.SpotifyAlbums = map[string][]byte{}
	fm.SpotifyTracks = map[string][]byte{}
	fm.SpotifySearchAlbums = map[string][]byte{}
	fm.SpotifySearchTracks = map[string][]byte{}

	fm.YandexAlbums = map[string][]byte{}
	fm.YandexTracks = map[string][]byte{}
	fm.YandexSearchAlbums = map[string][]byte{}
	fm.YandexSearchTracks = map[string][]byte{}

	fm.YoutubeAlbums = map[string][]byte{}
	fm.YoutubeTracks = map[string][]byte{}
	fm.YoutubeSearchAlbums = map[string][]byte{}
	fm.YoutubeSearchTracks = map[string][]byte{}
}
