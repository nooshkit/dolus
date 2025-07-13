# Dolus: Lightweight Backend Mocking for Frontend Developers

**Dolus** is a lightweight command-line tool written in Go that helps frontend developers simulate backend APIs using a simple YAML configuration file. No server setup required — just define your endpoints and responses in YAML and start mocking instantly.

---

## Features

- Define API routes and responses using a single YAML file
- Quickly mock endpoints with custom status codes, and JSON payloads
- Supports `GET`, `POST`, `PUT`, `DELETE`, and more
- Toggle verbose request logging to inspect incoming HTTP requests
- Perfect for frontend prototyping and integration testing

---

## Installation

- NOTE: Its recommended to pull the source and build the binary locally to ensure things work on your system

Download the latest release from [Releases](https://https://github.com/nooshkit/dolus/releases) or install via `go install`:

```bash
go install https://github.com/nooshkit/dolus/releases@Latest
```

## Usage
```bash
dolus -config ./config.yaml
```
### Flags
- `-config <path>`: By default, dolus will attempt to load `config.yaml` from the local directory
- `-v`: Enable verbose mode to log incoming HTTP requests

## Verbose Logging
```bash
dolus -config ./mock-api.yaml -v
```
- You'll see output like:
```
Method: GET
Path:   /api/v1/locations
Query Params:
Headers:
  Connection: [keep-alive]
  Accept: [*/*]
  Accept-Encoding: [gzip, deflate, br]
```

## Example Configuration
```yaml
server:
  port: 8080
  listenaddr: localhost
paths: # Paths can be defined as any arbitrary path desired. Must begin with '/'
  /item:
    get:
      response: 200
      content:
        itemid: 42
        itemName: "The Answer"
        itemDetails:
          price: 420
          count: 100
    post:
      response: 201
      content:
        status: "Item created"
  /api/v1/locations:
    get:
      response: 200
      content:
        franchises:
          - state: GA
            cities:
              - Atlanta
              - Savannah
              - Macon
          - state: TX
            cities:
              - Austin
              - Dallas
              - Houston
          - state: CA
            cities:
              - Los Angeles
              - San Franciso
              - San Diego
```

## Contributing

Pull requests are welcome! This was quickly thrown together as a P.O.C -- if there's desire for new features it may take some time to get around to it.
