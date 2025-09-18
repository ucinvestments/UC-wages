<script lang="ts">
	import { onMount } from 'svelte';
	import * as d3 from 'd3';
	import { campusThemes, getCampusTheme } from '$lib/utils/campusColors';
	import type { PageData } from './$types';
	import type { AggregatedWageData } from '$lib/types/wages';

	export let data: PageData;

	let heroContainer: HTMLDivElement;
	let totalWagesContainer: HTMLDivElement;
	let averageWagesContainer: HTMLDivElement;
	let employeeCountContainer: HTMLDivElement;
	let wagePyramidContainer: HTMLDivElement;
	let campusGridContainer: HTMLDivElement;

	// Legend containers
	let totalWagesLegend: HTMLDivElement;
	let averageWagesLegend: HTMLDivElement;
	let employeeCountLegend: HTMLDivElement;
	let wagePyramidLegend: HTMLDivElement;

	const { wageData, pyramidData, titleData, summary } = data;
	const { latestYear, totalEmployees, totalWages, averageWage, highestPaidCampus } = summary;

	// Filter and toggle states
	let selectedCampus = 'All Campuses';
	let showTotalLines = false;

	// Get unique campuses for filter dropdown
	$: campuses = [...new Set(wageData.map(d => d.location))];

	// Filter data based on selected campus
	$: filteredWageData = selectedCampus === 'All Campuses'
		? wageData
		: wageData.filter(d => d.location === selectedCampus);

	onMount(() => {
		createHeroVisualization();
		updateCharts();
	});

	// Reactive updates when filters change
	$: if (selectedCampus || showTotalLines !== undefined) {
		updateCharts();
	}

	function createChartLegend(container: HTMLDivElement, campuses: string[], colorScale: any, showTotal: boolean, totalLabel: string) {
		if (!container) return;

		const legendDiv = d3.select(container)
			.style('padding', '10px 0')
			.style('border-bottom', '1px solid #e5e7eb')
			.style('margin-bottom', '15px');

		const legendItems = legendDiv.append('div')
			.style('display', 'flex')
			.style('flex-wrap', 'wrap')
			.style('gap', '15px')
			.style('align-items', 'center')
			.style('justify-content', 'center');

		// Add total line item first if needed
		if (showTotal) {
			const totalItem = legendItems.append('div')
				.style('display', 'flex')
				.style('align-items', 'center')
				.style('gap', '6px');

			totalItem.append('div')
				.style('width', '20px')
				.style('height', '3px')
				.style('background', '#1f2937')
				.style('border', '2px dashed #1f2937')
				.style('opacity', '0.7');

			totalItem.append('span')
				.style('font-size', '12px')
				.style('font-weight', '600')
				.style('color', '#374151')
				.text(totalLabel);
		}

		// Add campus items
		campuses.forEach(campus => {
			const item = legendItems.append('div')
				.style('display', 'flex')
				.style('align-items', 'center')
				.style('gap', '6px');

			item.append('div')
				.style('width', '12px')
				.style('height', '12px')
				.style('background', colorScale(campus))
				.style('border-radius', '2px');

			item.append('span')
				.style('font-size', '12px')
				.style('color', '#374151')
				.text(campus);
		});
	}

	function createPyramidLegend(container: HTMLDivElement, sortedRanges: [string, number][], colorScale: any) {
		if (!container) return;

		const legendDiv = d3.select(container)
			.style('padding', '10px 0')
			.style('border-bottom', '1px solid #e5e7eb')
			.style('margin-bottom', '15px');

		// Add title
		legendDiv.append('div')
			.style('text-align', 'center')
			.style('font-size', '14px')
			.style('font-weight', 'bold')
			.style('color', '#1f2937')
			.style('margin-bottom', '10px')
			.text('Wage Ranges (Employee Count)');

		const legendItems = legendDiv.append('div')
			.style('display', 'flex')
			.style('flex-wrap', 'wrap')
			.style('gap', '12px')
			.style('align-items', 'center')
			.style('justify-content', 'center');

		// Add wage range items
		sortedRanges.forEach(([range, count], index) => {
			const item = legendItems.append('div')
				.style('display', 'flex')
				.style('align-items', 'center')
				.style('gap', '6px');

			item.append('div')
				.style('width', '12px')
				.style('height', '12px')
				.style('background', colorScale(index))
				.style('border-radius', '2px')
				.style('opacity', '0.9')
				.style('border', '1px solid #fff');

			item.append('span')
				.style('font-size', '11px')
				.style('color', '#374151')
				.text(`${range}: ${(count / 1000).toFixed(0)}K`);
		});
	}

	function updateCharts() {
		// Clear existing charts and legends
		if (totalWagesContainer) {
			d3.select(totalWagesContainer).selectAll('*').remove();
			d3.select(totalWagesLegend).selectAll('*').remove();
			createTotalWagesChart();
		}
		if (averageWagesContainer) {
			d3.select(averageWagesContainer).selectAll('*').remove();
			d3.select(averageWagesLegend).selectAll('*').remove();
			createAverageWagesChart();
		}
		if (employeeCountContainer) {
			d3.select(employeeCountContainer).selectAll('*').remove();
			d3.select(employeeCountLegend).selectAll('*').remove();
			createEmployeeCountChart();
		}
		if (wagePyramidContainer) {
			d3.select(wagePyramidContainer).selectAll('*').remove();
			d3.select(wagePyramidLegend).selectAll('*').remove();
			createWagePyramidChart();
		}
		if (campusGridContainer) {
			d3.select(campusGridContainer).selectAll('*').remove();
			createCampusGrid();
		}
	}

	function createHeroVisualization() {
		if (!heroContainer || !wageData.length) return;

		const width = heroContainer.clientWidth;
		const height = 200;

		// Create animated background pattern
		const svg = d3.select(heroContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height)
			.style('position', 'absolute')
			.style('top', 0)
			.style('left', 0)
			.style('z-index', -1);

		// Create flowing data points
		const campuses = [...new Set(wageData.map(d => d.location))];
		const points = campuses.flatMap(campus => {
			const theme = getCampusTheme(campus);
			return Array.from({ length: 5 }, (_, i) => ({
				campus,
				color: theme.primary,
				x: Math.random() * width,
				y: Math.random() * height,
				r: Math.random() * 8 + 3,
				vx: (Math.random() - 0.5) * 2,
				vy: (Math.random() - 0.5) * 2
			}));
		});

		const circles = svg.selectAll('.hero-dot')
			.data(points)
			.enter()
			.append('circle')
			.attr('class', 'hero-dot')
			.attr('cx', d => d.x)
			.attr('cy', d => d.y)
			.attr('r', 0)
			.attr('fill', d => d.color)
			.style('opacity', 0.6);

		// Animate in
		circles
			.transition()
			.delay((d, i) => i * 100)
			.duration(1000)
			.ease(d3.easeBackOut)
			.attr('r', d => d.r);

		// Continuous floating animation
		function animate() {
			circles
				.transition()
				.duration(3000)
				.ease(d3.easeLinear)
				.attr('cx', d => {
					d.x += d.vx;
					if (d.x < 0 || d.x > width) d.vx *= -1;
					return d.x;
				})
				.attr('cy', d => {
					d.y += d.vy;
					if (d.y < 0 || d.y > height) d.vy *= -1;
					return d.y;
				})
				.on('end', animate);
		}
		animate();
	}

	function createTotalWagesChart() {
		if (!totalWagesContainer || !filteredWageData.length) return;

		const width = totalWagesContainer.clientWidth;
		const height = 400;
		const margin = { top: 60, right: 50, bottom: 60, left: 100 };

		const svg = d3.select(totalWagesContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const activeCampuses = [...new Set(filteredWageData.map(d => d.location))];
		const years = [...new Set(filteredWageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		// Better color scheme
		const colorScale = d3.scaleOrdinal()
			.domain(activeCampuses)
			.range(['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#f97316', '#84cc16', '#ec4899', '#6366f1', '#14b8a6', '#f43f5e', '#8b5cf6']);

		// Calculate total line (sum of all campuses by year) - only if showing all campuses and toggle is on
		let totalData = [];
		if (showTotalLines && selectedCampus === 'All Campuses') {
			const totalByYear = d3.rollup(filteredWageData,
				v => d3.sum(v, d => d.totalWages),
				d => d.year
			);
			totalData = Array.from(totalByYear, ([year, total]) => ({ year, totalWages: total }))
				.sort((a, b) => a.year - b.year);
		}

		// Include total data in y-scale domain to ensure total line is visible
		const allWageData = [...filteredWageData.map(d => d.totalWages), ...totalData.map(d => d.totalWages)];
		const yScale = d3.scaleLinear()
			.domain(d3.extent(allWageData) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(filteredWageData, d => d.location);

		// Create line generator with correct typing
		const line = d3.line<any>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.totalWages))
			.curve(d3.curveCatmullRom);

		// Draw total line first (in background)
		if (totalData.length > 1) {
			const totalLine = g.append('path')
				.datum(totalData)
				.attr('fill', 'none')
				.attr('stroke', '#1f2937')
				.attr('stroke-width', 4)
				.attr('stroke-dasharray', '5,5')
				.attr('opacity', 0.7)
				.attr('d', line);

			const totalLength = totalLine.node()!.getTotalLength();
			totalLine
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.duration(2000)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);
		}

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const sortedValues = values.sort((a, b) => a.year - b.year);
			const campusIndex = activeCampuses.indexOf(campus);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', colorScale(campus))
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campusIndex * 200 + 500)
				.duration(1500)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);

			// Add dots
			g.selectAll(`.total-dot-${campus.replace(/\s+/g, '-')}`)
				.data(sortedValues)
				.enter()
				.append('circle')
				.attr('class', `total-dot-${campus.replace(/\s+/g, '-')}`)
				.attr('cx', d => xScale(d.year))
				.attr('cy', d => yScale(d.totalWages))
				.attr('r', 0)
				.attr('fill', colorScale(campus))
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campusIndex * 200 + 1500)
				.duration(800)
				.ease(d3.easeBackOut)
				.attr('r', 5);
		});

		// Add axes
		g.append('g')
			.attr('transform', `translate(0,${height - margin.top - margin.bottom})`)
			.call(d3.axisBottom(xScale).tickFormat(d3.format('d')));

		g.append('g')
			.call(d3.axisLeft(yScale).tickFormat(d => `$${(d / 1e9).toFixed(1)}B`));

		// Create external legend
		createChartLegend(totalWagesLegend, activeCampuses, colorScale, totalData.length > 0, 'System Total');
	}

	function createAverageWagesChart() {
		if (!averageWagesContainer || !filteredWageData.length) return;

		const width = averageWagesContainer.clientWidth;
		const height = 400;
		const margin = { top: 60, right: 50, bottom: 60, left: 100 };

		const svg = d3.select(averageWagesContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const activeCampuses = [...new Set(filteredWageData.map(d => d.location))];
		const years = [...new Set(filteredWageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		// Better color scheme
		const colorScale = d3.scaleOrdinal()
			.domain(activeCampuses)
			.range(['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#f97316', '#84cc16', '#ec4899', '#6366f1', '#14b8a6', '#f43f5e', '#8b5cf6']);

		// Calculate total line (weighted average across all campuses) - only if showing all campuses and toggle is on
		let totalData = [];
		if (showTotalLines && selectedCampus === 'All Campuses') {
			const totalByYear = d3.rollup(filteredWageData,
				v => {
					const totalWages = d3.sum(v, d => d.totalWages);
					const totalEmployees = d3.sum(v, d => d.employeeCount);
					return totalEmployees > 0 ? totalWages / totalEmployees : 0;
				},
				d => d.year
			);
			totalData = Array.from(totalByYear, ([year, averageWage]) => ({ year, averageWage }))
				.sort((a, b) => a.year - b.year);
		}

		// Include total data in y-scale domain to ensure total line is visible
		const allAverageData = [...filteredWageData.map(d => d.averageWage), ...totalData.map(d => d.averageWage)];
		const yScale = d3.scaleLinear()
			.domain(d3.extent(allAverageData) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(filteredWageData, d => d.location);

		// Create line generator
		const line = d3.line<any>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.averageWage))
			.curve(d3.curveCatmullRom);

		// Draw total line first (in background)
		if (totalData.length > 1) {
			const totalLine = g.append('path')
				.datum(totalData)
				.attr('fill', 'none')
				.attr('stroke', '#1f2937')
				.attr('stroke-width', 4)
				.attr('stroke-dasharray', '5,5')
				.attr('opacity', 0.7)
				.attr('d', line);

			const totalLength = totalLine.node()!.getTotalLength();
			totalLine
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.duration(2000)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);
		}

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const sortedValues = values.sort((a, b) => a.year - b.year);
			const campusIndex = campuses.indexOf(campus);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', colorScale(campus))
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campusIndex * 200 + 500)
				.duration(1500)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);

			// Add dots
			g.selectAll(`.avg-dot-${campus.replace(/\s+/g, '-')}`)
				.data(sortedValues)
				.enter()
				.append('circle')
				.attr('class', `avg-dot-${campus.replace(/\s+/g, '-')}`)
				.attr('cx', d => xScale(d.year))
				.attr('cy', d => yScale(d.averageWage))
				.attr('r', 0)
				.attr('fill', colorScale(campus))
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campusIndex * 200 + 1500)
				.duration(800)
				.ease(d3.easeBackOut)
				.attr('r', 5);
		});

		// Add axes
		g.append('g')
			.attr('transform', `translate(0,${height - margin.top - margin.bottom})`)
			.call(d3.axisBottom(xScale).tickFormat(d3.format('d')));

		g.append('g')
			.call(d3.axisLeft(yScale).tickFormat(d => `$${(d / 1000).toFixed(0)}K`));

		// Create external legend
		createChartLegend(averageWagesLegend, activeCampuses, colorScale, showTotalLines && selectedCampus === 'All Campuses' && totalData.length > 1, 'System Average');
	}

	function createEmployeeCountChart() {
		if (!employeeCountContainer || !filteredWageData.length) return;

		const width = employeeCountContainer.clientWidth;
		const height = 400;
		const margin = { top: 80, right: 50, bottom: 60, left: 100 };

		const svg = d3.select(employeeCountContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const activeCampuses = [...new Set(filteredWageData.map(d => d.location))];
		const years = [...new Set(filteredWageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		// Better color scheme
		const colorScale = d3.scaleOrdinal()
			.domain(activeCampuses)
			.range(['#3b82f6', '#ef4444', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#f97316', '#84cc16', '#ec4899', '#6366f1', '#14b8a6', '#f43f5e', '#8b5cf6']);

		// Calculate total line (sum of all employees by year) - only if showing all campuses and toggle is on
		let totalData = [];
		if (showTotalLines && selectedCampus === 'All Campuses') {
			const totalByYear = d3.rollup(filteredWageData,
				v => d3.sum(v, d => d.employeeCount),
				d => d.year
			);
			totalData = Array.from(totalByYear, ([year, employeeCount]) => ({ year, employeeCount }))
				.sort((a, b) => a.year - b.year);
		}

		// Include total data in y-scale domain to ensure total line is visible
		const allEmployeeData = [...filteredWageData.map(d => d.employeeCount), ...totalData.map(d => d.employeeCount)];
		const yScale = d3.scaleLinear()
			.domain(d3.extent(allEmployeeData) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(filteredWageData, d => d.location);

		// Create line generator with correct typing
		const line = d3.line<any>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.employeeCount))
			.curve(d3.curveCatmullRom);

		// Draw total line first (in background)
		if (totalData.length > 1) {
			const totalLine = g.append('path')
				.datum(totalData)
				.attr('fill', 'none')
				.attr('stroke', '#1f2937')
				.attr('stroke-width', 4)
				.attr('stroke-dasharray', '5,5')
				.attr('opacity', 0.7)
				.attr('d', line);

			const totalLength = totalLine.node()!.getTotalLength();
			totalLine
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.duration(2000)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);
		}

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const sortedValues = values.sort((a, b) => a.year - b.year);
			const campusIndex = activeCampuses.indexOf(campus);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', colorScale(campus))
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campusIndex * 200 + 500)
				.duration(1500)
				.ease(d3.easeQuadInOut)
				.attr('stroke-dashoffset', 0);

			// Add dots
			g.selectAll(`.emp-dot-${campus.replace(/\s+/g, '-')}`)
				.data(sortedValues)
				.enter()
				.append('circle')
				.attr('class', `emp-dot-${campus.replace(/\s+/g, '-')}`)
				.attr('cx', d => xScale(d.year))
				.attr('cy', d => yScale(d.employeeCount))
				.attr('r', 0)
				.attr('fill', colorScale(campus))
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campusIndex * 200 + 1500)
				.duration(800)
				.ease(d3.easeBackOut)
				.attr('r', 5);
		});

		// Add axes
		g.append('g')
			.attr('transform', `translate(0,${height - margin.top - margin.bottom})`)
			.call(d3.axisBottom(xScale).tickFormat(d3.format('d')));

		g.append('g')
			.call(d3.axisLeft(yScale).tickFormat(d => d.toLocaleString()));

		// Create external legend
		createChartLegend(employeeCountLegend, activeCampuses, colorScale, showTotalLines && selectedCampus === 'All Campuses' && totalData.length > 1, 'System Total');
	}

	function createWagePyramidChart() {
		if (!wagePyramidContainer || !pyramidData.length) return;

		const width = wagePyramidContainer.clientWidth;
		const height = 600;
		const margin = { top: 80, right: 50, bottom: 80, left: 50 };

		const svg = d3.select(wagePyramidContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${width / 2},${margin.top})`);

		// Use latest year data for pyramid visualization and apply campus filtering
		let filteredPyramidData = pyramidData.filter(d => d.year === latestYear);

		if (selectedCampus !== 'All Campuses') {
			filteredPyramidData = filteredPyramidData.filter(d => d.location === selectedCampus);
		}

		if (filteredPyramidData.length === 0) return;

		// Aggregate all brackets across filtered campuses for pyramid structure
		const wageRanges = {};
		filteredPyramidData.forEach(campus => {
			campus.brackets.forEach(bracket => {
				if (!wageRanges[bracket.range]) {
					wageRanges[bracket.range] = 0;
				}
				wageRanges[bracket.range] += bracket.count;
			});
		});

		// Sort wage ranges from low to high for pyramid structure
		const sortedRanges = Object.entries(wageRanges).sort((a, b) => {
			const getMin = (range) => parseInt(range.split('-')[0].replace(/[^\d]/g, ''));
			return getMin(a[0]) - getMin(b[0]);
		});

		const maxCount = Math.max(...Object.values(wageRanges));
		const pyramidWidth = Math.min(width - margin.left - margin.right, 500); // Limit max width
		const pyramidHeight = height - margin.top - margin.bottom;
		const levelHeight = pyramidHeight / sortedRanges.length;

		// Color gradient from bottom to top
		const colorScale = d3.scaleSequential()
			.domain([0, sortedRanges.length - 1])
			.interpolator(d3.interpolateViridis);


		// Create pyramid levels
		sortedRanges.forEach(([range, count], index) => {
			const levelY = index * levelHeight;
			const levelWidth = (count / maxCount) * pyramidWidth;
			const levelX = -levelWidth / 2;

			// Create trapezoid for pyramid effect - better proportions
			const tapering = 0.8 - (index * 0.02); // More gradual tapering
			const points = [
				[levelX, levelY + levelHeight],
				[levelX + levelWidth, levelY + levelHeight],
				[levelX + levelWidth * tapering, levelY],
				[levelX + levelWidth * (1 - tapering), levelY]
			];

			const levelGroup = g.append('g')
				.attr('class', `pyramid-level-${index}`);

			// Add the pyramid block
			const block = levelGroup.append('polygon')
				.attr('points', points.map(p => `${p[0]},${p[1]}`).join(' '))
				.attr('fill', colorScale(index))
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.attr('opacity', 0)
				.style('cursor', 'pointer');

			// Animate blocks growing from bottom
			block
				.transition()
				.delay(index * 150)
				.duration(1000)
				.ease(d3.easeBounceOut)
				.attr('opacity', 0.85);

			// Add hover effects
			block.on('mouseover', function() {
				d3.select(this)
					.transition()
					.duration(200)
					.attr('opacity', 1)
					.attr('stroke-width', 3);
			})
			.on('mouseout', function() {
				d3.select(this)
					.transition()
					.duration(200)
					.attr('opacity', 0.85)
					.attr('stroke-width', 2);
			});

			// Add count labels - prevent overflow
			if (count > maxCount * 0.03) { // Higher threshold to reduce clutter
				levelGroup.append('text')
					.attr('x', 0)
					.attr('y', levelY + levelHeight / 2)
					.attr('text-anchor', 'middle')
					.attr('dy', '0.35em')
					.style('font-size', Math.max(9, 12 - index * 0.3) + 'px')
					.style('font-weight', 'bold')
					.style('fill', 'white')
					.style('text-shadow', '1px 1px 2px rgba(0,0,0,0.7)')
					.style('opacity', 0)
					.text(count.toLocaleString())
					.transition()
					.delay(index * 150 + 800)
					.duration(600)
					.style('opacity', 1);
			}

			// Add wage range labels - only for levels with significant width
			if (levelWidth > pyramidWidth * 0.1) { // Only show labels for levels wide enough
				const isRightSide = index % 2 === 0;
				const labelX = isRightSide ? levelWidth / 2 + 15 : -(levelWidth / 2) - 15;
				const anchor = isRightSide ? 'start' : 'end';

				levelGroup.append('text')
					.attr('x', labelX)
					.attr('y', levelY + levelHeight / 2)
					.attr('text-anchor', anchor)
					.attr('dy', '0.35em')
					.style('font-size', '10px')
					.style('font-weight', '600')
					.style('fill', '#374151')
					.style('opacity', 0)
					.text(range)
					.transition()
					.delay(index * 150 + 1000)
					.duration(600)
					.style('opacity', 1);
			}
		});

		// Add floating animation for the entire pyramid
		function floatAnimation() {
			g.selectAll('polygon')
				.transition()
				.duration(4000)
				.ease(d3.easeSinInOut)
				.attr('transform', 'translate(0, -3)')
				.transition()
				.duration(4000)
				.ease(d3.easeSinInOut)
				.attr('transform', 'translate(0, 3)')
				.on('end', floatAnimation);
		}

		// Start floating animation after initial animation completes
		setTimeout(floatAnimation, sortedRanges.length * 150 + 1500);

		// Create external pyramid legend
		createPyramidLegend(wagePyramidLegend, sortedRanges, colorScale);

		// Add summary statistics at the bottom
		const totalEmployees = Object.values(wageRanges).reduce((sum, count) => sum + count, 0);
		svg.append('text')
			.attr('x', width / 2)
			.attr('y', height - 20)
			.attr('text-anchor', 'middle')
			.style('font-size', '14px')
			.style('font-weight', '600')
			.style('fill', '#6b7280')
			.text(`Total UC System Employees: ${totalEmployees.toLocaleString()}`);
	}

	function createCampusGrid() {
		if (!campusGridContainer || !filteredWageData.length) return;

		const campuses = [...new Set(wageData.map(d => d.location))];
		const cardWidth = 300;
		const cardHeight = 200;

		campuses.forEach(campus => {
			const campusData = wageData.filter(d => d.location === campus);
			const theme = getCampusTheme(campus);

			const card = d3.select(campusGridContainer)
				.append('div')
				.style('width', `${cardWidth}px`)
				.style('height', `${cardHeight}px`)
				.style('background', theme.gradient)
				.style('border-radius', '16px')
				.style('padding', '20px')
				.style('margin', '10px')
				.style('color', 'white')
				.style('display', 'inline-block')
				.style('vertical-align', 'top');

			card.append('h3')
				.text(theme.name)
				.style('margin', '0 0 10px 0')
				.style('font-size', '18px');

			const latestData = campusData.find(d => d.year === latestYear);
			if (latestData) {
				card.append('p')
					.html(`
						<strong>Employees:</strong> ${latestData.employeeCount.toLocaleString()}<br>
						<strong>Avg Wage:</strong> $${Math.round(latestData.averageWage).toLocaleString()}<br>
						<strong>Total Wages:</strong> $${(latestData.totalWages / 1e6).toFixed(1)}M
					`)
					.style('font-size', '14px')
					.style('line-height', '1.5');
			}

			// Add mini chart
			const svg = card.append('svg')
				.attr('width', cardWidth - 40)
				.attr('height', 80)
				.style('margin-top', '10px');

			if (campusData.length > 1) {
				const years = campusData.map(d => d.year).sort();
				const xScale = d3.scaleLinear()
					.domain(d3.extent(years) as [number, number])
					.range([0, cardWidth - 40]);

				const yScale = d3.scaleLinear()
					.domain(d3.extent(campusData, d => d.averageWage) as [number, number])
					.range([60, 0]);

				const line = d3.line<AggregatedWageData>()
					.x(d => xScale(d.year))
					.y(d => yScale(d.averageWage))
					.curve(d3.curveCatmullRom);

				svg.append('path')
					.datum(campusData.sort((a, b) => a.year - b.year))
					.attr('fill', 'none')
					.attr('stroke', 'rgba(255,255,255,0.8)')
					.attr('stroke-width', 2)
					.attr('d', line);
			}
		});
	}
</script>

<svelte:head>
	<title>UC Wage Explorer - Interactive Data Visualization</title>
	<meta name="description" content="Explore University of California employee wage data with interactive D3.js visualizations" />
</svelte:head>

<!-- Hero Section -->
<section class="hero" bind:this={heroContainer}>
	<div class="hero-content">
		<h1 class="hero-title">UC Wage Explorer</h1>
		<p class="hero-subtitle">Interactive visualization of University of California employee compensation data</p>
		<div class="hero-stats">
			{#if totalEmployees > 0}
				<div class="stat">
					<span class="stat-value">{totalEmployees.toLocaleString()}</span>
					<span class="stat-label">{latestYear} Total Employees</span>
				</div>
				<div class="stat">
					<span class="stat-value">${(totalWages / 1e9).toFixed(1)}B</span>
					<span class="stat-label">{latestYear} Total Wages</span>
				</div>
				<div class="stat">
					<span class="stat-value">${(averageWage / 1000).toFixed(0)}K</span>
					<span class="stat-label">{latestYear} Average Wage</span>
				</div>
			{/if}
		</div>
	</div>
</section>

<!-- Chart Controls -->
<section class="section">
	<div class="container">
		<div class="chart-controls">
			<div class="control-group">
				<label for="campus-filter">Filter by Campus:</label>
				<select id="campus-filter" bind:value={selectedCampus} class="control-select">
					<option value="All Campuses">All Campuses</option>
					{#each campuses as campus}
						<option value={campus}>{campus}</option>
					{/each}
				</select>
			</div>

			<div class="control-group">
				<label class="toggle-label">
					<input type="checkbox" bind:checked={showTotalLines} class="toggle-checkbox" />
					<span class="toggle-text">Show Total Lines</span>
				</label>
			</div>
		</div>
	</div>
</section>

<!-- Total Wages Section -->
<section class="section">
	<div class="container">
		<h2>Total Wages by Campus</h2>
		<div class="chart-legend" bind:this={totalWagesLegend}></div>
		<div class="chart-container" bind:this={totalWagesContainer}></div>
	</div>
</section>

<!-- Average Wages Section -->
<section class="section">
	<div class="container">
		<h2>Average Wages by Campus</h2>
		<div class="chart-legend" bind:this={averageWagesLegend}></div>
		<div class="chart-container" bind:this={averageWagesContainer}></div>
	</div>
</section>

<!-- Employee Count Section -->
<section class="section">
	<div class="container">
		<h2>Employee Count by Campus</h2>
		<div class="chart-legend" bind:this={employeeCountLegend}></div>
		<div class="chart-container" bind:this={employeeCountContainer}></div>
	</div>
</section>

<!-- Wage Pyramid Section -->
<section class="section">
	<div class="container">
		<h2>Wage Distribution Pyramid</h2>
		<div class="chart-legend" bind:this={wagePyramidLegend}></div>
		<div class="chart-container" bind:this={wagePyramidContainer}></div>
	</div>
</section>

<!-- Campus Grid -->
<section class="section">
	<div class="container">
		<h2>Campus Breakdown</h2>
		<div class="campus-grid" bind:this={campusGridContainer}></div>
	</div>
</section>

<style>
	.hero {
		position: relative;
		min-height: 70vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #1e293b 0%, #3b82f6 50%, #8b5cf6 100%);
		color: white;
		text-align: center;
		overflow: hidden;
	}

	.hero-content {
		z-index: 1;
		max-width: 800px;
		padding: 2rem;
	}

	.hero-title {
		font-size: 4rem;
		font-weight: 800;
		margin-bottom: 1rem;
		letter-spacing: -0.02em;
	}

	.hero-subtitle {
		font-size: 1.5rem;
		margin-bottom: 3rem;
		opacity: 0.9;
	}

	.hero-stats {
		display: flex;
		justify-content: center;
		gap: 4rem;
		margin-top: 2rem;
	}

	.stat {
		text-align: center;
	}

	.stat-value {
		display: block;
		font-size: 2.5rem;
		font-weight: 700;
		color: #fbbf24;
	}

	.stat-label {
		font-size: 0.875rem;
		opacity: 0.8;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.section {
		padding: 4rem 0;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 2rem;
	}

	h2 {
		font-size: 2.5rem;
		font-weight: 700;
		text-align: center;
		margin-bottom: 3rem;
		background: linear-gradient(135deg, #1e293b, #3b82f6);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.chart-container {
		background: white;
		border-radius: 16px;
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
		padding: 2rem;
		margin: 2rem 0;
		min-height: 450px;
	}

	.campus-grid {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 1rem;
		margin: 2rem 0;
	}

	.chart-controls {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 3rem;
		background: white;
		border-radius: 12px;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		padding: 1.5rem 2rem;
		margin-bottom: 1rem;
	}

	.control-group {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.control-group label {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
		white-space: nowrap;
	}

	.control-select {
		padding: 0.5rem 1rem;
		border: 2px solid #e5e7eb;
		border-radius: 8px;
		font-size: 0.875rem;
		background: white;
		color: #374151;
		min-width: 140px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.control-select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.toggle-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
		user-select: none;
	}

	.toggle-checkbox {
		width: 1rem;
		height: 1rem;
		cursor: pointer;
	}

	.toggle-text {
		font-weight: 500;
		color: #374151;
	}

	.chart-legend {
		background: #f8fafc;
		border-radius: 8px;
		margin: 0 0 20px 0;
		min-height: 40px;
	}

	@media (max-width: 768px) {
		.hero-title {
			font-size: 2.5rem;
		}

		.hero-subtitle {
			font-size: 1.2rem;
		}

		.hero-stats {
			flex-direction: column;
			gap: 2rem;
		}

		.stat-value {
			font-size: 2rem;
		}

		h2 {
			font-size: 2rem;
		}
	}
</style>
