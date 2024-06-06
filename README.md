# Library Management System

This is a simple library management system built with Go. It provides routes to manage books, borrows, and authors.

## Getting Started

To start the server, run the following command:

```sh
go run main.go
```

## API Routes

### Books

- **Get All Books**
  ```
  GET /books
  ```
  Example: [http://localhost:8080/books](http://localhost:8080/books)

- **Search Books by Title or Author**
  ```
  GET /books/search?q={search-value}
  ```
  Example: [http://localhost:8080/books/search?q=Harry](http://localhost:8080/books/search?q=Harry)

### Borrows

- **Return Book**
  ```
  POST /borrows/return/{bookId}
  ```
  Example: [http://localhost:8080/borrows/return/123](http://localhost:8080/borrows/return/123)

- **Find Books Borrowed by User**
  ```
  GET /borrows/search?u={search-value}
  ```
  Example: [http://localhost:8080/borrows/search?u=John](http://localhost:8080/borrows/search?u=John)

- **Register a Borrow**
  ```
  POST /borrows/register
  ```
  Example: [http://localhost:8080/borrows/register](http://localhost:8080/borrows/register)

  **Example Request Body**:
  ```json
  {
    "bookId": "cca31ff9-2d5f-4565-a31f-f92d5f6565d5",
    "userId": "85e773f6-66c7-42d2-a773-f666c7a2d2bc"
  }
  ```

### Authors

- **Get All Authors**
  ```
  GET /authors
  ```
  Example: [http://localhost:8080/authors](http://localhost:8080/authors)

## Configuration

Ensure you have the following environment variables set:

- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_PORT`: Database port
- `DB_HOST`: Database host
- `PORT`: Port on which the server runs (default is `8080`)

## Development

To run the application locally, ensure you have [Go](https://golang.org/dl/) installed and set up. Then, clone the
repository and run the following commands:

```sh
git clone https://github.com/juanPabloCano/library-management-system.git
cd library-management-system
go mod download
go run main.go