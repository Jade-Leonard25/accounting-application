import { sqliteTable, index, foreignKey, primaryKey, unique, text, integer } from "drizzle-orm/sqlite-core"
import { sql } from "drizzle-orm"

export const account = sqliteTable("account", {
	id: text().primaryKey(),
	accountId: text("account_id").notNull(),
	providerId: text("provider_id").notNull(),
	userId: text("user_id").notNull().references(() => user.id, { onDelete: "cascade" } ),
	accessToken: text("access_token"),
	refreshToken: text("refresh_token"),
	idToken: text("id_token"),
	accessTokenExpiresAt: text("access_token_expires_at"),
	refreshTokenExpiresAt: text("refresh_token_expires_at"),
	scope: text(),
	password: text(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").notNull(),
},
(table) => [index("account_userId_idx").on(table.userId),
]);

export const employee = sqliteTable("employee", {
	id: text().primaryKey(),
	userId: text("user_id").references(() => user.id, { onDelete: "set null" } ),
	firstName: text("first_name").notNull(),
	lastName: text("last_name").notNull(),
	email: text().notNull(),
	taxId: text("tax_id"),
	sssNumber: text("sss_number"),
	philhealthNumber: text("philhealth_number"),
	pagibigNumber: text("pagibig_number"),
	department: text(),
	status: text().default("active").notNull(),
	hireDate: text("hire_date").notNull(),
	payFrequency: text("pay_frequency").default("semi_monthly").notNull(),
	baseSalaryCents: integer("base_salary_cents").notNull(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").default(sql`datetime('now')`).notNull(),
},
(table) => [index("employee_status_idx").on(table.status),
index("employee_userId_idx").on(table.userId),
unique("employee_email_key").on(table.email),
]);

export const session = sqliteTable("session", {
	id: text().primaryKey(),
	expiresAt: text("expires_at").notNull(),
	token: text().notNull(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").notNull(),
	ipAddress: text("ip_address"),
	userAgent: text("user_agent"),
	userId: text("user_id").notNull().references(() => user.id, { onDelete: "cascade" } ),
},
(table) => [index("session_userId_idx").on(table.userId),
unique("session_token_key").on(table.token),
]);

export const user = sqliteTable("user", {
	id: text().primaryKey(),
	name: text().notNull(),
	email: text().notNull(),
	uniqeIdentification: text("uniqe_identification").notNull(),
	emailVerified: integer("email_verified", {"mode":"boolean"}).default(false).notNull(),
	image: text(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").default(sql`datetime('now')`).notNull(),
},
(table) => [unique("user_email_key").on(table.email),
]);

export const verification = sqliteTable("verification", {
	id: text().primaryKey(),
	identifier: text().notNull(),
	value: text().notNull(),
	expiresAt: text("expires_at").notNull(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").default(sql`datetime('now')`).notNull(),
},
(table) => [index("verification_identifier_idx").on(table.identifier),
]);

