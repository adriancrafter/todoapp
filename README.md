# Todo List Reference Implementation

This project serves as a basic Todo list application designed to be a reference implementation for future multi-tenant capable applications. Whether built as microservices, micro-monoliths, or monolith apps, this project provides a foundation for integrating key features.

## Key Features:

- **Authentication and Authorization:** The application includes a robust session-based authentication and authorization system for web features. API endpoints use tokens for secure authentication.

- **Role-Based Access Control (RBAC):** A simple RBAC system is implemented to control access to different functions within the application. This allows for fine-grained control over user permissions.

- **Web Views with SSR Rendering:** The web views are designed to support Server-Side Rendering (SSR) for improved performance. Additionally, dynamic content is managed using htmx for a seamless user experience.

## Layered Architecture:

The project follows a layered architecture with an emphasis on the use of interfaces. The graph of dependencies is assembled in the `App` package, including:

- **Config:** Configuration settings for the application.

- **Logger:** Logging functionality for monitoring and debugging.

- **Database Access:** Support for SQL databases, with a current focus on PostgreSQL.

- **Repos:** Repositories for data storage and retrieval, implementing the repository pattern.

- **SQL Queries:** Externalized queries, easily editable in custom files, seamlessly embedded in the binary.

- **Services:** Business logic services.

- **Controllers (Web and API):** Handles HTTP requests and invokes the appropriate service layer functions.

- **Templates:** Streamlined loading from an embedded file system with implemented caching. Future plans include the flexibility to switch to an external filesystem.

## Project Structure:

The project follows a feature-based organization approach, where each feature is encapsulated within its own package. This promotes modularity, reduces import lists, and ensures a clear separation of concerns. The organization includes components such as routes definition, controllers, services, models, view models, data access, and conversion utilities.

### Feature Packages:

- `Todo List`: The core feature package containing components for managing todo lists, items, and associated functionality.

### Future Features:

As the project evolves, additional features will be added, each following a similar distribution of file types within its own package. The consistent structure facilitates maintainability and scalability.

## Notes

## Notes

## Notes

After the foundational structure of the project is established in this exploratory phase, a comprehensive set of tests will then be developed to ensure thorough coverage of the codebase.
