# Awesome Blog

Awesome Blog is a full-stack web application for managing blog posts, comments, and users. It consists of a backend API built with Go and a frontend application built with React.

## Features

- User authentication (register, login, logout)
- Create, read, update, and delete blog posts
- Comment on blog posts
- User management

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/awesome-blog.git
   cd awesome-blog
   ```

2. Start the application using Docker Compose:
   ```
   docker-compose up --build
   ```

3. The application will be available at:
    - Backend API: http://localhost:8080
    - Frontend: http://localhost:3000

## Project Structure

```
awesome-blog/
├── backend/
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   ├── migrations/
│   ├── static/
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── public/
│   ├── src/
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yaml
└── README.md
```

## API Documentation

The API documentation is available in OpenAPI (Swagger) format. You can view it by accessing the `/swagger` endpoint of the backend server when it's running.

## Development

To run the project in development mode:

1. Start the database:
   ```
   docker-compose up db
   ```

2. Run the backend:
   ```
   cd backend
   go run cmd/server/main.go
   ```

3. Run the frontend:
   ```
   cd frontend
   npm start
   ```

## Testing

To run the tests for the backend:

```
cd backend
go test ./...
```

To run the tests for the frontend:

```
cd frontend
npm test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.