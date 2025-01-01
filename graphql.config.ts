import dotenv from "dotenv";

import { getEndpointsMap } from "./packages/graphql/codegen";

dotenv.config({
  path: ["./.env.local", "./.env"],
});

module.exports = {
  projects: getEndpointsMap().reduce((acc, [entry, url]) => {
    const [dir, file] = entry.split("/");
    const documents = [
      `./packages/graphql/src/modules/${dir}/${file ?? dir}.graphql`,
    ];
    if (entry === "dex/api") {
      documents.push("./packages/b-sdk/**/*.ts");
    }
    acc[entry] = {
      schema: url,
      documents: `./packages/graphql/src/modules/${dir}/${file ?? dir}.graphql`,
    };
    return acc;
  }, {}),
};
