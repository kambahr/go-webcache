# WebCache v2.0.0

A high-performance, thread-safe LRU (Least Recently Used) cache for Go applications. 

## Overview
WebCache is designed to handle web assets and byte arrays efficiently. Version 2.0 represents a complete architectural rewrite, moving from a slice-based approach to a **Map + Doubly Linked List** structure.

### Key Features
* **$O(1)$ Time Complexity:** Instant lookups regardless of cache size.
* **LRU Eviction:** Automatically drops the least accessed items when the limit is reached.
* **Thread-Safe:** Safe for concurrent use in high-traffic web servers via `sync.RWMutex`.
* **Automatic Expiration:** A background janitor periodically cleans up stale entries.

## Installation
```bash
go get [github.com/yourusername/webcache](https://github.com/yourusername/webcache)
