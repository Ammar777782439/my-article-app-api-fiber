# **Project: my-article-app \- An Article and Author Management System**

### **Abstract**

The my-article-app system constitutes a fully integrated web application, engineered utilizing the Go programming language to furnish a RESTful Application Programming Interface for the comprehensive management of articles and their respective authors. The implementation leverages the Fiber framework in conjunction with the GORM library for object-relational mapping to a PostgreSQL database, adhering to a paradigm of clean, scalable architecture predicated on the principle of separation of concerns.

## **1\. System Capabilities**

The core functionalities of the system are delineated as follows:

* **Comprehensive Article Entity Management:** Full lifecycle management of article entities, including creation, retrieval, modification, and termination.  
* **Comprehensive Author Entity Management:** Full lifecycle management of author entities, encompassing creation, retrieval, modification, and termination.  
* **Architectural Paradigm:** The system is structured upon a clean architectural model, featuring distinct layers for handlers, data transfer objects, business logic, and data persistence.  
* **Optimized Data Relationship Management:** An efficient Preloading strategy is employed for the retrieval of associated data entities, such as an author's articles, within a single database query. This approach is designed to obviate common performance deficiencies, most notably the N+1 query problem.  
* **Input Data Integrity Verification:** The go-playground/validator library is integrated for the programmatic validation of all incoming data, thereby ensuring its structural and syntactical integrity.  
* **System Deployment and Initialization:** The operational environment is streamlined through the use of Docker and Docker Compose, facilitating a simplified and reproducible deployment of the PostgreSQL database instance.

## **2\. Architectural Design**

The project's structural design adheres to a multi-layered architectural pattern, implemented to enforce a strict separation of concerns, thereby enhancing modularity and long-term maintainability of the codebase.

my-article-app/  
├── cmd/api/main.go              \# نقطة الدخول الرئيسية للتطبيق  
├── internal/  
│   ├── database/                \# إعداد اتصال قاعدة البيانات  
│   ├── dto/                     \# (DTOs) كائنات نقل البيانات للـ API  
│   ├── handlers/                \# معالجات طلبات HTTP (تتعامل مع DTOs)  
│   ├── models/                  \# نماذج GORM (تمثل جداول قاعدة البيانات)  
│   ├── repository/              \# طبقة الوصول إلى البيانات (تتعامل مع Models)  
│   └── usecase/                 \# طبقة منطق العمل (الجسر بين DTOs و Models)  
├── go.mod  
├── go.sum  
└── docker-compose.yml


### **2.1. Request Processing Flow**

The procedural sequence for request processing is delineated as follows:

1. An inbound HTTP request is received by the Fiber server process.  
2. The request is routed by the Fiber engine to the designated **Handler**.  
3. The Handler is responsible for parsing the request payload into its corresponding **Request Data Transfer Object (DTO)** and subsequently validating its contents.  
4. The Handler invokes the appropriate **UseCase** method, passing the validated DTO as an argument.  
5. The UseCase executes the core business logic, which includes the transformation of the DTO into a database **Model**.  
6. The UseCase delegates the data persistence operation to the **Repository** layer, providing the populated Model.  
7. The Repository interfaces with the database via the GORM library to execute the requisite query.  
8. The data retrieval or modification result is propagated back through the layers; the UseCase transforms the resultant Model into a **Response DTO**.  
9. The Handler receives the Response DTO, serializes it into a JSON format, and transmits it as the final HTTP response to the originating client.

## **3\. Technological Stack Utilized**

* **Programming Language:** Go (Golang), version 1.18 or subsequent.  
* **Web Application Framework:** Fiber, version 2\.  
* **Object-Relational Mapper (ORM):** GORM.  
* **Database Management System:** PostgreSQL.  
* **Data Validation Library:** go-playground/validator, version 10\.  
* **Containerization Technology:** Docker & Docker Compose.

## **4\. System Prerequisites and Configuration Protocol**

### **4.1. Prerequisites**

Prior to system initialization, the following software components must be installed and configured on the host machine:

* Go: Version 1.18 or a more recent release.  
* Docker & Docker Compose: For containerized database deployment.  
* Git: For version control and repository cloning.

### **4.2. Configuration Steps**

The protocol for system setup is as follows:

1. **Repository Cloning:** The source code repository is to be cloned via Git.  
   git clone \<repository-link\>  
   cd my-article-app

2. **Database Instance Initialization:** A pre-configured docker-compose.yml file is supplied to orchestrate the database container.  
   * The database service can be initiated by executing the command below from the project's root directory:  
     docker-compose up \-d

   * **Note:** It is imperative that the database connection string (DSN) specified within internal/database/gorm.go aligns with the configuration of the local or containerized PostgreSQL instance.  
3. **Dependency Installation:** Project dependencies must be resolved and installed.  
   go mod tidy

4. **Application Execution:** The main application binary can be run.  
   go run ./cmd/api/main.go

   Upon successful startup, a confirmation message will be logged, indicating that the server is listening for connections on port 3000\.

## **5\. Application Programming Interface (API) Endpoints**

The API employs Data Transfer Objects (DTOs) for both inbound requests and outbound responses to ensure security, abstraction, and clarity in data contracts.

### **5.1. Author Endpoints**

| Method | Path | Description | Request Body (Example) | Successful Response (Example) |
| ----: | ----: | ----: | ----: | ----: |
| POST | /api/v1/authors | Creates a new author entity. | {"name": "Ahmed", "email": "a@a.com"} | 201 Created with AuthorResponse |
| GET | /api/v1/authors | Retrieves a collection of all authors. | (None) | 200 OK with \[\]AuthorResponse |
| GET | /api/v1/authors/{id} | Retrieves a specific author entity. | (None) | 200 OK with AuthorDetailResponse |
| PUT | /api/v1/authors/{id} | Updates an existing author entity. | {"name": "Ahmed New"} | 200 OK with AuthorResponse |
| DELETE | /api/v1/authors/{id} | Deletes an author entity. | (None) | 204 No Content |

### **5.2. Article Endpoints**

| Method | Path | Description | Request Body (Example) | Successful Response (Example) |
| ----: | ----: | ----: | ----: | ----: |
| POST | /api/v1/articles | Creates a new article entity. | {"title": "New Article", "content": "...", "author\_id": 1} | 201 Created with ArticleResponse |
| GET | /api/v1/articles | Retrieves a collection of all articles. | (None) | 200 OK with \[\]ArticleResponse |
| GET | /api/v1/articles/{id} | Retrieves a specific article entity. | (None) | 200 OK with ArticleResponse |
| PUT | /api/v1/articles/{id} | Updates an existing article entity. | {"title": "Updated Title"} | 200 OK with ArticleResponse |
| DELETE | /api/v1/articles/{id} | Deletes an article entity. | (None) | 204 No Content |

## **6\. Exemplary cURL Invocations for Endpoint Verification**

1. **Author Creation:**  
   curl \-X POST \-H "Content-Type: application/json" \\  
        \-d '{"name": "Ammar", "email": "ammar@example.com"}' \\  
        http://localhost:3000/api/v1/authors

2. **Article Creation:** (Requires a pre-existing author with the specified author\_id)  
   curl \-X POST \-H "Content-Type: application/json" \\  
        \-d '{"title": "Article about Go", "content": "Content here...", "author\_id": 1}' \\  
        http://localhost:3000/api/v1/articles

3. **Article Collection Retrieval:** (Note the embedded author data in the response)  
   curl http://localhost:3000/api/v1/articles

4. **Author Retrieval:** (Note the embedded article data in the response)  
   curl http://localhost:3000/api/v1/authors/1

## **7\. Recommendations for Future System Enhancements**

It is posited that future development iterations could benefit from the incorporation of the following enhancements:

* **Authentication and Authorization:** Implementation of robust security mechanisms, such as JSON Web Tokens (JWT) or analogous protocols, is advisable to protect the API endpoints.  
* **Result Set Manipulation:** The introduction of support for pagination, sorting, and filtering functionalities for GET requests would arguably improve the API's utility.  
* **Automated Testing:** The development of a comprehensive suite of unit and integration tests would substantially increase code reliability and maintainability.  
* **Logging and Monitoring:** Integration of advanced logging libraries (e.g., logrus, zap) and monitoring systems (e.g., Prometheus, Grafana) is recommended for production environments.  
* **Interactive API Documentation:** Generation of interactive API documentation using standards such as OpenAPI (via tools like Swaggo) may improve the developer experience.  
* **Configuration Management:** The externalization of sensitive configuration parameters through the use of environment variables is a recommended security best practice.