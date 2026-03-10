# ZabScrap

A modern web application for extracting and visualizing data from ZabDesk.

## Features

- **Attendance Tracking** - Extract and view your attendance records from ZabDesk
- **Responsive Layout** - Works seamlessly on desktop and mobile devices
- **Real-time Data** - Fetch fresh data directly from ZabDesk
- **JSON API** - RESTful API endpoints for easy integration

## Tech Stack

### Backend
- **Go 1.25** - Fast and efficient server
- **Chi Router** - Lightweight HTTP router
- **Standard Library** - Minimal dependencies

### Frontend
- **Vanilla JavaScript** - No frameworks, just pure JS
- **HTML5 & CSS3** - Semantic markup and modern styling

### Design
- **Neobrutalism** - Bold borders, hard shadows, and vibrant colors
- **Accessibility** - Semantic HTML and clear visual hierarchy
- **No Build Step** - Static files served directly

## Installation

### Prerequisites
- Go 1.25 or higher
- Git

### Clone the Repository
```bash
git clone https://github.com/hammadmajid/zabscrap.git
cd zabscrap
```

### Install Dependencies
```bash
go mod download
```

### Set Environment Variables
Create a `.env` file or export the PORT variable:
```bash
export PORT=8080
```

### Build and Run
```bash
go build -o zabscrap ./cmd/zabscrap
./zabscrap
```

Or use [Air](https://github.com/air-verse/air) for development with hot reload:
```bash
air
```

The application will be available at `http://localhost:8080`

## Project Structure

```
zabscrap/
├── cmd/
│   └── zabscrap/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   └── handler.go        # HTTP handlers
│   ├── app/
│   │   └── app.go           # Application context
│   ├── models/
│   │   └── attendance.go    # Data models
│   ├── router/
│   │   └── router.go        # Route definitions
│   └── scraper/
│       └── scraper.go       # ZabDesk scraping logic
├── web/
│   ├── css/
│   │   └── styles.css       # Neobrutalist styles
│   ├── js/
│   │   └── app.js          # Frontend logic
│   └── index.html          # Main HTML page
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

### `POST /fetch`
Fetch attendance data from ZabDesk.

**Request Body:**
```json
{
  "username": "your_username",
  "password": "your_password"
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "courseName": "Course Name",
      "instructor": "Instructor Name",
      "records": [
        {
          "lecture": "1",
          "date": "2026-02-16",
          "status": "Present"
        }
      ]
    }
  ]
}
```

### `GET /api/build-info`
Get the latest commit information from GitHub.

**Response:**
```json
{
  "success": true,
  "data": {
    "hash": "abc1234",
    "message": "Latest commit message",
    "timeAgo": "2 days ago",
    "available": true
  }
}
```

### `GET /health`
Health check endpoint.

## Development

### Hot Reload with Air
The project includes Air configuration for development:
```bash
air
```

### Code Style
- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write meaningful commit messages

### Adding New Features
1. Create a new branch: `git checkout -b feature/your-feature`
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## Roadmap

- [ ] Marks extraction
- [ ] Dark mode toggle
- [ ] Caching for better performance

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

**Note**: This project is not affiliated with or endorsed by SZABIST. Use responsibly and in accordance with your institution's policies.

