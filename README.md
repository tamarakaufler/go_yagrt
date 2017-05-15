# go_router
Recreate a wheel to know how it works.

## Request journey

- request comes in and is intercepted by a router/multiplexor

- router hands it over, based on the request path, to the relevant handler

- handler is a function that takes (at least) response and request objects as arguments:
  * processes the request
  * writes out response headers and the body

- it is possible to use a default router provided by the net/http package or create one's own

- the default router is of type Handler, meaning it satisfies its interface, ie implements the ServeHTTP method:
  * ServeHTTP(http.ResponseWriter, *http.Request)

- any object can be a handler. If it satisfies the Handler interface then it can be used to handle requests         



## Implementation

- router:
  * requirements:
    * must be able to process incoming request routes:
      * it must be possible to register a request path with a provided handler
      * it must be possible to use middleware
    * provide some basic handlers

  * implementation:
    * use router.Params to store parameters defined during the path registration
      * use two different implementations to identify and store parameters defined during the path registration:
        * use a Trie type tree
        * use regexes
