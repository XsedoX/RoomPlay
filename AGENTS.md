# AI Agent Instructions for RoomPlay

This file provides architectural context, strict coding rules, and project conventions for AI coding agents working on the RoomPlay repository.

## 1. Project Context

- **Project Name:** RoomPlay
- **Repository Structure:** Monorepo containing a `/backend` (Go) and `/frontend` (Vue.js/TypeScript).
- **Domain:** A real-time room/music playing application featuring rooms, song queues, enqueued songs, default playlists, user voting, and host device management.

## 2. Tech Stack

### Backend (Go)

- **Language:** Go
- **Architecture:** Clean Architecture / Domain-Driven Design (DDD) with CQRS
- **Database:** PostgreSQL (with raw SQL initialization scripts)
- **API:** REST with Swagger/OpenAPI documentation

### Frontend (Vue.js)

- **Framework:** Vue.js (Composition API)
- **Language:** TypeScript
- **UI Framework:** Vuetify
- **State Management:** Pinia
- **Build Tool:** Vite
- **Package Manager:** Bun
- **Testing:** Vitest (Unit) and Playwright (E2E)

---

## 3. Backend Guidelines (Go)

### Architecture & Directory Structure

The backend strictly follows Clean Architecture. Dependencies point INWARD.

- `/backend/domain`: Core entities, aggregates, and value objects (e.g., `room`, `enqueued_song`, `user`). **Rule:** NO external dependencies or framework imports here.
- `/backend/application`: Use cases structured using **CQRS**. Contains Commands, Queries, and their respective Handlers. Contains interface contracts (`*_contracts/`).
- `/backend/infrastructure`: Implementations of application contracts (e.g., Repositories, DAOs, external Google OIDC services, caching).
- `/backend/presentation`: API controllers, custom middlewares (CORS, JWT cookies), and HTTP response formatting.

### Coding Rules

- **CQRS Pattern:** Every new feature must be split into a Command (for mutations) or a Query (for reads) and placed in the appropriate `application/<feature>` folder.
- **Error Handling:** Use custom domain errors (e.g., `domain_errors/validation_domain_error`) rather than generic Go errors.
- **Database Access:** Repositories handle aggregations, DAOs handle raw data queries. Always use the provided Unit of Work (`i_unit_of_work`) for transactional safety.
- **Testing:** Write unit tests for all Command/Query handlers (e.g., `*_test.go`). Use the established mocks in `backend/test_helpers/integration_tests/` for external dependencies.

---

## 4. Frontend Guidelines (Vue/TS)

### Architecture & Directory Structure

- `/frontend/src/pages`: Route-level components. Keep logic minimal; delegate state to stores.
- `/frontend/src/shared`: Highly reusable, dumb UI components.
- `/frontend/src/stores`: Pinia stores (e.g., `room_store.ts`, `user_store.ts`).
- `/frontend/src/infrastructure`: Data fetching, API clients, and service logic. UI components MUST NOT call `api_client.ts` directly; they must use service classes or stores.

### Coding Rules

- **TypeScript:** Strictly type all interfaces, models, and DTOs (e.g., `IUserListElementDto.ts`, `TSongState.ts`). Avoid `any`.
- **Vue Components:** Use Vue 3 Composition API (`<script setup>`).
- **Styling:** Rely on Vuetify components and the configured themes (`assets/themes.ts`) over custom CSS whenever possible.
- **State:** Use Pinia for cross-component state management.

---

## 5. Development & Verification Commands

When writing or modifying code autonomously, always run the relevant verification steps before completing your task.

**Frontend Commands (run inside `/frontend`):**

- Run Typecheck & Linter: `bun run lint` (assumes standard lint script setup)
- Run Unit Tests: `bun run test:unit`
- Run E2E Tests: `bun run test:e2e`

**Backend Commands (run inside `/backend`):**

- Run Tests: `go test ./...`
- Run formatting: `go fmt ./...`

**Environment Management:**

- Docker Compose is available at the root (`docker-compose.yml`) for spinning up necessary external services like the PostgreSQL database\*.
- Do not change the code unless told to do so. Always show code snippets/fragments in chat instead of making changes directly to the codebase.
