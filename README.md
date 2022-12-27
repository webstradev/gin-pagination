# gin-pagination
Simple pagination middleware for the gin framework

## Installation
``` bash
$ go get github.com/webstradev/gin-pagination
```

## Default Usage
This package comes with a default pagination handler. This uses query parameters `page` and `size` with default values of `1` and `10` and a maximum page size of `100`.

#### Using the middleware on a router will apply the it to all request on that router:
``` go
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination"
)

func main(){
  r := gin.Default()
  
  r.Use(pagination.Default())
  
  r.GET("/hello", func(c *gin.Context){
    c.Status(http.StatusOK)  
  })
  
  r.Run(":3000")
}
```

#### Using the middleware on a single route will only apply it to that route:
``` go
package main

import (
  "net/http"
  
  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination"
)

func main(){
  r := gin.Default()
  
  r.GET("/hello", pagination.Default(), func(c *gin.Context){
    page := c.GetInt("page")
  
    c.JSON(http.StatusOK, gin.H{"page":page})  
  })
  
  r.Run(":3000")
}
```
The `page` and `size` are now available in the gin context of a request and can be used to paginate your date (for example in an SQL query.

 
## Custom Usage
To create a pagination middleware with custom parameters use the `pagination.New()` function.
``` go
package main

import (
  "net/http"
  
  "github.com/gin-gonic/gin"
  "github.com/webstradev/gin-pagination"
)

func main(){
  r := gin.Default()
  
  paginator := pagination.New("page", "rowsPerPage", "1", "15", 5, 150)
  
  r.GET("/hello", paginator, func(c *gin.Context){
    c.Status(http.StatusOK)  
  })
  
  r.Run(":3000")
}
```
The custom middleware can also be used on an entire router object similar to the first example fo the Default Usage.
