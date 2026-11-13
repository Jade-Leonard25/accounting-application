import { z } from 'zod';

export const createEmployeeSchema = z.object({
  firstName: z.string().min(1),
  lastName: z.string().min(1),
  email: z.email().min(1),
  department: z.string().optional(),
  hireDate: z.string(), // ISO date string
  baseSalaryCents: z.number().int().positive(),
    
});