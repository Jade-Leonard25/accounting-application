import { Hono } from 'hono';
import { cors } from 'hono/cors';
import dotenv from 'dotenv';
import { register, login } from './user_auth/auth.controller';

dotenv.config();

const app = new Hono();
const port = Number(process.env.PORT)

app.use(cors({ origin: ["*"] }));

// Routes
app.get("/", (c) => {
  return c.json({ message: "Hello World", ok: true });
});

app.post("/api/auth/register", register);
app.post("/api/auth/login", login);

// Print active endpoints on startup
console.log(`🚀 Hono Auth Server running on port ${port}`);
console.log(`-----------------------------------------`);
console.log(`Available Endpoints:`);
console.log(`  GET  http://localhost:${port}/`);
console.log(`  POST http://localhost:${port}/api/auth/register`);
console.log(`  POST http://localhost:${port}/api/auth/login`);
console.log(`-----------------------------------------`);

export default { port, fetch: app.fetch };