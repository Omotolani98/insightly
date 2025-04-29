# Insightly: Time-Series Log Summarizer

Insightly is an intelligent log summarizer that ingests logs via gRPC, batches them in Redis Streams, summarizes each time window using a local LLM (Docker Model Runner), and stores human-readable digests in PostgreSQL. You can query summaries via gRPC or through a simple HTML interface.

## 📁 Folder Structure

```
├── Makefile                     # Commands to build/run components
├── README.md                    # Documentation
├── bin                          # Compiled binaries
│   ├── api                      # API binary
│   ├── ingest                   # Ingest binary
│   ├── query                    # Query binary
│   └── summarizer               # Summarizer binary
├── cmd                          # Application entrypoints
│   ├── api                      # Fiber HTTP API server
│   ├── ingest                   # gRPC ingest server
│   ├── query                    # gRPC query server
│   └── summarizer               # Log summarization worker
├── deployments                  # Dockerfiles
├── docker-compose.yml           # Compose file for Redis & PostgreSQL
├── internal                     # Application logic
│   ├── cache                    # Redis caching logic
│   ├── config                   # Application configuration
│   ├── db                       # PostgreSQL database connections
│   ├── ingest                   # Ingest handlers
│   ├── llm                      # Local LLM integration
│   ├── query                    # Query handlers
│   ├── storage                  # Database models and migrations
│   └── summarizer               # Summarization worker logic
├── proto                        # Protobuf definitions
│   ├── ingest                   # Compiled ingest protobuf
│   ├── ingest.proto             # Ingest protobuf definition
│   ├── query                    # Compiled query protobuf
│   └── query.proto              # Query protobuf definition
├── scripts                      # Migration scripts
├── views                        # HTML views (not currently used)
└── .env                         # Environment variables
```

## ⚙️ Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Redis & PostgreSQL (auto-deployed via Docker Compose)
- Docker Model Runner enabled on Docker Desktop (`ai/llama3.2` model)
- `protoc` and Go plugins (`protoc-gen-go`, `protoc-gen-go-grpc`)

## 🛠 Components

| Component             | Role                                                   | Startup           |
|-----------------------|--------------------------------------------------------|-------------------|
| **Ingest API**        | gRPC server: streams logs to Redis                     | `./start.sh`      |
| **Redis Stream**      | Buffers logs temporarily                               | Docker Compose    |
| **Summarizer Worker** | Processes logs with LLM, saves summaries to PostgreSQL | `./start.sh`      |
| **Docker‑LLM**        | Local LLM container (Docker Model Runner)              | Docker Desktop    |
| **PostgreSQL**        | Stores human-readable summaries                        | Docker Compose    |
| **Query API**         | gRPC server: fetches summaries                         | `./start.sh`      |
| **Fiber API**         | HTTP server for frontend summary fetching              | `./start.sh`      |

## 🚀 Running Locally

1. Clone the repository and navigate into the folder:

```bash
git clone [REPO_URL]
cd insightly
```

2. Install dependencies:

```bash
go mod tidy
```

3. Generate Protobuf files:

```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

4. Enable Docker Model Runner:

- Open Docker Desktop → Settings → Features In Development
- Check "Enable Docker Model Runner" & "Enable host-side TCP support"
- Confirm LLM runs at `localhost:12434`

5. Create a `.env` file:

```env
# PostgreSQL
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=insightly
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# LLM
LLM_HOST=localhost
LLM_PORT=12434
ENGINE_ID=llama.cpp
MODEL_NAME=ai/llama3.2

# gRPC Ports
INGEST_PORT=50051
QUERY_PORT=50052
```

6. Start the Insightly stack:

```bash
./start.sh
```

- To stop the application, use:

```bash
make down
```

## 📡 Testing with Postman

- Open Postman
- Create a gRPC request
  - Upload Proto file (`proto/ingest.proto`)
  - Select the `StreamLogs` method
- Send the message

Check query logs using another gRPC request (`proto/query.proto`) with the `GetSummaries` method.

## 🌐 Viewing Summaries

Create a simple frontend using HTML, Tailwind CSS, and JavaScript:

### `index.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Insightly</title>
  <script defer src="app.js"></script>
  <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter&display=swap" rel="stylesheet">
  <style>
    body { font-family: 'Inter', sans-serif; }
  </style>
</head>
<body class="bg-gray-50 p-8">
  <div class="max-w-2xl mx-auto bg-white rounded-lg shadow p-6">
    <h1 class="text-2xl font-bold mb-4">Insightly</h1>
    <input id="filterInput" type="text" placeholder="Filter by stream..." class="border p-2 rounded w-full">
    <ul id="summaryList" class="space-y-4 mt-4"></ul>
  </div>
</body>
</html>
```

### `app.js`

[JS content remains unchanged from previous version]

