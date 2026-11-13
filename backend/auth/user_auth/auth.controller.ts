import { Context } from 'hono';
import { auth } from '../auth'; // Import your Better Auth instance
import type { LoginInput, RegisterInput } from './schema';

// Helper for your custom ID format: 000-000-000
const generateCustomId = (): string => {
  const rand1to9 = () => Math.floor(Math.random() * 9) + 1;
  const rand0to9 = () => Math.floor(Math.random() * 10);
  const randLetter = () => {
    const letters = "abcdefghijklmnopqrstuvwxyz";
    return letters[Math.floor(Math.random() * letters.length)];
  };

  const block1 = `${rand1to9()}${rand1to9()}${rand1to9()}`;
  const block2 = `${rand0to9()}${rand0to9()}${rand0to9()}`;
  const block3 = `${rand0to9()}${rand0to9()}${randLetter()}`;

  return `${block1}-${block2}-${block3}`;
};

export const register = async (c: Context) => {
  // Properly type the validated JSON body
  const { name, email, password } = await c.req.json<RegisterInput>();

  // Generate your custom multi-block ID
  const customId = generateCustomId();

  try {
    // 1. Let Better Auth handle the user creation, passing your custom id
    const result = await auth.api.signUpEmail({
      body: {
        uniqe_identification: customId, // 👈 Passing the custom generated ID to Better Auth / database
        name,
        email,
        password,
      },
    });

    return c.json({ 
      message: "User registered successfully", 
      ok: true, 
      user: result.user 
    }, 201);
  } catch (error: any) {
    return c.json({ ok: false, error: error.message || "Registration failed" }, 400);
  }
};

export const login = async (c: Context) => {
  const { email, password } = await c.req.json<LoginInput>();

  try {
    // 1. Let Better Auth verify the password and create the session
    const session = await auth.api.signInEmail({
      body: {
        email,
        password,
      },
    });

    // 2. Return the session/token to the client
    return c.json({ 
      message: "Logged in successfully", 
      ok: true, 
      session: session // Contains the session token and user data
    });
  } catch (error: any) {
    return c.json({ ok: false, error: "Invalid email or password" }, 401);
  }
};