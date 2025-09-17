// UC Campus Brand Colors and Themes
export interface CampusTheme {
	primary: string;
	secondary: string;
	accent: string;
	light: string;
	dark: string;
	gradient: string;
	name: string;
}

export const campusThemes: Record<string, CampusTheme> = {
	'Berkeley': {
		primary: '#003262', // Berkeley Blue
		secondary: '#FDB515', // California Gold
		accent: '#859438', // Medalist
		light: '#C4820E',
		dark: '#001B44',
		gradient: 'linear-gradient(135deg, #003262 0%, #FDB515 100%)',
		name: 'UC Berkeley'
	},
	'Los Angeles': {
		primary: '#2774AE', // UCLA Blue
		secondary: '#FFD100', // UCLA Gold
		accent: '#8BB8E8',
		light: '#FFED4E',
		dark: '#003B5C',
		gradient: 'linear-gradient(135deg, #2774AE 0%, #FFD100 100%)',
		name: 'UCLA'
	},
	'San Diego': {
		primary: '#182B49', // UCSD Navy
		secondary: '#C69214', // UCSD Gold
		accent: '#6F99D3',
		light: '#E8E3D0',
		dark: '#0F1419',
		gradient: 'linear-gradient(135deg, #182B49 0%, #C69214 100%)',
		name: 'UC San Diego'
	},
	'Davis': {
		primary: '#022851', // Aggie Blue
		secondary: '#FFBF00', // Aggie Gold
		accent: '#B0D236',
		light: '#C6AA76',
		dark: '#001633',
		gradient: 'linear-gradient(135deg, #022851 0%, #FFBF00 100%)',
		name: 'UC Davis'
	},
	'Irvine': {
		primary: '#0064A4', // UCI Blue
		secondary: '#FFD200', // UCI Gold
		accent: '#1B365D',
		light: '#7FC8F8',
		dark: '#003D5C',
		gradient: 'linear-gradient(135deg, #0064A4 0%, #FFD200 100%)',
		name: 'UC Irvine'
	},
	'Santa Barbara': {
		primary: '#003660', // UCSB Navy
		secondary: '#DDD15A', // UCSB Gold
		accent: '#9CBEBE',
		light: '#F2E6CE',
		dark: '#001B2E',
		gradient: 'linear-gradient(135deg, #003660 0%, #DDD15A 100%)',
		name: 'UC Santa Barbara'
	},
	'Santa Cruz': {
		primary: '#003C6C', // UCSC Blue
		secondary: '#F29813', // UCSC Gold/Orange
		accent: '#7BA7BC',
		light: '#FFEAA7',
		dark: '#002147',
		gradient: 'linear-gradient(135deg, #003C6C 0%, #F29813 100%)',
		name: 'UC Santa Cruz'
	},
	'Riverside': {
		primary: '#003DA5', // UCR Blue
		secondary: '#FFB81C', // UCR Gold
		accent: '#005581',
		light: '#B7CDF1',
		dark: '#002952',
		gradient: 'linear-gradient(135deg, #003DA5 0%, #FFB81C 100%)',
		name: 'UC Riverside'
	},
	'San Francisco': {
		primary: '#052049', // UCSF Navy
		secondary: '#90BD31', // UCSF Green
		accent: '#18A3AC',
		light: '#6DACE4',
		dark: '#031329',
		gradient: 'linear-gradient(135deg, #052049 0%, #90BD31 100%)',
		name: 'UCSF'
	},
	'Merced': {
		primary: '#002856', // UCM Blue
		secondary: '#FFB310', // UCM Gold
		accent: '#8A8B8C',
		light: '#B3A369',
		dark: '#001B3D',
		gradient: 'linear-gradient(135deg, #002856 0%, #FFB310 100%)',
		name: 'UC Merced'
	},
	'UCOP': {
		primary: '#1B4685', // UC System Blue
		secondary: '#DAAA00', // UC System Gold
		accent: '#767F8B',
		light: '#8BB8E8',
		dark: '#0D2142',
		gradient: 'linear-gradient(135deg, #1B4685 0%, #DAAA00 100%)',
		name: 'UC Office of the President'
	},
	'ASUCLA': {
		primary: '#2774AE', // UCLA colors (ASUCLA is UCLA-affiliated)
		secondary: '#FFD100',
		accent: '#8BB8E8',
		light: '#FFED4E',
		dark: '#003B5C',
		gradient: 'linear-gradient(135deg, #2774AE 0%, #FFD100 100%)',
		name: 'ASUCLA'
	},
	'UC SF Law': {
		primary: '#052049', // UCSF colors
		secondary: '#90BD31',
		accent: '#18A3AC',
		light: '#6DACE4',
		dark: '#031329',
		gradient: 'linear-gradient(135deg, #052049 0%, #90BD31 100%)',
		name: 'UC San Francisco Law'
	}
};

export function getCampusTheme(campus: string): CampusTheme {
	return campusThemes[campus] || {
		primary: '#1B4685',
		secondary: '#DAAA00',
		accent: '#767F8B',
		light: '#8BB8E8',
		dark: '#0D2142',
		gradient: 'linear-gradient(135deg, #1B4685 0%, #DAAA00 100%)',
		name: campus
	};
}

export function getCampusColorScale() {
	const colors = Object.values(campusThemes).map(theme => theme.primary);
	const campuses = Object.keys(campusThemes);
	return { colors, campuses };
}