import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { ucWages } from '$lib/server/db/schema';
import { eq, and, sql, desc, asc } from 'drizzle-orm';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url }) => {
	try {
		const location = url.searchParams.get('location');
		const year = url.searchParams.get('year');
		const limit = parseInt(url.searchParams.get('limit') || '1000');
		const offset = parseInt(url.searchParams.get('offset') || '0');
		const aggregated = url.searchParams.get('aggregated') === 'true';

		if (aggregated) {
			// Return aggregated data for visualization
			const query = db
				.select({
					location: ucWages.location,
					year: ucWages.year,
					totalWages: sql<number>`SUM(${ucWages.grosspay})::numeric`,
					averageWage: sql<number>`AVG(${ucWages.grosspay})::numeric`,
					employeeCount: sql<number>`COUNT(*)::integer`,
					maxWage: sql<number>`MAX(${ucWages.grosspay})::numeric`,
					minWage: sql<number>`MIN(${ucWages.grosspay})::numeric`
				})
				.from(ucWages)
				.groupBy(ucWages.location, ucWages.year)
				.orderBy(asc(ucWages.location), asc(ucWages.year));

			// Apply filters if provided
			let conditions = [];
			if (location) {
				conditions.push(eq(ucWages.location, location));
			}
			if (year) {
				conditions.push(eq(ucWages.year, parseInt(year)));
			}

			if (conditions.length > 0) {
				query.where(and(...conditions));
			}

			const result = await query;

			// Convert string decimals to numbers
			const processedResult = result.map(row => ({
				...row,
				totalWages: parseFloat(row.totalWages.toString()),
				averageWage: parseFloat(row.averageWage.toString()),
				maxWage: parseFloat(row.maxWage.toString()),
				minWage: parseFloat(row.minWage.toString())
			}));

			return json(processedResult);
		} else {
			// Return individual wage records
			let query = db
				.select()
				.from(ucWages)
				.limit(limit)
				.offset(offset)
				.orderBy(desc(ucWages.grosspay));

			// Apply filters if provided
			let conditions = [];
			if (location) {
				conditions.push(eq(ucWages.location, location));
			}
			if (year) {
				conditions.push(eq(ucWages.year, parseInt(year)));
			}

			if (conditions.length > 0) {
				query = query.where(and(...conditions));
			}

			const result = await query;

			// Convert decimal strings to numbers for frontend
			const processedResult = result.map(record => ({
				...record,
				basepay: parseFloat(record.basepay?.toString() || '0'),
				overtimepay: parseFloat(record.overtimepay?.toString() || '0'),
				adjustpay: parseFloat(record.adjustpay?.toString() || '0'),
				grosspay: parseFloat(record.grosspay?.toString() || '0')
			}));

			return json(processedResult);
		}
	} catch (error) {
		console.error('Error fetching wage data:', error);
		return json({ error: 'Failed to fetch wage data' }, { status: 500 });
	}
};

export const POST: RequestHandler = async ({ request }) => {
	try {
		const { records, location, year } = await request.json();

		if (!records || !Array.isArray(records)) {
			return json({ error: 'Records array is required' }, { status: 400 });
		}

		// Convert records to the correct format
		const wageRecords = records.map((record: any) => ({
			location: location || record.location,
			year: year || record.year,
			employeeId: record.employee_id || record.employeeId || record.id,
			firstname: record.firstname,
			lastname: record.lastname,
			title: record.title,
			basepay: record.basepay ? parseFloat(record.basepay.toString()) : 0,
			overtimepay: record.overtimepay ? parseFloat(record.overtimepay.toString()) : 0,
			adjustpay: record.adjustpay ? parseFloat(record.adjustpay.toString()) : 0,
			grosspay: record.grosspay ? parseFloat(record.grosspay.toString()) : 0,
			scrapedAt: record.scraped_at ? new Date(record.scraped_at) : new Date()
		}));

		// Insert records with upsert logic
		const insertQuery = db
			.insert(ucWages)
			.values(wageRecords)
			.onConflictDoUpdate({
				target: [ucWages.location, ucWages.year, ucWages.employeeId],
				set: {
					firstname: sql`EXCLUDED.firstname`,
					lastname: sql`EXCLUDED.lastname`,
					title: sql`EXCLUDED.title`,
					basepay: sql`EXCLUDED.basepay`,
					overtimepay: sql`EXCLUDED.overtimepay`,
					adjustpay: sql`EXCLUDED.adjustpay`,
					grosspay: sql`EXCLUDED.grosspay`,
					scrapedAt: sql`EXCLUDED.scraped_at`,
					uploadedAt: sql`CURRENT_TIMESTAMP`
				}
			});

		await insertQuery;

		return json({
			success: true,
			message: `Successfully uploaded ${wageRecords.length} records for ${location} ${year}`
		});

	} catch (error) {
		console.error('Error uploading wage data:', error);
		return json({ error: 'Failed to upload wage data' }, { status: 500 });
	}
};