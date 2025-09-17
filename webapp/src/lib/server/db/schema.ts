import { pgTable, serial, integer, varchar, decimal, timestamp, text, index, unique } from 'drizzle-orm/pg-core';

export const ucWages = pgTable('uc_wages', {
	id: serial('id').primaryKey(),
	location: varchar('location', { length: 50 }).notNull(),
	year: integer('year').notNull(),
	employeeId: integer('employee_id'),
	firstname: varchar('firstname', { length: 100 }),
	lastname: varchar('lastname', { length: 100 }),
	title: varchar('title', { length: 200 }),
	basepay: decimal('basepay', { precision: 12, scale: 2 }).default('0.00'),
	overtimepay: decimal('overtimepay', { precision: 12, scale: 2 }).default('0.00'),
	adjustpay: decimal('adjustpay', { precision: 12, scale: 2 }).default('0.00'),
	grosspay: decimal('grosspay', { precision: 12, scale: 2 }).default('0.00'),
	scrapedAt: timestamp('scraped_at'),
	uploadedAt: timestamp('uploaded_at').defaultNow()
}, (table) => ({
	uniqueLocationYearEmployee: unique().on(table.location, table.year, table.employeeId),
	locationIdx: index('idx_uc_wages_location').on(table.location),
	yearIdx: index('idx_uc_wages_year').on(table.year),
	locationYearIdx: index('idx_uc_wages_location_year').on(table.location, table.year),
	grosspayIdx: index('idx_uc_wages_grosspay').on(table.grosspay),
	titleIdx: index('idx_uc_wages_title').on(table.title)
}));

export const uploadProgress = pgTable('upload_progress', {
	id: serial('id').primaryKey(),
	location: varchar('location', { length: 50 }).notNull(),
	year: integer('year').notNull(),
	totalRecords: integer('total_records'),
	uploadedRecords: integer('uploaded_records').default(0),
	status: varchar('status', { length: 20 }).default('pending'), // pending, processing, completed, failed
	startedAt: timestamp('started_at').defaultNow(),
	completedAt: timestamp('completed_at'),
	errorMessage: text('error_message')
}, (table) => ({
	uniqueLocationYear: unique().on(table.location, table.year)
}));
