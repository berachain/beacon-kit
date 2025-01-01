<br>

<p align="center">
  <a href="https://wagmi.sh">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://res.cloudinary.com/duv0g402y/image/upload/v1713381289/monobera_color_alt_fgny7b.svg">
      <img alt="wagmi logo" src="https://res.cloudinary.com/duv0g402y/image/upload/v1713381289/monobera_color_alt2_ppo8o6.svg" width="auto" height="200">
    </picture>
  </a>
</p>
<p align="center">
    Monorepo for maintaining Berachain Applications & Libraries
<p>

![CI](https://github.com/berachain/monobera/actions/workflows/quality.yml/badge.svg?branch=v2)

## Installation

In order to setup your local environment, run

```
git submodule init
git submodule update
pnpm i
pnpm setenv bartio
```

You'll also need to run

```
cp .env.local.example .env.local
```

After that, create your own `NEXT_PUBLIC_DYNAMIC_API_KEY` without any CORS restrictions to write into `.env.local`.


All Berachain dapps are built to be single chain applications.

| Environment Variables | Environment |
| --------------------- | ---------------------------------------------------------------------------------------------- |
| `.env.bartio` | Environment variables for Berachain bartio testnet |

## Development Tips

Sometimes your VSCode will not correctly pick up type information from `/packages`.  

```bash
pnpm build:pkg # This will build all packages
pnpm format    # This will format all packages
```

Then in VSCode press `CMD+Shift+P` to open the command pallette and type `>` then type `Restart TS Server` and press enter. This will 
restart the typescript server and should pick up the new types.

## Commands

Monobera requires node 22.12+.

| Script                   | Description                                                                                              |
| ------------------------ | -------------------------------------------------------------------------------------------------------- |
| `pnpm i`                 | Installs packages for all apps & packages                                                                |
| `pnpm build`             | Builds all packages and apps. Not recommended as it takes large amounts of memory                        |
| `pnpm setenv <environment>`     | Copies `.env.<environment>` into `.env`. Don't do it manually.                                                  |
| `pnpm build:hub`         | Builds only the `Hub` and related packages.                                                              |
| `pnpm build:honey`       | Builds only the `Honey` and related packages.                                                            |
| `pnpm build:lend`        | Builds only the `Bend` and related packages.                                                             |
| `pnpm build:perp`        | Builds only the `Berps` and related packages.                                                            |
| `pnpm build:berajs-docs` | Builds only the `Berajs Docs` and related packages.                                                      |
| `pnpm build:pkg`         | Builds all packages.                                                                                     |
| `pnpm dev:pkg`           | Runs all packages in dev mode.                                                                           |
| `pnpm clean`             | Cleans the project using turbo clean and removes untracked files with git clean, including node_modules. |
| `pnpm pullenv`           | Pulls production environment variables from Vercel. Requires Vercel Login                                |
| `pnpm check-types`       | Runs type-checking across all apps and packages.                                                         |
| `pnpm lint`              | Lints all apps and packages.                                                                             |
| `pnpm format:check`      | Checks the formatting of all apps and packages without making changes.                                   |
| `pnpm format`            | Formats the apps and packages and writes the changes.                                                    |
| `pnpm check`             | Performs a comprehensive check of all apps and packages, including linting and type-checking.            |
| `pnpm prepare`           | Installs Husky, setting up Git hooks for the project.                                                    |
| `pnpm upsertenv`         | Runs a script to upsert environment variables in Vercel for the project.                                 |
| `pnpm knip`              | Executes the knip command to exclude binaries from operations.                                           |

To run Hub for example, run `pnpm i && pnpm dev:hub`

## Apps

| App                  | Description                                    |
| -------------------- | ---------------------------------------------- |
| `app/hub`            | `Hub` application code                         |
| `app/honey`          | `Honey` application code                       |
| `app/lend`           | `Bend` application code                        |
| `app/perp`           | `Berps` application code                       |

## Packages

| Package                 | Description                                                                                               |
| ----------------------- | --------------------------------------------------------------------------------------------------------- |
| `packages/berajs`       | A Typescript package for interacting with Berachain. [View Docs](https://berajsdocs.vercel.app/)          |
| `packages/wagmi`        | A package to create a shared wagmi / dynamic config for web3 applications                                 |
| `packages/config`       | A package to store shared config variables across applications                                            |
| `packages/graphql`      | A package to store appolo clients / gql subgraph queries                                                  |
| `packages/proto`        | A package to generate e2e typing & protobuf for interacting with Cosmos-SDK                               |
| `packages/shared-ui`    | A package of built UI widgets made from `packages/ui` component                                           |
| `packages/ui`           | A package of [shadcn](https://ui.shadcn.com/) components                                                  |



## Tooling & Libraries

A short list of tooling and libraries we use across all apps and packages.

- [biomejs](https://biomejs.dev/)
- [knip](https://knip.dev/)
- [turbo](https://turbo.build/)
- [next](https://nextjs.org/)
- [wagmi](https://wagmi.sh/)
- [viem](https://viem.sh/)
- [swr](https://swr.vercel.app/)
- [vocs](https://vocs.dev/)
- [shadcn](https://ui.shadcn.com/)
- [tailwind](https://tailwindcss.com/)

## CI

We have set up a caching mechanism for turbo builds to speed up CI times. This is done by using an open source github action called [rharkor/caching-for-turbo@v1.5](https://github.com/rharkor/caching-for-turbo).

## Dapps banner management

Banners serve as an essential tool for communicating urgent messages or event-related information to users across the site. The management of these banners is centralized in the Bannerconfig component, located within the `packages/shared-ui` directory. This allows for effective global notification during scenarios like RPC issues or network congestion.
For targeted communications, banners can be configured to appear on specific pages by listing the desired paths in the hrefs field. For example, to display a banner only on the "Pools," "Swap," and homepage in BEX, you would set `hrefs` to `["/pools", "/swap", "/"]`.
To modify the banner configuration, submit a PR with changes to the `enabled` field in the `bannerConfig`. This will update the banner's active status and display it as specified.
