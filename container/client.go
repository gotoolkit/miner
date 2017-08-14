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
		fmt.Printf(`%s        %s
		`, container.ID[:10], container.Image)
	}
}

func (c *client) Images() string {
	images, err := c.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return fmt.Sprintln("error", err)
	}
	iList := ""
	for _, image := range images {
		iList += fmt.Sprintf(`%s        %s      %s
		`, Time(image.Created), Bytes(uint64(image.Size)), image.RepoTags)
	}
	return iList
}

func (c *client) Version() string {
	sVersion, err := c.cli.ServerVersion(context.Background())
	if err != nil {
		return fmt.Sprintln("error", err)
	}
	return fmt.Sprintf(`
Server:
 Version:      %s
 API version:  %s
 Go version:   %s
 Git commit:   %s
 Built:        %s
 OS/Arch:      %s
 Experimental: %t
`, sVersion.Version, sVersion.APIVersion, sVersion.GoVersion, sVersion.GitCommit, sVersion.BuildTime, sVersion.Os, sVersion.Experimental)
}

func (c *client) Info() string {
	info, err := c.cli.Info(context.Background())
	if err != nil {
		return fmt.Sprintln("error", err)
	}
	return fmt.Sprintf(`
Containers: %d
 Running: %d
 Paused: %d
 Stopped: %d
Images: %d
Server Version: %s
Kernel Version: %s
Operating System: %s
OSType: %s
Architecture: %s
CPUs: %d
Total Memory: %s
Name: %s
Experimental: %t`, info.Containers, info.ContainersRunning, info.ContainersPaused, info.ContainersStopped, info.Images, info.ServerVersion, info.KernelVersion, info.OperatingSystem, info.OSType, info.Architecture, info.NCPU, Bytes(uint64(info.MemTotal)), info.Name, info.ExperimentalBuild)
}

func (c *client) PS() string {
	containers, err := c.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return fmt.Sprintln("error", err)
	}
	cList := ""
	for _, container := range containers {
		cList += fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
	}
	return cList
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
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("container", "container"),
			tgbotapi.NewInlineKeyboardButtonData("image", "image"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("network", "network"),
			tgbotapi.NewInlineKeyboardButtonData("volume", "volume"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("info", "info"),
			tgbotapi.NewInlineKeyboardButtonData("version", "version"),
		),
	)
	// tgbotapi.NewInlineKeyboardButtonURL("dash", "https://dash.gokit.info/dashboard/#/")),
	for update := range updates {

		if update.CallbackQuery != nil {
			callBackConf := tgbotapi.CallbackConfig{
				CallbackQueryID: update.CallbackQuery.ID,
			}
			if _, err := bot.AnswerCallbackQuery(callBackConf); err != nil {
				log.Println(err)
			}
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      int64(update.CallbackQuery.From.ID),
					MessageID:   update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &markup,
				},
			}
			switch update.CallbackQuery.Data {
			case "container":
				edit.Text = c.PS()
			case "image":
				edit.Text = c.Images()
			case "network":
				edit.Text = "network"
			case "volume":
				edit.Text = "volume"
			case "info":
				edit.Text = c.Info()
			case "version":
				edit.Text = c.Version()
			}
			_, err = bot.Send(edit)
			continue
		}
		// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, c.PS())
		// msg.ReplyMarkup = markup

		// bot.Send(msg)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(int64(update.Message.From.ID), update.Message.Text+"\n Group: https://t.me/dockertutorial"+"\n Author: @paultian")
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

func (c *client) StartInlineQuery() {
	bot := c.bot
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error on update chan: ", err)
	}

	for update := range updates {
		if update.InlineQuery == nil { // if no inline query, ignore it
			continue
		}

		article := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
		article.Description = update.InlineQuery.Query

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     0,
			Results:       []interface{}{article},
		}

		if _, err := bot.AnswerInlineQuery(inlineConf); err != nil {
			log.Println(err)
		}
	}
}
