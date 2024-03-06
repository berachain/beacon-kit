// Constants
// ========================================================
import { defineConfig, type SidebarItem } from "vocs";
import { theme } from "./vocs.theme.beaconkit";

// Vocs Config
// ========================================================
export default defineConfig({
  title: "BeaconKit Docs",
  theme,
  topNav: [
    {
      text: 'Quick Start',
      link: '/'
    },
    {
      text: 'Concepts',
      link: '/concepts'
    },
    {
      text: 'Guides',
      link: '/guides'
    },
    {
      text: 'Changelog',
      link: 'https://github.com/berachain/beacon-kit/blob/main/CHANGELOG'
    }
  ],
  sidebar: {
    '/': [
      {
        text: 'What Is BeaconKit?',
        link: '/'
      },
      {
        text: 'Quick Start: Run A Node',
        link: '/quick-start'
      },
      {
        text: 'Installation Guides',
        items: [
          {
            text: 'Docker', 
            link: '/installation-guides/docker'
          },
          {
            text: 'Build From Source',
            link: '/installation-guides/build-from-source'
          }
        ]
      },
      {
        text: 'Node Configurations',
        items: [
          {
            text: 'Validator', 
            link: '/node-configurations/validator'
          },
          {
            text: 'RPC',
            link: '/node-configurations/rpc'
          },
          {
            text: 'Archive Node',
            link: '/node-configurations/archive-node'
          },
          {
            text: 'Snapshot',
            link: '/node-configurations/snapshot'
          }
        ]
      },
      {
        text: 'CLI',
        link: '/cli'
      },
      {
        text: 'API Reference',
        link: '/api'
      },
      {
        text: 'FAQ',
        link: '/faq'
      },
      {
        text: 'Glossary',
        link: '/glossary'
      },
    ],
    '/concepts': [
      {
        text: "Nodes",
        items: [
          {
            text: 'Nodes & Network',
            link: '/concepts/nodes-network'
          },
          {
            text: 'Validator Lifecycle',
            link: '/concepts/validator-lifecycle'
          },
          {
            text: 'Key, Wallets, & Accounts',
            link: '/concepts/key-wallets-accounts'
          },
          {
            text: 'Slashing',
            link: '/concepts/slashing'
          },
        ]
      },
      {
        text: 'Protocol',
        items: [
          {
            text: 'Proof Of Stake',
            link: '/concepts/proof-of-stake'
          },
          {
            text: 'Rewards',
            link: '/concepts/rewards'
          },
          {
            text: 'Finality / Settlement',
            link: '/concepts/finality-settlement'
          },
          {
            text: 'Engine API',
            link: '/concepts/engine-api'
          },
          {
            text: 'Blobs',
            link: '/concepts/blobs'
          },
        ]
      }
    ],
    '/guides': [
      {
        text: 'Voting',
        link: '/guides/voting'
      },
      {
        text: 'Upgrading',
        link: '/guides/upgrading'
      },
      {
        text: 'Unjailing',
        link: '/guides/unjailing'
      },
      {
        text: 'Resyncing',
        link: '/guides/resyncing'
      },
      {
        text: 'Snapshots',
        link: '/guides/snapshots'
      },
      {
        text: 'Kurtosis',
        link: '/guides/kurtosis'
      },
      {
        text: 'Custom Rewards',
        link: '/guides/custom-rewards'
      },
    ]
  },
  socials: [
    { 
      icon: 'github', 
      link: 'https://github.com/berachain/beacon-kit', 
    }, 
    { 
      icon: 'discord', 
      link: 'https://discord.com/invite/berachain', 
    }, 
  ],
  editLink: {
    pattern: 'https://github.com/berachain/beacon-kit/tree/main/docs/docs/pages/:path',
    text: 'Edit on GitHub'
  }
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
