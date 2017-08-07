package container

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"gopkg.in/telegram-bot-api.v4"
)

type client struct {
	cli *dockerclient.Client
	bot *tgbotapi.BotAPI
}

func NewClient(botToken string) *client {

	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	return &client{
		cli: cli,
		bot: bot,
	}
}

func (c *client) ListContainers() {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println("error", err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (c *client) StartBot(authorID int) {
	bot := c.bot
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error on update chan: ", err)
	}
	// markup := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("test1"), tgbotapi.NewKeyboardButton("test2")))
	markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("test", "@ps"), tgbotapi.NewInlineKeyboardButtonURL("dash", "https://dash.gokit.info/dashboard/#/")))
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			log.Println("button: ", update.CallbackQuery.Data)
		}
		
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyMarkup = markup

		bot.Send(msg)
	}
}

func (c *client) StartWebHook() {
	bot := c.bot
	updates := bot.ListenForWebhook("/test")
	go http.ListenAndServe(":8443", nil)
	for update := range updates {
		log.Printf("%+v\n", update)

		msg := tgbotapi.NewMessage(230431366, "test")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonSwitch("test", "http://localhost:8080")))
		// msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
