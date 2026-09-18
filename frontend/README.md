# Frontend - URL Shortener

React single-page application built with [Vite](https://vitejs.dev/) for the URL Shortener service.

---

## Project Layout

```
frontend/
├── .env.example          # Template for environment configuration
├── .gitignore            # Frontend-specific gitignore rules
├── README.md             # Frontend architecture and instructions
├── index.html            # Main HTML document template
├── package.json          # Node.js dependencies and build scripts
├── vite.config.js        # Vite build tool configuration
├── public/               # Static assets served as-is
└── src/
    ├── App.css           # Global application styles
    ├── App.jsx           # Root React component placeholder
    ├── index.css         # Baseline CSS resets
    ├── main.jsx          # React DOM entry point
    ├── assets/           # Project images, icons, and fonts
    ├── components/       # Reusable presentational & container components
    ├── context/          # React Context providers for global state
    ├── hooks/            # Custom React hooks
    ├── pages/            # Page-level route views
    ├── services/         # API clients and HTTP communication layer
    └── utils/            # Shared formatting and validation helpers
```

---

## Folder Purposes

- **`src/components/`**: Reusable modular UI elements (e.g., URL input bar, link cards, statistics widgets, copy buttons).
- **`src/pages/`**: High-level page components (e.g., Home/Dashboard page, Analytics page, 404 page).
- **`src/services/`**: API abstraction layer (e.g., Axios/Fetch wrappers) to communicate with the Go backend API.
- **`src/hooks/`**: Custom reusable React hooks (e.g., `useClipboard`, `useShortenerApi`).
- **`src/context/`**: React Context instances for cross-cutting state (e.g., notifications, user preferences).
- **`src/utils/`**: Helper utilities (e.g., URL validation, date formatting, string helpers).
- **`src/assets/`**: Static media assets such as SVG icons and branding images.

---

## Getting Started

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start development server:
   ```bash
   npm run dev
   ```

4. Build for production:
   ```bash
   npm run build
   ```
