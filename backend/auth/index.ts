import { Hono } from 'hono';
import { cors } from 'hono/cors';
import dotenv from 'dotenv';
import { register, login } from './user_auth/auth.controller';
import { auth } from './auth'; 
dotenv.config();

const app = new Hono();
const port = Number(process.env.PORT);

// Fixed: changed "localhosht" to "localhost"
app.use(
  '/api/auth/*',
  cors({
    origin: 'http://localhost:5173', // Must be the exact frontend origin (cannot be '*')
    allowHeaders: ['Content-Type', 'Authorization'],
    allowMethods: ['POST', 'GET', 'OPTIONS'],
    exposeHeaders: ['Content-Length'],
    maxAge: 600,
    credentials: true, // 👈 Fixes the Access-Control-Allow-Credentials error
  })
)
// Routes
app.get("/", (c) => {
  return c.json({ message: "Hello World", ok: true });
});
app.on(['POST', 'GET'], '/api/auth/*', (c) => {
  return auth.handler(c.req.raw)
})
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
