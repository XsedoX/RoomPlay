// Aspire TypeScript AppHost
// For more information, see: https://aspire.dev

import { createBuilder } from "./.aspire/modules/aspire.mjs";

const builder = await createBuilder();

const postgres = await builder
  .addPostgres("database")
  .withEnvironment("postgres", "postgres")
  .withHostPort({ port: 5432 })
  .withDataVolume()
  .addDatabase("postgres");

var postgresConnectionString = await postgres.uriExpression();
var api = await builder
  .addGoApp("api", "../backend")
  .withEnvironment("PORT", "7654")
  .withEnvironment("DATABASE_CONNECTION_STRING", postgresConnectionString)
  .waitFor(postgres);

await builder
  .addViteApp("frontend", "../frontend/")
  .withBun()
  .withHttpEndpoint({ port: 5173 })
  .publishAsStaticWebsite({ apiPath: "/api/v1", apiTarget: api })
  .waitFor(api);

await builder.build().run();
