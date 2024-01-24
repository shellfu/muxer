# Tutorial: Building a RESTful API with the `muxer` Package

In this tutorial, we will explore how to build a RESTful API using the
`muxer` package. The `muxer` package provides an HTTP request multiplexer,
allowing you to define routes, handle different HTTP methods, extract path
parameters, and use middleware functions. We will go through the process
of setting up a basic API, defining routes, handling requests, and using
middleware functions for authentication.

## Prerequisites

Before we begin, make sure you have the following prerequisites:

- Basic understanding of Go programming language.
- Go installed on your system. You can download it from the official Go website: https://golang.org/dl/

## Step 1: Create a new Go module

Let's start by creating a new Go module for our project. Open your
terminal or command prompt and run the following command:

```sh
go mod init myapi
```

This will initialize a new Go module named `myapi` in the current directory.

## Step 2: Install the `muxer` package

Next, we need to install the `muxer` package. Run the following command in your terminal:

```sh
go get github.com/shellfu/muxer
```

This will download and install the `muxer` package into your Go module.

## Step 3: Create the main API file

In the root directory of your project, create a new file named `main.go`.
This file will contain the main entry point of our API.

Open the `main.go` file in a text editor and add the following code:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/shellfu/muxer"
)

func main() {
    router := muxer.NewRouter()

    // Define your routes here

    // Start the HTTP server
    fmt.Println("Server listening on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Step 4: Define routes

Now, let's define some routes for our API. We will create routes for handling GET, POST, PUT, and DELETE requests.

Add the following code inside the `main()` function in `main.go`:

```go
func main() {
    router := muxer.NewRouter()

    // Define routes
    router.HandleFunc("http.MethodGet, "/users", getUsersHandler)
    router.HandleFunc("http.MethodPost, "/users", createUserHandler)
    router.HandleFunc("http.MethodPut, "/users/:id", updateUserHandler)
    router.HandleFunc("http.MethodDelete, "/users/:id", deleteUserHandler)

    // Start the HTTP server
    fmt.Println("Server listening on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func getUsersHandler(w http.ResponseWriter, req *http.Request) {
    // Handle GET /users request
    fmt.Println("Handling GET /users request")
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
    // Handle POST /users request
    fmt.Println("Handling POST /users request")
}

func updateUserHandler(w http.ResponseWriter, req *http.Request) {
    // Handle PUT /users/:id request
    fmt.Println("Handling PUT /users request")
}

func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
    // Handle DELETE /users/:id request
    fmt.Println("Handling DELETE /users request")
}
```

Here, we have defined four routes: `GET /users`, `POST /users`, `PUT
/users/:id`, and `DELETE /users/:id`. Replace the empty handler functions
with your own implementation for handling the corresponding requests.

## Step 5: Handling requests

Let's now implement the handler functions for our routes. These functions
will be executed when a matching request is received by the server.

For example, in the `getUsersHandler` function, you can fetch a list of
users from a database and return it as a JSON response. Here's an example
implementation:

```go
import (
    "encoding/json"
    "net/http"
)

func getUsersHandler(w http.ResponseWriter, req *http.Request) {
    // Fetch users from the database (replace with your own implementation)
    users := []User{
        {ID: 1, Name: "John"},
        {ID: 2, Name: "Jane"},
    }

    // Convert users to JSON
    jsonData, err := json.Marshal(users)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Set response headers
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    // Write JSON response
    w.Write(jsonData)
}
```

Implement similar handler functions for the remaining routes
(`createUserHandler`, `updateUserHandler`, and `deleteUserHandler`) based
on your API requirements.

## Step 6: Starting the API server

Now that we have defined our routes and implemented the handler functions,
it's time to start the API server and test our endpoints.

In the terminal, navigate to the root directory of your project and run the following command:

```go
go run main.go
```

This will start the API server on `http://localhost:8080`.

## Step 7: Testing the API

With the API server running, you can now test the defined routes using a tool like cURL or Postman.

For example, to test the `GET /users` route, you can use the following cURL command:

```sh
curl http://localhost:8080/users
```

You should receive a JSON response containing the list of users.

Similarly, you can test the other routes (`POST /users`, `PUT /users/:id`,
and `DELETE /users/:id`) by sending requests to the corresponding URLs.

Congratulations! You have successfully created a basic RESTful API using the `muxer` package.

## Conclusion

In this tutorial, we have learned how to use the `muxer` package to build
a RESTful API in Go. We explored how to define routes, handle different
HTTP methods, extract path parameters, and use middleware functions. You
can now expand upon this foundation to build more complex APIs with
additional features.

Remember to refer to the official documentation of the `muxer` package for
more details on its features and capabilities.

Happy coding!
