# DockerMonitoringApp

<div align="center">
    <img src="imgs/logo.png" width=40%>
</div>

<br>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
![GitHub Issues](https://img.shields.io/github/issues/k6zma/DockerMonitoringApp.svg)
![GitHub Pull Requests](https://img.shields.io/github/issues-pr/k6zma/DockerMonitoringApp.svg)

</div>

<div align="center">
  <img src="https://img.shields.io/badge/Go-black?style=for-the-badge&logo=go&logoColor=#00ADD8"/>
  <img src="https://img.shields.io/badge/Next.js-black?style=for-the-badge&logo=next.js&logoColor=#000000"/>
  <img src="https://img.shields.io/badge/TypeScript-black?style=for-the-badge&logo=typescript&logoColor=#3178C6"/>
  <img src="https://img.shields.io/badge/PostgreSQL-black?style=for-the-badge&logo=postgresql&logoColor=#4169E1"/>
  <img src="https://img.shields.io/badge/Docker-black?style=for-the-badge&logo=docker&logoColor=#2496ED"/>
  <img src="https://img.shields.io/badge/Nginx-black?style=for-the-badge&logo=nginx&logoColor=green"/>
  <img src="https://img.shields.io/badge/github-black?style=for-the-badge&logo=github&logoColor=#000000"/>
</div>

<br>


---

## **Starting the Project**
Before starting the application, create a .env file in the project root. You can use the provided .env.example as a reference for the required environment variables

The project is deployed using **Docker Compose**. All services are containerized and can be started with a single command:

```bash
docker-compose -f dev.docker-compose.yml up --build -d
```

After a successful startup, the system will be available at the following addresses:
- **Frontend**: [http://localhost](http://localhost)
- **Backend API**: [http://localhost/api/v1](http://localhost/api/v1)
- **Swagger Documentation**: [http://localhost/swagger](http://localhost/swagger)


---

## **Technical Specification**

### **Description**  
The task is to develop an application for continuous monitoring of Docker container statuses. The application should retrieve container IP addresses, periodically ping them, and store the data in a database. The retrieved data should be displayed on a dynamically generated web page.


## **Tasks**  
The application should be developed using Go and JavaScript (TypeScript) and perform the following functions:  
- Retrieve the IP addresses of running Docker containers.  
- Periodically ping them at a specified interval.  
- Store the obtained data in a database.  
- Provide access to container status data through a dynamically updated web page.

## **Services**  
The project requires the development of four services:

1. **Backend Service**  
   - Provides a RESTful API for retrieving data from the database and adding new records.  

2. **Frontend Service**  
   - Developed in JavaScript using any UI library (preferably React).  
   - Fetches data through the Backend API.  
   - Displays data in a table format with columns: IP address, ping time, and the date of the last successful attempt.  
   - For data presentation in HTML, Bootstrap, Ant Design, or similar libraries can be used.  

3. **Database**  
   - PostgreSQL.  

4. **Pinger Service**  
   - Retrieves a list of all running Docker containers.  
   - Pings them.  
   - Sends the collected data to the database via the Backend API.  

Additionally, the following complexity enhancements are possible:  
- Adding Nginx.  
- Implementing a message queue service.  
- Using netns.  
- Separate configuration for the service with verification.  

## **Result**  
- Each service must have its own **Dockerfile**.  
- A common **Docker Compose** file should be created to build and run all services.  
- After launching the system, container status data should be accessible via HTTP after the first polling cycle.  
- The code must be hosted in a **GitHub/GitLab** repository with a **README.md** file containing a brief description of functionality and setup instructions.  

---

## **Project Description**

## **Backend Service**

### **1. Configuration**
The service configuration is located in the `config.json` file, which contains database connection settings, API configurations, and other parameters. The **Viper** package is used for convenient configuration management

Path to the configuration file:  
```
backend/config.json
```

Example `config.json`:
```json
{
  "server": {
    "port": 8080
  },
  "database": {
    "host": "postgres_db",
    "port": 5432,
    "user": "admin",
    "password": "password",
    "dbname": "docker_monitoring"
  },
  "auth": {
    "api_key": "your-api-key"
  }
}
```

### **2. Architecture (Layered Model)**
The backend service is built using **Clean (Layered) Architecture**, where the code is divided into several layers:

```md
backend/
├── cmd/server            # Entry point
│   └── main.go
├── internal/
│   ├── application/      # Business logic
│   │   ├── dto/          # Data structures for handling containers
│   │   ├── usecases/     # Core business logic
│   │   ├── repositories/ # Database access
│   ├── domain/           # Core system entities
│   ├── infrastructure/   # Configuration, logging, database interactions
│   │   ├── config/       # Configuration management
│   │   ├── db/postgres/  # Database handling logic
│   │   ├── migrations/   # Database migrations
|   |   ├── flags/        # Flags managment
│   ├── presentation/     # User interaction
│   │   ├── handlers/     # HTTP request processing
│   │   ├── middlewares/  # Authentication, CORS, logging
│   │   ├── routes/       # API routing
│   │   ├── mapper/       # DTO mapper
│   │   ├── server/       # Server implementation
├── docs/                 # API documentation (Swagger)
├── migrations/           # SQL migration files
├── pkg/                  # Utility functions (logging, IP collection, etc.)
├── config.json           # Main configuration file
├── Dockerfile            # Dockerfile for backend service
├── go.mod                # Go module dependencies
└── go.sum                # Go dependencies lock file

```

### **3. Handling HTTP Requests (REST API)**  

The backend service implements a **REST-API** that allows managing container statuses. All request handlers are located in:  
```
backend/internal/presentation/handlers/
```

### **Main API Endpoints**  

| Method     | Endpoint                                  | Description                                   |
|------------|-------------------------------------------|-----------------------------------------------|
| **GET**    | `/api/v1/container_status`                | Retrieve a list of containers (with filters)  |
| **POST**   | `/api/v1/container_status`                | Create a new container entry                  |
| **PATCH**  | `/api/v1/container_status/{container_id}` | Update a container by ID                      |
| **DELETE** | `/api/v1/container_status/{container_id}` | Delete a container by ID                      |


### **Detailed API Description**  

#### **1. Retrieve a List of Containers**  
##### **GET** `/api/v1/container_status`  

Returns a list of containers with optional filtering parameters.  

##### **Query Parameters (Optional Filters):**  
| Parameter         | Type      | Description                                           |
|------------------|----------|------------------------------------------------------|
| `container_id`   | `string`  | Filter by container ID                               |
| `ip`            | `string`  | Filter by IP address                                |
| `name`          | `string`  | Filter by container name                           |
| `status`        | `string`  | Filter by status (running, exited, etc.)           |
| `ping_time_min` | `number`  | Minimum ping time                                  |
| `ping_time_max` | `number`  | Maximum ping time                                  |
| `created_at_gte` | `string`  | Filter by creation date (≥, RFC3339 format)         |
| `created_at_lte` | `string`  | Filter by creation date (≤, RFC3339 format)         |
| `updated_at_gte` | `string`  | Filter by last update date (≥, RFC3339 format)      |
| `updated_at_lte` | `string`  | Filter by last update date (≤, RFC3339 format)      |
| `limit`         | `integer` | Limit the number of returned records               |

##### **Response:**  
```json
[
    {
        "container_id": "abc123",
        "ip_address": "192.168.1.10",
        "name": "nginx-container",
        "status": "running",
        "ping_time": 15.2,
        "last_successful_ping": "2025-02-09T12:34:56Z",
        "created_at": "2025-02-08T10:00:00Z",
        "updated_at": "2025-02-09T12:35:00Z"
    }
]
```


#### **2. Create a New Container Entry**  
##### **POST** `/api/v1/container_status`  

Adds a new container to the database.  

##### **Request Body:**  
```json
{
    "container_id": "abc123",
    "ip_address": "192.168.1.10",
    "name": "nginx-container",
    "status": "running",
    "ping_time": 15.2,
    "last_successful_ping": "2025-02-09T12:34:56Z"
}
```

##### **Response:**  
```json
{
    "container_id": "abc123",
    "ip_address": "192.168.1.10",
    "name": "nginx-container",
    "status": "running",
    "ping_time": 15.2,
    "last_successful_ping": "2025-02-09T12:34:56Z",
    "created_at": "2025-02-08T10:00:00Z",
    "updated_at": "2025-02-09T12:35:00Z"
}
```

##### **Possible Responses:**  
- **`201 Created`** - Container added successfully  
- **`400 Bad Request`** - Invalid input data  
- **`500 Internal Server Error`** - Server-side issue  


#### **3. Update a Container by ID**  
##### **PATCH** `/api/v1/container_status/{container_id}`  

Updates the information of a container (partial update).  

##### **Path Parameter:**  
| Parameter      | Type    | Description           |
|--------------|--------|----------------------|
| `container_id` | `string` | ID of the container |

##### **Request Body (only include fields to update):**  
```json
{
    "name": "updated-nginx-container",
    "status": "restarting",
    "ping_time": 20.5
}
```

##### **Response:**  
- **`204 No Content`** - Updated successfully  
- **`400 Bad Request`** - Invalid input data  
- **`500 Internal Server Error`** - Server-side issue  

#### **4. Delete a Container by ID**  
##### **DELETE** `/api/v1/container_status/{container_id}`  

Removes a container record from the database.  

##### **Path Parameter:**  
| Parameter      | Type    | Description           |
|--------------|--------|----------------------|
| `container_id` | `string` | ID of the container |

##### **Response:**  
- **`204 No Content`** - Deleted successfully  
- **`404 Not Found`** - Container not found  
- **`500 Internal Server Error`** - Server-side issue  


### **Authentication & Security**  
All endpoints require authentication via API Key. Clients must include the following HTTP header in requests:  
```http
X-Api-Key: your-api-key
```
Without this key, the server will return **`401 Unauthorized`**.

### **4. Database Interaction**  

The backend service uses **PostgreSQL** as the database, with the **pgx** driver for efficient database interactions

#### **Migrations**  
Database migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate) and are stored in:  
```
backend/migrations/
```
Migrations are executed **directly in the application code**, and the migration behavior is configurable via the **config file**:  
```json
"migrations": {
  "path": "/app/migrations",
  "type": "apply"
}
```
The `type` field allows manipulation of migration execution behavior. 

#### **Database Schema**  

The **`container_status`** table is used to store container information:  

```sql
CREATE TABLE container_status (
    container_id TEXT PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    ip_address INET NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'created',
    ping_time DOUBLE PRECISION NULL,
    last_successful_ping TIMESTAMP,
    updated_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
```

#### **Indexes for Query Optimization**  
To enhance query performance, the following indexes are applied:  
```sql
CREATE INDEX idx_last_successful_ping ON container_status(last_successful_ping);
CREATE INDEX idx_updated_at ON container_status(updated_at);
```
These indexes optimize retrieval of records based on recent updates and successful pings


### **5. Swagger Documentation**
**Swagger** is used for API documentation. Documentation files are located in:
```
backend/docs/
```

To view the documentation in a browser:
[http://localhost/swagger](http://localhost/swagger)

---

### **5. Pinger Service**  

### **Configuration**  

The Pinger service is configurable via `config.json`, which defines key parameters such as the ping interval, Docker socket path, and backend API connection details.  

```json
{
  "ping": {
    "ping_interval": "5s"
  },
  "docker": {
    "socket_path": "/var/run/docker.sock"
  },
  "backend": {
    "url": "http://backend_service:8080",
    "api_key": "your-api-key"
  }
}
```
- **`ping_interval`** – Defines how often the service pings active containers
- **`socket_path`** – Specifies the path to the Docker daemon socket for retrieving container information
- **`backend.url`** – API endpoint of the Backend Service where ping results are sent
- **`backend.api_key`** – Authentication key for the Backend API

---

### **Architecture**  

The service follows a **layered (onion) architecture**, ensuring a separation of concerns and maintainability

```
pinger/
├── cmd/pinger
│   └── main.go              # Entry point (main.go)
├── internal/
│   ├── application/         # Business logic
│   │   ├── repositories/    # Interfaces for data sources
│   │   ├── usecases/        # Core pinging logic
│   ├── domain/              # Core system entities
│   ├── infrastructure/      # External service interactions
│   │   ├── backend/         # Communication with backend API
│   │   ├── config/          # Configuration management
│   │   ├── docker/          # Interaction with Docker API
│   │   ├── flags/           # Command-line flag parsing
│   └── pkg/
│       └── utils/           # Logging utilities
```

---

### **How It Works**  

1. **Retrieving Container Data**  
   - The service connects to the **Docker daemon** via sock path.
   - It fetches all running containers and extracts their **IP addresses**
   - This logic is implemented in `internal/infrastructure/docker/container_repository.go`

2. **Pinging Containers**  
   - The service uses [`pro-bing`](https://github.com/prometheus-community/pro-bing) to perform ping requests 
   - Pings are executed at the interval defined in `ping_interval`
   - The **ping results** (latency, success/failure) are processed and formatted
   - The core pinging logic is implemented in `internal/application/usecases/pinger_usecase.go`

3. **Sending Data to the Backend**  
   - After each ping, results are **sent via REST API** to the **Backend Service**.  
   - API interaction is handled in `internal/infrastructure/backend/status_repository.go`.  
   - The service authenticates using the **API key** configured in `config.json`.  

