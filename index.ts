import { Hono } from "hono";
import dotenv from 'dotenv';
import {cors} from "hono/cors";
dotenv.config();
const app = new Hono().basePath("/api");
app.use(cors({
  origin:["*"]
}));
app.get("/health",(c) =>{
  return c.json({message:"Hello World"});
  
});




