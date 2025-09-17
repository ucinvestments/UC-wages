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
	let campusGridContainer: HTMLDivElement;

	const { wageData, summary } = data;
	const { latestYear, totalEmployees, totalWages, averageWage, highestPaidCampus } = summary;

	onMount(() => {
		createHeroVisualization();
		createTotalWagesChart();
		createAverageWagesChart();
		createEmployeeCountChart();
		createCampusGrid();
	});

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
		if (!totalWagesContainer || !wageData.length) return;

		const width = totalWagesContainer.clientWidth;
		const height = 400;
		const margin = { top: 40, right: 150, bottom: 60, left: 100 };

		const svg = d3.select(totalWagesContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const campuses = [...new Set(wageData.map(d => d.location))];
		const years = [...new Set(wageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		const yScale = d3.scaleLinear()
			.domain(d3.extent(wageData, d => d.totalWages) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(wageData, d => d.location);

		// Create line generator
		const line = d3.line<AggregatedWageData>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.totalWages))
			.curve(d3.curveCatmullRom);

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const theme = getCampusTheme(campus);
			const sortedValues = values.sort((a, b) => a.year - b.year);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', theme.primary)
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campuses.indexOf(campus) * 200)
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
				.attr('fill', theme.primary)
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campuses.indexOf(campus) * 200 + 1000)
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

		// Add title
		g.append('text')
			.attr('x', (width - margin.left - margin.right) / 2)
			.attr('y', -10)
			.attr('text-anchor', 'middle')
			.style('font-size', '16px')
			.style('font-weight', 'bold')
			.text('Total Wages by Campus');

		// Add legend
		const legend = svg.append('g')
			.attr('transform', `translate(${width - margin.right + 20}, ${margin.top})`);

		campuses.forEach((campus, i) => {
			const theme = getCampusTheme(campus);
			const legendItem = legend.append('g')
				.attr('transform', `translate(0, ${i * 25})`);

			legendItem.append('rect')
				.attr('width', 18)
				.attr('height', 3)
				.attr('fill', theme.primary);

			legendItem.append('text')
				.attr('x', 25)
				.attr('y', 2)
				.attr('dy', '0.35em')
				.style('font-size', '12px')
				.text(campus);
		});
	}

	function createAverageWagesChart() {
		if (!averageWagesContainer || !wageData.length) return;

		const width = averageWagesContainer.clientWidth;
		const height = 400;
		const margin = { top: 40, right: 150, bottom: 60, left: 100 };

		const svg = d3.select(averageWagesContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const campuses = [...new Set(wageData.map(d => d.location))];
		const years = [...new Set(wageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		const yScale = d3.scaleLinear()
			.domain(d3.extent(wageData, d => d.averageWage) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(wageData, d => d.location);

		// Create line generator
		const line = d3.line<AggregatedWageData>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.averageWage))
			.curve(d3.curveCatmullRom);

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const theme = getCampusTheme(campus);
			const sortedValues = values.sort((a, b) => a.year - b.year);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', theme.primary)
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campuses.indexOf(campus) * 200)
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
				.attr('fill', theme.primary)
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campuses.indexOf(campus) * 200 + 1000)
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

		// Add title
		g.append('text')
			.attr('x', (width - margin.left - margin.right) / 2)
			.attr('y', -10)
			.attr('text-anchor', 'middle')
			.style('font-size', '16px')
			.style('font-weight', 'bold')
			.text('Average Wages by Campus');

		// Add legend
		const legend = svg.append('g')
			.attr('transform', `translate(${width - margin.right + 20}, ${margin.top})`);

		campuses.forEach((campus, i) => {
			const theme = getCampusTheme(campus);
			const legendItem = legend.append('g')
				.attr('transform', `translate(0, ${i * 25})`);

			legendItem.append('rect')
				.attr('width', 18)
				.attr('height', 3)
				.attr('fill', theme.primary);

			legendItem.append('text')
				.attr('x', 25)
				.attr('y', 2)
				.attr('dy', '0.35em')
				.style('font-size', '12px')
				.text(campus);
		});
	}

	function createEmployeeCountChart() {
		if (!employeeCountContainer || !wageData.length) return;

		const width = employeeCountContainer.clientWidth;
		const height = 400;
		const margin = { top: 40, right: 150, bottom: 60, left: 100 };

		const svg = d3.select(employeeCountContainer)
			.append('svg')
			.attr('width', width)
			.attr('height', height);

		const g = svg.append('g')
			.attr('transform', `translate(${margin.left},${margin.top})`);

		const campuses = [...new Set(wageData.map(d => d.location))];
		const years = [...new Set(wageData.map(d => d.year))].sort();

		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width - margin.left - margin.right]);

		const yScale = d3.scaleLinear()
			.domain(d3.extent(wageData, d => d.employeeCount) as [number, number])
			.range([height - margin.top - margin.bottom, 0]);

		// Group data by campus
		const campusData = d3.group(wageData, d => d.location);

		// Create line generator
		const line = d3.line<AggregatedWageData>()
			.x(d => xScale(d.year))
			.y(d => yScale(d.employeeCount))
			.curve(d3.curveCatmullRom);

		// Draw lines for each campus
		campusData.forEach((values, campus) => {
			const theme = getCampusTheme(campus);
			const sortedValues = values.sort((a, b) => a.year - b.year);

			// Add line
			const path = g.append('path')
				.datum(sortedValues)
				.attr('fill', 'none')
				.attr('stroke', theme.primary)
				.attr('stroke-width', 3)
				.attr('d', line);

			// Animate line drawing
			const totalLength = path.node()!.getTotalLength();
			path
				.attr('stroke-dasharray', totalLength + ' ' + totalLength)
				.attr('stroke-dashoffset', totalLength)
				.transition()
				.delay(campuses.indexOf(campus) * 200)
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
				.attr('fill', theme.primary)
				.attr('stroke', 'white')
				.attr('stroke-width', 2)
				.transition()
				.delay(campuses.indexOf(campus) * 200 + 1000)
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

		// Add title
		g.append('text')
			.attr('x', (width - margin.left - margin.right) / 2)
			.attr('y', -10)
			.attr('text-anchor', 'middle')
			.style('font-size', '16px')
			.style('font-weight', 'bold')
			.text('Employee Count by Campus');

		// Add legend
		const legend = svg.append('g')
			.attr('transform', `translate(${width - margin.right + 20}, ${margin.top})`);

		campuses.forEach((campus, i) => {
			const theme = getCampusTheme(campus);
			const legendItem = legend.append('g')
				.attr('transform', `translate(0, ${i * 25})`);

			legendItem.append('rect')
				.attr('width', 18)
				.attr('height', 3)
				.attr('fill', theme.primary);

			legendItem.append('text')
				.attr('x', 25)
				.attr('y', 2)
				.attr('dy', '0.35em')
				.style('font-size', '12px')
				.text(campus);
		});
	}

	function createCampusGrid() {
		if (!campusGridContainer || !wageData.length) return;

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
					<span class="stat-label">Total Employees</span>
				</div>
				<div class="stat">
					<span class="stat-value">${(totalWages / 1e9).toFixed(1)}B</span>
					<span class="stat-label">Total Wages</span>
				</div>
				<div class="stat">
					<span class="stat-value">{latestYear}</span>
					<span class="stat-label">Latest Data</span>
				</div>
			{/if}
		</div>
	</div>
</section>

<!-- Total Wages Section -->
<section class="section">
	<div class="container">
		<h2>Total Wages by Campus</h2>
		<div class="chart-container" bind:this={totalWagesContainer}></div>
	</div>
</section>

<!-- Average Wages Section -->
<section class="section">
	<div class="container">
		<h2>Average Wages by Campus</h2>
		<div class="chart-container" bind:this={averageWagesContainer}></div>
	</div>
</section>

<!-- Employee Count Section -->
<section class="section">
	<div class="container">
		<h2>Employee Count by Campus</h2>
		<div class="chart-container" bind:this={employeeCountContainer}></div>
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
