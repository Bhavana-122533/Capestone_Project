# Production-Ready Order Management System (Week 5)

Go + Gin REST API, GraphQL read endpoint, PostgreSQL/RDS support, Docker, tests, observability and CI. Built against the uploaded Week-5 case study.

## 1. Prerequisites
- Go 1.23+
- Docker Desktop with Docker Compose
- Optional: PostgreSQL 15+ if not using Compose

## 2. Run with Docker (recommended)
1. Copy `.env.example` to `.env`.
2. Set `DB_PASSWORD` in `.env` to a strong local development password. For Compose, also set `DB_USER=orderuser`, `DB_NAME=orderdb`, `DB_HOST=postgres`, `DB_SSLMODE=disable`. Compose injects the database user/password for the app as `orderuser` and `${DB_PASSWORD}`.
3. Start services:
   ```bash
   docker compose up --build
   ```
4. API is at `http://localhost:8080`; PostgreSQL is on port 5432.
5. Stop with `docker compose down`; add `-v` only if you intentionally want to delete local database data.

If `API_KEY` is set, send header `X-API-Key: <your value>` to `/api/v1/*`. Do not commit `.env`.

## 3. Run directly with Go
- `cp .env.example .env` (Windows PowerShell: `Copy-Item .env.example .env`)
- Leave `DB_HOST` blank for an in-memory demo, or configure PostgreSQL connection values.
- `go mod tidy`
- `go run .`
The in-memory mode is for local learning only; data disappears when the process stops. Use PostgreSQL/Docker for persistence.

## 4. REST endpoints
When `API_KEY` is configured, add `X-API-Key` to all `/api/v1` calls.

### Create order
`POST /api/v1/orders`
```json
{"customerId":101,"productId":501,"quantity":2}
```
The Product Service must return a positive price as JSON `{"price": 1200}` or `{"unitPrice": 1200}`. Set `PRODUCT_SERVICE_URL` to its base URL (the client requests `/products/{productId}`). For local API testing without a product service, create a mock service or configure one; the app intentionally fails with 502 instead of inventing a price.

### Get order
`GET /api/v1/orders/1`

### List orders
`GET /api/v1/orders?limit=20&offset=0`

### Update status
`PUT /api/v1/orders/1/status`
```json
{"status":"CONFIRMED"}
```
Valid statuses: `CREATED`, `CONFIRMED`, `CANCELLED`.

## 5. Health, metrics and GraphQL
- `GET /health/live` — process is alive.
- `GET /health/ready` — repository/database is reachable.
- `GET /metrics` — local request count, latency summary and error count in Prometheus-style text format for easy inspection without a Prometheus/Grafana stack.
- `POST /graphql` — GraphQL read query:
```json
{"query":"{ orders { id customerId productId quantity totalAmount status } }"}
```
or
```json
{"query":"{ order(id: \"1\") { id status totalAmount } }"}
```

## 6. Tests and checks
```bash
go mod tidy
gofmt -w .
go build ./...
go test ./...
go test -race ./...
docker build -t order-management:local .
```
The tests use a thread-safe in-memory repository and do not need AWS or PostgreSQL.

## 7. Configuration and security notes
- All database credentials are environment variables; never put secrets in source control.
- `API_KEY` is a minimal demonstration of authentication. Production should use a managed identity provider (e.g. OIDC/JWT), validate issuer/audience/signature, and enforce role/permission policies.
- `DB_SSLMODE=disable` is for local Docker only. For RDS, use TLS (`verify-full` with the appropriate CA bundle) and the RDS endpoint.
- Terminate HTTPS at an Application Load Balancer with an ACM certificate when available.
- Restrict the RDS security group inbound rule to the application security group on TCP 5432. Do not expose RDS publicly.
- Use least-privilege IAM roles. If Secrets Manager is unavailable, inject secrets through the deployment environment temporarily; do not bake them into the image or Git.
- CloudWatch integration can be added through the deployment log driver/agent when available. The application writes structured key-value logs to stdout for container log collection.
- Distributed tracing is represented by a request ID propagated in `X-Request-ID` and logged across the request path. For full distributed tracing, add OpenTelemetry instrumentation/exporter when your telemetry backend is available.

## 8. Product Service resilience
The external product client has per-request timeout, bounded retries (0–5, default 2), exponential backoff, context cancellation, and a circuit breaker. When the service is unavailable or breaker is open, create order returns a controlled 502 response. The product service contract is `GET /products/{id}` and a JSON positive `price` or `unitPrice`.

## 9. AWS deployment guide (only when services are available)
1. Create/use a VPC with public subnets for the ALB and private subnets for app tasks/instances and RDS.
2. Create an application security group allowing inbound only from the ALB security group; permit outbound to RDS and required endpoints.
3. Create an RDS PostgreSQL database in private subnets; apply `database/schema.sql` using a controlled migration step.
4. Create a least-privilege IAM role for EC2/ECS task execution. Grant only required access to pull the image/logs and read the specific secret when Secrets Manager is available.
5. Build and push the Docker image to ECR when available. If ECR is unavailable, retain the image locally and do not claim deployment completed.
6. Deploy two app replicas behind an ALB when ECS/EC2 and the load balancer are available. Configure health check path `/health/ready`.
7. Configure environment variables from Secrets Manager or the platform's secret injection mechanism. Use RDS TLS settings.
8. Configure CloudWatch logs/alarms when available. Confirm `/health/live`, `/health/ready`, API calls, and logs before calling the deployment complete.

This repository contains development artifacts and deployment instructions. It does not claim that any AWS resources have been created.

## 10. Known scope boundaries
- Four required order functions are implemented.
- GraphQL is read-only for orders; the required business mutations are REST.
- Database error metric is registered as a metric placeholder; database failures should be incremented at the repository/error instrumentation boundary before production rollout.
- The product service is external and must be supplied/configured. A missing URL is reported as a controlled dependency failure, not silently bypassed.
