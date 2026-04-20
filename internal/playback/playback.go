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
	"github.com/gopxl/beep/flac"
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
		isFlac := strings.Contains(contentType, "flac")
		isAAC := strings.Contains(contentType, "aac") || strings.Contains(contentType, "m4a")
		isHLS := strings.Contains(contentType, "mpegurl") || strings.Contains(contentType, "apple.mpegurl")

		// Route to native decoders if possible
		if isOgg {
			streamer, format, err = vorbis.Decode(res.Body)
		} else if isMp3 {
			streamer, format, err = mp3.Decode(res.Body)
		} else if isFlac {
			streamer, format, err = flac.Decode(res.Body)
		}

		// If it's a known FFmpeg-only format, or native decoding failed, fallback.
		if isAAC || isHLS || (err != nil && (isOgg || isMp3 || isFlac)) {
			res.Body.Close()
			return p.playWithFFmpeg(url)
		}

		// If we still don't have a streamer and it wasn't a known codec,
		// we can try one last greedy attempt with FFmpeg.
		if streamer == nil {
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
