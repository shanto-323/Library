# Library API v2.0.0
### This is a self-hosted backend service for managing a digital library, designed with scalability and security in mind. It supports user-based operations, full-text book search, and robust access control using JWT authentication.

## Requirements
-  ``go``  
-  ``docker`` 

## Version
-  `` 1.0.0 ``  Basic CRUD operations • Pagination (limit/offset) • Docker support. ✅ Released
-  `` 2.0.0 ``  JWT authentication • Role-based access • Secure password hashing. ✅ Released
-  `` 3.0.0 ``  Redis rate limiting • Response caching. Planned

### Features

-  **JWT Authentication** with access & refresh tokens
-  **Full-text Search** powered by Elasticsearch
-  **Relational Data Storage** using PostgreSQL
-  **Password Hashing** with bcrypt for secure credentials
-  **Role-based Access** (ADMIN/USER)

### Tech Used
- **Go (Golang)** – Backend server
- **PostgreSQL** – Primary database
- **Elasticsearch** – Search engine
- **bcrypt** – Secure password hashing
- **JWT** – Auth tokens for session management

## Run Locally

Clone the project

```bash
  git clone https://github.com/shanto-323/Library.git
```

Go to the project directory

```bash
  cd project/path/Library
```

Start the server

```bash
  docker-compose up --build -d
```


## API Reference


### Overview
- All routes are under the base path `/library/v2`.
- JWT authentication is required for protected routes.
- Role-based access control: 'USER' and 'ADMIN'.
- Access tokens are **Short-Lived(20 munite)** and provided via cookie.
- Refresh tokens are stored securely in the database.
- All passwords are hashed using bcrypt before being stored.
- Refresh token support is available via `/user/token/{id}`.
- No CORS support
<br><br><br>

### User Endpoints

| Method | Endpoint                   | Description                       | Auth Required | User Type |
|--------|----------------------------|-----------------------------------|----------------|------------|
| POST   | `/user/signup`             | Register a new user               | No             | -          |
| POST   | `/user/login`              | Log in with email and password    | No             | -          |
| GET    | `/user/token/{id}`         | Get a new access token using refresh token (include header `x-jwt-token`)           | No             | -          |
| POST   | `/user/logout/{id}`        | Log out and invalidate tokens     | Yes            | USER/ADMIN |
| PATCH  | `/user/{id}`               | Update user info                  | Yes            | USER/ADMIN |
| DELETE | `/user/delete/{id}`        | Delete user account               | Yes            | USER/ADMIN |
| GET    | `/user/admin/all`          | Get all users (admin only)        | Yes            | ADMIN      |
---

### Book Endpoints

| Method | Endpoint                   | Description                       | Auth Required | User Type |
|--------|----------------------------|-----------------------------------|----------------|------------|
| GET    | `/books/search?query=...`  | Search for books                  | Yes            | USER/ADMIN |
| GET    | `/books/{isbn}`            | Get book by ISBN                  | Yes            | USER/ADMIN |
| GET    | `/books`                   | Get all books                     | Yes            | USER/ADMIN |
| POST   | `/books/admin`             | Add a new book                    | Yes            | ADMIN      |
| PATCH  | `/books/admin/{isbn}`      | Update book information           | Yes            | ADMIN      |
| DELETE | `/books/admin/{isbn}`      | Delete a book                     | Yes            | ADMIN      |

---