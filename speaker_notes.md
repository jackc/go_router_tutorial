# Intro

This is a tutorial project for creating a router in Go. It will route based on query path and HTTP method.

# s1 - Failing test for not found

* httptest.ResponseRecorder

# s2 - Implement 404 for everything

# s3 - Failing test for basic string matching

# s4 - Implement basic string matching

# s5 - Refactor tests

* Extract stubHandler is straightforeward
* Extract testRequest is more nuanced. The inline version had the advantage of not having to pass in the *Router and *testing.T due to it being a closure. The extracted version has the advantage of the test body being cleaner and it being usable in multiple tests.

# s6 - Failing test for parameter matching

* Updated stub handler to print query parameters

# s7 - Implement parameter matching

Unfortunately, this is a big step. Fundamental change is breaking URLs into segments and the router being defined as a tree.

* Router is defined recursively. Each router is responsible for one segment.
* AddRoute now uses segmentizePath then addRouteFromSegments (which does most of the work) to construct routes
* segmentizePath splits strings by '/' and removes any empty segments causes by a '/' being the first or last character
* addRouteFromSegments recursively calls itself, creating new routers along the way as needed
* ServeHTTP now uses segmentizePath and findEndpoint to traverse the tree and find the correct endpoint
* findEndpoint traverses the routing tree recursively

At this point we are routing to the correct endpoint, but we are not giving it the parameters.

# s8 - Implement parameter capturing

Unfortunately, Go does not make it easy to attach context to an http.Request and pass it to the next handler.

Possible approaches:

* Global map indexed by request - must make currency safe (gorilla/mux)
* Break compatibility with http.Handler by adding arguments to handler signature
* Rewrite URL query string

Rewrite URL query string

* endpoint is new structure that stores handler and ordered names of parameters
* Use endpoint instead of http.Handler throughout router
* AddRoute now extracts parameter names
* findEndpoint now records the values of the parameter segments it matches
* addRouteArgumentsToRequest takes the parameter names from the endpoint and the parameter values from findEndpoint and rewrites the query string to include them

# s9 - Failing method matching test

# s10 - Implement method matching

* Add map for method to endpoint

# s11 - Add convenience methods

* go test -cover # 100% :)

# s12 - Demo app

* Fix - Use same package name as directory name so Go can import package

# s13 - Benchmarks

# Final thoughts

Very simple, very fast. We rebuilt https://github.com/JackC/quo_vadis which I am using on multiple projects. See https://github.com/JackC/tpr for example.
