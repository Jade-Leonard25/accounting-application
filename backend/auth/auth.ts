import { betterAuth } from 'better-auth';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { db } from './db';
import * as schema from './drizzle/schema';

const generateCustomId = (): string => {
  const rand1to9 = () => Math.floor(Math.random() * 9) + 1;
  const rand0to9 = () => Math.floor(Math.random() * 10);
  const randLetter = () => {
    const letters = "abcdefghijklmnopqrstuvwxyz";
    return letters[Math.floor(Math.random() * letters.length)];
  };

  // 2 digits + 1 letter for Block 1
  const block1 = `${rand0to9()}${rand0to9()}${randLetter()}`;
  // 2 digits + 1 letter for Block 2
  const block2 = `${rand0to9()}${rand0to9()}${randLetter()}`;
  // 2 digits + 1 letter for Block 3
  const block3 = `${rand0to9()}${rand0to9()}${randLetter()}`;

  return `${block1}-${block2}-${block3}`; // Result looks like: 49v-82m-31k
};
export const auth = betterAuth({
  baseURL:"http://localhost:4002",
trustedOrigins:["http://localhost:5173"],
  database: drizzleAdapter(db, {
    provider: 'sqlite', // 👈 Update provider to sqlite since you migrated to Turso/SQLite
    schema,
  }),
  user: {
    additionalFields: {
      uniqeIdentification: {
        type: "string",
        required: true,
        defaultValue: () => generateCustomId(), 
        input: true, 
      },
    },
  },
  emailAndPassword: {
    enabled: true,
  },
});
