/**
 * SPDX-License-Identifier: MIT
 *
 * Copyright (c) 2024 Berachain Foundation
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation
 * files (the "Software"), to deal in the Software without
 * restriction, including without limitation the rights to use,
 * copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following
 * conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

// Constants
// ========================================================
import { defineConfig } from "vocs";
import { theme } from "./vocs.theme.beaconkit";

// Vocs Config
// ========================================================
export default defineConfig({
  title: "BeaconKit Docs",
  theme,
  topNav: [
    {
      text: "Quick Start",
      link: "/what-is-beaconkit",
    },
    {
      text: "Concepts",
      link: "/concepts",
    },
    {
      text: "Guides",
      link: "/guides",
    },
    {
      text: "Changelog",
      link: "https://github.com/berachain/beacon-kit/blob/main/CHANGELOG",
    },
  ],
  sidebar: {
    "/": [
      {
        text: "What Is BeaconKit?",
        link: "/what-is-beaconkit",
      },
      {
        text: "Quick Start: Run A Node",
        link: "/quick-start",
      },
      {
        text: "Installation Guides",
        items: [
          {
            text: "Docker",
            link: "/installation-guides/docker",
          },
          {
            text: "Build From Source",
            link: "/installation-guides/build-from-source",
          },
        ],
      },
      {
        text: "Node Configurations",
        items: [
          {
            text: "Validator",
            link: "/node-configurations/validator",
          },
          {
            text: "RPC",
            link: "/node-configurations/rpc",
          },
          {
            text: "Archive Node",
            link: "/node-configurations/archive-node",
          },
          {
            text: "Snapshot",
            link: "/node-configurations/snapshot",
          },
        ],
      },
      {
        text: "CLI",
        link: "/cli",
      },
      {
        text: "API Reference",
        link: "/api",
      },
      {
        text: "FAQ",
        link: "/faq",
      },
      {
        text: "Glossary",
        link: "/glossary",
      },
    ],
    "/concepts": [
      {
        text: "Nodes",
        items: [
          {
            text: "Nodes & Network",
            link: "/concepts/nodes-network",
          },
          {
            text: "Validator Lifecycle",
            link: "/concepts/validator-lifecycle",
          },
          {
            text: "Key, Wallets, & Accounts",
            link: "/concepts/key-wallets-accounts",
          },
          {
            text: "Slashing",
            link: "/concepts/slashing",
          },
        ],
      },
      {
        text: "Protocol",
        items: [
          {
            text: "Proof Of Stake",
            link: "/concepts/proof-of-stake",
          },
          {
            text: "Rewards",
            link: "/concepts/rewards",
          },
          {
            text: "Finality / Settlement",
            link: "/concepts/finality-settlement",
          },
          {
            text: "Engine API",
            link: "/concepts/engine-api",
          },
          {
            text: "Blobs",
            link: "/concepts/blobs",
          },
        ],
      },
    ],
    "/guides": [
      {
        text: "Voting",
        link: "/guides/voting",
      },
      {
        text: "Upgrading",
        link: "/guides/upgrading",
      },
      {
        text: "Unjailing",
        link: "/guides/unjailing",
      },
      {
        text: "Resyncing",
        link: "/guides/resyncing",
      },
      {
        text: "Snapshots",
        link: "/guides/snapshots",
      },
      {
        text: "Kurtosis",
        link: "/guides/kurtosis",
      },
      {
        text: "Custom Rewards",
        link: "/guides/custom-rewards",
      },
    ],
  },
  socials: [
    {
      icon: "github",
      link: "https://github.com/berachain/beacon-kit",
    },
    {
      icon: "discord",
      link: "https://discord.com/invite/berachain",
    },
  ],
  editLink: {
    pattern:
      "https://github.com/berachain/beacon-kit/tree/main/docs/docs/pages/:path",
    text: "Edit on GitHub",
  },
  sponsors: [
    {
      name: "Collaborator",
      height: 120,
      items: [
        [
          {
            name: "Paradigm",
            link: "https://paradigm.xyz",
            image:
              "https://raw.githubusercontent.com/wevm/.github/main/content/sponsors/paradigm-light.svg",
          },
        ],
      ],
    },
    {
      name: "Large Enterprise",
      height: 60,
      items: [
        [
          {
            name: "WalletConnect",
            link: "https://walletconnect.com",
            image:
              "https://raw.githubusercontent.com/wevm/.github/main/content/sponsors/walletconnect-light.svg",
          },
          {
            name: "Stripe",
            link: "https://www.stripe.com",
            image:
              "https://raw.githubusercontent.com/wevm/.github/main/content/sponsors/stripe-light.svg",
          },
        ],
      ],
    },
  ],
  // sidebar: [
  //   {
  //     text: "Introduction to Beacon Kit",
  //     link: "/getting-started",
  //   },
  //   {
  //     text: "APIs",
  //     collapsed: false,
  //     items: [
  //       {
  //         text: "Pets",
  //         collapsed: false,
  //         link: "/pets",
  //         items: [
  //           {
  //             text: "get",
  //             link: "/pets/get",
  //           },
  //           {
  //             text: "post",
  //             link: "/pets/post",
  //           },
  //         ],
  //       },
  //       {
  //         text: "Streams",
  //         collapsed: false,
  //         items: [
  //           {
  //             text: "post",
  //             link: "/streams/post",
  //           },
  //         ],
  //       },
  //     ],
  //   },
  // ],

});
