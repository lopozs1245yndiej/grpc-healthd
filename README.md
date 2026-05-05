# grpc-healthd

Lightweight sidecar daemon that exposes gRPC health check endpoints for containerized services.

---

## Installation

```bash
go install github.com/yourorg/grpc-healthd@latest
```

Or pull the Docker image:

```bash
docker pull ghcr.io/yourorg/grpc-healthd:latest
```

## Usage

Run `grpc-healthd` alongside your service, pointing it at the target gRPC server:

```bash
grpc-healthd --port 50051 --target localhost:9090 --service my.Service
```

### Docker Compose Example

```yaml
services:
  app:
    image: my-app:latest

  healthd:
    image: ghcr.io/yourorg/grpc-healthd:latest
    command: ["--target", "app:9090", "--service", "my.Service"]
    ports:
      - "50051:50051"
```

### Flags

| Flag        | Default       | Description                        |
|-------------|---------------|------------------------------------|
| `--port`    | `50051`       | Port to expose the health endpoint |
| `--target`  | `localhost:80`| Target gRPC service address        |
| `--service` | `""`          | Service name to check              |
| `--interval`| `10s`         | Health check poll interval         |

## Configuration

`grpc-healthd` also supports configuration via environment variables. Prefix any flag with `HEALTHD_` (e.g., `HEALTHD_TARGET=localhost:9090`).

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.

## License

[MIT](LICENSE)