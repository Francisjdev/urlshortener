URL Shortener (Go + PostgreSQL)

A simple URL shortener built in Go using PostgreSQL for persistence.

This project was built to practice backend fundamentals without using frameworks or ORMs. The focus is on clean architecture, interface-based design, and proper database handling.

Features

Create short URLs from long URLs

Redirect using short codes

Track hit counts

Persist data in PostgreSQL

In-memory repository available for testing

Tech Stack

Go (standard library net/http)

PostgreSQL

pgx driver (via database/sql)

Raw SQL (no ORM)

Project Structure

The application follows a layered architecture:

HTTP Handler → Service → Repository Interface → Repository Implementation

Handler Layer

Handles HTTP routing

Parses and validates JSON

Returns proper HTTP responses

Service Layer

Contains business logic

Generates short codes

Coordinates repository operations

Repository Layer

Defined as an interface:

type URLRepository interface {
    Create(ctx context.Context, url *model.URL) error
    GetByCode(ctx context.Context, code string) (*model.URL, error)
    IncrementHitCount(ctx context.Context, code string) error
}


Two implementations exist:

In-memory repository (for development/testing)

PostgreSQL repository (persistent storage)

The service depends only on the interface.
The concrete repository implementation is injected in main (dependency injection).

Database Schema
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NULL,
    hit_count INT NOT NULL DEFAULT 0
);


The unique constraint on code ensures short codes cannot collide.

Running Locally
1. Ensure PostgreSQL is running

Create a database named:

urlshortener

2. Run the server

From the project root:

go run .


Server runs on:

http://localhost:8080

Example Usage
Create a short URL
curl -X PUT \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.example.com"}' \
  http://localhost:8080/create


Example response:

{
  "code": "AbC123",
  "short_url": "http://localhost:8080/AbC123"
}

Redirect
curl -v http://localhost:8080/AbC123


Returns:

HTTP/1.1 302 Found
Location: https://www.example.com

Check Hit Count in PostgreSQL
SELECT code, hit_count FROM urls WHERE code = 'AbC123';

Design Decisions

No ORM — raw SQL via database/sql

Database constraints enforce uniqueness

Database-specific errors (e.g. unique violation 23505) are mapped to domain errors

Clean separation between HTTP, business logic, and persistence

Repository pattern allows swapping storage implementations

What I Learned

Structuring backend applications with clear boundaries

Using interfaces to decouple layers

Handling database errors correctly

Implementing dependency injection manually in Go

Working directly with PostgreSQL and SQL queries

Future Improvements

URL expiration logic

Statistics endpoint

Environment-based configuration

Docker setup

Deployment

Automated tests
