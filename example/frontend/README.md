# sv

Everything you need to build a Svelte project, powered by [`sv`](https://github.com/sveltejs/cli).

## Creating a project

If you're seeing this, you've probably already done this step. Congrats!

```sh
# create a new project
pnpm dlx sv create my-app
```

To recreate this project with the same configuration:

```sh
# recreate this project
pnpm dlx sv@0.15.3 create --template minimal --types ts --add tailwindcss="plugins:none" --install pnpm frontend
```

## Developing

## Environment

The frontend uses SvelteKit and Vite's native `.env` loading. Copy the example file before running the app:

```sh
cp .env.example .env
```

Adjust `SYNCRA_API_BASE_URL` if your local API is not running on `http://localhost:8080`. Keep these values server-only unless they intentionally use SvelteKit's `PUBLIC_` prefix.

Once you've created a project, install dependencies:

```sh
pnpm install
```

Then start a development server:

```sh
pnpm dev

# or start the server and open the app in a new browser tab
pnpm dev -- --open
```

## Building

To create a production version of your app:

```sh
pnpm build
```

You can preview the production build with:

```sh
pnpm preview
```

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
