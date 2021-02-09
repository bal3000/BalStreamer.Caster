package infrastructure

import (
	"errors"
	vlc "github.com/adrg/libvlc-go/v3"
)

const microdns = "microdns_renderer"

type Streamer interface {
	DiscoverChromecasts(itemAdded chan string, itemDeleted chan string) error
	GetChromecast(key string) *vlc.Renderer
	StartCasting(url string, rendererItem *vlc.Renderer) error
	StopCasting(rendererItem *vlc.Renderer) error
	CloseAndCleanUp()
}

type chromecastStreamer struct {
	discoverer  *vlc.RendererDiscoverer
	player      *vlc.Player
	chromecasts map[string]vlc.Renderer
}

func NewChromecastStreamer() (Streamer, error) {
	// Initialize libVLC. Additional command line arguments can be passed in
	// to libVLC by specifying them in the Init function.
	if err := vlc.Init("--no-audio"); err != nil {
		return nil, err
	}

	discoverer, err := getDiscoverer()
	if err != nil {
		return nil, err
	}

	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
		return nil, err
	}

	return &chromecastStreamer{discoverer: discoverer, player: player, chromecasts: make(map[string]vlc.Renderer)}, nil
}

func getDiscoverer() (*vlc.RendererDiscoverer, error) {
	descriptors, err := vlc.ListRendererDiscoverers()
	if err != nil {
		return nil, err
	}

	for _, descriptor := range descriptors {
		if descriptor.Name != microdns {
			continue
		}

		discoverer, err := vlc.NewRendererDiscoverer(descriptor.Name)
		if err != nil {
			return nil, err
		}

		return discoverer, nil
	}

	return nil, errors.New("could not find discovery service")
}

// DiscoverChromecasts notifies when chromecasts are found or lost on my network
func (s *chromecastStreamer) DiscoverChromecasts(itemAdded chan string, itemDeleted chan string) error {
	// Start renderer discovery.
	stop := make(chan error)

	callback := func(event vlc.Event, r *vlc.Renderer) {
		// NOTE: the discovery service cannot be stopped or released from
		// the callback function. Doing so will result in undefined behavior.

		switch event {
		case vlc.RendererDiscovererItemAdded:
			// New renderer (`r`) found.
			rendererType, err := r.Type()
			if err != nil {
				stop <- err
			}
			if rendererType == vlc.RendererChromecast {
				name, err := r.Name()
				if err != nil {
					stop <- err
				}
				s.chromecasts[name] = *r
				itemAdded <- name
			}
		case vlc.RendererDiscovererItemDeleted:
			// The renderer (`r`) is no longer available.
			rendererType, err := r.Type()
			if err != nil {
				stop <- err
			}
			if rendererType == vlc.RendererChromecast {
				name, err := r.Name()
				if err != nil {
					stop <- err
				}
				delete(s.chromecasts, name)
				itemDeleted <- name
			}
		}
	}

	if err := s.discoverer.Start(callback); err != nil {
		return err
	}

	if err := <-stop; err != nil {
		return err
	}
	if err := s.discoverer.Stop(); err != nil {
		return err
	}

	return nil
}

func (s *chromecastStreamer) GetChromecast(key string) *vlc.Renderer {
	if renderer, ok := s.chromecasts[key]; ok {
		return &renderer
	}

	return nil
}

func (s *chromecastStreamer) StartCasting(url string, rendererItem *vlc.Renderer) error {
	media, err := s.player.LoadMediaFromURL(url)
	if err != nil {
		return err
	}
	defer media.Release()

	// Set renderer to the given chromecast
	if err := s.player.SetRenderer(rendererItem); err != nil {
		return err
	}
	// Set media to play
	if err := s.player.SetMedia(media); err != nil {
		return err
	}

	// Start media playback.
	if err = s.player.Play(); err != nil {
		return err
	}
	return nil
}

func (s *chromecastStreamer) StopCasting(rendererItem *vlc.Renderer) error {
	if s.player.IsPlaying() {
		err := s.player.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// CloseAndCleanUp closes all connections and disposes resources
func (s *chromecastStreamer) CloseAndCleanUp() {
	if s.player.IsPlaying() {
		s.player.Stop()
	}
	s.player.Release()
	defer vlc.Release()
	defer s.discoverer.Release()
}
