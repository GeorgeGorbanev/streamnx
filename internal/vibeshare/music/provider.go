package music

type Provider struct {
	Name string
	Code string
}

var (
	Apple = &Provider{
		Name: "Apple",
		Code: "ap",
	}
	Spotify = &Provider{
		Name: "Spotify",
		Code: "sf",
	}
	Yandex = &Provider{
		Name: "Yandex",
		Code: "ya",
	}
	Youtube = &Provider{
		Name: "Youtube",
		Code: "yt",
	}

	Providers = []*Provider{
		Apple,
		Spotify,
		Yandex,
		Youtube,
	}
)

func FindProviderByCode(code string) *Provider {
	for _, provider := range Providers {
		if provider.Code == code {
			return provider
		}
	}
	return nil
}
