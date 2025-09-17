<script lang="ts">
	import { onMount } from 'svelte';
	import * as d3 from 'd3';
	import type { AggregatedWageData } from '$lib/types/wages';

	export let data: AggregatedWageData[];
	export let metric: 'totalWages' | 'averageWage' | 'employeeCount' = 'totalWages';

	let chartContainer: HTMLDivElement;
	let selectedCampuses: Set<string> = new Set();
	let svg: d3.Selection<SVGSVGElement, unknown, null, undefined>;
	let g: d3.Selection<SVGGElement, unknown, null, undefined>;
	let isLogarithmic = false;
	let isAnimating = false;

	// Get unique campuses and years
	$: campuses = [...new Set(data.map(d => d.location))].sort();
	$: years = [...new Set(data.map(d => d.year))].sort();

	// Initialize with all campuses selected
	onMount(() => {
		selectedCampuses = new Set(campuses);
		initializeChart();
		drawChart(true); // Initial draw with animation
	});

	// Redraw chart when data or selections change
	$: if (chartContainer && selectedCampuses.size > 0 && svg) {
		drawChart(false); // Redraw without full animation
	}

	// Animate when metric changes
	$: if (chartContainer && svg && !isAnimating) {
		drawChart(false, true); // Animate metric change
	}

	function toggleCampus(campus: string) {
		if (selectedCampuses.has(campus)) {
			selectedCampuses.delete(campus);
		} else {
			selectedCampuses.add(campus);
		}
		selectedCampuses = new Set(selectedCampuses);
	}

	function toggleLogScale() {
		isLogarithmic = !isLogarithmic;
		drawChart(false, true);
	}

	function initializeChart() {
		if (!chartContainer) return;

		// Clear previous chart
		d3.select(chartContainer).selectAll('*').remove();

		// Chart dimensions
		const margin = { top: 20, right: 140, bottom: 60, left: 100 };
		const width = 900 - margin.left - margin.right;
		const height = 550 - margin.top - margin.bottom;

		// Create SVG
		svg = d3
			.select(chartContainer)
			.append('svg')
			.attr('width', width + margin.left + margin.right)
			.attr('height', height + margin.top + margin.bottom)
			.attr('class', 'wage-chart');

		// Add gradient definitions
		const defs = svg.append('defs');

		// Create gradients for each campus
		campuses.forEach((campus, i) => {
			const gradient = defs
				.append('linearGradient')
				.attr('id', `gradient-${i}`)
				.attr('gradientUnits', 'userSpaceOnUse')
				.attr('x1', 0).attr('y1', height)
				.attr('x2', 0).attr('y2', 0);

			const color = d3.schemeCategory10[i % 10];
			gradient.append('stop')
				.attr('offset', '0%')
				.attr('stop-color', color)
				.attr('stop-opacity', 0.1);
			gradient.append('stop')
				.attr('offset', '100%')
				.attr('stop-color', color)
				.attr('stop-opacity', 0.8);
		});

		g = svg.append('g').attr('transform', `translate(${margin.left},${margin.top})`);

		// Add chart background
		g.append('rect')
			.attr('width', width)
			.attr('height', height)
			.attr('fill', 'transparent')
			.attr('stroke', 'none');
	}

	function drawChart(initialAnimation = false, metricAnimation = false) {
		if (!chartContainer || !svg || !g) return;

		isAnimating = metricAnimation;

		// Filter data for selected campuses
		const filteredData = data.filter(d => selectedCampuses.has(d.location));
		if (filteredData.length === 0) {
			// Clear chart if no data
			g.selectAll('.chart-element').remove();
			isAnimating = false;
			return;
		}

		// Chart dimensions
		const margin = { top: 20, right: 140, bottom: 60, left: 100 };
		const width = 900 - margin.left - margin.right;
		const height = 550 - margin.top - margin.bottom;

		// Scales
		const xScale = d3.scaleLinear()
			.domain(d3.extent(years) as [number, number])
			.range([0, width]);

		// Create appropriate Y scale based on logarithmic setting
		const extent = d3.extent(filteredData, d => d[metric]) as [number, number];
		const yScale = isLogarithmic
			? d3.scaleLog()
				.domain([Math.max(1, extent[0]), extent[1]])
				.range([height, 0])
				.nice()
			: d3.scaleLinear()
				.domain(extent)
				.range([height, 0])
				.nice();

		// Enhanced color scale with better colors
		const colorScale = d3.scaleOrdinal()
			.domain(campuses)
			.range([
				'#2563eb', '#dc2626', '#059669', '#d97706', '#7c3aed',
				'#db2777', '#0891b2', '#65a30d', '#dc2626', '#4338ca'
			]);

		// Smooth curve line generator
		const line = d3
			.line<AggregatedWageData>()
			.x(d => xScale(d.year))
			.y(d => yScale(Math.max(1, d[metric])))
			.curve(d3.curveCatmullRom.alpha(0.5)); // Smoother curves

		// Group data by campus
		const campusData = d3.group(filteredData, d => d.location);

		// Update or create axes with animation
		const transitionDuration = metricAnimation ? 800 : initialAnimation ? 1200 : 300;

		// X Axis
		const xAxis = g.selectAll('.x-axis').data([null]);
		const xAxisEnter = xAxis.enter()
			.append('g')
			.attr('class', 'x-axis chart-element')
			.attr('transform', `translate(0,${height})`);

		xAxisEnter.merge(xAxis)
			.transition()
			.duration(transitionDuration)
			.call(d3.axisBottom(xScale)
				.tickFormat(d3.format('d'))
				.tickSize(-height)
				.tickPadding(10))
			.selectAll('.tick line')
			.attr('stroke', '#e5e7eb')
			.attr('stroke-dasharray', '2,2');

		// Y Axis
		const yAxis = g.selectAll('.y-axis').data([null]);
		const yAxisEnter = yAxis.enter()
			.append('g')
			.attr('class', 'y-axis chart-element');

		// Custom dollar formatting function
		const formatDollars = (value: number) => {
			if (value >= 1e9) {
				return `$${(value / 1e9).toFixed(1)}B`;
			} else if (value >= 1e6) {
				return `$${(value / 1e6).toFixed(1)}M`;
			} else if (value >= 1e3) {
				return `$${(value / 1e3).toFixed(0)}K`;
			} else {
				return `$${value.toFixed(0)}`;
			}
		};

		// Use different formatting based on metric
		const yAxisFormat = metric === 'employeeCount'
			? d3.format(',')
			: formatDollars;

		yAxisEnter.merge(yAxis)
			.transition()
			.duration(transitionDuration)
			.call(d3.axisLeft(yScale)
				.tickFormat(yAxisFormat)
				.tickSize(-width)
				.tickPadding(10))
			.selectAll('.tick line')
			.attr('stroke', '#e5e7eb')
			.attr('stroke-dasharray', '2,2');

		// Style axis text
		g.selectAll('.x-axis text, .y-axis text')
			.style('fill', '#6b7280')
			.style('font-size', '12px')
			.style('font-weight', '400');

		g.selectAll('.x-axis path, .y-axis path')
			.style('stroke', '#d1d5db');

		// Draw lines with animation
		const lines = g.selectAll('.line-path')
			.data([...campusData.entries()], d => d[0]);

		lines.exit()
			.transition()
			.duration(300)
			.style('opacity', 0)
			.remove();

		const linesEnter = lines.enter()
			.append('path')
			.attr('class', 'line-path chart-element')
			.attr('fill', 'none')
			.attr('stroke-width', 3)
			.attr('stroke-linejoin', 'round')
			.attr('stroke-linecap', 'round')
			.style('opacity', 0);

		const linesUpdate = linesEnter.merge(lines)
			.attr('stroke', d => colorScale(d[0]))
			.attr('d', d => line(d[1].sort((a, b) => a.year - b.year)));

		// Sort years for chronological animation
		const sortedYears = [...years].sort((a, b) => a - b);

		if (initialAnimation) {
			// Smooth interpolated line drawing animation
			linesUpdate.each(function(d, campusIndex) {
				const path = d3.select(this);
				const campusData = d[1].sort((a, b) => a.year - b.year);
				const totalLength = (this as SVGPathElement).getTotalLength();

				// Set up for smooth drawing animation
				path
					.attr('stroke-dasharray', totalLength + ' ' + totalLength)
					.attr('stroke-dashoffset', totalLength)
					.style('opacity', 1)
					.transition()
					.delay(campusIndex * 300)
					.duration(2000)
					.ease(d3.easeCubicInOut)
					.attr('stroke-dashoffset', 0)
					.on('end', function() {
						d3.select(this).attr('stroke-dasharray', 'none');
					});
			});
		} else {
			linesUpdate
				.style('opacity', 1)
				.transition()
				.duration(transitionDuration)
				.ease(d3.easeCubicInOut)
				.attr('d', d => line(d[1].sort((a, b) => a.year - b.year)));
		}

		// Draw enhanced dots
		const allPoints = [...campusData.entries()].flatMap(([campus, values]) =>
			values.map(d => ({ ...d, campus }))
		);

		const dots = g.selectAll('.data-dot')
			.data(allPoints, d => `${d.campus}-${d.year}`);

		dots.exit()
			.transition()
			.duration(300)
			.attr('r', 0)
			.style('opacity', 0)
			.remove();

		const dotsEnter = dots.enter()
			.append('circle')
			.attr('class', 'data-dot chart-element')
			.attr('r', 0)
			.style('opacity', 0);

		dotsEnter.merge(dots)
			.attr('cx', d => xScale(d.year))
			.attr('cy', d => yScale(Math.max(1, d[metric])))
			.attr('fill', d => colorScale(d.campus))
			.attr('stroke', 'white')
			.attr('stroke-width', 2)
			.style('cursor', 'pointer')
			.transition()
			.delay(initialAnimation ? (d, i) => {
				const campusIndex = campuses.indexOf(d.campus);
				return campusIndex * 300 + 1000; // Appear after lines start drawing
			} : 0)
			.duration(initialAnimation ? 1200 : transitionDuration)
			.ease(d3.easeCubicOut)
			.attr('r', 5)
			.style('opacity', 0.95);

		// Add hover effects
		g.selectAll('.data-dot')
			.on('mouseover', function(event, d) {
				d3.select(this)
					.transition()
					.duration(150)
					.attr('r', 6)
					.style('opacity', 1);

				// Add tooltip
				const tooltip = g.append('g')
					.attr('class', 'tooltip')
					.attr('transform', `translate(${xScale(d.year)}, ${yScale(Math.max(1, d[metric]))})`);

				// Format tooltip value properly
				const tooltipValue = metric === 'employeeCount'
					? d[metric].toLocaleString()
					: formatDollars(d[metric]);

				const rect = tooltip.append('rect')
					.attr('x', -60)
					.attr('y', -35)
					.attr('width', 120)
					.attr('height', 25)
					.attr('fill', 'rgba(0,0,0,0.8)')
					.attr('rx', 4);

				tooltip.append('text')
					.attr('text-anchor', 'middle')
					.attr('y', -15)
					.attr('fill', 'white')
					.style('font-size', '12px')
					.text(tooltipValue);
			})
			.on('mouseout', function(event, d) {
				d3.select(this)
					.transition()
					.duration(150)
					.attr('r', 4)
					.style('opacity', 0.9);

				g.select('.tooltip').remove();
			});

		// Update axis labels
		const yLabel = g.selectAll('.y-label').data([getYAxisLabel(metric)]);
		const yLabelEnter = yLabel.enter()
			.append('text')
			.attr('class', 'y-label chart-element')
			.attr('transform', 'rotate(-90)')
			.attr('y', -margin.left + 20)
			.attr('x', -height / 2)
			.style('text-anchor', 'middle')
			.style('font-size', '14px')
			.style('font-weight', '500')
			.style('fill', '#374151');

		yLabelEnter.merge(yLabel)
			.transition()
			.duration(transitionDuration)
			.tween('text', function(d) {
				const i = d3.interpolateString(this.textContent, d);
				return function(t) {
					this.textContent = i(t);
				};
			});

		const xLabel = g.selectAll('.x-label').data(['Year']);
		const xLabelEnter = xLabel.enter()
			.append('text')
			.attr('class', 'x-label chart-element')
			.attr('transform', `translate(${width / 2}, ${height + margin.bottom - 10})`)
			.style('text-anchor', 'middle')
			.style('font-size', '14px')
			.style('font-weight', '500')
			.style('fill', '#374151');

		xLabelEnter.merge(xLabel).text('Year');

		// Enhanced legend
		const legend = svg.selectAll('.legend').data([null]);
		const legendEnter = legend.enter()
			.append('g')
			.attr('class', 'legend')
			.attr('transform', `translate(${width + margin.left + 20}, ${margin.top + 20})`);

		const legendItems = legendEnter.merge(legend)
			.selectAll('.legend-item')
			.data([...selectedCampuses]);

		legendItems.exit().remove();

		const legendItemsEnter = legendItems.enter()
			.append('g')
			.attr('class', 'legend-item')
			.style('cursor', 'pointer');

		legendItemsEnter.append('rect')
			.attr('width', 16)
			.attr('height', 3)
			.attr('rx', 1.5);

		legendItemsEnter.append('text')
			.attr('x', 22)
			.attr('y', 2)
			.attr('dy', '0.35em')
			.style('font-size', '13px')
			.style('font-weight', '500')
			.style('fill', '#374151');

		legendItemsEnter.merge(legendItems)
			.attr('transform', (d, i) => `translate(0, ${i * 25})`)
			.on('click', (event, d) => toggleCampus(d));

		legendItemsEnter.merge(legendItems).select('rect')
			.attr('fill', d => colorScale(d));

		legendItemsEnter.merge(legendItems).select('text')
			.text(d => d);

		if (metricAnimation) {
			setTimeout(() => {
				isAnimating = false;
			}, transitionDuration);
		}
	}

	function getYAxisLabel(metric: string): string {
		switch (metric) {
			case 'totalWages':
				return 'Total Wages ($)';
			case 'averageWage':
				return 'Average Wage ($)';
			case 'employeeCount':
				return 'Employee Count';
			default:
				return '';
		}
	}

	function formatMetricLabel(metric: string): string {
		switch (metric) {
			case 'totalWages':
				return 'Total Wages';
			case 'averageWage':
				return 'Average Wage';
			case 'employeeCount':
				return 'Employee Count';
			default:
				return metric;
		}
	}
</script>

<div class="chart-wrapper">
	<div class="controls">
		<div class="metric-selector">
			<label for="metric-select">Metric:</label>
			<select id="metric-select" bind:value={metric}>
				<option value="totalWages">Total Wages</option>
				<option value="averageWage">Average Wage</option>
				<option value="employeeCount">Employee Count</option>
			</select>
		</div>

		<div class="scale-selector">
			<label class="scale-toggle">
				<input
					type="checkbox"
					bind:checked={isLogarithmic}
					on:change={toggleLogScale}
				/>
				<span class="toggle-slider"></span>
				Logarithmic Scale
			</label>
		</div>

		<div class="campus-selector">
			<h4>Select Campuses:</h4>
			<div class="campus-checkboxes">
				{#each campuses as campus}
					<label class="campus-checkbox">
						<input
							type="checkbox"
							checked={selectedCampuses.has(campus)}
							on:change={() => toggleCampus(campus)}
						/>
						{campus}
					</label>
				{/each}
			</div>
		</div>
	</div>

	<div class="chart-title">
		<h3>UC Wage Data: {formatMetricLabel(metric)} Over Time</h3>
	</div>

	<div bind:this={chartContainer} class="chart-container"></div>
</div>

<style>
	.chart-wrapper {
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(248, 250, 252, 0.95) 100%);
		backdrop-filter: blur(20px);
		border-radius: 24px;
		padding: 2.5rem;
		box-shadow:
			0 20px 25px -5px rgba(0, 0, 0, 0.1),
			0 10px 10px -5px rgba(0, 0, 0, 0.04),
			0 0 0 1px rgba(255, 255, 255, 0.5);
		border: 1px solid rgba(229, 231, 235, 0.3);
		margin: 3rem 0;
		position: relative;
		overflow: hidden;
	}

	.chart-wrapper::before {
		content: '';
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 2px;
		background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4, #10b981);
		background-size: 400% 100%;
		animation: gradient-shift 8s ease infinite;
	}

	@keyframes gradient-shift {
		0%, 100% { background-position: 0% 50%; }
		50% { background-position: 100% 50%; }
	}

	.controls {
		display: flex;
		gap: 2rem;
		margin-bottom: 3rem;
		flex-wrap: wrap;
		align-items: flex-start;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.4);
		border-radius: 16px;
		border: 1px solid rgba(229, 231, 235, 0.2);
		backdrop-filter: blur(10px);
	}

	.metric-selector {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.metric-selector label {
		font-weight: 500;
		color: var(--text-primary);
	}

	.metric-selector select {
		padding: 0.875rem 1.25rem;
		border: 2px solid rgba(229, 231, 235, 0.3);
		border-radius: 12px;
		background: rgba(255, 255, 255, 0.8);
		backdrop-filter: blur(10px);
		color: var(--text-primary);
		cursor: pointer;
		font-weight: 600;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		min-width: 180px;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
	}

	.metric-selector select:hover {
		border-color: #3b82f6;
		background: rgba(255, 255, 255, 0.95);
		transform: translateY(-1px);
		box-shadow: 0 8px 15px -3px rgba(0, 0, 0, 0.1);
	}

	.metric-selector select:focus {
		outline: none;
		border-color: #3b82f6;
		background: rgba(255, 255, 255, 1);
		box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1), 0 8px 15px -3px rgba(0, 0, 0, 0.1);
		transform: translateY(-1px);
	}

	.scale-selector {
		display: flex;
		align-items: center;
	}

	.scale-toggle {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		cursor: pointer;
		font-weight: 500;
		color: var(--text-primary);
		user-select: none;
	}

	.scale-toggle input[type="checkbox"] {
		display: none;
	}

	.toggle-slider {
		position: relative;
		width: 52px;
		height: 28px;
		background: linear-gradient(135deg, #e5e7eb, #d1d5db);
		border-radius: 14px;
		transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.toggle-slider::before {
		content: '';
		position: absolute;
		top: 2px;
		left: 2px;
		width: 24px;
		height: 24px;
		background: linear-gradient(135deg, #ffffff, #f8fafc);
		border-radius: 50%;
		transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15), 0 1px 3px rgba(0, 0, 0, 0.1);
	}

	.scale-toggle input:checked + .toggle-slider {
		background: linear-gradient(135deg, #3b82f6, #1d4ed8);
		box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1), 0 0 0 2px rgba(59, 130, 246, 0.2);
	}

	.scale-toggle input:checked + .toggle-slider::before {
		transform: translateX(24px);
		background: linear-gradient(135deg, #ffffff, #f1f5f9);
	}

	.campus-selector h4 {
		margin: 0 0 0.75rem 0;
		color: var(--text-primary);
		font-size: 0.875rem;
		font-weight: 600;
	}

	.campus-checkboxes {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: 0.75rem;
		max-width: 600px;
	}

	.campus-checkbox {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--text-primary);
		cursor: pointer;
		padding: 0.5rem 0.75rem;
		border-radius: 8px;
		transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
		background: rgba(255, 255, 255, 0.3);
		border: 1px solid rgba(229, 231, 235, 0.2);
	}

	.campus-checkbox:hover {
		background: rgba(255, 255, 255, 0.6);
		transform: translateY(-1px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}

	.campus-checkbox input {
		cursor: pointer;
		width: 16px;
		height: 16px;
		accent-color: #3b82f6;
	}

	.chart-title {
		text-align: center;
		margin-bottom: 2rem;
		position: relative;
	}

	.chart-title h3 {
		margin: 0;
		background: linear-gradient(135deg, #1e293b, #3b82f6, #8b5cf6);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		font-size: 1.75rem;
		font-weight: 700;
		letter-spacing: -0.025em;
		line-height: 1.2;
	}

	.chart-container {
		overflow-x: auto;
		display: flex;
		justify-content: center;
	}

	.chart-container :global(svg) {
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.8) 0%, rgba(248, 250, 252, 0.9) 100%);
		border-radius: 16px;
		overflow: visible;
		backdrop-filter: blur(10px);
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.chart-container :global(.wage-chart) {
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}

	.chart-container :global(.data-dot) {
		filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
		transition: all 0.2s ease;
	}

	.chart-container :global(.data-dot:hover) {
		filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2));
	}

	.chart-container :global(.line-path) {
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
	}

	.chart-container :global(.legend-item:hover rect) {
		filter: brightness(1.1);
		transform: scale(1.05);
	}

	.chart-container :global(.legend-item:hover text) {
		fill: var(--pri);
	}

	.chart-container :global(.tooltip) {
		pointer-events: none;
	}

	.chart-container :global(.x-axis .tick line),
	.chart-container :global(.y-axis .tick line) {
		opacity: 0.3;
	}

	.chart-container :global(.domain) {
		stroke: #d1d5db;
		stroke-width: 1;
	}

	@media (max-width: 768px) {
		.chart-wrapper {
			padding: 1rem;
			margin: 1rem 0;
		}

		.controls {
			flex-direction: column;
			gap: 1rem;
		}

		.campus-checkboxes {
			grid-template-columns: repeat(auto-fit, minmax(100px, 1fr));
		}

		.chart-container :global(svg) {
			max-width: 100%;
			height: auto;
		}
	}
</style>