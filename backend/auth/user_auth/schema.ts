import { z } from 'zod';

// Used for Sign Up
export const registerSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  email: z.email({error:"Enter a valid email"}),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export const loginSchema = z.object({
  email: z.email({error:"Enter a valid email"}),
  password: z.string().min(1, "Password is required"),
});

// Type inference for TypeScript
export type RegisterInput = z.infer<typeof registerSchema>;
export type LoginInput = z.infer<typeof loginSchema>;