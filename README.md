# Go Fullstack Website

This is a Go-based fullstack web application with HTMX and Tailwind CSS for building modern, dynamic UIs. This project serves as a template for developing Go-based websites, and it includes reusable components like a `<head>` template for every page.

## Table of Contents

- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Setup](#development-setup)
- [Running the Application](#running-the-application)
- [Building and Deploying](#building-and-deploying)
- [Contributing](#contributing)
- [License](#license)

---

## Getting Started

To get started with developing this Go web application, clone the repository and set up your environment:

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/go-fullstack-your-project.git
   cd go-fullstack-your-project
   ```

2. **Set Up Environment Variables:**

   Copy `.env.example` to `.env` and fill in the required values for your local development environment.

   ```bash
   cp .env.example .env
   ```

   Make sure the `.env` file contains the correct environment variables for Postgres, Redis, and any other configurations.

---

## Project Structure

The project structure is organized as follows:

```
/go-fullstack-your-project
├── docker-compose.yml        # Docker configuration
├── go.mod                    # Go module file
├── go.sum                    # Go checksum file
├── main.go                   # Main Go application file
├── README.md                 # Project README
├── routes                    # Contains route handlers
│   ├── auth.go               # Authentication routes
│   ├── index.go              # Home page routes
│   ├── routes.go             # Router setup
│   └── user.go               # User-related routes
├── services                  # Contains service logic
│   ├── database              # Database related logic
│   │   └── database.go       # Database connection
│   └── redis                 # Redis-related logic
│       └── redis.go          # Redis connection
└── templates                 # HTML templates
    ├── error.html            # Error page template
    ├── index.html            # Home page template
    ├── login.html            # Login page template
    └── register.html        # Registration page template
```

- **`/static/`**: Holds the static assets such as your Tailwind CSS and JavaScript files.
- **`/templates/`**: Contains the HTML templates, including the reusable `head.html`, `layout.html`, and page-specific templates.
- **`main.go`**: Main Go application that renders templates and serves the web application.
- **`.env`**: Stores the configuration for environment variables (e.g., Redis, Postgres, etc.).

---

## Development Setup

To get started with local development, follow these steps:

### 1. Install Go

If you don't have Go installed, follow the official Go installation guide:
https://golang.org/doc/install

### 2. Install Dependencies

Make sure to install the required Go dependencies:

```bash
go mod tidy
```

### 3. Set Up Your Database and Redis

Ensure you have a working Postgres database and Redis server. You can use Docker to quickly set up both services:

```bash
docker-compose up
```

Make sure your `.env` file has the correct connection details for both services.

---

## Running the Application

Once you've set up everything, you can run the application using Go:

```bash
go run main.go
```

This will start the application on `http://localhost:8080`. You should see the home page served from `home.html`, with Tailwind CSS styling applied.

### Available Endpoints:

- `/`: The home page.
- Any other endpoints you add will follow the same template rendering pattern.

---

## Building and Deploying

To build the project for production, you can compile it using Go's build command:

```bash
go build -o go-fullstack-your-project
```

After building, you can deploy it to any hosting service (e.g., DigitalOcean, Heroku, etc.). You can also containerize the app using Docker:

---

## Contributing

We welcome contributions! If you have any improvements or bug fixes, please follow these steps:

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a pull request.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy coding!**
