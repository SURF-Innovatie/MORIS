# MORIS

MORIS (Modular Research Information System) is a comprehensive backend and frontend solution designed to streamline the management of research project information. It facilitates the integration of project details, contributors, and related identifiers, serving as a central hub for research metadata.

## Table of Contents
- [Project Overview](#project-overview)
- [Prerequisites](#prerequisites)
- [Development Environment Setup](#development-environment-setup)
  - [Initial Setup](#initial-setup)
  - [Backend Configuration](#backend-configuration)
  - [Frontend Configuration](#frontend-configuration)
  - [Infrastructure & Seeding](#infrastructure--seeding)
  - [Running the Application](#running-the-application)
- [Production Environment Setup](#production-environment-setup)
- [Database Migrations](#database-migrations)
- [Contributing](#contributing)

## Project Overview

MORIS connects various research infrastructure services to provide a unified view of research projects. The system is split into a Go-based backend API (`apps/backend`) and a React-based frontend (`apps/frontend`), managed as a monorepo using Turbo.

## Prerequisites

Before you begin, ensure you have the following installed:
- [nvm](https://github.com/nvm-sh/nvm) (for managing Node.js versions)
- [pnpm](https://pnpm.io/) (Package Manager)
- [Go](https://go.dev/) (1.25.4 or later)
- [wgo](https://github.com/bokwoon95/wgo) (Live reload for Go development)
- [swag](https://github.com/swaggo/swag) (Swagger documentation generator)
- [turbo](https://turbo.build/) (Repo management)
- [podman](https://podman.io/) (Container management)
- [atlas](https://atlasgo.io/) (Database migration tool)

## Development Environment Setup

Follow these steps to set up a local development environment.

### Initial Setup

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd MORIS
    ```

2.  **Setup Node.js environment:**
    ```sh
    nvm use
    pnpm install
    ```

3.  **Initialize Backend:**
    Download Go dependencies and generate Ent definitions and Swag documentation:
    ```sh
    cd apps/backend
    go get ./...
    pnpm run api:generate
    ```

### Backend Configuration

Create the backend environment file and configure the defaults.

1.  **Create `.env`:**
    ```sh
    # Assuming you are still in apps/backend
    cp .env.example .env
    ```

2.  **Configure `.env`:**
    Update the file with the following development defaults:

    ```bash
    # App Configuration
    APP_ENV=dev
    PORT=8080
    JWT_SECRET=this0is1a2secret

    # Database (PostgreSQL)
    DB_HOST=localhost
    DB_PORT=8765
    DB_USER=moris
    DB_PASSWORD=moris
    DB_NAME=moris

    # Redis Cache
    CACHE_HOST=localhost
    CACHE_PORT=6380
    CACHE_PASSWORD=moris
    CACHE_USER=

    # ORCID Integration
    ORCID_CLIENT_ID=
    ORCID_CLIENT_SECRET=
    ORCID_REDIRECT_URL=http://127.0.0.1:3000/orcid-callback
    ORCID_SANDBOX=true

    # Crossref Integration
    CROSSREF_USER_AGENT=MORIS/1.0 (https://github.com/SURF-Innovatie/MORIS)
    CROSSREF_MAILTO=your_email@example.com
    CROSSREF_BASE_URL=https://api.crossref.org

    # RAiD Integration
    RAID_API_URL=https://api.demo.raid.org.au/
    RAID_AUTH_URL=https://auth.demo.raid.org.au/realms/RAiD/protocol/openid-connect/token
    RAID_USERNAME=
    RAID_PASSWORD=

    # Zenodo Integration
    ZENODO_CLIENT_ID=
    ZENODO_CLIENT_SECRET=
    ZENODO_REDIRECT_URL=http://127.0.0.1:3000/zenodo-callback
    ZENODO_SANDBOX=true

    # SURFconext (OpenID Connect) Integration
    SURFCONEXT_ISSUER_URL=
    SURFCONEXT_CLIENT_ID=
    SURFCONEXT_CLIENT_SECRET=
    SURFCONEXT_REDIRECT_URL=http://127.0.0.1:3000/surfconext-callback
    SURFCONEXT_SCOPES=
    ```

### Frontend Configuration

Install dependencies and configure the frontend environment.

1.  **Initialize Frontend:**
    ```sh
    cd ../frontend
    pnpm install
    pnpm run api:generate
    ```

2.  **Configure `.env`:**
    ```sh
    cp .env.example .env
    ```

    Update the file content:
    ```bash
    # App Name
    VITE_APP_NAME="MORIS"

    # API URL (Backend endpoint)
    VITE_API_BASE_URL=/api
    ```

### Infrastructure & Seeding

1.  **Start Services:**
    Spin up the database and Redis containers using Podman from the root directory:
    ```sh
    cd ../..
    podman compose up -d
    ```

2.  **Seed Database (Optional):**
    If you need dummy data for development:
    ```sh
    cd apps/backend
    pnpm run db:seed
    cd ../..
    ```

3. **Add Admin User (if you skipped seeding the database)**
    ```
    cd apps/backend
    pnpm run db:add-admin --email="verify_admin@example.com" --name="Verify Admin" --password="password123"
    cd ../..
    ```

### Running the Application

Start the development server. This runs both the backend (with hot-reload via `wgo`) and the frontend (via Vite) concurrently:

```sh
# From the root directory
pnpm run dev
```

-   **Frontend:** `http://127.0.0.1:3000`
-   **Backend API:** `http://localhost:8080`
-   **API Docs:** `http://localhost:8080/swagger/index.html`

---

## Production Environment Setup

> **Coming Soon**
>
> The production environment setup is currently under development. It will utilize optimized Docker images for both the backend and frontend services.

---

## Database Migrations

MORIS uses [Atlas](https://atlasgo.io/) for versioned database migrations with [Ent](https://entgo.io/).

### Prerequisites

- [Atlas CLI](https://atlasgo.io/getting-started#installation) installed (`brew install ariga/tap/atlas` or `curl -sSf https://atlasgo.sh | sh`)
- Docker/Podman running (for the dev database used in schema diffing)

### Generating Migrations

After modifying ent schemas in `apps/backend/ent/schema/`, generate a new migration:

```sh
cd apps/backend
pnpm run db:migrate:diff <migration_name>
```
This creates timestamped SQL files in `apps/backend/ent/migrate/migrations/`.

### Applying Migrations

Apply pending migrations to your database:

```sh
cd apps/backend
pnpm run db:migrate:apply
```

### Checking Migration Status

View which migrations have been applied:

```sh
cd apps/backend
pnpm run db:migrate:status
```

## Contributing

Please read our contributing guidelines before submitting a Pull Request.

### Reporting Issues
-   Use the [Bug Report](.github/ISSUE_TEMPLATE/bug_report.md) template for bugs.
-   Use the [Feature Request](.github/ISSUE_TEMPLATE/feature_request.md) template for suggestions.

### Pull Requests
-   Use the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md) when opening a PR.
-   Ensure all checks pass (`pnpm lint`, `pnpm build`).