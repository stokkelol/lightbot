package bot

import (
	"fmt"
	"github.com/dgraph-io/badger"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stokkelol/lightbot/pkg/cache"
	"log/slog"
	"os"
	"time"
)

const dbPath = "/db"

const helpCommand = "help"
const checkCommand = "check"
const help = `
	Hey there! I'm SvitlaBot and I'm here to help you with the one simple questions - "Світло є чи нема?"
`

type commandHandler interface {
	handle() error
}

var commandsList = map[string]commandHandler{}

// Bot is a struct that represents a telegram bot
type Bot struct {
	backend string
	bot     *telegram.BotAPI
	token   string
	db      *badger.DB

	cache *cache.Cache
}

// New creates a new bot
func New(token string, cache *cache.Cache) (*Bot, error) {
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:   bot,
		token: token,
		db:    db,
		cache: cache,
	}, nil
}

func (b *Bot) Run() {
	b.bot.Debug = true
	if err := b.setWebhook(); err != nil {
		return
	}

	b.run(b.bot.ListenForWebhook(fmt.Sprintf("/%s/handle", os.Getenv("TELEGRAM_TOKEN"))))

	return
}
func (b *Bot) run(ch telegram.UpdatesChannel) {
	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for {
		select {
		//case <-ticker.C:
		//	diff := time.Now().Sub(b.cache.GetLastTimestamp())
		//
		//	if diff > 60*time.Second {
		//		msg, err := b.bot.Send(telegram.NewMessage(0, "Світло є чи нема?"))
		//		if err != nil {
		//			slog.Error("send message", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		//
		//			continue
		//		}
		//
		//		slog.Info("send message", slog.Attr{Key: "message", Value: slog.StringValue(msg.Text)})
		//	}
		case update := <-ch:
			switch {
			case update.Message == nil:
				continue
			case update.Message.IsCommand():
				b.handleCommand(update)
			}
		}
	}
}

func (b *Bot) setWebhook() error {
	info, err := b.bot.GetWebhookInfo()
	if err != nil {
		return err
	}

	if info.URL != "" {
		return nil
	}

	_, err = b.bot.SetWebhook(telegram.NewWebhook(webhookAddr()))
	return err
}

func (b *Bot) handleCommand(update telegram.Update) {
	var txt string
	switch update.Message.Command() {
	case helpCommand:
		txt = "TODO"
	case checkCommand:
		txt = "Світло є чи нема?"
		if diff := time.Now().Sub(b.cache.GetLastTimestamp()); diff > 60*time.Second {
			txt += "Світла нема."
			txt += fmt.Sprintf("Останній раз було о %s", b.cache.GetLastTimestamp().Format("15:04:05"))
		} else {
			txt += "Світло є."
		}
	}

	msg := telegram.NewMessage(update.Message.Chat.ID, help)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("send message", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}
}

func webhookAddr() string {
	return fmt.Sprintf("%s/%s/handle", os.Getenv("DOMAIN_NAME"), os.Getenv("TELEGRAM_KEY"))
}
