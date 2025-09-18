import { db } from '$lib/server/db';
import { wageSummaries, wagePyramids, titleAnalysis } from '$lib/server/db/schema';
import { sql } from 'drizzle-orm';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	try {
		// Get precalculated wage summaries for fast visualization
		const summaryData = await db
			.select({
				location: wageSummaries.location,
				year: wageSummaries.year,
				employeeCount: wageSummaries.employeeCount,
				totalGrossPay: wageSummaries.totalGrossPay,
				avgGrossPay: wageSummaries.avgGrossPay,
				medianPay: wageSummaries.medianPay,
				stdDev: wageSummaries.stdDev,
				minPay: wageSummaries.minPay,
				maxPay: wageSummaries.maxPay,
				percentiles: wageSummaries.percentiles
			})
			.from(wageSummaries)
			.orderBy(wageSummaries.location, wageSummaries.year);

		// Get precalculated wage pyramids for pyramid visualizations
		const pyramidData = await db
			.select({
				location: wagePyramids.location,
				year: wagePyramids.year,
				totalEmployees: wagePyramids.totalEmployees,
				totalPay: wagePyramids.totalPay,
				brackets: wagePyramids.brackets
			})
			.from(wagePyramids)
			.orderBy(wagePyramids.location, wagePyramids.year);

		// Get title analysis data for additional insights
		const titleData = await db
			.select({
				location: titleAnalysis.location,
				year: titleAnalysis.year,
				uniqueTitles: titleAnalysis.uniqueTitles,
				topTitles: titleAnalysis.topTitles
			})
			.from(titleAnalysis)
			.orderBy(titleAnalysis.location, titleAnalysis.year);

		// Convert data for frontend compatibility
		const processedSummaries = summaryData.map(row => ({
			location: row.location,
			year: row.year,
			employeeCount: row.employeeCount,
			totalWages: parseFloat(row.totalGrossPay.toString()),
			averageWage: parseFloat(row.avgGrossPay.toString()),
			medianWage: parseFloat(row.medianPay.toString()),
			stdDev: parseFloat(row.stdDev.toString()),
			minWage: parseFloat(row.minPay.toString()),
			maxWage: parseFloat(row.maxPay.toString()),
			percentiles: typeof row.percentiles === 'string' ? JSON.parse(row.percentiles) : row.percentiles
		}));

		const processedPyramids = pyramidData.map(row => ({
			location: row.location,
			year: row.year,
			totalEmployees: row.totalEmployees,
			totalPay: parseFloat(row.totalPay.toString()),
			brackets: typeof row.brackets === 'string' ? JSON.parse(row.brackets) : row.brackets
		}));

		const processedTitles = titleData.map(row => ({
			location: row.location,
			year: row.year,
			uniqueTitles: row.uniqueTitles,
			topTitles: typeof row.topTitles === 'string' ? JSON.parse(row.topTitles) : row.topTitles
		}));

		// Get summary statistics from precalculated data
		const latestYear = processedSummaries.length > 0 ? Math.max(...processedSummaries.map(d => d.year)) : new Date().getFullYear();
		const latestData = processedSummaries.filter(d => d.year === latestYear);

		const totalEmployees = latestData.reduce((sum, d) => sum + d.employeeCount, 0);
		const totalWages = latestData.reduce((sum, d) => sum + d.totalWages, 0);
		const averageWage = totalEmployees > 0 ? totalWages / totalEmployees : 0;
		const highestPaidCampus = latestData.length > 0
			? latestData.reduce((max, d) => d.averageWage > max.averageWage ? d : max)
			: null;

		return {
			wageData: processedSummaries,
			pyramidData: processedPyramids,
			titleData: processedTitles,
			summary: {
				latestYear,
				totalEmployees,
				totalWages,
				averageWage,
				highestPaidCampus
			}
		};
	} catch (error) {
		console.error('Error loading precalculated wage data:', error);

		// Return empty data if database query fails
		return {
			wageData: [],
			pyramidData: [],
			titleData: [],
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