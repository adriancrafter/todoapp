# `auth` package

This package follows a feature-based organization approach, where all components related to a specific feature are grouped together within a single package. Unlike the traditional layered structure with separate packages for controllers, models, and other components, this approach keeps all parts of a feature within its own package. The files within the package are organized to provide a clear structure and separation of concerns.

## Feature-Based Organization Advantages

- **Modularity**: Each feature is encapsulated within its own package, promoting modularity. This makes it easier to understand, modify, and extend individual features without affecting other parts of the system.

- **Reduced Import Lists**: With feature-based organization, you only need to import the specific components related to the feature you are working on. This leads to shorter and more focused import lists in your code.

- **Clear Separation of Concerns**: Despite having all components in the same package, there is still a clear separation of concerns. Files are organized based on their responsibilities, making it easy to locate and understand different parts of the feature.

- **Easier Navigation**: Developers can easily navigate through the codebase, focusing on the feature they are currently working on. This reduces cognitive load and helps maintain a clear mental model of the codebase.

## Package Structure

### - `webroutes.go`

Defines route definitions for web endpoints related to the todo feature. Maps routes to corresponding web controller handlers.

### - `apiroutes.go`

Defines routes for API endpoints specific to the todo feature. Maps routes to corresponding API controller handlers.

### - `controller.go`

Handles HTTP requests and invokes the appropriate service layer functions for the todo feature.

### - `service.go`

Contains business logic related to the todo feature. Interacts with the repository for data access and performs necessary operations.

### - `repo.go`

Implements the repository pattern for data storage and retrieval related to the todo feature. The default implementation is for PostgreSQL, but it's designed to be extendable to other databases.

### - `list.go`

Defines the model for todo lists, including fields such as list name, creation date, etc.

### - `item.go`

Model for individual todo list items, containing information like task description, due date, completion status, etc.

### - `listvm.go`

View model for rendering todo lists. It encapsulates data and behavior needed for rendering lists in the UI.

### - `listitemvm.go`

View model for rendering individual todo list items. Similar to `listvm.go`, but specific to list items.

### - `listda.go`

Data access layer for todo lists. Business structs operate with Golang's standard types, and flattened structs based on SQL null-types are used for data access. This design avoids the complexity of using null types in business models.

### - `listitemda.go`

Data access layer for todo list items. Similar to `listda.go`, this file manages data access for individual list items.

### - `criteriavm.go`

View model for criteria related to searching and filtering todo lists.

### - `convvm.go`

Conversion utilities for converting view models to models and vice versa. Helps in maintaining separation of concerns between layers.

### - `convda.go`

Conversion utilities for converting between persistence data access and business models. Useful for transforming data when interacting with the database.

### - `interface.go`

Contains interface definitions for API and web controllers, service, repository, and other components related to the todo feature.

## Note

As we continue to add new features to the project, you can expect new packages following a similar distribution of file types (routes, controllers, service, models, vm, da, etc.). The names of these packages may vary, reflecting the models and functionality specific to each feature.
