## Quick Start Guide

### Prerequisites:
- Go (≥ 1.24.6);
- Postgresql;
- Goose (optional);
- Docker (optional).

### Running Locally:

- **Step 1**: Create a `config.env` file with environment variables, for example:

```bash
BIND_IP=127.0.0.1
LISTEN_PORT=8001
PSQL_HOST=your_db_host
PSQL_PORT=your_db_port
PSQL_NAME=your_db_name
PSQL_USER=your_db_user
PSQL_PASSWORD=your_db_password
```

- **Step 2**: Install `goose` migration tool (optional):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- **Step 3**: Apply database migrations (optional):

```bash
goose -dir=migrations postgres \
"host=your_db_host port=your_db_port dbname=your_db_name user=your_db_user password=your_db_password sslmode=disable" up
```

- **Step 4**: Build and run the server:

```bash
make build
make run
```

---

### Running with Docker:

- **Step 1**: Create a `config.env` file with environment variables, for example:

```bash
BIND_IP=0.0.0.0
LISTEN_PORT=8001
PSQL_HOST=task4-db
PSQL_PORT=5432
PSQL_NAME=postgres
PSQL_USER=postgres
PSQL_PASSWORD=postgres
DOCKER_SERVICE_PORT=8001
DOCKER_PSQL_PORT=25432
```

- **Step 2**: Start the containerized application:

```bash
make docker-compose-up
```

### Web Interface:

For better visualization and hands-on testing of the API, a web interface has been developed. It allows you to perform full CRUD operations (Create, Read, Update, Delete) and test complex filtering/pagination logic in real-time.

- **Location**: The web interface files are located in the web/ directory (main file: web/index.html).
- **Configuration**: By default, the frontend is configured to connect to the server at 127.0.0.1:8001 (set in web/config.js).
- **Important**: If you change the BIND_IP or LISTEN_PORT in your config.env, you must update the BASE_URL in web/config.js accordingly to ensure the interface can reach the API.
