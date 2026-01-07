# MORIS

MORIS (Modular Research Information System) is a comprehensive backend and frontend solution designed to streamline the management of research project information. It facilitates the integration of project details, contributors, and related identifiers, serving as a central hub for research metadata.

## Table of Contents
- [Project Overview](#project-overview)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Configuration](#configuration)
  - [Building the Project](#building-the-project)
  - [Running the Project](#running-the-project)
- [Usage](#usage)
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

## Setup

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

3.  **Setup Go environment:**
    Ensure you are in the `apps/backend` directory or have the workspace configured correctly for Go tools.

### Configuration

The project uses `.env` files for configuration. Copy the example files to create your local configuration:

-   **Backend:**
    ```sh
    cp apps/backend/.env.example apps/backend/.env
    ```
    Edit `apps/backend/.env` to set your database credentials, API secrets (ORCID, Crossref), and JWT secret.

-   **Frontend:**
    ```sh
    cp apps/frontend/.env.example apps/frontend/.env
    ```
    Edit `apps/frontend/.env` to point to your backend API URL.

### Building the Project

Use Turbo to build the entire monorepo:
```sh
pnpm build
```
Or build specific parts:
```sh
turbo run build --filter=backend
turbo run build --filter=frontend
```

### Running the Project

1.  **Start Dependencies (Database, Redis, etc.):**
    ```sh
    pnpm db:start
    ```
    This uses `podman compose` to start the required services.

2.  **Run Development Server:**
    ```sh
    pnpm dev
    ```
    This will start both the backend (using `wgo` for hot reload) and the frontend (using Vite).

## Usage

-   **Frontend:** Accessed via `http://127.0.0.1:3000` (default Vite port).
-   **Backend API:** Accessed via `http://localhost:8080` (or configured port).
-   **API Documentation:** Usage of `swag` implies Swagger UI is available, typically at `/swagger/index.html` on the backend URL.

## Contributing

Please read our contributing guidelines before submitting a Pull Request.

### Reporting Issues
-   Use the [Bug Report](.github/ISSUE_TEMPLATE/bug_report.md) template for bugs.
-   Use the [Feature Request](.github/ISSUE_TEMPLATE/feature_request.md) template for suggestions.

### Pull Requests
-   Use the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md) when opening a PR.
-   Ensure all checks pass (`pnpm lint`, `pnpm build`).