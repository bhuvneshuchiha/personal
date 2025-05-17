
# Go Learning Journey: Roadmap & Progress

## ✅ What You’ve Covered So Far

### 🧠 Core Fundamentals

* Variables, constants, types
* Functions and error handling
* Control flow: `if`, `switch`, `for`, etc.
* Pointers, slices, maps, structs
* Interfaces

### ⚙️ Concurrency

* Goroutines
* Channels (unbuffered, buffered)
* `sync.WaitGroup`
* `select` statement
* Timeouts with `context.Context`
* Cancelation and values using `context.WithTimeout`, `WithValue`, etc.

### 🧪 Testing

* `testing` package
* Using `t.Run()` and table-driven tests
* Assertion with `testify/assert`
* Structuring test cases
* Pretty test output using `gotest`

---

## ὓ9 What’s Next in Testing

### ✅ Completed

* Basic assertions with `testify`
* Table-driven tests
* Handling multiple test cases

### ⏳ Upcoming

* Testing error-returning functions
* Subtests for edge cases
* Benchmark tests (optional)
* Mocks & interfaces (when testing dependencies)

---

## 💃 What’s Coming After Testing

### 🚀 Project Phase – Practical Application of Everything

#### 📁 Backend Application in Go

1. Project setup with Go modules
2. HTTP server using `net/http` or `chi`
3. Routing, middlewares, handlers
4. CRUD operations
5. Working with JSON and request/response cycles
6. File handling (optional)
7. Integrating context properly
8. Unit testing + integration tests

#### 🗃️ Databases

* SQLite/Postgres (starting with `database/sql`)
* Struct-to-DB mapping
* Migrations & querying
* Testing DB code

#### 🧐 Bonus Topics

* Logging
* Environment configuration
* Graceful shutdowns
* Deploying or sharing the app
