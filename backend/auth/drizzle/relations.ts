import { defineRelations } from "drizzle-orm";
import * as schema from "./schema";

export const relations = defineRelations(schema, (r) => ({
	account: {
		user: r.one.user({
			from: r.account.userId,
			to: r.user.id
		}),
	},
	user: {
		accounts: r.many.account(),
		employees: r.many.employee(),
		sessions: r.many.session(),
	},
	employee: {
		user: r.one.user({
			from: r.employee.userId,
			to: r.user.id
		}),
	},
	session: {
		user: r.one.user({
			from: r.session.userId,
			to: r.user.id
		}),
	},
}))