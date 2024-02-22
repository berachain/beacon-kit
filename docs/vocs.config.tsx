import { defineConfig } from "vocs";

export default defineConfig({
  title: "BeaconKit Docs",
  sidebar: [
    {
      text: "Introduction to Beacon Kit",
      link: "/getting-started",
    },
    {
      text: "APIs",
      collapsed: false,
      items: [
        {
          text: "Pets",
          collapsed: false,
          link: "/pets",
          items: [
            {
              text: "get",
              link: "/pets/get",
            },
            {
              text: "post",
              link: "/pets/post",
            },
          ],
        },
        {
          text: "Streams",
          collapsed: false,
          items: [
            {
              text: "post",
              link: "/streams/post",
            },
          ],
        },
      ],
    },
  ],
});
