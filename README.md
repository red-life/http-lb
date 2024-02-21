# HTTP Load Balancer
A lightweight, extensible Go-based load balancer with extensibility, robust configurability, seamless scalability, and integrated health checks, ensuring optimal performance and reliability for high-demand applications.

## Features
- **Configuration Flexibility**: Easily configurable through YAML files, enabling fine-tuning of load balancing behavior to suit specific requirements.
- **Efficient Containerization**: Fully containerized with Docker, ensuring seamless deployment and portability across diverse environments.
- **Simplified Execution**: Streamlined execution with a user-friendly Makefile, facilitating straightforward setup and operation.
- **Versatile Load Balancing**: Supports a wide array of static load balancing algorithms, providing flexibility to optimize traffic distribution.
- **Extensibility**: Built with extensibility in mind, allowing for easy extension and customization of each component to adapt to evolving needs.
- **Robust Health Monitoring**: Empowered with built-in health checks, ensuring continuous monitoring of backend server health for enhanced reliability.
- **Enhanced Monitoring Capabilities**: Equipped with comprehensive logging features, making it well-suited for integration with monitoring tools to track performance metrics and diagnose issues effectively.

## Usage
### Requirements
- Make
- Docker
### Steps:
1. Copy the `config.yaml` and configure it.
2. Run it using `make run` \
**If you are on a development server, Be cautious and run it using `make run_dev`**
3. That's it :)

## Configuration
Here is the sample configuration file in YAML (available at [config.yaml](./config.yaml)):
```yaml
algorithm: round-robin
log_level: debug

frontend:
    listen: 0.0.0.0:8000
    tls:
      cert: cert.ctr
      key: key.key

backend:
  - address: http://127.0.0.1:5001
    timeout: 1s # ms, s, min, h
    keep_alive:
      max_idle_connections: 100 # 0 means no limit
      idle_connection_timeout: 30s # maximum amount of time an idle connection will remain idle before closing itself. (0 means no limit)
    
  - address: http://127.0.0.1:5002
    timeout: 2s
    keep_alive:
      max_idle_connections: 50
      idle_connection_timeout: 10s # ms, s, min, h

health_check:
  endpoint: /health_check
  expected_status_code: 200 # HTTP Status codes like 200 or 304
  interval: 10s # ms, s, min, h
  timeout: 2s
```
- `algorithm`: Determines the backend selection method for client connections. Available options include:
    - Round-robin (`round-robin`)
    - Sticky round-robin (`sticky-round-robin`)
    - URL Hash (`url-hash`)
    - IP Hash (`ip-hash`)
    - Random (`random`) \
    **If you want to know how they work, read this article [HERE](https://blog.bytebytego.com/i/103707419/what-are-the-common-load-balancing-algorithms)**
- `log_level`: Sets the logging level for monitoring purposes. Options range from `info` to `debug`.
- `frontend`: It's responsible for accepting the incoming requests
    - `listen`: Specifies the IP and port for the HTTP server to listen on. Sync the port with the `PORT` variable in [Makefile](./Makefile)
    - **OPTIONAL** `tls`: If you want the front-end only accepts tls requests (**Recommended**) \
        - `cert`: The path to the certificate file.
        - `key`: The path to the private key file. \
- `backend`: The most important field. You define your backends in an array.
    - `address`: URL of the backend server.
    - `timeout`: Timeout for requests to the backend.
    - **OPTIONAL** `keep_alive`: If you don't want the reverse proxy creates connection for each request keep the keep_alive field.
        - `max_idle_connections`: Maximum number of idle connections (0 for no limit).
        - `idle_connection_timeout`: The maximum amount of time an idle connection will remain idle before closing itself. (0 means no limit) \
- `health_check`: Monitors backend health and adjusts routing accordingly.
    - `endpoint`: Path for the health check to verify backend availability.
    - `expected_status_code`: Expected HTTP status code indicating backend health.
    - `interval`: Interval for checking backend health.
    - `timeout`: Timeout for health check requests.

## Extensibility
**HTTP-LB** is designed with extensibility in mind, with every component built upon abstraction. This architecture enables seamless development and customization of each component for effortless integration and utilization.
- **Algorithms**: Based on the `LoadBalancingAlgorithm` interface defined in [definitions.go](./definitions.go), allowing for easy implementation of custom load balancing algorithms to suit specific requirements.
- **Request Forwarder**: This crucial component handles incoming requests from the frontend, invoking the selected algorithm and forwarding requests to the reverse proxy seamlessly.
- **Health Checker**:  Built upon the `HealthChecker` interface, ensuring continuous monitoring of backend server health for enhanced reliability.\
**These components serve as the backbone of the load balancer, with additional components following similar abstraction principles for further extensibility. Explore [definitions.go](./definitions.go) to discover more abstraction-based components ready for customization and expansion.**


## Contribution
Contributions to HTTP-LB are warmly welcomed. To contribute, please follow these steps:
1. Fork the repository.
2. Make your changes, ensuring clear and concise commit messages.
3. Push your branch to your fork.
4. Submit a pull request (PR) to the main repository, detailing the changes and any relevant information.
5. Your PR will be reviewed, and upon approval, it will be merged into the main branch.
