# API Flow and Architecture

## Project Structure

The project follows a modular design, with different functionalities separated into distinct packages (e.g., `auth`, `customer`, `dashboard`, `inbound`, `outbound`). This promotes code organization and maintainability.

It utilizes a `factory` pattern for creating repositories and services. This pattern decouples the instantiation of objects from their usage, making it easier to manage dependencies and switch implementations if needed.

## API Flow

1.  **HTTP Request Handling:** Incoming HTTP requests are first received by the `chi` router, which is configured in `server/server.go`. The router is responsible for matching the request path and method to the appropriate handler.
2.  **Middleware:** Before reaching the specific handler, requests pass through several middleware functions. These include:
    *   `cors`: Handles Cross-Origin Resource Sharing.
    *   `middleware.RequestID`, `middleware.RealIP`, `middleware.Logger`, `middleware.Recoverer`, `middleware.StripSlashes`, `middleware.Timeout`: These are standard `chi` middleware for request ID generation, IP address identification, logging, panic recovery, URL path cleaning, and request timeout handling.
    *   `jwtauth.Verifier` and `jwtauth.Authenticator`: These middleware handle JWT-based authentication, ensuring that protected routes are accessed only by authenticated users.
3.  **Routing to Domain Handlers:** Based on the route, the request is dispatched to a specific domain handler (e.g., `customerHandler` for `/customers` routes, `mawbHandler` for `/mawb` routes). These handlers are also defined in `server/server.go`.
4.  **Service Layer:** The domain handler then calls the appropriate method in the corresponding service (e.g., `customerService.CreateCustomer()`). Services contain the business logic of the application. They interact with repositories to access and manipulate data.
5.  **Repository Layer:** Repositories are responsible for data access. They interact with the database (PostgreSQL in this case) to perform CRUD (Create, Read, Update, Delete) operations.
6.  **Response:** After processing the request, the service returns the result to the handler, which then sends an HTTP response back to the client.

## Key Components

*   `main.go`: The entry point of the application. It initializes configurations, database connections, GCS client, factories, and the HTTP server.
*   `server/server.go`: Configures the `chi` router, middleware, and defines routes for different domains. It also includes utility functions for handling API responses and errors.
*   `factory/`: Contains factory functions for creating repository and service instances. This promotes loose coupling and dependency injection.
*   **Domain Packages** (e.g., `customer/`, `mawb/`, `auth/`): Each domain package typically contains:
    *   `*.go` (e.g., `customer.go`): Defines the data structures (structs) for that domain.
    *   `service.go`: Implements the business logic for the domain.
    *   `repository.go`: Implements data access logic for the domain, interacting with the database.
    *   Some domains might also have `logging.go` for specific logging needs.
*   `database/`: Contains code for establishing database connections (e.g., `postgresql.go`).
*   `config/`: Handles application configuration loading.
*   `utils/`: Contains utility functions used across the application.

This architecture separates concerns, making the codebase more organized, testable, and easier to maintain.

# Guide on Adding a New API

To add a new API, you would typically follow these steps:

1.  **Define Your New Domain/Functionality:**
    *   Decide what new data or operations you want to expose via an API. For example, let's say you want to add a "products" API to manage product information.

2.  **Create or Update Domain Package:**
    *   **New Domain:** If your API introduces a completely new concept (like "products"), create a new directory for it (e.g., `products/`).
        *   Inside `products/`, create `products.go` to define the `Product` struct (e.g., with fields like `ID`, `Name`, `Price`, `Description`).
        *   Create `service.go` for the `ProductService` interface and its implementation. This service will contain methods like `CreateProduct`, `GetProductByID`, `UpdateProduct`, `DeleteProduct`.
        *   Create `repository.go` for the `ProductRepository` interface and its implementation. This repository will handle database operations for products.
    *   **Existing Domain:** If your new API extends an existing domain (e.g., adding a new operation to the `customer` domain), you'll modify the existing files in that domain package.

3.  **Implement Service Logic:**
    *   In the `service.go` file for your domain, implement the business logic for your new API endpoints.
    *   For example, in `products/service.go`, the `CreateProduct` method in `ProductServiceImpl` would take product data, validate it, and then call the corresponding repository method to save it to the database.

4.  **Implement Repository Logic:**
    *   In the `repository.go` file, implement the database interaction logic.
    *   For example, in `products/repository.go`, the `CreateProduct` method in `ProductRepositoryImpl` would contain the SQL query or ORM calls to insert a new product record into the `products` table.
    *   You might need to create a new database table (e.g., `products`) if it doesn't exist. This would typically be handled via database migration scripts, which are not explicitly shown in the provided file structure but are a common practice.

5.  **Update Factories:**
    *   **Repository Factory (`factory/repository.go`):**
        *   Add a new method to create an instance of your new repository (e.g., `NewProductRepository`).
        *   Update the `RepositoryFactory` struct to include your new repository interface.
    *   **Service Factory (`factory/service.go`):**
        *   Add a new method to create an instance of your new service (e.g., `NewProductService`). This method will typically take the `RepositoryFactory` as a dependency to get the required repositories.
        *   Update the `ServiceFactory` struct to include your new service interface.
        *   In `main.go`, when `factory.NewServiceFactory` is called, ensure your new repository is correctly passed and initialized if needed.

6.  **Create API Handler:**
    *   In the `server/` directory, create a new handler file for your domain, similar to `server/customer.go` or `server/mawb.go`. For our example, this would be `server/products.go`.
    *   This file will define a struct (e.g., `productsHandler`) that embeds your new service (e.g., `products.ProductService`).
    *   It will also have a `router()` method that defines the specific routes for your API (e.g., `POST /products`, `GET /products/{id}`).
    *   Implement methods in this handler for each route. These methods will:
        *   Parse request data (e.g., from JSON body or URL parameters).
        *   Call the appropriate methods on your service.
        *   Use `render.Respond` to send the HTTP response (success or error).

7.  **Mount New Routes in `server/server.go`:**
    *   In `server/server.go`, within the `New()` function:
        *   Create an instance of your new handler, passing the corresponding service from the `svcFactory`. For example:
            ```go
            productsSvc := productsHandler{s.svcFactory.ProductSvc} // Assuming ProductSvc is added to ServiceFactory
            r.Mount("/products", productsSvc.router())
            ```
        *   Decide if your new API endpoints should be under the authenticated group (`/v1`) or the public group. Add the `r.Mount(...)` call in the appropriate `r.Group(...)` section.

8.  **Add Database Migrations (If Applicable):**
    *   If you added new tables or modified existing ones, create database migration scripts to apply these changes. The project doesn't show a migration system, but tools like `golang-migrate/migrate` or `GORM's` auto-migration features are commonly used.

9.  **Write Tests:**
    *   It's crucial to write unit tests for your new repository methods, service methods, and API handlers to ensure they function correctly.

**Hypothetical Example: Adding a "GET /products/{id}" API:**

1.  **Domain (`products/`):**
    *   `products/products.go`: Define `Product` struct.
    *   `products/service.go`:
        *   `ProductService` interface with `GetProductByID(ctx context.Context, id string) (*Product, error)`.
        *   `ProductServiceImpl` implementing `GetProductByID`.
    *   `products/repository.go`:
        *   `ProductRepository` interface with `FindProductByID(ctx context.Context, id string) (*Product, error)`.
        *   `ProductRepositoryImpl` implementing `FindProductByID` (fetches product from DB).
2.  **Factories:**
    *   Update `factory/repository.go` and `factory/service.go` to include `ProductRepository` and `ProductService`.
3.  **Handler (`server/products.go`):**
    *   `productsHandler` struct.
    *   `router()` method: `r.Get("/{id}", ph.handleGetProductByID)`
    *   `handleGetProductByID(w http.ResponseWriter, r *http.Request)` method:
        *   Extract `id` using `chi.URLParam(r, "id")`.
        *   Call `ph.ProductService.GetProductByID(r.Context(), id)`.
        *   Respond with product data or an error.
4.  **Server (`server/server.go`):**
    *   `productSvc := productsHandler{s.svcFactory.ProductSvc}`
    *   `r.Mount("/v1/products", productSvc.router())` (assuming it's an authenticated route)

By following this structure, you can integrate new APIs into the existing architecture in a consistent and maintainable way. Remember to also update `main.go` if your new service or repository requires specific initialization parameters beyond what the factories provide.
