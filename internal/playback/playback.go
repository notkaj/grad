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
	current    io.Closer
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

		isOgg := strings.Contains(contentType, "ogg") || strings.Contains(contentType, "vorbis")
		isMp3 := strings.Contains(contentType, "mpeg")

		if isOgg {
			streamer, format, err = vorbis.Decode(res.Body)
		} else if isMp3 {
			streamer, format, err = mp3.Decode(res.Body)
		}

		// If it's a known format but decode failed, or it's an unknown format (AAC, HLS, etc.)
		// we fallback to FFmpeg.
		if err != nil || (!isOgg && !isMp3) {
			res.Body.Close()
			return p.playWithFFmpeg(url)
		}

		resampled := beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
		wrapped := &closeWrapper{Streamer: resampled, closer: streamer}

		speaker.Lock()
		p.current = wrapped
		p.mixer.Add(wrapped)
		speaker.Unlock()

		return PlayerStartedMsg(fmt.Sprintf("Playing %s", contentType))
	}
}

func (p *Player) playWithFFmpeg(url string) tea.Msg {
	s, err := newFFmpegStreamer(url, p.sampleRate)
	if err != nil {
		return PlayerErrorMsg("ffmpeg error: " + err.Error())
	}

	speaker.Lock()
	p.current = s
	p.mixer.Add(s)
	speaker.Unlock()

	return PlayerStartedMsg("Playing via FFmpeg (fallback)")
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
