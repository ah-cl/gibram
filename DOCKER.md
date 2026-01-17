# GibRAM - Graph in-Buffer Retrieval & Associative Memory

[![Docker Image Size](https://img.shields.io/docker/image-size/gibramio/gibram/latest)](https://hub.docker.com/r/gibramio/gibram)
[![Docker Pulls](https://img.shields.io/docker/pulls/gibramio/gibram)](https://hub.docker.com/r/gibramio/gibram)

High-performance in-memory knowledge graph server for GraphRAG applications. Combines vector search (HNSW) with graph storage for fast, context-aware retrieval.

## Quick Start

```bash
# Run server
docker run -d -p 6161:6161 --name gibram gibramio/gibram:latest

# With custom config
docker run -d -p 6161:6161 \
  -v $(pwd)/config.yaml:/etc/gibram/config.yaml:ro \
  gibramio/gibram:latest

# With data persistence
docker run -d -p 6161:6161 \
  -v gibram-data:/var/lib/gibram/data \
  gibramio/gibram:latest
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GIBRAM_LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `GIBRAM_VECTOR_DIM` | `1536` | Vector embedding dimension |
| `TZ` | `UTC` | Timezone |

### Ports

| Port | Protocol | Description |
|------|----------|-------------|
| `6161` | TCP | GibRAM Protocol (Protobuf) |

### Volumes

| Path | Description |
|------|-------------|
| `/etc/gibram/config.yaml` | Server configuration file |
| `/var/lib/gibram/data` | Data directory for persistence |
| `/var/lib/gibram/certs` | TLS certificates (optional) |
| `/var/log/gibram` | Log files (optional) |

## Docker Compose

```yaml
version: '3.8'

services:
  gibram:
    image: gibramio/gibram:latest
    container_name: gibram-server
    restart: unless-stopped
    
    ports:
      - "6161:6161"
    
    volumes:
      - gibram-data:/var/lib/gibram/data
    
    environment:
      - GIBRAM_VECTOR_DIM=1536
      - TZ=UTC
    
    healthcheck:
      test: ["CMD", "nc", "-z", "localhost", "6161"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  gibram-data:
```

## Health Check

The container includes a built-in health check that verifies port 6161 is accessible.

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' gibram
```

## Usage with Python SDK

```bash
# Install SDK
pip install gibram

# Connect to Docker container
python -c "
from gibram import GibRAMIndexer

indexer = GibRAMIndexer(
    session_id='test',
    host='localhost',
    port=6161
)
print('Connected!')
"
```

## Supported Platforms

- `linux/amd64`
- `linux/arm64`

## Security

The image runs as non-root user `gibram` (UID/GID 1000) for enhanced security.

For production deployments:
- Use custom TLS certificates
- Enable authentication via config file
- Run behind reverse proxy (nginx, traefik)

## Links

- **Documentation**: https://github.com/gibram-io/gibram
- **Python SDK**: https://pypi.org/project/gibram/
- **Issues**: https://github.com/gibram-io/gibram/issues

## License

MIT
