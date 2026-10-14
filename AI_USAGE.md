# AI usage

Entries 1–5 are prompts collected during development. Their wording is preserved. Formatting, such as lists and code blocks, was cleaned up. They are not a complete timeline.

Entries 6–8 are summaries of later work. They are not verbatim prompts. The exact wording of those requests was not kept in this file.

## 1. Go backend implementation

### Verbatim prompt: Go Calculator Backend

I am building a full-stack calculator application.

#### Objective

Build a clean, maintainable Go REST API for a traditional calculator. The frontend will be developed separately using React and TypeScript.

#### 1. Project structure

Set up the following structure inside the existing project root. Do not create another nested `calculator-service` directory.

```text
calculator-service/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── calculator/
│   │   │   ├── calculator.go
│   │   │   └── calculator_test.go
│   │   └── handler/
│   │       ├── handler.go
│   │       └── handler_test.go
│   ├── go.mod
│   └── .gitignore
├── frontend/
├── postman/
├── README.md
├── AI_USAGE.md
└── .gitignore
```

Create the backend files and the empty `frontend/` and `postman/` directories. Do not implement the frontend yet. The root README may initially contain a short project outline, but do not invent completed features.

#### 2. Technology and design

- Use Go and the standard `net/http` package.
- Avoid external dependencies unless genuinely necessary.
- Keep calculation logic independent of HTTP handlers.
- Use small, readable functions and idiomatic Go.
- Use the `internal` directory to organize backend implementation.
- Do not introduce a database, authentication, elaborate interfaces, or unnecessary abstractions.
- Read configuration such as the server port from an environment variable, defaulting to `8080`.

#### 3. API contract

Implement one endpoint:

`POST /api/calculate`

Accept JSON requests containing an operation and its operands.

Example request:

```json
{
  "operation": "add",
  "a": 19,
  "b": 9
}
```

Example successful response:

```json
{
  "result": 28
}
```

Support these operations:

- `add`: addition
- `subtract`: subtraction
- `multiply`: multiplication
- `divide`: division
- `power`: exponentiation
- `sqrt`: square root of `a`
- `percentage`: calculate `b` percent of `a`

For example, percentage with `a = 200` and `b = 15` should return `30`.

Document the percentage convention in the code and API documentation.

#### 4. Validation and error handling

- Reject malformed JSON and unsupported operations.
- Validate required operands according to the operation.
- Reject division by zero.
- Reject square roots of negative numbers.
- Handle overflow and non-finite calculation results.
- Return consistent JSON error responses with a machine-readable error code and a useful message.
- Use appropriate HTTP status codes: `200` for success, `400` for invalid requests, and `405` for unsupported HTTP methods.
- Set the appropriate JSON content type.
- Avoid silently accepting trailing JSON data or unexpected request bodies.

Do not add complicated validation frameworks.

#### 5. Automated tests

Write table-driven unit tests for calculation logic, covering:

- Every supported operation
- Positive and negative operands
- Zero and decimal values
- Division by zero
- Negative square roots
- Invalid operations
- Overflow or non-finite results

Write HTTP handler tests using Go's `net/http/httptest` package, covering:

- Successful requests
- Malformed JSON
- Missing operands
- Unsupported operations
- Invalid mathematical inputs
- Unsupported HTTP methods
- Correct HTTP status codes and JSON response structures

Do not claim tests pass unless they have actually been executed.

#### 6. Maintainability and documentation

- Use meaningful names and comments only where they clarify behavior.
- Keep business logic separate from HTTP transport logic.
- Add a health endpoint such as `GET /health` only if it can be implemented simply and documented.
- Provide commands to start the server and run the tests.
- Include example curl requests for success and error scenarios.
- Document design choices and assumptions.
- Do not implement Docker, the React frontend, deployment, or the Postman collection in this step.

#### 7. Completion checklist

After implementation:

1. Explain the purpose of each created file.
2. Show the final backend directory tree.
3. Run `gofmt` on Go source files.
4. Run `go test ./...` from the backend directory and report the actual result.
5. Fix any failing tests.
6. Show the exact commands needed to start the server and test the API.
7. Explain any remaining limitations or unfinished items.

Work within the existing project root. Prioritize correctness, clarity, and maintainability over extra features.

## 2. Backend improvements

### Verbatim prompt

Review the existing Go calculator implementation and make only these targeted improvements:

1. Percentage calculation: avoid unnecessary intermediate overflow in `(a * b) / 100` where practical. Validate the final result using the existing non-finite result handling. Add tests for large operands, overflow, and normal percentage calculations. Do not claim that every mathematically representable result can be recovered from floating-point intermediate overflow.
2. Power operation: explicitly decide how to handle `0^0`. For this calculator API, reject it as an invalid mathematical input with a consistent JSON error response. Add a corresponding unit test and handler test.
3. Keep all existing API contracts and passing tests intact. Do not refactor unrelated code, add dependencies, or introduce new features.

Run `gofmt` and `go test ./...` from the backend directory. Report the actual test results and summarize the changes.

## 3. React frontend

### Verbatim prompt

Build the React frontend for my existing full-stack calculator assignment in the `frontend/` directory.

First inspect the existing Go backend and its API implementation. Do not modify the working backend unless a change is necessary and explain why.

Requirements:

- Use React, TypeScript, and Vite.
- Build a polished, responsive, traditional calculator interface with a display and buttons for digits, decimal, addition, subtraction, multiplication, division, clear, delete, and equals.
- Also support power, square root, and percentage.
- Connect to the existing `POST /api/calculate` endpoint. Inspect the backend's actual request and response formats rather than assuming them.
- All arithmetic calculations must be performed by the Go API, not duplicated in frontend calculation logic.
- Use a Vite development proxy to forward `/api` requests to `http://localhost:8080`, avoiding CORS configuration for local development.
- Support sequential calculations, such as pressing `19`, `+`, `9`, `−`, `6`, `=`, using the previous result. This is a sequential calculator, not a full expression parser with mathematical operator precedence.
- Handle loading, API errors, division by zero, invalid inputs, and repeated button presses gracefully.
- Define and document the percentage button's behavior consistently with the backend's implementation.
- Keep components and API calls reasonably organized. Avoid unnecessary libraries and overengineering.
- Add frontend unit tests using Vitest and React Testing Library, including basic rendering, button interactions, API success, and API error handling.
- Add clear frontend setup and test commands to the README, preserving existing backend documentation.
- Record the actual prompts used in `AI_USAGE.md`. Do not invent prompts or claim tests passed without running them.
- Run the frontend tests and production build. Fix failures and report the commands run and their results.

Do not add authentication, a database, Docker, or hosting requirements.

## 4. UI refinements

### Verbatim prompt

In `frontend/src/Calculator.tsx`:

1. Remove the percentage explanation displayed below the calculator buttons: "Percentage is b percent of a. Enter 200, %, 15, = to get 30."
2. Keep the division-by-zero error message exactly as "division by zero" (lowercase). Do not change it to "Cannot divide by zero" or any other wording.
3. Keep the percentage button and all calculator functionality unchanged.
4. Do not modify the Go backend or any unrelated UI elements.

Fix the calculator display in `frontend/Calculator.tsx` (or wherever the display component is implemented).

Currently, when dividing by zero, the display shows a large `0` with a smaller `division by zero` message above it. This is misleading because users may interpret zero as the calculation result.

Required changes:

1. When a calculation fails, show the error message `division by zero` prominently in the main display, replacing the large numeric value.
2. Hide the large numeric result while the error is active. Do not show `0` as if it were the result of the failed calculation.
3. Preserve the exact lowercase error text `division by zero`.
4. When the user presses `C` or starts a new calculation, clear the error and restore the normal numeric display.
5. Keep successful calculation behavior, API integration, existing styling, and all other operations unchanged.
6. Add or update a frontend test verifying that division by zero displays the error instead of a large `0`.
7. Run the frontend tests and production build, and report the actual results.

Do not modify the Go backend or change the API error response.

## 5. Calculation behavior

### Verbatim prompt

Fix the inconsistent calculation timing in the React calculator in `frontend/`.

Expected behavior:

1. Binary operations (`+`, `−`, `×`, `÷`, `xʸ`) must not calculate or show the final result until the user presses `=`.
2. Pressing a binary operator should store the first operand and selected operation, then allow the user to enter the second operand.
3. Pressing `=` should send the operation and both operands to the existing Go API and display the returned result.
4. Square root should calculate immediately when `√` is pressed, using the Go API.
5. Keep percentage behavior consistent with its existing implementation and document how it works.
6. Preserve clear (`C`), delete (`DEL`), decimal input, sequential calculations, and API error handling. Division by zero must display exactly `division by zero`.
7. Do not perform arithmetic calculations in the frontend; use the Go API.
8. Do not modify the Go backend or unrelated styling.

Inspect the current state management and fix the root cause rather than applying a cosmetic workaround. Add or update tests to verify that binary operations do not call the API until `=` is pressed, and that pressing `=` displays the API result. Run the frontend tests and production build, and report the actual results.

## 6. Server behavior

### Summary, not a verbatim prompt

The API was extended without changing arithmetic or existing error codes. Unknown routes return `404` JSON `{"code":"not_found","message":"not found"}`. `/health` and `/api/calculate` stay as they were. `cmd/server` handles `SIGINT` and `SIGTERM` with `http.Server.Shutdown` and an 8-second limit, and treats `http.ErrServerClosed` as a normal stop. `ReadHeaderTimeout` stays 5 seconds. `ReadTimeout` is 30 seconds, `WriteTimeout` is 15 seconds, and `IdleTimeout` is 120 seconds. Handler tests cover the JSON `404`. `cmd/server` still has no tests.

## 7. Docker, Nginx, and Postman

### Summary, not a verbatim prompt

Docker and a Postman Collection v2.1 file were added. `docker-compose.yml` builds both services, publishes the frontend on host port 3000, and does not publish the backend port. Nginx serves `/`, `/index.html`, `/favicon.svg`, and `/assets/`, and proxies `/api/`, `/health`, and other paths to the Compose service `backend`.

A check of `GET http://localhost:3000/no-such-route` first returned the calculator HTML with status 200, because Nginx sent unknown paths to `index.html`. `docker/nginx.conf` was changed so those paths reach the API. A later check returned `404` and `{"code":"not_found","message":"not found"}`.

The collection has 21 requests. Every URL uses `{{baseUrl}}`. `baseUrl` is `http://localhost:3000`.

## 8. Documentation and ignore rules

### Summary, not a verbatim prompt

`README.md` was revised to match the repository: the public URL `https://github.com/imloshmitha/calculator-service`, Docker and Postman instructions, server timeouts, the JSON `404`, and a Numeric Precision subsection. A sentence that named an employer and a take-home assignment was removed from the README. Handler coverage in the README was updated from 91.2% to 91.4% after `go test -cover ./...`.

Root `.gitignore` gained `frontend/coverage/`. In `backend/.gitignore`, the pattern `server` also ignored `backend/cmd/server/main.go`. It was replaced with `/server`, which ignores only a binary in the backend root.

## Review and verification

These checks were run in this workspace. This file does not claim they were also run by hand outside that work.

| Check | Result |
| --- | --- |
| `gofmt -l` on the changed Go files | No output. The files were already formatted. |
| `go test ./...` from `backend/` | Passed. `internal/calculator` and `internal/handler` ok. `cmd/server` has no test files. |
| `go vet ./...` from `backend/` | Passed, with no findings. |
| `go test -cover ./...` from `backend/` | `cmd/server` 0.0%. `internal/calculator` 95.6%. `internal/handler` 91.4%. |
| `npm test` from `frontend/` | Passed. 15 tests. |
| `npm run build` from `frontend/` | Passed. |
| `npm run test:coverage` from `frontend/` | Passed. 15 tests. Statements 64.54% (142/220), branches 56.29% (85/151), functions 88.57% (31/35), lines 64.05% (139/217). |
| `docker compose config` | Succeeded. |
| `docker compose up --build -d` | Both containers started. The backend health check reported healthy. |
| HTTP checks on port 3000 after the Nginx fix | `GET /health` returned 200 `{"status":"ok"}`. `POST /api/calculate` with add `2` and `3` returned 200 `{"result":5}`. Divide `1` by `0` returned 400 `{"code":"division_by_zero","message":"division by zero"}`. `GET /no-such-route` returned 404 `{"code":"not_found","message":"not found"}`. `/`, `/index.html`, `/favicon.svg`, and the built JavaScript file returned 200. |
| Newman | `newman run postman/calculator-api.postman_collection.json` passed: 21 requests, 45 assertions, 0 failures, using `http://localhost:3000`. |
| `git diff --check` | Exited 0 for the README edits checked with it. |

`git check-ignore -v backend/cmd/server/main.go` exited 1 after the ignore-rule fix, which means that source file is not ignored.

Issues corrected during these checks are the ignore rule for `cmd/server`, the README coverage figure, and the Nginx fallback that hid the JSON `404`. The containers were left running after the successful port 3000 checks. Whether the UI was exercised manually after that is not recorded here.
