## Template Manager

The template manager is designed to efficiently handle and organize templates for different controllers and handlers. It takes care of loading, parsing, and caching the templates at startup, ensuring that they are pre-processed and ready to use when needed by the controllers. This guide explains how to use the Template Manager, detailing file placement, layout fallbacks, and the template preparation process.

### Overview

The Template Manager automates the process of loading and caching HTML templates for web applications. It supports:
- Controller-specific templates
- Handler-specific templates within controllers
- Shared partial templates
- Layouts with fallback mechanisms

### File Structure

Your templates should be organized in the `assets/templates/web/` directory with the following structure:

``` 
assets/templates/web/
├── {controller}
│ ├── layouts
│ │ ├── layout.tmpl (optional)
│ │ └── {handler}_layout.tmpl (optional)
│ ├── partials
│ │ ├── _{partial}.tmpl
│ │ └── ...
│ ├── {handler}.tmpl
│ └── ...
└── layout
└── base.tmpl
```

### How It Works

1. **Controller:** Each controller of your application should have its own directory under `assets/templates/web/`. For example, templates for the `auth` controller should be in `assets/templates/web/auth/`.

2. **Handler:** Place the main templates for each handler in the respective controller directory. For instance, the main template for the `signin` handler of the `auth` controller should be `assets/templates/web/auth/signin.tmpl`.

3. **Partials:** Shared partial templates used across different handlers of a controller should be placed in the `partial` subdirectory.

4. **Layout:**
    - The general layout for the entire application is `assets/templates/web/layout/base_layout.tmpl`.
    - Optional controller-specific layouts can be placed in `assets/templates/web/{controller}/layout/layout.tmpl`.
    - Optional handler-specific layouts can be placed in `assets/templates/web/{controller}/layout/{handler}_layout.tmpl`.

### Layout Fallback Mechanism

The Template Manager uses a fallback mechanism for layouts:
- It first checks for a handler-specific layout (`{handler}_layout.tmpl`).
- If not found, it looks for a controller-specific layout (`layout.tmpl` within the controller's `layouts` directory).
- If neither is found, it falls back to the default application layout (`base.tmpl`).

### Usage

When rendering a template for a specific handler, the Template Manager automatically selects the appropriate layout based on the available templates and the fallback mechanism described above. This process is handled internally, ensuring that the correct combination of layout, main template, and partials is used.
