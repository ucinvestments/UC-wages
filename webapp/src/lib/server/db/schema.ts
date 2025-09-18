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

// Precalculated summary statistics
export const wageSummaries = pgTable('wage_summaries', {
	id: serial('id').primaryKey(),
	location: varchar('location', { length: 50 }).notNull(),
	year: integer('year').notNull(),
	employeeCount: integer('employee_count').notNull(),
	totalGrossPay: decimal('total_gross_pay', { precision: 15, scale: 2 }).notNull(),
	avgGrossPay: decimal('avg_gross_pay', { precision: 12, scale: 2 }).notNull(),
	medianPay: decimal('median_pay', { precision: 12, scale: 2 }).notNull(),
	stdDev: decimal('std_dev', { precision: 12, scale: 2 }).notNull(),
	minPay: decimal('min_pay', { precision: 12, scale: 2 }).notNull(),
	maxPay: decimal('max_pay', { precision: 12, scale: 2 }).notNull(),
	// Percentiles stored as JSON
	percentiles: text('percentiles').notNull(), // JSON string
	// Pay components
	totalBase: decimal('total_base', { precision: 15, scale: 2 }).notNull(),
	totalOvertime: decimal('total_overtime', { precision: 15, scale: 2 }).notNull(),
	totalAdjustments: decimal('total_adjustments', { precision: 15, scale: 2 }).notNull(),
	avgBase: decimal('avg_base', { precision: 12, scale: 2 }).notNull(),
	avgOvertime: decimal('avg_overtime', { precision: 12, scale: 2 }).notNull(),
	avgAdjustments: decimal('avg_adjustments', { precision: 12, scale: 2 }).notNull(),
	generatedAt: timestamp('generated_at').notNull(),
	uploadedAt: timestamp('uploaded_at').defaultNow()
}, (table) => ({
	uniqueLocationYear: unique().on(table.location, table.year),
	locationIdx: index('idx_wage_summaries_location').on(table.location),
	yearIdx: index('idx_wage_summaries_year').on(table.year),
	locationYearIdx: index('idx_wage_summaries_location_year').on(table.location, table.year)
}));

// Wage pyramid brackets
export const wagePyramids = pgTable('wage_pyramids', {
	id: serial('id').primaryKey(),
	location: varchar('location', { length: 50 }).notNull(),
	year: integer('year').notNull(),
	totalEmployees: integer('total_employees').notNull(),
	totalPay: decimal('total_pay', { precision: 15, scale: 2 }).notNull(),
	// Store brackets as JSON array
	brackets: text('brackets').notNull(), // JSON string
	generatedAt: timestamp('generated_at').notNull(),
	uploadedAt: timestamp('uploaded_at').defaultNow()
}, (table) => ({
	uniqueLocationYear: unique().on(table.location, table.year),
	locationIdx: index('idx_wage_pyramids_location').on(table.location),
	yearIdx: index('idx_wage_pyramids_year').on(table.year),
	locationYearIdx: index('idx_wage_pyramids_location_year').on(table.location, table.year)
}));

// Job title analysis
export const titleAnalysis = pgTable('title_analysis', {
	id: serial('id').primaryKey(),
	location: varchar('location', { length: 50 }).notNull(),
	year: integer('year').notNull(),
	uniqueTitles: integer('unique_titles').notNull(),
	// Store top titles as JSON array
	topTitles: text('top_titles').notNull(), // JSON string
	generatedAt: timestamp('generated_at').notNull(),
	uploadedAt: timestamp('uploaded_at').defaultNow()
}, (table) => ({
	uniqueLocationYear: unique().on(table.location, table.year),
	locationIdx: index('idx_title_analysis_location').on(table.location),
	yearIdx: index('idx_title_analysis_year').on(table.year),
	locationYearIdx: index('idx_title_analysis_location_year').on(table.location, table.year)
}));
