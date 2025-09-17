// Real wage data types based on the UC database schema
export interface WageRecord {
	id: number;
	location: string;
	year: number;
	employeeId: number | null;
	firstname: string | null;
	lastname: string | null;
	title: string | null;
	basepay: number;
	overtimepay: number;
	adjustpay: number;
	grosspay: number;
	scrapedAt: Date | null;
	uploadedAt: Date | null;
}

export interface AggregatedWageData {
	location: string;
	year: number;
	totalWages: number;
	averageWage: number;
	employeeCount: number;
	maxWage: number;
	minWage: number;
}

export interface UploadProgress {
	id: number;
	location: string;
	year: number;
	totalRecords: number | null;
	uploadedRecords: number;
	status: 'pending' | 'processing' | 'completed' | 'failed';
	startedAt: Date | null;
	completedAt: Date | null;
	errorMessage: string | null;
}

export interface WageFileData {
	location: string;
	year: number;
	scraped_at: string;
	total_records: number;
	records: Array<{
		id?: number;
		employee_id?: number;
		firstname?: string;
		lastname?: string;
		title?: string;
		basepay?: string | number;
		overtimepay?: string | number;
		adjustpay?: string | number;
		grosspay?: string | number;
	}>;
}

// UC Campus list
export const ucCampuses = [
	'ASUCLA',
	'Berkeley',
	'Davis',
	'UC SF Law',
	'Irvine',
	'Los Angeles',
	'Merced',
	'Riverside',
	'San Diego',
	'San Francisco',
	'Santa Barbara',
	'Santa Cruz',
	'UCOP'
] as const;

export type UCCampus = typeof ucCampuses[number];