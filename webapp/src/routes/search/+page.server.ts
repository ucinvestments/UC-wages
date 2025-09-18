import { db } from '$lib/server/db';
import { ucWages } from '$lib/server/db/schema';
import { sql, like, and, eq, desc } from 'drizzle-orm';
import type { PageServerLoad } from './$types';

const ITEMS_PER_PAGE = 50;

export const load: PageServerLoad = async ({ url }) => {
	try {
		const searchParams = url.searchParams;
		const name = searchParams.get('name') || '';
		const job = searchParams.get('job') || '';
		const location = searchParams.get('location') || '';
		const year = searchParams.get('year') ? parseInt(searchParams.get('year')!) : null;
		const page = Math.max(1, parseInt(searchParams.get('page') || '1'));
		const offset = (page - 1) * ITEMS_PER_PAGE;

		// Build search conditions
		const conditions = [];

		if (name) {
			// Search in both firstname and lastname
			conditions.push(
				sql`(${ucWages.firstname} ILIKE ${`%${name}%`} OR ${ucWages.lastname} ILIKE ${`%${name}%`} OR CONCAT(${ucWages.firstname}, ' ', ${ucWages.lastname}) ILIKE ${`%${name}%`})`
			);
		}
		if (job) {
			conditions.push(like(ucWages.title, `%${job}%`));
		}
		if (location) {
			conditions.push(eq(ucWages.location, location));
		}
		if (year) {
			conditions.push(eq(ucWages.year, year));
		}

		// Get total count for pagination
		const [countResult] = await db
			.select({ count: sql<number>`count(*)::integer` })
			.from(ucWages)
			.where(conditions.length > 0 ? and(...conditions) : undefined);

		const totalItems = countResult.count;
		const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

		// Get paginated results
		const results = await db
			.select({
				name: sql<string>`CONCAT(${ucWages.firstname}, ' ', ${ucWages.lastname})`,
				jobtitle: ucWages.title,
				location: ucWages.location,
				year: ucWages.year,
				grosspay: ucWages.grosspay,
				basePay: ucWages.basepay,
				overtimePay: ucWages.overtimepay,
				otherPay: ucWages.adjustpay
			})
			.from(ucWages)
			.where(conditions.length > 0 ? and(...conditions) : undefined)
			.orderBy(desc(ucWages.grosspay))
			.limit(ITEMS_PER_PAGE)
			.offset(offset);

		// Get available filters for dropdowns
		const locations = await db
			.selectDistinct({ location: ucWages.location })
			.from(ucWages)
			.orderBy(ucWages.location);

		const years = await db
			.selectDistinct({ year: ucWages.year })
			.from(ucWages)
			.orderBy(desc(ucWages.year));

		return {
			employees: results.map(emp => ({
				...emp,
				grosspay: parseFloat(emp.grosspay.toString()),
				basePay: parseFloat(emp.basePay.toString()),
				overtimePay: parseFloat(emp.overtimePay.toString()),
				otherPay: parseFloat(emp.otherPay.toString())
			})),
			pagination: {
				currentPage: page,
				totalPages,
				totalItems,
				itemsPerPage: ITEMS_PER_PAGE,
				hasNext: page < totalPages,
				hasPrev: page > 1
			},
			searchParams: {
				name,
				job,
				location,
				year
			},
			filters: {
				locations: locations.map(l => l.location),
				years: years.map(y => y.year)
			}
		};
	} catch (error) {
		console.error('Error loading employee search data:', error);

		return {
			employees: [],
			pagination: {
				currentPage: 1,
				totalPages: 0,
				totalItems: 0,
				itemsPerPage: ITEMS_PER_PAGE,
				hasNext: false,
				hasPrev: false
			},
			searchParams: {
				name: '',
				job: '',
				location: '',
				year: null
			},
			filters: {
				locations: [],
				years: []
			}
		};
	}
};