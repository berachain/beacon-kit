{
  "name": "monobera",
  "private": true,
  "engines": {
    "node": ">=v22.0.0",
    "pnpm": ">=9"
  },
  "pnpm": {
    "overrides": {
      "react": "<19",
      "react-dom": "<19",
      "@types/react": "<19"
    }
  },
  "packageManager": "pnpm@9.1.4",
  "scripts": {
    "build": "turbo build",
    "clean": "turbo clean && git clean -xdf node_modules",
    "pullenv": "vercel env pull --environment=production",
    "setenv": "node ./.scripts/set-env.js",
    "setenv:bartio": "pnpm run setenv bartio",
    "build:ipfs": "NEXT_PUBLIC_HOST=ipfs turbo build --filter=hub --filter=honey --filter='./packages/*' --concurrency=20",
    "build:hub": "turbo build --filter=hub --filter='./packages/*' --concurrency=20",
    "build:honey": "turbo build --filter=honey --filter='./packages/*' --concurrency=20",
    "build:storybook": "turbo build --filter=storybook --filter='./packages/*' --concurrency=20  --filter='!b-sdk'",
    "build:lend": "turbo build --filter=lend --filter='./packages/*' --concurrency=20 --filter='!b-sdk'",
    "build:perp": "turbo build --filter=perp --filter='./packages/*' --concurrency=20 --filter='!b-sdk'",
    "build:berajs-docs": "turbo build --filter=berajs-docs --filter='./packages/*' --concurrency=20",
    "build:pkg": "turbo build --filter='./packages/*' --concurrency=20",
    "dev:pkg": "turbo dev --filter='./packages/*' --concurrency=20",
    "check-types": "turbo check-types",
    "lint": "biome lint .",
    "format:check": "biome format .",
    "format": "biome format --write .",
    "postformat": "cd ./packages/b-sdk && pnpm format",
    "check": "biome check .",
    "prepare": "husky install",
    "upsertenv": "node ./scripts/upsertVercelEnv.js",
    "knip": "knip --exclude binaries"
  },
  "dependencies": {
    "@ianvs/prettier-plugin-sort-imports": "^3.7.2",
    "@types/prettier": "^2.7.2",
    "commander": "^12.1.0",
    "eslint": "^8.39.0",
    "husky": "^8.0.0",
    "prettier": "^2.8.8",
    "prettier-plugin-tailwindcss": "^0.2.8",
    "turbo": "^1.10.15",
    "typescript": "^5.0.4",
    "vocs": "latest"
  },
  "devDependencies": {
    "@biomejs/biome": "1.5.3",
    "jscpd": "^3.5.10",
    "knip": "^5.7.0",
    "next": "^14.2.11",
    "yargs": "^17.7.2"
  },
  "resolutions": {
    "rpc-websockets": "^7.11.1"
  }
}
