# gin-pagination
[![Run Tests](https://github.com/webstradev/gin-pagination/actions/workflows/test.yml/badge.svg)](https://github.com/webstradev/gin-pagination/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/webstradev/gin-pagination/branch/master/graph/badge.svg?token=C2D4QHYHI4)](https://codecov.io/gh/webstradev/gin-pagination)
[![Go Reference](https://pkg.go.dev/badge/github.com/webstradev/gin-pagination.svg)](https://pkg.go.dev/github.com/webstradev/gin-pagination/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/webstradev/gin-pagination)](https://goreportcard.com/report/github.com/webstradev/gin-pagination)
[![CodeQL](https://github.com/webstradev/gin-pagination/actions/workflows/codeql.yml/badge.svg)](https://github.com/webstradev/gin-pagination/actions/workflows/codeql.yml)

Simple pagination middleware for the gin framework. Allows for the usage of url parameters like `?page=1&size=25` to paginate data on your API.

## Installation
```bash
$ go get github.com/webstradev/gin-pagination/v2
```

## Default Usage
This package comes with a default pagination handler. This uses query parameters `page` and `size` with default values of `1` and `10` and a maximum page size of `100`.

#### Using the middleware on a router will apply the it to all requests on that router:
```go
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination/v2/pkg/pagination"
)

func main(){
  r := gin.Default()
  
  r.Use(pagination.New())
  
  r.GET("/hello", func(c *gin.Context){
    c.Status(http.StatusOK)  
  })
  
  r.Run(":3000")
}
```

#### Using the middleware on a single route will only apply it to that route:
```go
package main

import (
  "net/http"
  
  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination/v2/pkg/pagination"
)

func main(){
  r := gin.Default()
  
  r.GET("/hello", pagination.New(), func(c *gin.Context){
    page := c.GetInt("page")
  
    c.JSON(http.StatusOK, gin.H{"page" : page})  
  })
  
  r.Run(":3000")
}
```
The `page` and `size` are now available in the gin context of a request and can be used to paginate your data (for example in an SQL query).

 
## Custom Usage
To create a pagination middleware with custom parameters the New() function supports various custom options provided as functions that overwrite the default value.
All the options can be seen in the example below.
```go
package main

import (
  "net/http"
  
  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination/v2/pkg/pagination"
)

func main(){
  r := gin.Default()
  
  paginator := pagination.New(
    pagination.WithPageText("page"), 
    pagination.WithSizeText("rowsPerPage"),
    pagination.WithDefaultPage(1),
    pagination.WithDefaultPageSize(15),
    pagination.WithMinPageSize(5),
    pagination.WithMaxPageSize(15),
    pagination.WithHeaderPrefix(""),
  )
  
  r.GET("/hello", paginator, func(c *gin.Context){
    c.Status(http.StatusOK)  
  })
  
  r.Run(":3000")
}
```

The custom middleware can also be used on an entire router object similarly to the first example fo the Default Usage.
