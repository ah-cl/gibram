# GibRAM

**Graph in-Buffer Retrieval & Associative Memory**

High-performance knowledge graph server for GraphRAG applications. Store entities, relationships, communities with semantic search using HNSW vector indexing.

## Features

- **Vector Search**: HNSW-based semantic search (sub-millisecond queries)
- **Graph Storage**: Entities, relationships, and communities with TTL-based eviction
- **Protocol**: Binary protobuf protocol for low-latency communication
- **Language Support**: Python SDK
- **Production Ready**: Tested with real-world workloads, memory efficient

## Quick Start

### Install via Binary

```bash
# Install via script
curl -fsSL https://gibram.io/install.sh | sh

# Run server
gibram-server --insecure
```

Server runs on port **6161** by default.

### Install via Docker

```bash
# Run server
docker run -p 6161:6161 gibramio/gibram:latest

# With custom config
docker-compose up -d
```

### Python SDK

```bash
pip install gibram
```

**Basic Usage:**

```python
from gibram import GibRAMIndexer

# Initialize indexer
indexer = GibRAMIndexer(
    server_host="localhost",
    server_port=6161,
    openai_api_key="your-api-key"
)

# Index documents
stats = indexer.index_documents([
    "Python is a programming language created by Guido van Rossum.",
    "JavaScript was created by Brendan Eich at Netscape in 1995."
])

print(f"Entities: {stats.total_entities}")
print(f"Relationships: {stats.total_relationships}")

# Query
results = indexer.query("Who created JavaScript?", top_k=3)
for result in results:
    print(f"{result.entity.name}: {result.score}")
```

**Custom Components:**

```python
from gibram import GibRAMIndexer
from gibram.chunkers import TokenChunker
from gibram.extractors import OpenAIExtractor
from gibram.embedders import OpenAIEmbedder

indexer = GibRAMIndexer(
    chunker=TokenChunker(chunk_size=512, overlap=50),
    extractor=OpenAIExtractor(model="gpt-4o", api_key="..."),
    embedder=OpenAIEmbedder(model="text-embedding-3-small", api_key="...")
)
```

## Performance

- Query latency: **P50 = 2ms**, P95 = 5ms
- Write throughput: **100K ops/s**
- Memory per entity: **~3.2KB**
- Supports: **1-5M entities** per session

## Documentation

- [Server Documentation](./server/README.md)
- [Python SDK Documentation](./sdk/python/README.md)
- [Design Decisions](./docs/design-decisions/)
- [Performance Guide](./docs/performance/)

## License

MIT
