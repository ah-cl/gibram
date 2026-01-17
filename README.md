# GibRAM

**Graph in-Buffer Retrieval & Associative Memory**

- **Graph in-Buffer**: Graph structure (entities + relationships) stored in RAM
- **Retrieval**: Query mechanism for retrieving relevant context in RAG workflows  
- **Associative Memory**: Traverse between associated nodes via relationships, all accessed from memory

GibRAM is an in-memory knowledge graph server designed for retrieval augmented generation (RAG) workflows. It combines a lightweight graph store with vector search so that related pieces of information remain connected in memory. This makes it easier to retrieve related regulations, articles or other text when a query mentions specific subjects.

## Why GibRAM?
- In memory and Ephemeral: Data lives in RAM with a configurable time to live. It is meant for short lived analysis and exploration rather than persistent storage.
- Graph and Vectors Together: Stores named entities, relationships and document chunks alongside their embeddings in the same structure.
- Graph aware Retrieval: Supports traversal over entities and relations as well as semantic search, helping you pull in context that would be missed by vector similarity alone.
- Python SDK: Provides a GraphRAG style workflow for indexing documents and running queries with minimal code. Components such as chunker, extractor and embedder can be swapped out.

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
    session_id="my-project",
    host="localhost",
    port=6161,
    llm_api_key="sk-..."  # or set OPENAI_API_KEY env
)

# Index documents
stats = indexer.index_documents([
    "Python is a programming language created by Guido van Rossum.",
    "JavaScript was created by Brendan Eich at Netscape in 1995."
])

print(f"Entities: {stats.entities_extracted}")
print(f"Relationships: {stats.relationships_extracted}")

# Query
results = indexer.query("Who created JavaScript?", top_k=3)
for entity in results.entities:
    print(f"{entity.title}: {entity.score}")
```

**Custom Components:**

```python
from gibram import GibRAMIndexer
from gibram.chunkers import TokenChunker
from gibram.extractors import OpenAIExtractor
from gibram.embedders import OpenAIEmbedder

indexer = GibRAMIndexer(
    session_id="custom-project",
    chunker=TokenChunker(chunk_size=512, chunk_overlap=50),
    extractor=OpenAIExtractor(model="gpt-4o", api_key="..."),
    embedder=OpenAIEmbedder(model="text-embedding-3-small", api_key="...")
)
```

## License

MIT
