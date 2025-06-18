package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	telebot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	api_client "github.com/snakehunterr/hacs_dbapi_client"
	types "github.com/snakehunterr/hacs_dbapi_types"
	api_errors "github.com/snakehunterr/hacs_dbapi_types/errors"
)

var (
	client *api_client.APIClient
	csh    = make(ClientStateHandler)
	ch     = make(map[int64]*types.Client)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	opts := []telebot.Option{
		telebot.WithDebug(),
	}

	bot := prepare(opts)
	bot.Start(ctx)
}

func DefaultHandler(ctx context.Context, bot *telebot.Bot, update *models.Update) {
	var id int64
	switch {
	case update.Message != nil:
		id = update.Message.From.ID

	case update.CallbackQuery != nil:
		id = update.CallbackQuery.From.ID
	}
	state := csh.Get(id)

	if state == StateRegisterClient {
		EndClientRegistration(ctx, bot, id)
	}

	c, err := client.ClientGetByID(id)
	if err != nil {
		if api_errors.IsChildErr(err, api_errors.ErrSQLNoRows) {
			StartClientRegistration(ctx, bot, id)
			return
		}

		log.Println("ClientGetByID() err:", err)
		SendError(ctx, bot, id)
		return
	}

	switch csh.Get(id) {
	case NoState:
		ShowMainMenu(ctx, bot, c)
		return

	case StateMainMenu:
		if update.Message != nil {
			_, err := bot.DeleteMessage(ctx, &telebot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			})
			if err != nil {
				log.Println("bot.DeleteMessage() err:", err)
			}
		}
	}
}

func prepare(opts []telebot.Option) *telebot.Bot {
	c := api_client.New(
		os.Getenv("DBAPI_SERVER_HOST"),
		os.Getenv("DBAPI_SERVER_PORT"),
	)
	client = &c

	bot, err := telebot.New(
		os.Getenv("TELEBOT_KEY"),
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return bot
}

type ClientState uint8

const (
	NoState = iota
	StateRegisterClient
	StateMainMenu
)

type ClientStateHandler map[int64]ClientState

func (csh ClientStateHandler) Set(id int64, state ClientState) {
	csh[id] = state
}

func (csh ClientStateHandler) Get(id int64) ClientState {
	s, ok := csh[id]
	if !ok {
		csh[id] = NoState
		return NoState
	}
	return s
}
