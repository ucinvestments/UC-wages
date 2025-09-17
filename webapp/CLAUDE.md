# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

UC Wages Webapp - A SvelteKit web application for visualizing and analyzing University of California employee wage data. Features data ingestion, interactive D3.js visualizations, and comprehensive wage analytics. Built with modern web technologies including SvelteKit, TypeScript, TailwindCSS, and Drizzle ORM with PostgreSQL database.

## Architecture

### Tech Stack
- **Frontend**: SvelteKit 2 with Svelte 5, TypeScript
- **Styling**: TailwindCSS 4 with custom components
- **Database**: PostgreSQL with Drizzle ORM
- **Deployment**: Vercel (configured with @sveltejs/adapter-vercel)
- **Build Tool**: Vite 7

### Project Structure
```
src/
├── app.html                 # Main HTML template
├── app.css                  # Global styles import
├── app.d.ts                 # TypeScript app declarations
├── routes/                  # SvelteKit routes
│   ├── +layout.svelte      # Root layout component
│   └── +page.svelte        # Homepage
├── lib/
│   ├── components/         # Reusable Svelte components
│   ├── layouts/           # Layout components
│   ├── assets/            # Static assets (favicon, etc.)
│   └── server/
│       └── db/            # Database configuration and schema
│           ├── index.ts   # Drizzle database connection
│           └── schema.ts  # Database schema definitions
└── styles/                # CSS modules and utilities
```

### Database Setup
- Uses Neon serverless PostgreSQL
- Configured with Drizzle ORM for type-safe database operations
- Schema defined in `src/lib/server/db/schema.ts`
- Two main tables: `uc_wages` (wage records) and `upload_progress` (upload tracking)
- Environment variable `DATABASE_URL` required for database connection

### API Routes
- `GET/POST /api/wages` - Query and upload wage data
- `GET/POST /api/upload` - File upload and progress tracking
- Supports aggregated data queries for visualization
- Batch processing with upsert logic for data integrity

## Commands

### Development
```bash
npm run dev                 # Start development server
npm run dev -- --open      # Start dev server and open browser
npm run build              # Build for production
npm run preview            # Preview production build
```

### Code Quality
```bash
npm run check              # Run svelte-check for TypeScript validation
npm run check:watch        # Run svelte-check in watch mode
npm run format             # Format code with Prettier
npm run lint               # Check code formatting with Prettier
```

### Database Operations
```bash
npm run db:generate        # Generate database migrations
npm run db:push            # Push schema changes to database
npm run db:migrate         # Run database migrations
npm run db:studio          # Open Drizzle Studio for database management
```

## Environment Configuration

Required environment variables (see `.env.example`):
- `DATABASE_URL`: PostgreSQL connection string for Neon database

## Key Development Notes

- Uses Svelte 5 with the new runes syntax (`$props()`, `$state()`, etc.)
- TailwindCSS 4 configured with Vite plugin for styling
- Database operations should use the configured Drizzle instance from `$lib/server/db`
- Components follow SvelteKit conventions with TypeScript
- Custom font (AldotheApache.ttf) included in src directory