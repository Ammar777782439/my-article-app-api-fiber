# my-article-app: Article and Author Management System (Go, Fiber, GORM)

`my-article-app` is an integrated web application built with Go that provides a RESTful API for managing articles and authors. The application is built using the Fiber framework and the GORM library for interacting with PostgreSQL, following a clean and scalable architecture.

---

## ✨ Key Features

- **Comprehensive Article Management:**
  - Create new articles and link them to an author.
  - Read all articles or a specific article (including author data).
  - Update article data.
  - Delete articles.
- **Comprehensive Author Management:**
  - Create new authors.
  - Read all authors or a specific author (including their list of articles).
  - Update author data.
  - Delete authors (considering associated articles).
- **Strong Data Relationships:**
  - One-to-many relationship between authors and articles.
  - Preloading of related data for improved performance.
- **Data Validation:**
  - Use of the `go-playground/validator` library to ensure the integrity of input data.
- **Clean Application Architecture:**
  - Handlers, Use Cases, Repositories, Models.
- **Easy Setup and Operation:**
  - Use of Docker and Docker Compose to easily run a PostgreSQL database.

---

## 🛠️ Technologies Used

- **Programming Language:** Go (Golang) 1.18+
- **Web Framework:** Fiber v2
- **ORM:** GORM
- **Database:** PostgreSQL
- **Data Validation:** `go-playground/validator` v10
- **Containers:** Docker & Docker Compose

---

## 🏛️ Project Structure


 * my-article-app/ (Root directory of your project)
   * main.go (Main application entry point)
   * article.go (Article model and related logic - e.g., handlers, repository)
   * author.go (Author model and related logic)
   * database.go (Database connection and setup)
   * go.mod (Manages project dependencies)
   * go.sum (Checksums of dependencies)
   * README.md (Project description and setup instructions)
   * docker-compose.yml (Optional: for Docker setup)



**Request Flow:**
1. An HTTP request arrives at the Fiber server.
2. Fiber routes the request to the appropriate Handler.
3. The Handler parses the request and validates the input.
4. The Handler calls the appropriate UseCase.
5. The UseCase coordinates operations with the Repository.
6. The Repository executes database operations using GORM.
7. Results are returned through the layers to the Handler, which sends the response to the client.

---

## 🚀 Prerequisites and Setup

### Prerequisites:
- **Go:** Version 1.18 or later
- **Docker and Docker Compose:** To run the PostgreSQL database (optional)
- **Git:** To clone the repository

### Setup Steps:

1.  **Clone the repository:**
    ```bash
    git clone <repository-link>
    cd my-article-app
    ```
2.  **Set up the database (using Docker Compose):**
    - A pre-configured `docker-compose.yml` file is provided:
    ```yaml
    version: '3.8'
    services:
      postgres_db:
        image: postgres:13-alpine
        container_name: my_article_app_db
        environment:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: article_db
        ports:
          - "5432:5432"
        volumes:
          - postgres_data:/var/lib/postgresql/data
    volumes:
      postgres_data:
    ```
    - To run the database:
    ```bash
    docker-compose up -d
    ```
    - **Note:** Ensure that the connection settings in `internal/database/gorm.go` (the `dsn` variable) match your Docker or local PostgreSQL settings.

3.  **Install dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Run the application:**
    ```bash
    go run ./cmd/api/main.go
    ```
    A message will appear indicating that the server is running on port `3000`.

---

## 📡 API Endpoints

### Authors

| Method | Path              | Description            | Request Body (Example)                               | Successful Response (Example) |
|--------|-------------------|------------------------|------------------------------------------------------|-----------------------------|
| POST   | /authors          | Create a new author    | `{ "name": "Author Name", "email": "email@example.com" }` | 201 Created                 |
| GET    | /authors          | Fetch all authors      | (None)                                               | 200 OK                      |
| GET    | /authors/{id}     | Fetch a specific author| (None)                                               | 200 OK                      |
| PUT    | /authors/{id}     | Update author data     | `{ "name": "New Name", "email": "email@domain.com" }`  | 200 OK                      |
| DELETE | /authors/{id}     | Delete an author       | (None)                                               | 204 No Content              |

### Articles

| Method | Path              | Description            | Request Body (Example)                                   | Successful Response (Example) |
|--------|-------------------|------------------------|----------------------------------------------------------|-----------------------------|
| POST   | /articles         | Create a new article   | `{ "title": "Title", "content": "...", "author_id": 1 }` | 201 Created                 |
| GET    | /articles         | Fetch all articles     | (None)                                                 | 200 OK                      |
| GET    | /articles/{id}    | Fetch a specific article| (None)                                                 | 200 OK                      |
| PUT    | /articles/{id}    | Update article data    | `{ "title": "New Title", "content": "...", "author_id": 1 }` | 200 OK                      |
| DELETE | /articles/{id}    | Delete an article      | (None)                                                 | 204 No Content              |

### Health Check

| Method | Path    | Description                |
|--------|---------|----------------------------|
| GET    | /health | Check if the application is running |

---

## 🧪 `curl` Examples for Testing Endpoints

1.  **Create a new author:**
    ```bash
    curl -X POST -H "Content-Type: application/json" \
         -d '{"name": "Author Name", "email": "email@example.com"}' \
         http://localhost:3000/api/v1/authors
    ```

2.  **Create a new article (with author_id):**
    ```bash
    curl -X POST -H "Content-Type: application/json" \
         -d '{"title": "Article Title", "content": "Content...", "author_id": 1}' \
         http://localhost:3000/api/v1/articles
    ```

3.  **Fetch all articles:**
    ```bash
    curl http://localhost:3000/api/v1/articles
    ```

4.  **Fetch a specific article:**
    ```bash
    curl http://localhost:3000/api/v1/articles/1
    ```

---

## 💡 Suggested Future Improvements and Developments

- **Authentication and Authorization:** Add JWT or similar to protect endpoints.
- **Pagination, Sorting, and Filtering:** Support for these features.
- **Testing:** Write unit and integration tests.
- **Logging and Monitoring:** Integrate logrus/zap and Prometheus/Grafana.
- **Interactive API Documentation:** Use Swagger/OpenAPI (Swaggo).
- **Improved Error Handling:** Provide standard error codes and descriptions.
- **Use Environment Variables:** To manage sensitive configurations.

---


## 🤝 Contributing

1.  Fork the Project.
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

---

## 📄 License

This project is licensed under the MIT License.


