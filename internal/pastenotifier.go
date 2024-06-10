package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gen2brain/beeep"
	"github.com/jimyag/log"
	"golang.design/x/clipboard"
)

type PasteNotifier struct {
	handlers []Handler
	clipText <-chan []byte
}

func (pn *PasteNotifier) Register(h Handler) {
	pn.handlers = append(pn.handlers, h)
}

var PN *PasteNotifier

func init() {
	clipText := clipboard.Watch(context.Background(), clipboard.FmtText)
	PN = &PasteNotifier{
		clipText: clipText,
	}
}

func Run() {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGSEGV, syscall.SIGABRT, syscall.SIGTERM)
	signal.Ignore(syscall.SIGPIPE, syscall.SIGHUP)
	for {
		select {
		case <-signalC:
			return
		case clipText := <-PN.clipText:
			for _, h := range PN.handlers {
				title, message, err := h.Handle(string(clipText))
				if err != nil {
					log.Error(err).Msg("handle error")
				}
				if title != "" && message != "" {
					if err = beeep.Notify(title, message, ""); err != nil {
						log.Error(err).Str("title", title).Str("message", message).Msg("notify error")
					}

				}
			}
		}
	}
}
