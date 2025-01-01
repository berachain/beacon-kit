import type { KnipConfig } from "knip";

const config: KnipConfig = {
  ignore: [
    "packages/proto/**",
    "apps/perp/public/static/**",
    "**/tsup.config.ts",
  ],
  project: ["apps/dex/**", "apps/honey/**", "apps/lend/**", "apps/hub/**"],
};

export default config;
