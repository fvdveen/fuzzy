package fuzzy

import (
	"io"
	"sync"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/fuzzy/internal/queue"
)

// VoiceHandlerGenerator creates a voicehandler
type VoiceHandlerGenerator func(sess *discordgo.Session, bot *Bot, gid, voiceChanID, textChanID string, vi VoiceItem) VoiceHandler

// VoiceHandler handles voice commands
type VoiceHandler interface {
	Play(VoiceItem)

	Skip()
	Stop()
	Pause()
	Resume()
	Loop()
	Repeat()
}

// VoiceItem is a item that should be played in the VoiceHandler
// it should be replayable
type VoiceItem interface {
	// OpusFrame must return io.EOF when it is done
	OpusFrame() ([]byte, error)
	ResetPlayback()
}

type defaultVoiceHandler struct {
	bot                          *Bot
	sess                         *discordgo.Session
	gid, voiceChanID, textChanID string
	log                          Logger
	voiceConn                    *discordgo.VoiceConnection
	mu                           sync.Mutex

	queue *queue.Queue

	loop, repeat, paused atomic.Value

	skipChan   chan interface{}
	stopChan   chan interface{}
	pauseChan  chan interface{}
	resumeChan chan interface{}
	loopChan   chan interface{}
	repeatChan chan interface{}
}

// DefaultVoiceHandler creates the default voice handler
func DefaultVoiceHandler(sess *discordgo.Session, bot *Bot, gid, voiceChanID, textChanID string, vi VoiceItem) VoiceHandler {
	vh := &defaultVoiceHandler{
		gid:         gid,
		voiceChanID: voiceChanID,
		textChanID:  textChanID,
		bot:         bot,
		sess:        sess,
		queue:       queue.New(),
		log:         bot.Generator().Logger(bot.Config().LogLevel),
		skipChan:    make(chan interface{}, 1),
		stopChan:    make(chan interface{}, 1),
		pauseChan:   make(chan interface{}, 1),
		resumeChan:  make(chan interface{}, 1),
		loopChan:    make(chan interface{}, 1),
		repeatChan:  make(chan interface{}, 1),
	}
	vh.queue.PushBack(vi)

	vh.loop.Store(false)
	vh.repeat.Store(false)
	vh.paused.Store(false)

	go vh.handle()

	return vh
}

func (vh *defaultVoiceHandler) Play(vi VoiceItem) {
	vh.queue.PushBack(vi)
}

func (vh *defaultVoiceHandler) Skip() {
	vh.skipChan <- 0
}

func (vh *defaultVoiceHandler) Stop() {
	vh.stopChan <- 0
}

func (vh *defaultVoiceHandler) Pause() {
	vh.pauseChan <- 0
}

func (vh *defaultVoiceHandler) Resume() {
	vh.resumeChan <- 0
}

func (vh *defaultVoiceHandler) Loop() {
	vh.loopChan <- 0
}

func (vh *defaultVoiceHandler) Repeat() {
	vh.repeatChan <- 0
}

func (vh *defaultVoiceHandler) handle() {
	if vh.queue == nil {
		return
	}

	var err error
	vh.voiceConn, err = vh.sess.ChannelVoiceJoin(vh.gid, vh.voiceChanID, false, true)
	if err != nil {
		vh.log.Errorf("Could not open voice connection: %v", err)
		vh.bot.DeleteVoiceHandler(vh.gid)
		return
	}

	if err := vh.voiceConn.Speaking(true); err != nil {
		vh.log.Errorf("Could not send speaking packet: %v", err)
		vh.bot.DeleteVoiceHandler(vh.gid)
		return
	}

	for {
		if vh.queue == nil || vh.queue.Length() == 0 {
			if err := vh.voiceConn.Disconnect(); err != nil {
				vh.log.Errorf("Could not send disconnect from voice channel: %v", err)
			}
			vh.bot.DeleteVoiceHandler(vh.gid)

			return
		}
		vi, ok := vh.queue.PopFront().(VoiceItem)
		if !ok {
			vh.log.Errorf("Could not convert voice item of type: %T to VoiceItem", vi)
			continue
		}

		if err := vh.playItem(vi); err != nil {
			vh.log.Errorf("Could not play voice item: %v", err)
		}

		if vh.repeat.Load().(bool) {
			vi.ResetPlayback()
			vh.queue.PushFront(vi)
		} else if vh.loop.Load().(bool) {
			vi.ResetPlayback()
			vh.queue.PushBack(vi)
		}
	}
}

func (vh *defaultVoiceHandler) playItem(vi VoiceItem) error {
	for {
		f, err := vi.OpusFrame()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		for vh.paused.Load().(bool) {
			select {
			case <-vh.pauseChan:
				vh.paused.Store(true)
			case <-vh.resumeChan:
				vh.paused.Store(false)
			case <-vh.loopChan:
				vh.loop.Store(!vh.loop.Load().(bool))
			case <-vh.repeatChan:
				vh.repeat.Store(!vh.repeat.Load().(bool))
			case <-vh.skipChan:
				vh.paused.Store(false)
				vi.ResetPlayback()
				return nil
			case <-vh.stopChan:
				vh.mu.Lock()
				vh.stop()
				vh.mu.Unlock()
				return nil
			}
		}
		select {
		case <-vh.resumeChan:
			vh.paused.Store(false)
		case <-vh.pauseChan:
			vh.paused.Store(true)
		case <-vh.loopChan:
			vh.loop.Store(!vh.loop.Load().(bool))
		case <-vh.repeatChan:
			vh.repeat.Store(!vh.repeat.Load().(bool))
		case <-vh.skipChan:
			vh.paused.Store(false)
			vi.ResetPlayback()
			return nil
		case <-vh.stopChan:
			vh.mu.Lock()
			vh.stop()
			vh.mu.Unlock()
			return nil
		case vh.voiceConn.OpusSend <- f:
		}
	}
}

func (vh *defaultVoiceHandler) stop() {
	vh.paused.Store(false)
	vh.loop.Store(false)
	vh.repeat.Store(false)
	vh.queue = nil
}
