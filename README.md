
![Configo Logo](assets/logo-name.png)

### Configo - Configuration Management System
A flexible, dynamic configuration management system with a Go API backend and Flutter web interface that allows you to easily register, modify and retrieve configuration values for your services.

## Overview
Configo provides a centralized configuration management solution with the following features:

* Register configuration objects with validation and type information
* Web UI for managing configurations
* Swagger API documentation
* Automatic configuration refresh
* Support for multiple configuration value types
* Simple repository interface for storage backends

## How to use

### Integration Guide
Integrating Configo into your Go application is straightforward:

#### 1. Create a Repository Implementation
First, implement the Repository interface for your preferred storage backend:

```
type Repo struct {
    cache cache.Cache  // Or your preferred storage method
}

func NewRepo() *Repo {
    return &Repo{
        cache: memory.NewCache(), // Or Redis, etc.
    }
}

func (r *Repo) SaveConfig(ctx context.Context, cfg *configo.Config) error {
    return r.cache.Set(ctx, cfg.Key, cfg)
}

func (r *Repo) GetConfig(ctx context.Context, key string) (*configo.Config, error) {
    var cfg configo.Config
    err := r.cache.GetV(ctx, key, &cfg)
    if err != nil {
        return nil, err
    }
    if cfg.Key == "" {
        return nil, nil
    }
    return &cfg, nil
}
```

#### 2. Initialize the Config Manager

```
configo, err := configo.NewConfigManager(ctx, &configo.ConfigManagerConfig{
    ServiceName:        "Your Service",
    ServiceDescription: "Description of your service",
    ConfigRefreshInSecs: 30, // Optional, defaults to 10
}, NewRepo())
if err != nil {
    panic(err)
}
```

#### 3. Define Configuration Structures
Create structs for your configurations with appropriate field tags:

```
type MyConfig struct {
    DatabaseURL string `name:"database_url" type:"string" description:"URL for the database" required:"true"`
    MaxConn     int    `name:"max_connections" type:"number" description:"Maximum connections" required:"true"`
    EnableCache bool   `name:"enable_cache" type:"boolean" description:"Enable caching" required:"false"`
}

```
#### 4. Register your configurations

```
myConfig := &MyConfig{
    DatabaseURL: "postgres://default:password@localhost:5432/mydb",
    MaxConn: 10,
}
err = configo.RegisterConfig(ctx, myConfig)
if err != nil {
    log.Fatalf("Failed to register config: %v", err)
}
```

#### 5. Register Routes with Your Web Framework
For Gin:

```
router := gin.New()
group := router.Group("/your-service")

// Register configo routes
err = configo.RegisterRoute(ctx, func(method, path string, handler http.HandlerFunc) error {
    ginHandler := func(c *gin.Context) {
        handler(c.Writer, c.Request)
    }
    group.Handle(method, path, ginHandler)
    return nil
})
if err != nil {
    panic(err)
}

router.Run(":8085")
```

Note: You can register with any of the router you are using.

### Configuration Types
Configo supports the following configuration types via field tags:

| Type | Description | UI Element |
|------|-------------|------------|
| string | Basic string values | Text input |
| number | Integer values | Number input |
| boolean | Boolean values (true/false) | Checkbox |
| json | JSON objects | JSON editor textarea |
| bigText | Multiline text | Large textarea |
| choice | Selection from predefined options | Dropdown |
| parent | Parent configuration | Container |
| list | List of values | List widget |


### Web UI
Access the web UI at:
```
http://localhost:8085/your-service/configo/web/
```
### API Documentation
Swagger documentation is available at:

```
http://localhost:8085/your-service/configo/swagger/index.html
```

#### API Endpoints
* `GET /configo/metadata` - Get all configuration keys
* `GET /configo/config?key=YOUR_KEY` - Get specific configuration
* `POST /configo/config` - Update configuration
* `GET /configo/swagger/*any` - Swagger documentation
* `GET /configo/web/*any` - Web UI


## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.





## Controbute

### Installation
#### Prerequisites
* Go 1.16 or later
* Git
* Flutter
#### From Source
```
# Clone the repository
git clone https://github.com/gofreego/configo.git
cd configo

# Install dependencies
go mod tidy

# Run the example
go run main.go
```

License
This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

Contact
For any inquiries, please contact:

* Developer: https://github.com/pavanyewale
* Linkedin: https://www.linkedin.com/in/pavanyewale/
* Email: pavanyewale1996@gmail.com
* GitHub: https://github.com/gofreego
