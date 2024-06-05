package youtube

type snippet struct {
	Title                  string `json:"title"`
	ChannelTitle           string `json:"channelTitle"`
	Description            string `json:"description"`
	VideoOwnerChannelTitle string `json:"videoOwnerChannelTitle"`
}

func (s *snippet) ownerChannelTitle() string {
	if s.VideoOwnerChannelTitle != "" {
		return s.VideoOwnerChannelTitle
	}
	return s.ChannelTitle
}
