# Golem

A robust Layer 7 (HTTP) Load Balancer written in Go.

## What is Golem?

**Golem** is a lightweight, custom-built Layer 7 (Application Layer) load balancer. It was designed to explore concepts and concurrency + traffic management in Go.

Golem acts as a reverse proxy, distributing incoming HTTP traffic across a pool of user configured backend servers. It ensures high availability by actively monitoring server health and automatically routing traffic away from failing nodes.

## Architecture

Golem sits between the client and backend services. It accepts requests, selects an optimal backend based on the configured strategy, and forwards the response back to the client.

## Key Features

- **Layer 7 Load Balancing:** Distribution of HTTP requests.
    
- **Multiple Strategies:**
    
    - **Round Robin:** Rotates requests sequentially across all healthy servers.
        
    - **Weighted Round Robin:** Prioritizes servers based on their weight.
        
    - **Least Connections:** Prioritizes servers with fewer active connections.
    
    - **IP Hash:** Distributes requests based on the client's IP address. Used HRW (Hash Ring Weighted) algorithm for consistent hashing.
        
- **Active Health Monitoring:** Periodically pings backends to verify status. Unhealthy servers are removed from rotation instantly based on the configured threshold.
    
- **Concurrency:** Built on Go's goroutines to handle high throughput with low overhead.
    
- **Configurable:** Easy to read and setup YAML configuration for defining host, port, intervals, timeouts and health paths.
- **Connection Pooling:** Efficiently manages connections to backend servers.
- **Timeouts:** Sets timeouts for connection establishment, request processing, and response transmission.
    

## To Run 

### Prerequisites

- Go 1.21 or higher
    

### Installation

1. **Clone the repository**
    
    ```
    git clone [https://github.com/Sahil-796/golem.git](https://github.com/Sahil-796/golem.git)
    cd golem
    ```
2. **Configure**: Edit the `config.yaml` file in the root directory.
3. **Run main.go**
    
    ```
    go run main.go
    ```
    
    

## Configuration

Open the `config.yaml` file in the root directory.

## Usage

1. **Start your backend servers** (for testing, you can use any simple servers):
    
2. **Run Golem**:
    
3. **Send a request**:
    
    You should see the response from one of your backends according to the configured load balancing algorithm. Repeated requests will cycle through them.
    

## What is a Load Balancer?

A load balancer is a device or service that acts as a traffic cop for your servers. It distributes network traffic across multiple servers to:

- **Optimize resource use:** Prevent one server from doing all the work while others sit idle.
    
- **Maximize throughput:** Process more requests per second.
    
- **Ensure reliability:** If one server crashes, the load balancer stops sending it traffic.

## Wisdom
The name golem is inspired by a clash royale card named golem which symbolizes strength and resilience. It's a fitting name for a load balancer that can handle heavy loads and keep your servers running smoothly.
Also the go in golem stands for golang. 
