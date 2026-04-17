import { db } from '$lib/server/db';
import { ucWages } from '$lib/server/db/schema';
import { sql, like, and, eq, desc, asc } from 'drizzle-orm';
import type { PageServerLoad } from './$types';

const ITEMS_PER_PAGE = 50;
const FILTER_CACHE_TTL_MS = 1000 * 60 * 30; // 30 minutes
const SORT_COLUMNS = {
	name: sql`LOWER(CONCAT(${ucWages.firstname}, ' ', ${ucWages.lastname}))`,
	jobtitle: ucWages.title,
	location: ucWages.location,
	year: ucWages.year,
	grosspay: ucWages.grosspay,
	basePay: ucWages.basepay,
	overtimePay: ucWages.overtimepay,
	otherPay: ucWages.adjustpay
} as const;
const DEFAULT_SORT_COLUMN = 'grosspay' as const;
const DEFAULT_SORT_DIRECTION = 'desc' as const;

type Filters = {
	locations: string[];
	years: number[];
};
type SortableColumn = keyof typeof SORT_COLUMNS;
type SortDirection = 'asc' | 'desc';

let filterCache: { data: Filters; expires: number } | null = null;

async function getFilters(): Promise<Filters> {
	const now = Date.now();
	if (filterCache && filterCache.expires > now) {
		return filterCache.data;
	}

	const [locations, years] = await Promise.all([
		db.selectDistinct({ location: ucWages.location }).from(ucWages).orderBy(ucWages.location),
		db.selectDistinct({ year: ucWages.year }).from(ucWages).orderBy(desc(ucWages.year))
	]);

	const data = {
		locations: locations.map(l => l.location),
		years: years.map(y => y.year)
	};

	filterCache = {
		data,
		expires: now + FILTER_CACHE_TTL_MS
	};

	return data;
}

export const load: PageServerLoad = async ({ url }) => {
	try {
		const searchParams = url.searchParams;
		const filterKeys = ['name', 'job', 'location', 'year'] as const;
		const hasFilterParam = filterKeys.some(key => searchParams.has(key));
		const name = searchParams.get('name') || '';
		const job = searchParams.get('job') || '';
		const location = searchParams.get('location') || '';
		const year = searchParams.get('year')
			? parseInt(searchParams.get('year')!)
			: hasFilterParam
				? null
				: 2024;
		const page = Math.max(1, parseInt(searchParams.get('page') || '1'));
		const requestedSort = searchParams.get('sort') as SortableColumn | null;
		const requestedDirection = searchParams.get('direction') === 'asc' ? 'asc' : 'desc';
		const sortColumn: SortableColumn =
			requestedSort && requestedSort in SORT_COLUMNS ? requestedSort : DEFAULT_SORT_COLUMN;
		const sortDirection: SortDirection =
			requestedSort && requestedSort in SORT_COLUMNS ? requestedDirection : DEFAULT_SORT_DIRECTION;
		const orderExpression = SORT_COLUMNS[sortColumn];
		const orderByClause = sortDirection === 'asc' ? asc(orderExpression) : desc(orderExpression);
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

		const whereClause = conditions.length > 0 ? and(...conditions) : undefined;

		const [countResult, results, filters] = await Promise.all([
			db
				.select({ count: sql<number>`count(*)::integer` })
				.from(ucWages)
				.where(whereClause),
			db
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
				.where(whereClause)
				.orderBy(orderByClause)
				.limit(ITEMS_PER_PAGE)
				.offset(offset),
			getFilters()
		]);

		const totalItems = countResult[0]?.count ?? 0;
		const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

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
				year,
				sort: sortColumn,
				direction: sortDirection
			},
			filters
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
				year: 2024,
				sort: DEFAULT_SORT_COLUMN,
				direction: DEFAULT_SORT_DIRECTION
			},
			filters: {
				locations: [],
				years: []
			}
		};
	}
};
