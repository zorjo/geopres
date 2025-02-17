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
	password := os.Getenv("PASSWORD")
	connStr := fmt.Sprintf("postgresql://sourjyendra:%s@ep-hidden-limit-01274492-pooler.us-east-2.aws.neon.tech/neondb?sslmode=require", password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

		} else if update.Message.Location != nil && State[int(update.Message.Chat.ID)] == "addoffice" {
			var officeID int64
			err = db.QueryRow("INSERT INTO office_locations (latitude,longitude,employer_id) VALUES ($1, $2,$3) RETURNING id", update.Message.Location.Latitude, update.Message.Location.Longitude, update.Message.Chat.ID).Scan(&officeID)
			if err != nil {
				msg.Text = "Error adding office"
			} else {
				msg.Text = "Office added successfully on ID : "
			}
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
						tgbotapi.NewKeyboardButton("/checkin"),
					),
				)
			case "checkout":
				State[int(update.Message.Chat.ID)] = "checkout"
				msg.Text = "You have checked out"
			case "addoffice":
				//add office location to the database and return the office id
				State[int(update.Message.Chat.ID)] = "addoffice"
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButtonLocation("Share Office Location"),
					),
				)
				msg.Text = "Please share the office location"
			case "subscribe":
				args := update.Message.CommandArguments()
				if args == "" {
					msg.Text = "Please provide an office ID. Usage: /subscribe <office_id>"
				} else {
					username := ""
					if update.Message.From.UserName != "" {
						username = update.Message.From.UserName
					}

					_, err = db.Exec("INSERT INTO users (telegram_id,telegram_name,office_id) VALUES ($1, $2, $3)", update.Message.Chat.ID, username, args)
					if err != nil {
						msg.Text = "Error subscribing to office"
						fmt.Println(err)
					} else {
						msg.Text = fmt.Sprintf("You have subscribed to office %s updates", args)
					}
				}
			case "deloffice":
				args := update.Message.CommandArguments()
				if args == "" {
					msg.Text = "Please provide an office ID. Usage: /deloffice <office_id>"
				} else {
					State[int(update.Message.Chat.ID)] = "deloffice"
					msg.Text = fmt.Sprintf("Deleting office %s", args)
				}
			case "listoffice":
				State[int(update.Message.Chat.ID)] = "listoffice"
				rows, err := db.Query("SELECT id, latitude, longitude FROM offices")
				if err != nil {
					msg.Text = "Error fetching offices"
				} else {
					var officeList string = "*Registered Offices:*\n"
					for rows.Next() {
						var id int
						var lat, long float64
						if err := rows.Scan(&id, &lat, &long); err != nil {
							continue
						}
						officeList += fmt.Sprintf("ID: `%d` | Location: `%f, %f`\n", id, lat, long)
					}
					msg.Text = officeList
					msg.ParseMode = "MarkdownV2"
				}
				defer rows.Close()
			case "status":
				State[int(update.Message.Chat.ID)] = "status"
				msg.Text = "Current status:"
			case "history":
				State[int(update.Message.Chat.ID)] = "history"
				msg.Text = "Your attendance history:"
			case "help":
				msg.Text = "Available commands: \n/checkin - Check in to the office \n/checkout - Check out of the office \n/addoffice - Add office location\n/status - Show your current status"
			default:
				msg.Text = "Unsupported command, check /help"
			}

		}
		bot.Send(msg)
	}
}
