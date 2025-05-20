# Library API v1.0.0
### A self-hosted book management service with Elasticsearch backend

## Requirements
-  ``go``  
-  ``docker`` 

## Version
-  `` 1.0.0 ``  Basic CRUD operations • Pagination (limit/offset) • Docker support. ✅ Released
-  `` 2.0.0 ``  JWT authentication • Role-based access • Secure password hashing. In Progress
-  `` 3.0.0 ``  Redis rate limiting • Response caching. Planned



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


# API Reference

## Create Book
```http
  POST library/v1/books
```

```bash
  curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "isbn": "978-0134190440",
    "writer": "Alan A. A. Donovan & Brian W. Kernighan",
    "total_copies": 19,
    "on_loan": 3
  }' \
  "localhost:8080/library/v1/books"
```

#### Book Object Structure

| Field          | Type     | Required | Description                          | Example Value                                 |
|----------------|----------|----------|--------------------------------------|-----------------------------------------------|
| `title`        | `string` | Yes      | Title of the book                    | `"The Go Programming Language"`               |
| `isbn`         | `string` | Yes      | 10/13-digit ISBN                     | `"978-0134190440"`                            |
| `writer`       | `string` | Yes      | Author(s)                            | `"Alan A. A. Donovan & Brian W. Kernighan"`   |
| `total_copies` | `integer`| Yes      | Total inventory count (≥0)           | `19`                                          |
| `on_loan`      | `integer`| Yes      | Currently borrowed copies (≤total)   | `3`                                           |
| `created_at`   | `string` | Auto     | ISO-8601 creation timestamp          | `"2025-05-18T18:19:27Z"`                      |
| `updated_at`   | `string` | Auto     | ISO-8601 last update timestamp       | `"2025-05-20T09:15:42Z"`                      |

## Delete Book By ISBN
```http
  DELETE library/v1/books/{isbn}
```
```bash
  curl -X DELETE "localhost:8080/library/v1/books/978-0134190440"
```
| UrlParam     | Type     | Description |
| :----------  | :------  | :---------- |  
| `bookIsbn`   | `string` | book isbn number |

## Update Book By ISBN
```http
  PATCH library/v1/books/{isbn}
```
```bash
curl -X PATCH \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language - Updated Edition",
    "writer": "Alan A. A. Donovan & Brian W. Kernighan",
    "total_copies": 25,
    "on_loan": 5
  }' \
  "http://localhost:8080/library/v1/books/978-0134190440"
```

| Field          | Type     |
|----------------|----------|
| `title`        | `string` | 
| `writer`       | `string` | 
| `total_copies` | `integer`| 
| `on_loan`      | `integer`| 

## Get Book By ISBN
```http
  GET library/v1/books/{isbn}
```
```bash
  curl -X GET "localhost:8080/library/v1/books/978-0134190440"
```
| UrlParam     | Type     | Description |
| :----------  | :------  | :---------- |  
| `bookIsbn`   | `string` | book isbn number |



## Get all Books
```http
  GET library/v1/books
```
```bash
  curl -X GET "localhost:8080/library/v1/books?limit=10&offfset=0"
```
| Query     | Type     | Description |
| :-------  | :------  | :---------- |  
| `limit`   | `uint64` | per page item |
| `offset`  | `uint64` | starting index |



## Search Books
```http
  GET library/v1/books/search
```
```bash
  curl -X GET "localhost:8080/library/v1/books/search?query=Go&limit=100&offfset=0"
```

| Query     | Type     | Description |
| :-------  | :------  | :---------- |  
| `query`   | `string` | search item |
| `limit`   | `uint64` | per page item |
| `offset`  | `uint64` | starting index |


# Output

- ### For 'CreateBooks' request and 'GetBooksByISBN'  expected output
```bash
  {
    "title": "The Go Programming Language",
    "isbn": "978-0134190440",
    "writer": "Alan A. A. Donovan & Brian W. Kernighan",
    "total_copies": 19,
    "on_loan": 3,
    "created_at": "2024-05-21T10:15:30Z",
    "updated_at": "2024-05-21T14:30:45Z"
 }
```



- ### For 'AllBooks' request and 'SearchBooks' expected out
```bash
{
  "total_page": 1,
  "total_book": 1,
  "books": [
    {
      "title": "The Go Programming Language",
      "isbn": "978-0134190440",
      "writer": "Alan A. A. Donovan & Brian W. Kernighan",
      "total_copies": 19,
      "on_loan": 3,
      "created_at": "2024-05-21T10:15:30Z",
      "updated_at": "2024-05-21T14:30:45Z"
    }
  ]
}
```
- ### For 'Delete' request and 'Update' expected out
```bash
{
  "msg":string
}
```

