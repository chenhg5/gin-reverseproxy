# Gin Reverse Proxy Middleware

proxy different domain request to different host through specified proxy rules <br>
no need nginx just go

## usage

```

import proxy "github.com/chenhg5/gin-reverseproxy"

router.Use(proxy.ReverseProxy(map[string] string {
    "www.xyz.com" : "localhost:4001",
    "www.abc.com" : "localhost:4003",
}))


```


## todo

- [ ] Load balance
- [ ] Error handle
- [ ] specified route