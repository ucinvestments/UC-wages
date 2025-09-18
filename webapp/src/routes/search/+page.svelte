<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	export let data: PageData;

	$: ({ employees, pagination, searchParams, filters } = data);

	let searchForm = {
		name: searchParams?.name || '',
		job: searchParams?.job || '',
		location: searchParams?.location || '',
		year: searchParams?.year?.toString() || '2024'
	};

	function handleSearch(event: SubmitEvent) {
		event.preventDefault();
		const params = new URLSearchParams();

		if (searchForm.name.trim()) params.set('name', searchForm.name.trim());
		if (searchForm.job.trim()) params.set('job', searchForm.job.trim());
		if (searchForm.location) params.set('location', searchForm.location);
		if (searchForm.year) params.set('year', searchForm.year);

		// Reset to page 1 for new search
		params.set('page', '1');

		goto(`/search?${params.toString()}`);
	}

	function clearFilters() {
		searchForm = { name: '', job: '', location: '', year: '' };
		goto('/search');
	}

	function goToPage(pageNum: number) {
		const params = new URLSearchParams($page.url.searchParams);
		params.set('page', pageNum.toString());
		goto(`/search?${params.toString()}`);
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 0,
			maximumFractionDigits: 0
		}).format(amount);
	}

	// Auto-search for 2024 data on page load if no search params
	onMount(() => {
		const hasSearchParams = $page.url.searchParams.toString();
		if (!hasSearchParams) {
			const params = new URLSearchParams();
			params.set('year', '2024');
			params.set('page', '1');
			goto(`/search?${params.toString()}`);
		}
	});
</script>

<svelte:head>
	<title>Search UC Employees - UC Wage Explorer</title>
	<meta name="description" content="Search University of California employee wage records by name, job title, campus, and year." />
</svelte:head>

<div class="search-page">
	<div class="container">
		<div class="page-header">
			<h1>Search UC Employees</h1>
			<p>Search through University of California employee wage records</p>
		</div>

		<!-- Search Form -->
		<div class="search-form-container">
			<form onsubmit={handleSearch} class="search-form">
				<div class="search-row">
					<div class="search-field">
						<label for="name">Employee Name</label>
						<input
							id="name"
							type="text"
							bind:value={searchForm.name}
							placeholder="Enter name..."
							class="search-input"
						/>
					</div>
					<div class="search-field">
						<label for="job">Job Title</label>
						<input
							id="job"
							type="text"
							bind:value={searchForm.job}
							placeholder="Enter job title..."
							class="search-input"
						/>
					</div>
				</div>

				<div class="search-row">
					<div class="search-field">
						<label for="location">Campus</label>
						<select id="location" bind:value={searchForm.location} class="search-select">
							<option value="">All Campuses</option>
							{#each filters?.locations || [] as location}
								<option value={location}>{location}</option>
							{/each}
						</select>
					</div>
					<div class="search-field">
						<label for="year">Year</label>
						<select id="year" bind:value={searchForm.year} class="search-select">
							<option value="">All Years</option>
							{#each filters?.years || [] as year}
								<option value={year.toString()}>{year}</option>
							{/each}
						</select>
					</div>
				</div>

				<div class="search-actions">
					<button type="submit" class="search-button primary">
						Search
					</button>
					<button type="button" class="search-button secondary" onclick={clearFilters}>
						Clear Filters
					</button>
				</div>
			</form>
		</div>

		<!-- Results Summary -->
		{#if pagination?.totalItems > 0}
			<div class="results-summary">
				<p>
					Showing {((pagination.currentPage - 1) * pagination.itemsPerPage) + 1}-{Math.min(pagination.currentPage * pagination.itemsPerPage, pagination.totalItems)}
					of {pagination.totalItems.toLocaleString()} employees
				</p>
			</div>
		{:else}
			<div class="results-summary">
				<p>No employees found matching your criteria.</p>
			</div>
		{/if}

		<!-- Results Table -->
		{#if employees?.length > 0}
			<div class="results-container">
				<div class="results-table-container">
					<table class="results-table">
						<thead>
							<tr>
								<th>Name</th>
								<th>Job Title</th>
								<th>Campus</th>
								<th>Year</th>
								<th class="currency">Gross Pay</th>
								<th class="currency">Base Pay</th>
								<th class="currency">Overtime</th>
								<th class="currency">Other Pay</th>
							</tr>
						</thead>
						<tbody>
							{#each employees as employee}
								<tr class="employee-row">
									<td class="employee-name">{employee.name}</td>
									<td class="job-title">{employee.jobtitle}</td>
									<td class="location">{employee.location}</td>
									<td class="year">{employee.year}</td>
									<td class="currency gross-pay">{formatCurrency(employee.grosspay)}</td>
									<td class="currency">{formatCurrency(employee.basePay)}</td>
									<td class="currency">{formatCurrency(employee.overtimePay)}</td>
									<td class="currency">{formatCurrency(employee.otherPay)}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>

			<!-- Pagination -->
			{#if pagination?.totalPages > 1}
				<div class="pagination">
					<div class="pagination-info">
						Page {pagination.currentPage} of {pagination.totalPages}
					</div>
					<div class="pagination-controls">
						{#if pagination?.hasPrev}
							<button class="pagination-button" onclick={() => goToPage(1)}>
								First
							</button>
							<button class="pagination-button" onclick={() => goToPage(pagination.currentPage - 1)}>
								Previous
							</button>
						{/if}

						<!-- Page numbers -->
						{#each Array.from({ length: Math.min(5, pagination?.totalPages || 0) }, (_, i) => {
							const start = Math.max(1, (pagination?.currentPage || 1) - 2);
							return start + i;
						}).filter(p => p <= (pagination?.totalPages || 0)) as pageNum}
							<button
								class="pagination-button {pageNum === pagination.currentPage ? 'current' : ''}"
								onclick={() => goToPage(pageNum)}
							>
								{pageNum}
							</button>
						{/each}

						{#if pagination?.hasNext}
							<button class="pagination-button" onclick={() => goToPage(pagination.currentPage + 1)}>
								Next
							</button>
							<button class="pagination-button" onclick={() => goToPage(pagination.totalPages)}>
								Last
							</button>
						{/if}
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>

<style>
	.search-page {
		min-height: 100vh;
		background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
		padding: 2rem 0;
	}

	.container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 0 2rem;
	}

	.page-header {
		text-align: center;
		margin-bottom: 3rem;
	}

	.page-header h1 {
		font-size: 3rem;
		font-weight: 800;
		background: linear-gradient(135deg, #1e293b, #3b82f6);
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		margin-bottom: 1rem;
	}

	.page-header p {
		font-size: 1.2rem;
		color: #64748b;
	}

	.search-form-container {
		background: white;
		border-radius: 16px;
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
		padding: 2rem;
		margin-bottom: 2rem;
	}

	.search-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.search-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}

	.search-field {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.search-field label {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
	}

	.search-input,
	.search-select {
		padding: 0.75rem;
		border: 2px solid #e5e7eb;
		border-radius: 8px;
		font-size: 1rem;
		transition: border-color 0.2s, box-shadow 0.2s;
	}

	.search-input:focus,
	.search-select:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.search-actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
		margin-top: 1rem;
	}

	.search-button {
		padding: 0.75rem 2rem;
		border-radius: 8px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		border: 2px solid transparent;
	}

	.search-button.primary {
		background: linear-gradient(135deg, #3b82f6, #1d4ed8);
		color: white;
	}

	.search-button.primary:hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 20px rgba(59, 130, 246, 0.3);
	}

	.search-button.secondary {
		background: white;
		color: #6b7280;
		border-color: #e5e7eb;
	}

	.search-button.secondary:hover {
		border-color: #d1d5db;
		background: #f9fafb;
	}

	.results-summary {
		background: white;
		border-radius: 8px;
		padding: 1rem;
		margin-bottom: 1rem;
		color: #6b7280;
		font-weight: 500;
	}

	.results-container {
		background: white;
		border-radius: 16px;
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
		margin-bottom: 2rem;
		overflow: hidden;
	}

	.results-table-container {
		overflow-x: auto;
	}

	.results-table {
		width: 100%;
		border-collapse: collapse;
	}

	.results-table th {
		background: #f8fafc;
		padding: 1rem;
		text-align: left;
		font-weight: 600;
		color: #374151;
		border-bottom: 2px solid #e5e7eb;
		white-space: nowrap;
	}

	.results-table th.currency {
		text-align: right;
	}

	.results-table td {
		padding: 1rem;
		border-bottom: 1px solid #f1f5f9;
	}

	.results-table td.currency {
		text-align: right;
		font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
	}

	.employee-row:hover {
		background: #f8fafc;
	}

	.employee-name {
		font-weight: 600;
		color: #1f2937;
	}

	.job-title {
		color: #4b5563;
		max-width: 200px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.gross-pay {
		font-weight: 600;
		color: #059669;
	}

	.pagination {
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: white;
		border-radius: 8px;
		padding: 1rem 2rem;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
	}

	.pagination-info {
		color: #6b7280;
		font-weight: 500;
	}

	.pagination-controls {
		display: flex;
		gap: 0.5rem;
	}

	.pagination-button {
		padding: 0.5rem 1rem;
		border: 1px solid #e5e7eb;
		background: white;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
		font-weight: 500;
	}

	.pagination-button:hover {
		border-color: #3b82f6;
		background: #eff6ff;
	}

	.pagination-button.current {
		background: #3b82f6;
		color: white;
		border-color: #3b82f6;
	}

	@media (max-width: 768px) {
		.search-row {
			grid-template-columns: 1fr;
		}

		.search-actions {
			flex-direction: column;
		}

		.pagination {
			flex-direction: column;
			gap: 1rem;
		}

		.pagination-controls {
			flex-wrap: wrap;
			justify-content: center;
		}

		.page-header h1 {
			font-size: 2rem;
		}

		.results-table {
			font-size: 0.875rem;
		}

		.results-table th,
		.results-table td {
			padding: 0.75rem 0.5rem;
		}
	}
</style>