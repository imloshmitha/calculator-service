# Calculator

A full stack calculator with both basic and advanced calc support. The interface is a React and TypeScript client. The calculation API is a Go REST service. Every arithmetic result is produced by the Go API. The browser sends operands and displays the response.

The application is not deployed. It does not use a database, authentication, or Docker.

## Technology Stack

Versions in the "Used in this checkout" column were read from the lockfile or from the tools on this machine. Ranges in the "Project requirement" column come from `go.mod` or from dependency engine fields. `package.json` does not pin Node.js or npm.

| Technology | Used in this checkout | Project requirement | Purpose |
| --- | --- | --- | --- |
| Go | 1.27.2 | 1.22 or newer (`backend/go.mod`) | Backend language and REST API |
| React | 19.3.0 (`package-lock.json`; declared `^19.2.8`) | `^19.2.8` | Frontend UI |
| TypeScript | 6.0.3 (declared `~6.0.2`) | `~6.0.2` | Frontend language |
| Vite | 8.3.4 (declared `^8.3.0`) | Node.js `^20.19.0` or `>=22.12.0` for the Vite package | Development server and production build |
| Vitest | 5.0.4, with `@vitest/coverage-v8` 5.0.4 | Node.js `^22.12.0`, `^24.0.0`, or `>=26.0.0` | Frontend test runner and V8 coverage |
| React Testing Library | `@testing-library/react` 16.3.3 | `^16.3.3` | Frontend component and behavior tests |
| Go standard library `net/http` | Included with Go 1.27.2 | Part of the Go toolchain | HTTP server and routing |
| Go `testing` and `net/http/httptest` | Included with Go 1.27.2 | Part of the Go toolchain | Backend unit and handler tests |
| HTTP and JSON | `POST /api/calculate` with `Content-Type: application/json` | Same contract in the Go handler and `frontend/src/api.ts` | Frontend/backend communication |
| npm | 10.9.7 | Not pinned | Frontend package management |
| Git | 2.50.1 | Not pinned | Source control |

This working copy is not a Git repository and has no remote URL.

## Features

The keypad implements:

- Addition, subtraction, multiplication, and division
- Exponentiation, labeled `x^y` and sent as `power`
- Square root, labeled `√` and sent as `sqrt`
- Percentage: **b percent of a**, `(a * b) / 100`. `200 % 15 =` displays `30`
- Clear (`C`), delete (`DEL`), and a single decimal point
- Validation and API error handling, including division by zero
- A layout that fits a narrow screen: below 380px, padding and key size shrink (`frontend/src/Calculator.css`)

`+`, `−`, `×`, `÷`, `x^y`, and `%` record the operation and wait until `=` before calling the API. Square root calls the API as soon as `√` is pressed.

## Architecture

`frontend/src/useCalculator.ts` stores the number on screen and, for a binary operation, the left operand and operation name. It does not perform the arithmetic.

`frontend/src/api.ts` sends the request:

```http
POST /api/calculate
Content-Type: application/json
```

`backend/cmd/server` listens with `net/http` and passes the request to `backend/internal/handler`. The handler checks the method, decodes one JSON object, and calls `backend/internal/calculator`. That package returns a number or a sentinel error and does not know about HTTP. The handler maps the error to a JSON body. The client prints `result`, or replaces the main display with the error `message`.

In development, Vite proxies `/api` to `http://localhost:8080` (`frontend/vite.config.ts`). `npm run build` writes a static client to `frontend/dist`. That build requests `/api/calculate` on the same origin. The proxy is not part of the build. Serving the static files from another origin requires the API on that same origin or a reverse proxy. The Go server does not send CORS headers. Nothing in this repository deploys either part.

## Prerequisites

- Go 1.22 or newer. This checkout used Go 1.27.2.
- Node.js `^22.12.0`, `^24.0.0`, or `>=26.0.0`, so both Vite and Vitest can run. This checkout used Node.js v22.22.2.
- npm, to install frontend packages and run the scripts. This checkout used npm 10.9.7. The repository does not pin npm.
- Git, to clone the repository. This machine has Git 2.50.1. The project files do not require a specific Git version.

## Getting Started

The directory that contains `backend/` and `frontend/` is the repository root. This copy is not yet a Git repository, so replace the placeholder with the real GitHub URL after the remote exists.

```bash
git clone <repository-url>
cd <repository-directory>
```

Start the API and the client in separate terminals.

Terminal 1, from the repository root:

```bash
cd backend
go run ./cmd/server
```

The process listens on port `8080` unless `PORT` is an integer from 1 through 65535.

```bash
PORT=9090 go run ./cmd/server
```

Terminal 2, from the repository root:

```bash
cd frontend
npm install
npm run dev
```

Open the URL Vite prints. With the default Vite port that is `http://localhost:5173/`. The client calls `/api/calculate` on the Vite origin, and Vite forwards `/api` to `http://localhost:8080`. No `.env` file is required for the default ports. If `PORT` is not `8080`, calculations fail until the proxy target in `frontend/vite.config.ts` is changed to match.

## API Documentation

### `POST /api/calculate`

Requests and responses use `Content-Type: application/json`.

Request:

```json
{
  "operation": "add",
  "a": 19,
  "b": 9
}
```

Success, `200`:

```json
{
  "result": 28
}
```

Division by zero, `400`:

```json
{
  "code": "division_by_zero",
  "message": "division by zero"
}
```

`operation` is required. `a` and `b` are JSON numbers, decoded as `float64`. `null` or an omitted required operand is treated as missing. `0` is a valid operand. Names are lowercase and case-sensitive. The body must be a single JSON object. Unknown fields, a second JSON value, and a body larger than 1 MiB are rejected.

| Operation | Required operands | Meaning |
| --- | --- | --- |
| `add` | `a`, `b` | `a + b` |
| `subtract` | `a`, `b` | `a - b` |
| `multiply` | `a`, `b` | `a * b` |
| `divide` | `a`, `b` | `a / b`. `b` must not be zero. |
| `power` | `a`, `b` | `a` raised to `b` with `math.Pow`. `0^0` is rejected. |
| `sqrt` | `a` | Square root of `a`. A supplied `b` is ignored. |
| `percentage` | `a`, `b` | `b` percent of `a`: `(a * b) / 100`. `a = 200` and `b = 15` returns `30`. |

| Status | When |
| --- | --- |
| `200` | The calculation succeeded. |
| `400` | The JSON, operation, or operands are invalid, or the result is not finite. |
| `405` | The method is not `POST`. The response includes `Allow: POST`. |
| `500` | An unexpected calculator error. Ordinary validation does not use this status. |

| Code | Message |
| --- | --- |
| `invalid_json` | `request body must be one JSON object`, `request body contains an unknown field`, `request body must not contain trailing data`, or `request body exceeds 1 MiB` |
| `missing_operation` | `operation is required` |
| `unsupported_operation` | `operation is not supported: "<name>"` |
| `missing_operand` | `operand a is required` or `operand b is required` |
| `division_by_zero` | `division by zero` |
| `negative_sqrt` | `square root of a negative number is not defined` |
| `undefined_power` | `zero raised to the power of zero is undefined` |
| `non_finite_result` | `operands must be finite numbers` or `calculation result is not a finite number` |
| `method_not_allowed` | `method not allowed` |
| `internal_error` | `internal error` |

```bash
curl -sS -D - -X POST http://localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":19,"b":9}'
```

```bash
curl -sS -D - -X POST http://localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":1,"b":0}'
```

### `GET /health`

Success, `200`:

```json
{
  "status": "ok"
}
```

Any other method receives `405` with `Allow: GET`.

Any other path receives the standard library's plain-text `404`.

## Testing

Run each command from the directory shown.

Backend tests, from `backend/`:

```bash
cd backend
go test ./...
```

Frontend tests, from `frontend/`:

```bash
cd frontend
npm test
```

`npm test` is `vitest run`. The suite has 15 tests.

Backend coverage, from `backend/`:

```bash
cd backend
go test -cover ./...
```

Frontend coverage, from `frontend/`:

```bash
cd frontend
npm run test:coverage
```

`npm run test:coverage` is `vitest run --coverage` using the V8 provider configured in `frontend/vite.config.ts`.

Frontend production build, from `frontend/`:

```bash
cd frontend
npm run build
```

`npm run build` is `tsc -b && vite build`. Output goes to `frontend/dist/`, which Git ignores.

## Test Coverage

These are the reported results for the commands above. Package percentages are not combined into one backend number.

Frontend, from `npm run test:coverage` (15 tests passed):

| Metric | Coverage |
| --- | --- |
| Statements | 64.54% (142/220) |
| Branches | 56.29% (85/151) |
| Functions | 88.57% (31/35) |
| Lines | 64.05% (139/217) |

Backend, from `go test -cover ./...`:

| Package | Statement coverage |
| --- | --- |
| `internal/calculator` | 95.6% |
| `internal/handler` | 91.2% |
| `cmd/server` | 0% |

`cmd/server` is the process entry point (`main` and port parsing). It currently has no direct test coverage.

## Design Decisions and Limitations

- The client is a sequential calculator. It does not parse a full expression or apply operator precedence. One binary operation is pending. `=` sends `operation`, `a`, and `b`. After a result, the next operator uses that result as the next left operand, and the next `=` sends a separate request. Pressing `=` again repeats the last completed binary operation. Another operator pressed before a right-hand number replaces the pending operation. If a right-hand number is already displayed, another operator still does not call the API: it keeps the original left operand, switches the operation, and the next digit starts a new right-hand number.
- Percentage means `b` percent of `a`, `(a * b) / 100`. It does not mean "`a` plus `b` percent." `200 % 15 =` sends `a = 200` and `b = 15` and shows `30` when the API returns that result.
- Division by zero is rejected before IEEE division. The API returns `400` with code `division_by_zero` and message `division by zero`. The main display shows that message instead of a number. `C` or a new entry clears it.
- Operands and results are IEEE-754 `float64`. A sum such as `0.1 + 0.2` can differ by normal binary rounding. Percentage divides the larger-magnitude operand by 100 before multiplying, which avoids some intermediate overflow, then rejects a non-finite final value. That does not recover every result that would mathematically fit in a `float64`.

### Numeric Precision

The calculator uses Go's `float64` type for arithmetic operations. Floating-point numbers cannot represent every decimal value exactly, so some calculations may produce small rounding differences. For example, `0.1 + 0.2` may produce `0.30000000000000004` rather than exactly `0.3`.

This is an intentional trade-off for a basic calculator. Results are returned as JSON numbers without artificial rounding. The frontend may format results for display, but the API preserves the calculated numeric value.

- The handler rejects malformed JSON, unknown fields, trailing JSON, missing operands, unsupported operations, negative square roots, `0^0`, and non-finite results. Those responses use the codes and messages in the API section.
- The development proxy only exists while `npm run dev` is running and only forwards `/api` to port `8080`. The static build has no proxy and no CORS headers from the API.
- There is no authentication, persistence, request logging, or rate limiting. `postman/` contains no collection.

## AI Usage

The prompt record is in [AI_USAGE.md](AI_USAGE.md).
