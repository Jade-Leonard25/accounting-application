import { defineConfig } from 'drizzle-kit';
import dotenv from 'dotenv';
dotenv.config();

export default defineConfig({
  schema: './drizzle/schema.ts',
  out: './drizzle',
  dialect: 'turso',
  dbCredentials: {
    url: process.env.DATABASE_URL!,
    authToken: process.env.DATABASE_TOKEN!,
  },
  introspect: {
    casing: 'camel',
  },
});