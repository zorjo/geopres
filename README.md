# GeoPresence Bot 🌍

A Telegram bot for managing employee attendance and office locations using geolocation data.

## 🌟 Features

- **Check-in/Check-out System**: Employees can check in and out using their current location
- **Office Management**: Add and manage multiple office locations
- **Location Verification**: Verify employee presence within office premises
- **Subscription System**: Users can subscribe to specific office locations
- **Status Tracking**: Monitor current attendance status
- **History View**: Access attendance history records

## 🚀 Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL database
- Telegram Bot Token
- Neon.tech account (or any other PostgreSQL provider)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
TELEGRAM_KEY=your_telegram_bot_token
PASSWORD=your_database_password
```

### Database Setup

The application requires the following Postgres migration script to run
`db/migrations/seed.sql`
### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/geopres.git
cd geopres
```

2. Install dependencies:
```bash
go mod download
```

3. Build and run:
```bash
go build
./geopres
```

## 📱 Usage

### Available Commands

- `/checkin` - Check in at current location
- `/checkout` - Check out from current location
- `/addoffice` - Add a new office location
- `/subscribe <office_id>` - Subscribe to an office
- `/deloffice <office_id>` - Delete an office location
- `/listoffice` - List all registered offices
- `/status` - View current status
- `/history` - View attendance history
- `/help` - Display help message

## 🔒 Security

- SSL mode is required for database connections
- Employee locations are verified against registered office coordinates
- Telegram's built-in security features ensure secure communication

## 📦 Dependencies

- [github.com/go-telegram-bot-api/telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [github.com/lib/pq](https://github.com/lib/pq)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - *zorjo*

## 🙏 Acknowledgments

- Telegram Bot API team
- PostgreSQL team
- Neon.tech for database hosting

## 📞 Support

For support, email sourjyo@protonmail.com or create an issue in the repository.

---
Made with ❤️ by [zorjo]
