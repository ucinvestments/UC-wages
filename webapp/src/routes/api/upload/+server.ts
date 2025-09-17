import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { ucWages, uploadProgress } from '$lib/server/db/schema';
import { eq, and, sql } from 'drizzle-orm';
import type { RequestHandler } from './$types';

// Handle file upload for wage data
export const POST: RequestHandler = async ({ request }) => {
	try {
		const formData = await request.formData();
		const file = formData.get('file') as File;
		const location = formData.get('location') as string;
		const year = parseInt(formData.get('year') as string);

		if (!file) {
			return json({ error: 'No file provided' }, { status: 400 });
		}

		if (!location || !year) {
			return json({ error: 'Location and year are required' }, { status: 400 });
		}

		// Read file content
		const fileContent = await file.text();

		try {
			const wageData = JSON.parse(fileContent);

			if (!wageData.records || !Array.isArray(wageData.records)) {
				return json({ error: 'Invalid file format: records array not found' }, { status: 400 });
			}

			// Start progress tracking
			await db
				.insert(uploadProgress)
				.values({
					location,
					year,
					totalRecords: wageData.records.length,
					uploadedRecords: 0,
					status: 'processing'
				})
				.onConflictDoUpdate({
					target: [uploadProgress.location, uploadProgress.year],
					set: {
						totalRecords: wageData.records.length,
						uploadedRecords: 0,
						status: 'processing',
						startedAt: sql`CURRENT_TIMESTAMP`,
						errorMessage: null
					}
				});

			// Process records in batches
			const batchSize = 1000;
			let uploadedRecords = 0;

			for (let i = 0; i < wageData.records.length; i += batchSize) {
				const batch = wageData.records.slice(i, i + batchSize);

				// Convert batch to the correct format
				const wageRecords = batch.map((record: any) => ({
					location: wageData.location || location,
					year: wageData.year || year,
					employeeId: record.id || record.employee_id,
					firstname: record.firstname || '',
					lastname: record.lastname || '',
					title: record.title || '',
					basepay: parseFloat(record.basepay?.toString().replace(/,/g, '') || '0'),
					overtimepay: parseFloat(record.overtimepay?.toString().replace(/,/g, '') || '0'),
					adjustpay: parseFloat(record.adjustpay?.toString().replace(/,/g, '') || '0'),
					grosspay: parseFloat(record.grosspay?.toString().replace(/,/g, '') || '0'),
					scrapedAt: wageData.scraped_at ? new Date(wageData.scraped_at) : new Date()
				}));

				// Insert batch with upsert logic
				await db
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

				uploadedRecords += wageRecords.length;

				// Update progress
				await db
					.update(uploadProgress)
					.set({
						uploadedRecords,
						status: uploadedRecords >= wageData.records.length ? 'completed' : 'processing',
						completedAt: uploadedRecords >= wageData.records.length ? sql`CURRENT_TIMESTAMP` : undefined
					})
					.where(and(
						eq(uploadProgress.location, location),
						eq(uploadProgress.year, year)
					));
			}

			return json({
				success: true,
				message: `Successfully uploaded ${uploadedRecords} records for ${location} ${year}`,
				uploadedRecords,
				totalRecords: wageData.records.length
			});

		} catch (parseError) {
			// Update progress with error
			await db
				.update(uploadProgress)
				.set({
					status: 'failed',
					errorMessage: `JSON parse error: ${parseError}`,
					completedAt: sql`CURRENT_TIMESTAMP`
				})
				.where(and(
					eq(uploadProgress.location, location),
					eq(uploadProgress.year, year)
				));

			return json({ error: 'Invalid JSON file format' }, { status: 400 });
		}

	} catch (error) {
		console.error('Error uploading file:', error);
		return json({ error: 'Failed to upload file' }, { status: 500 });
	}
};

// Get upload progress
export const GET: RequestHandler = async ({ url }) => {
	try {
		const location = url.searchParams.get('location');
		const year = url.searchParams.get('year');

		let query = db.select().from(uploadProgress);

		if (location && year) {
			query = query.where(and(
				eq(uploadProgress.location, location),
				eq(uploadProgress.year, parseInt(year))
			));
		} else if (location) {
			query = query.where(eq(uploadProgress.location, location));
		}

		const result = await query;
		return json(result);

	} catch (error) {
		console.error('Error fetching upload progress:', error);
		return json({ error: 'Failed to fetch upload progress' }, { status: 500 });
	}
};