-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE `account` (
	`id` text,
	`account_id` text NOT NULL,
	`provider_id` text NOT NULL,
	`user_id` text NOT NULL,
	`access_token` text,
	`refresh_token` text,
	`id_token` text,
	`access_token_expires_at` text,
	`refresh_token_expires_at` text,
	`scope` text,
	`password` text,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text NOT NULL,
	CONSTRAINT `account_pk` PRIMARY KEY(`id`),
	CONSTRAINT `fk_account_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
);
--> statement-breakpoint
CREATE TABLE `employee` (
	`id` text,
	`user_id` text,
	`first_name` text NOT NULL,
	`last_name` text NOT NULL,
	`email` text NOT NULL CONSTRAINT `employee_email_key` UNIQUE,
	`tax_id` text,
	`sss_number` text,
	`philhealth_number` text,
	`pagibig_number` text,
	`department` text,
	`status` text DEFAULT 'active' NOT NULL,
	`hire_date` text NOT NULL,
	`pay_frequency` text DEFAULT 'semi_monthly' NOT NULL,
	`base_salary_cents` integer NOT NULL,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text DEFAULT datetime('now') NOT NULL,
	CONSTRAINT `employee_pk` PRIMARY KEY(`id`),
	CONSTRAINT `fk_employee_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE SET NULL
);
--> statement-breakpoint
CREATE TABLE `session` (
	`id` text,
	`expires_at` text NOT NULL,
	`token` text NOT NULL CONSTRAINT `session_token_key` UNIQUE,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text NOT NULL,
	`ip_address` text,
	`user_agent` text,
	`user_id` text NOT NULL,
	CONSTRAINT `session_pk` PRIMARY KEY(`id`),
	CONSTRAINT `fk_session_user_id_user_id_fk` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE
);
--> statement-breakpoint
CREATE TABLE `user` (
	`id` text,
	`name` text NOT NULL,
	`email` text NOT NULL CONSTRAINT `user_email_key` UNIQUE,
	`email_verified` integer DEFAULT false NOT NULL,
	`image` text,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text DEFAULT datetime('now') NOT NULL,
	CONSTRAINT `user_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `verification` (
	`id` text,
	`identifier` text NOT NULL,
	`value` text NOT NULL,
	`expires_at` text NOT NULL,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text DEFAULT datetime('now') NOT NULL,
	CONSTRAINT `verification_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE INDEX `account_userId_idx` ON `account` (`user_id`);--> statement-breakpoint
CREATE INDEX `employee_status_idx` ON `employee` (`status`);--> statement-breakpoint
CREATE INDEX `employee_userId_idx` ON `employee` (`user_id`);--> statement-breakpoint
CREATE INDEX `session_userId_idx` ON `session` (`user_id`);--> statement-breakpoint
CREATE INDEX `verification_identifier_idx` ON `verification` (`identifier`);
*/