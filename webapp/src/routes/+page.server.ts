import { db } from '$lib/server/db';
import { ucWages } from '$lib/server/db/schema';
import { sql } from 'drizzle-orm';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		// Get aggregated wage data for visualization
		const aggregatedData = await db
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
			.orderBy(ucWages.location, ucWages.year);

		// Convert decimal strings to numbers for frontend
		const processedData = aggregatedData.map(row => ({
			...row,
			totalWages: parseFloat(row.totalWages.toString()),
			averageWage: parseFloat(row.averageWage.toString()),
			maxWage: parseFloat(row.maxWage.toString()),
			minWage: parseFloat(row.minWage.toString())
		}));

		// Get summary statistics
		const latestYear = processedData.length > 0 ? Math.max(...processedData.map(d => d.year)) : new Date().getFullYear();
		const latestData = processedData.filter(d => d.year === latestYear);

		const totalEmployees = latestData.reduce((sum, d) => sum + d.employeeCount, 0);
		const totalWages = latestData.reduce((sum, d) => sum + d.totalWages, 0);
		const averageWage = totalEmployees > 0 ? totalWages / totalEmployees : 0;
		const highestPaidCampus = latestData.length > 0
			? latestData.reduce((max, d) => d.averageWage > max.averageWage ? d : max)
			: null;

		return {
			wageData: processedData,
			summary: {
				latestYear,
				totalEmployees,
				totalWages,
				averageWage,
				highestPaidCampus
			}
		};
	} catch (error) {
		console.error('Error loading wage data:', error);

		// Return empty data if database query fails
		return {
			wageData: [],
			summary: {
				latestYear: new Date().getFullYear(),
				totalEmployees: 0,
				totalWages: 0,
				averageWage: 0,
				highestPaidCampus: null
			}
		};
	}
};