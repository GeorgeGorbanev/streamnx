package vibeshare

import (
	"fmt"
	"log"
	"strings"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/translit"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"

	"github.com/tucnak/telebot"
)

type Vibeshare struct {
	spotifyClient  *spotify.Client
	telegramRouter *telegram.Router
	telegramSender telegram.Sender
	ymusicClient   *ymusic.Client
}

type Input struct {
	SpotifyClient  *spotify.Client
	TelegramSender telegram.Sender
	YmusicClient   *ymusic.Client
}

func NewVibeshare(input *Input) *Vibeshare {
	vs := Vibeshare{
		spotifyClient:  input.SpotifyClient,
		telegramSender: input.TelegramSender,
		ymusicClient:   input.YmusicClient,
	}

	vs.telegramRouter = vs.makeRouter()

	return &vs
}

func (vs *Vibeshare) HandleText(inMsg *telebot.Message) {
	vs.telegramRouter.RouteMessage(inMsg)
}

func (vs *Vibeshare) makeRouter() *telegram.Router {
	router := telegram.NewRouter()

	router.Register(spotify.OpenTrackRe, vs.spotifyTrack())
	router.Register(ymusic.TrackURLRegExp, vs.yMusicTrack())
	router.RegisterNotFound(vs.notFound())

	return router
}

func (vs *Vibeshare) respond(inMsg *telebot.Message, text string) (*telebot.Message, error) {
	return vs.telegramSender.Send(
		&telegram.Message{
			To:   inMsg.Sender,
			Text: text,
		},
	)
}

func (vs *Vibeshare) yMusicSearch(spotifyTrack *spotify.Track) (*ymusic.Track, error) {
	artistName := strings.ToLower(spotifyTrack.Artists[0].Name)
	trackName := strings.ToLower(spotifyTrack.Name)

	yMusicTrack, err := vs.ymusicClient.SearchTrack(artistName, trackName)
	if err != nil {
		return nil, fmt.Errorf("failed to find ymusic track: %w", err)
	}
	if yMusicTrack != nil {
		foundLowcasedArtist := strings.ToLower(yMusicTrack.Artists[0].Name)
		if artistName == foundLowcasedArtist {
			return yMusicTrack, nil
		}

		translitedArtist := translit.CyrillicToLatin(foundLowcasedArtist)
		if artistName == translitedArtist {
			return yMusicTrack, nil
		}
		return nil, nil
	}

	if spotifyTrack.NameContainsRussianLetters() {
		translitedArtist := translit.LatinToCyrillic(artistName)
		yMusicTrack, err = vs.ymusicClient.SearchTrack(translitedArtist, trackName)
		if err != nil {
			return nil, fmt.Errorf("failed to find ymusic track: %w", err)
		}
	}

	return yMusicTrack, nil
}

func (vs *Vibeshare) spotifyTrack() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		trackID := spotify.DetectTrackID(inMsg.Text)
		spotifyTrack, err := vs.spotifyClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if spotifyTrack == nil {
			outMsg, err := vs.respond(inMsg, "track not found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		yMusicTrack, err := vs.yMusicSearch(spotifyTrack)
		if err != nil {
			log.Printf("failed to search ymusic: %s", err)
			return
		}
		if yMusicTrack == nil {
			outMsg, err := vs.respond(inMsg, "no ym track found")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		outMsg, err := vs.respond(inMsg, yMusicTrack.URL())
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}

func (vs *Vibeshare) yMusicTrack() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		trackID := ymusic.ParseTrackID(inMsg.Text)
		yMusicTrack, err := vs.ymusicClient.GetTrack(trackID)
		if err != nil {
			log.Printf("error fetching track: %s", err)
			return
		}
		if yMusicTrack == nil {
			outMsg, err := vs.respond(inMsg, "track not found in yandex music")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		spotifyTrack, err := vs.spotifyClient.SearchTrack(yMusicTrack.Artists[0].Name, yMusicTrack.Title)
		if err != nil {
			log.Printf("failed to search spotify: %s", err)
			return
		}
		if spotifyTrack == nil {
			outMsg, err := vs.respond(inMsg, "no track found in spotify")
			if err != nil {
				log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
				return
			}
			log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
			return
		}

		outMsg, err := vs.respond(inMsg, spotifyTrack.URL())
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}

func (vs *Vibeshare) notFound() telegram.HandlerFunc {
	return func(inMsg *telebot.Message) {
		outMsg, err := vs.respond(inMsg, "no track link found")
		if err != nil {
			log.Printf("failed to send message to %s: %s", inMsg.Sender.Username, err)
			return
		}
		log.Printf("Sent message to %s: %s", inMsg.Sender.Username, outMsg.Text)
	}
}
