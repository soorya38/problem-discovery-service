# Problem Discovery Service
Problem Discovery Service provides a structured API to fetch, filter, and sort programming problems. Results are cached in memory to reduce external API calls and improve response times.
- Fetch problems by single or multiple tags
- Exact tag filtering (single-tag-only and multi-tag-only)
- In-memory caching for repeated queries
- Problems sorted by rating

---

## Running the Binary
The service runs as a standalone binary.

### Linux
```bash
chmod +x problemsvc
./problemsvc
````

### macOS

```bash
chmod +x problemsvc
./problemsvc
```

If blocked by Gatekeeper:

```bash
xattr -d com.apple.quarantine problemsvc
```

### Windows

```powershell
problemsvc.exe
```

Or double-click the executable.

---

## Default Configuration

* **Port:** `49160`
* **Protocol:** HTTP
* **Startup Mode:** Foreground process

---

## API Endpoints

### Health & Metadata

* `GET /health`
* `GET /tags`

### Problem Discovery

* `GET /problems?tag={topic}`
* `GET /problems/multi?tags={topic},{topic}`
* `GET /problems/only?tag={topic}`
* `GET /problems/multi/only?tags={topic},{topic}`
* `GET /tags`

---
