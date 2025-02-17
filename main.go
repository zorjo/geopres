package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // Distance in kilometers
}
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
			// Get user's office_id
			var officeID int64
			err = db.QueryRow("SELECT office_id FROM users WHERE telegram_id = $1", update.Message.Chat.ID).Scan(&officeID)
			if err != nil {
				msg.Text = "Error: You're not subscribed to any office"
				return
			}

			// Get office location
			var officeLat, officeLong float64
			err = db.QueryRow("SELECT latitude, longitude FROM office_locations WHERE id = $1", officeID).Scan(&officeLat, &officeLong)
			if err != nil {
				msg.Text = "Error: Office location not found"
				return
			}

			// Calculate distance between user and office
			distance := calculateDistance(update.Message.Location.Latitude, update.Message.Location.Longitude, officeLat, officeLong)
			if distance <= 0.1 { // Within 100 meters
				_, err = db.Exec("INSERT INTO attendance (employee_id, check_in, office_location_id) VALUES ($1, NOW(), $2)",
					update.Message.Chat.ID, officeID)
				if err != nil {
					msg.Text = "Error recording attendance"
					log.Printf("Error: %v", err)
				} else {
					msg.Text = "You have successfully checked in"
				}
			} else {
				msg.Text = fmt.Sprintf("You are too far from the office (%.2f km away)", distance)
			}
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
				// First check if user exists
				var exists bool
				err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)", update.Message.Chat.ID).Scan(&exists)
				if !exists {
					msg.Text = "You need to subscribe to an office first. Use /subscribe <office_id>"
				} else {
					msg.Text = "Please share your location"
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButtonLocation("Share Location"),
						),
					)
				}

			case "addoffice":
				// Verify if user is an employer
				var userType string
				err = db.QueryRow("SELECT type FROM users WHERE telegram_id = $1", update.Message.Chat.ID).Scan(&userType)
				if err != nil || userType != "employer" {
					msg.Text = "Only employers can add office locations"
				} else {
					State[int(update.Message.Chat.ID)] = "addoffice"
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButtonLocation("Share Office Location"),
						),
					)
					msg.Text = "Please share the office location"
				}

			case "subscribe":
				args := update.Message.CommandArguments()
				if args == "" {
					msg.Text = "Please provide an office ID. Usage: /subscribe <office_id>"
				} else {
					username := update.Message.From.UserName
					// Check if office exists
					var exists bool
					err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM office_locations WHERE id = $1)", args).Scan(&exists)
					if !exists {
						msg.Text = "Office ID does not exist"
					} else {
						_, err = db.Exec("INSERT INTO users (telegram_id, telegram_name, office_id, type) VALUES ($1, $2, $3, 'employee') ON CONFLICT (telegram_id) DO UPDATE SET office_id = $3",
							update.Message.Chat.ID, username, args)
						if err != nil {
							msg.Text = "Error subscribing to office"
							log.Printf("Error: %v", err)
						} else {
							msg.Text = fmt.Sprintf("You have subscribed to office %s", args)
						}
					}
				}

			case "listoffice":
				rows, err := db.Query("SELECT id, name, latitude, longitude FROM office_locations")
				if err != nil {
					msg.Text = "Error fetching offices"
					log.Printf("Error: %v", err)
				} else {
					var officeList string = "Registered Offices:\n"
					for rows.Next() {
						var id int64
						var name string
						var lat, long float64
						if err := rows.Scan(&id, &name, &lat, &long); err != nil {
							continue
						}
						officeList += fmt.Sprintf("ID: %d | Name: %s | Location: %f, %f\n", id, name, lat, long)
					}
					msg.Text = officeList
				}
				defer rows.Close()

			case "status":
				var lastCheckIn time.Time
				var officeName string
				err = db.QueryRow(`
        SELECT a.check_in, o.name
        FROM attendance a
        JOIN office_locations o ON a.office_location_id = o.id
        WHERE a.employee_id = $1
        ORDER BY a.check_in DESC
        LIMIT 1`,
					update.Message.Chat.ID).Scan(&lastCheckIn, &officeName)

				if err == sql.ErrNoRows {
					msg.Text = "No check-ins recorded"
				} else if err != nil {
					msg.Text = "Error fetching status"
					log.Printf("Error: %v", err)
				} else {
					msg.Text = fmt.Sprintf("Last check-in: %s at %s", lastCheckIn.Format("2006-01-02 15:04:05"), officeName)
				}

			case "history":
				rows, err := db.Query(`
        SELECT a.check_in, o.name
        FROM attendance a
        JOIN office_locations o ON a.office_location_id = o.id
        WHERE a.employee_id = $1
        ORDER BY a.check_in DESC
        LIMIT 10`,
					update.Message.Chat.ID)

				if err != nil {
					msg.Text = "Error fetching history"
					log.Printf("Error: %v", err)
				} else {
					var history string = "Your last 10 check-ins:\n"
					for rows.Next() {
						var checkIn time.Time
						var officeName string
						if err := rows.Scan(&checkIn, &officeName); err != nil {
							continue
						}
						history += fmt.Sprintf("%s at %s\n", checkIn.Format("2006-01-02 15:04:05"), officeName)
					}
					msg.Text = history
				}
				defer rows.Close()
			}

		}
		bot.Send(msg)
	}
}
