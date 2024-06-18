
# Request Counter

## Description
This project implements a simple HTTP server that tracks and returns the number of requests received within a configurable time window.


## Dependencies

* Go ([https://go.dev/](https://go.dev/))

## Building

```bash
make build
```


## Running

```bash
make run
```


## Testing

### Unit Tests

```bash
make test
```

### Concurrency Tests

```bash
make test-race
```


## Linting

```bash
make lint
```


## Configuration

The Request Counter server can be configured using environment variables:

### TTL_SEC (Default: 60) 
Defines the window period (in seconds) for counting requests.
### PORT (Default: 8080) 
Specifies the port on which the server listens for incoming requests.

## Example Usage

1. Set environment variables (optional):

   ```bash
   export TTL_SEC=30  # Set window period to 30 seconds
   export PORT=8080  # Set port to 8080
   ```

2. Build and run the server:

   ```bash
   make build-run
   ```

   This builds the executable and starts the server listening on port 8080 with a window period of 30 seconds (or 60 seconds if not set).

3. Send requests to the server:

   Use tools like `curl` or a web browser to send HTTP requests to the server address (e.g., `http://localhost:8080`).

4. Access request count:

   Send a GET request to the root path (`/`) of the server. The response will be the current request count within the defined window period.

