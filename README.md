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




