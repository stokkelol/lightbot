package bot

import (
	"fmt"
	"github.com/dgraph-io/badger"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
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
}

// New creates a new bot
func New(token string) (*Bot, error) {
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
	for update := range ch {
		switch {
		case update.Message == nil:
			continue
		case update.Message.IsCommand():
			b.handleCommand(update)
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
	switch update.Message.Command() {
	case helpCommand:
		msg := telegram.NewMessage(update.Message.Chat.ID, help)
		msg.ReplyToMessageID = update.Message.MessageID
		b.bot.Send(msg)
	case checkCommand:

	}
}

func webhookAddr() string {
	return fmt.Sprintf("%s/%s/handle", os.Getenv("DOMAIN_NAME"), os.Getenv("TELEGRAM_KEY"))
}
