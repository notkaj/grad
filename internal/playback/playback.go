// Package playback provides audio playback support
package playback

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

type Player struct {
	URL        string
	mixer      *beep.Mixer
	current    beep.StreamCloser
	sampleRate beep.SampleRate
}

func NewPlayer() *Player {
	rate := beep.SampleRate(44100)
	speaker.Init(rate, rate.N(time.Second/10))
	mixer := &beep.Mixer{}
	speaker.Play(mixer)
	return &Player{sampleRate: rate, mixer: mixer}
}

func (p *Player) Play(url string) tea.Cmd {
	p.URL = url
	return func() tea.Msg {
		p.Stop()

		res, err := http.Get(url)
		if err != nil {
			return PlayerErrorMsg(err.Error())
		}

		contentType := strings.ToLower(res.Header.Get("Content-Type"))
		var (
			streamer beep.StreamSeekCloser
			format   beep.Format
		)

		// Select decoder based on Content-Type
		if strings.Contains(contentType, "ogg") || strings.Contains(contentType, "vorbis") {
			streamer, format, err = vorbis.Decode(res.Body)
		} else {
			// Default to MP3 for mpeg or unknown types
			streamer, format, err = mp3.Decode(res.Body)
		}

		if err != nil {
			res.Body.Close()
			return PlayerErrorMsg(fmt.Sprintf("decode error (%s): %v", contentType, err))
		}

		resampled := beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
		wrapped := &closeWrapper{Streamer: resampled, closer: streamer}

		speaker.Lock()
		p.current = wrapped
		p.mixer.Add(wrapped)
		speaker.Unlock()

		return PlayerStartedMsg("Player Playing")
	}
}

func (p *Player) Stop() {
	speaker.Lock()
	if p.current != nil {
		p.current.Close()
		p.current = nil
	}
	p.mixer.Clear()
	speaker.Unlock()
}

func (p *Player) IsPlaying() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return p.current != nil
}

type (
	PlayerErrorMsg   string
	PlayerStartedMsg string
)

type closeWrapper struct {
	beep.Streamer
	closer io.Closer
}

func (w *closeWrapper) Close() error {
	return w.closer.Close()
}
