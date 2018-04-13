package gin_reverseproxy

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func ReverseProxy(domains map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("host: " + c.Request.Host)
		if c.Request.Host != "localhost:4000" {
			log.Println("continue for next")
			c.Next()
		}

		// we need to buffer the body if we want to read it here and send it
		// in the request.
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// you can reassign the body if you need to parse it as multipart
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		url := fmt.Sprintf("%s://%s%s", "http", domains[c.Request.Host], c.Request.RequestURI)

		proxyReq, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))

		// We may want to filter some headers, otherwise we could just use a shallow copy
		// proxyReq.Header = c.Request.Header
		proxyReq.Header = make(http.Header)
		for h, val := range c.Request.Header {
			proxyReq.Header[h] = val
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		c.Writer.Write(bodyContent)
		for h, val := range resp.Header {
			c.Writer.Header()[h] = val
		}
		c.Abort()
	}
}
