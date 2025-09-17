// Mock wage data for UC campuses over the years
export interface WageRecord {
	campus: string;
	year: number;
	totalWages: number;
	averageWage: number;
	employeeCount: number;
}

export const ucCampuses = [
	'Berkeley',
	'Davis',
	'Irvine',
	'Los Angeles',
	'Merced',
	'Riverside',
	'San Diego',
	'San Francisco',
	'Santa Barbara',
	'Santa Cruz',
	'UCOP'
];

export const years = Array.from({ length: 15 }, (_, i) => 2010 + i); // 2010-2024

// Generate mock data - in reality this would come from your database
export function generateMockWageData(): WageRecord[] {
	const data: WageRecord[] = [];

	for (const campus of ucCampuses) {
		for (const year of years) {
			// Base values that vary by campus
			const baseMultiplier = getBaseMultiplier(campus);
			const yearGrowth = Math.pow(1.03, year - 2010); // 3% yearly growth

			// Add some randomness for realistic variation
			const variation = 0.9 + Math.random() * 0.2; // ±10% variation

			const employeeCount = Math.floor(baseMultiplier * 5000 * variation);
			const averageWage = Math.floor(baseMultiplier * 75000 * yearGrowth * variation);
			const totalWages = employeeCount * averageWage;

			data.push({
				campus,
				year,
				totalWages,
				averageWage,
				employeeCount
			});
		}
	}

	return data;
}

function getBaseMultiplier(campus: string): number {
	// Different campuses have different scales
	const multipliers: Record<string, number> = {
		'Berkeley': 2.5,
		'Los Angeles': 3.0,
		'San Diego': 2.2,
		'San Francisco': 2.8, // Medical school inflates wages
		'Davis': 2.0,
		'Irvine': 1.8,
		'Santa Barbara': 1.5,
		'Santa Cruz': 1.2,
		'Riverside': 1.4,
		'Merced': 0.8, // Newest campus
		'UCOP': 0.5 // Administrative office
	};

	return multipliers[campus] || 1.0;
}

export const mockWageData = generateMockWageData();