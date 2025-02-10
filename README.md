# Go API for User Authentication with JWT

This project is a REST API built with **Go** and the **Gin Gonic** framework, providing user authentication services. The API allows the following operations:

- **User Creation** (with email and password)
- **Login** (returns an access token and a refresh token)
- **Token Refresh** (returns new access and refresh tokens)

## Features

- **User Creation**: Allows users to register with an email and password. Passwords are securely hashed using **Argon2id**.
- **Login**: Authenticates a user and returns an **access token** and **refresh token** as **JWT** (JSON Web Tokens).
- **Token Refresh**: Allows a user to refresh their access token by providing the refresh token.
- **JWT Authentication**: Access and refresh tokens are generated and validated using **JWT** for secure authentication.
- **PostgreSQL Database**: All user data is stored in a **PostgreSQL** database.

## Libraries and Technologies Used

- **Gin Gonic**: A fast HTTP web framework for Go.
- **PostgreSQL**: A relational database for storing user credentials and tokens.
- **Argon2id**: A secure password-hashing algorithm to protect user passwords.
- **JWT (jsonwebtoken)**: Used to generate and validate access tokens and refresh tokens.
- **Air**: A live-reload tool to improve the development workflow.

## Prerequisites

Before running the project, you need to have the following installed:

- **Go**: Version 1.22 or later.
- **PostgreSQL**: Version 12 or later.
- **Air**: For live reloading during development. (Optional)
- **Make**: To use the `Makefile` commands (Optional).

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/pedrotunin/go-jwt-auth.git
    cd yourproject
    ```

2. Install Go dependencies:

    ```bash
    go mod tidy
    ```

3. Set up your PostgreSQL database:
    - Create a PostgreSQL database and user.
    - Configure the connection in your application by editing the `.env` file with the necessary database credentials.

    Example `.env` file:

    ```env
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=your_database
    HMAC_SECRET=
    PORT=8080
    MODE=DEBUG # DEBUG or PRODUCTION
    ```

4. **(Optional)** Install Air for live-reloading during development:

    ```bash
    go install github.com/cosmtrek/air@latest
    ```

5. **(Optional)** If you want to use the provided `Makefile`, you can run the following:

    ```bash
    make run
    ```

## Running the API

You can start the API server in different ways:

- **Using Air (Development Mode)**:

    If you installed **Air**, simply run:

    ```bash
    air
    ```

- **Without Air (Production Mode)**:

    Run the Go application directly:

    ```bash
    go run cmd/api/api.go
    ```

This will start the API server on port `8080` by default.