package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Replace with your bot token
	password := os.Getenv("PASSWORD");connStr := fmt.Sprintf("postgresql://sourjyendra:%s@ep-hidden-limit-01274492-pooler.us-east-2.aws.neon.tech/neondb?sslmode=require", password)
	db, err := sql.Open("postgres", connStr)
		if err != nil { log.Fatal(err) }; defer db.Close();

	State := make(map[int]string)
	val := os.Getenv("TELEGRAM_KEY")
	bot, err := tgbotapi.NewBotAPI(val)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//	log.Printf("%s %s", update.Message.From.FirstName, update.Message.Text)
		//Print out the shared gps location from the message
		if update.Message.Location != nil && State[int(update.Message.Chat.ID)] == "checkin" {
			log.Printf("Location: %f %f", update.Message.Location.Latitude, update.Message.Location.Longitude)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = "You have checked in"

		}

		//create commands to send location
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "checkin":
				State[int(update.Message.Chat.ID)] = "checkin"
				msg.Text = "Please share your location"
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButtonLocation("Share Location"),
					),
				)
			case "checkout":
				State[int(update.Message.Chat.ID)] = "checkout"
				msg.Text = "You have checked out"

			default:
				msg.Text = "Unsupported command, check /help"
			}

			bot.Send(msg)
		}
	}
}
