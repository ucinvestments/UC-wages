<script lang="ts">
		import { goto } from '$app/navigation';
		import { page } from '$app/stores';
		import { navigating } from '$app/stores';
		import Icon from '@iconify/svelte';
		import type { PageData } from './$types';

	export let data: PageData;

	$: ({ employees, pagination, searchParams, filters } = data);
	$: isLoading = !!$navigating;

		type SortableColumn =
			| 'name'
			| 'jobtitle'
			| 'location'
			| 'year'
			| 'grosspay'
			| 'basePay'
			| 'overtimePay'
			| 'otherPay';

		type Employee = {
			name: string;
			jobtitle: string | null;
			location: string;
			year: number;
			grosspay: number;
			basePay: number;
			overtimePay: number;
			otherPay: number;
		};

		let selectedEmployee: Employee | null = null;

		const sortableColumns: { key: SortableColumn; label: string; numeric?: boolean }[] = [
			{ key: 'name', label: 'Name' },
			{ key: 'grosspay', label: 'Gross Pay', numeric: true },
			{ key: 'jobtitle', label: 'Job Title' },
			{ key: 'location', label: 'Campus' },
			{ key: 'year', label: 'Year' },
			{ key: 'basePay', label: 'Base Pay', numeric: true },
			{ key: 'overtimePay', label: 'Overtime', numeric: true },
			{ key: 'otherPay', label: 'Other Pay', numeric: true }
		];

		let searchForm = {
			name: searchParams?.name || '',
			job: searchParams?.job || '',
			location: searchParams?.location || '',
			year: searchParams?.year ? searchParams.year.toString() : ''
		};

		$: currentSort = (searchParams?.sort as SortableColumn) || 'grosspay';
		$: currentDirection = (searchParams?.direction as 'asc' | 'desc') || 'desc';

		function appendSortParams(params: URLSearchParams) {
			const sort = $page.url.searchParams.get('sort');
			const direction = $page.url.searchParams.get('direction');
			if (sort) params.set('sort', sort);
			if (direction) params.set('direction', direction);
		}

	function handleSearch(event: SubmitEvent) {
		event.preventDefault();
		const params = new URLSearchParams();

		if (searchForm.name.trim()) params.set('name', searchForm.name.trim());
		if (searchForm.job.trim()) params.set('job', searchForm.job.trim());
		if (searchForm.location) params.set('location', searchForm.location);
			if (searchForm.year) params.set('year', searchForm.year);

			appendSortParams(params);

		// Reset to page 1 for new search
		params.set('page', '1');

		goto(`/search?${params.toString()}`);
	}

	function clearFilters() {
			searchForm = { name: '', job: '', location: '', year: '' };
			const params = new URLSearchParams();
			appendSortParams(params);
			goto(params.toString() ? `/search?${params.toString()}` : '/search');
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

		function toggleSort(column: SortableColumn) {
			const params = new URLSearchParams($page.url.searchParams);
			const activeColumn = params.get('sort') as SortableColumn | null;
			const activeDirection = (params.get('direction') as 'asc' | 'desc') || 'desc';

			if (activeColumn === column) {
				params.set('direction', activeDirection === 'asc' ? 'desc' : 'asc');
			} else {
				params.set('sort', column);
				params.set('direction', column === 'name' || column === 'jobtitle' || column === 'location' ? 'asc' : 'desc');
			}

			params.set('page', '1');
			goto(`/search?${params.toString()}`);
		}

		function ariaSort(column: SortableColumn) {
			if (currentSort !== column) return 'none';
			return currentDirection === 'asc' ? 'ascending' : 'descending';
		}

		function openModal(employee: Employee) {
			selectedEmployee = employee;
		}

		function closeModal() {
			selectedEmployee = null;
		}

		function handleModalKeydown(event: KeyboardEvent) {
			if (event.key === 'Escape') closeModal();
		}

		function slugify(value: string): string {
			return value
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/(^-|-$)/g, '') || 'employee';
		}

		function triggerDownload(content: string, filename: string, mime: string) {
			const blob = new Blob([content], { type: mime });
			const url = URL.createObjectURL(blob);
			const link = document.createElement('a');
			link.href = url;
			link.download = filename;
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			URL.revokeObjectURL(url);
		}

		function downloadJSON(employee: Employee) {
			const filename = `${slugify(employee.name)}-${employee.year}.json`;
			triggerDownload(JSON.stringify(employee, null, 2), filename, 'application/json');
		}

		function downloadCSV(employee: Employee) {
			const headers = ['name', 'jobtitle', 'location', 'year', 'grosspay', 'basePay', 'overtimePay', 'otherPay'];
			const escape = (v: string | number) => {
				const s = String(v ?? '');
				return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
			};
			const row = headers.map(h => escape((employee as Record<string, string | number>)[h]));
			const csv = `${headers.join(',')}\n${row.join(',')}\n`;
			const filename = `${slugify(employee.name)}-${employee.year}.csv`;
			triggerDownload(csv, filename, 'text/csv');
		}
	</script>

	<svelte:window onkeydown={handleModalKeydown} />

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

		<!-- Loading Indicator -->
		{#if isLoading}
			<div class="loading-container">
				<div class="loading-spinner"></div>
				<p class="loading-text">Searching employee data...</p>
			</div>
		{:else}
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
								{#each sortableColumns as column (column.key)}
									<th class={column.numeric ? 'currency' : ''} aria-sort={ariaSort(column.key)}>
										<button
											type="button"
											class="sort-button"
											onclick={() => toggleSort(column.key)}
											aria-label={`Sort by ${column.label}`}
										>
											<span>{column.label}</span>
											<span class={`sort-icon ${currentSort === column.key ? 'active' : ''}`}>
												{#if currentSort === column.key}
													<Icon
														icon={currentDirection === 'asc'
															? 'mdi:arrow-up'
															: 'mdi:arrow-down'}
														class="sort-arrow"
													/>
												{:else}
													<Icon icon="mdi:unfold-more-horizontal" class="sort-arrow" />
												{/if}
											</span>
										</button>
									</th>
								{/each}
							</tr>
						</thead>
						<tbody>
							{#each employees as employee}
								<tr
									class="employee-row"
									role="button"
									tabindex="0"
									onclick={() => openModal(employee)}
									onkeydown={(e) => {
										if (e.key === 'Enter' || e.key === ' ') {
											e.preventDefault();
											openModal(employee);
										}
									}}
								>
									<td class="employee-name">{employee.name}</td>
									<td class="currency gross-pay">{formatCurrency(employee.grosspay)}</td>
									<td class="job-title">{employee.jobtitle}</td>
									<td class="location">{employee.location}</td>
									<td class="year">{employee.year}</td>
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
					<div class="pagination bottom">
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
		{/if}
	</div>
</div>

{#if selectedEmployee}
	{@const employee = selectedEmployee}
	<div
		class="modal-backdrop"
		role="button"
		tabindex="0"
		onclick={closeModal}
		onkeydown={(e) => {
			if (e.key === 'Enter' || e.key === ' ') {
				e.preventDefault();
				closeModal();
			}
		}}
		aria-label="Close modal"
	>
		<div
			class="modal"
			role="dialog"
			aria-modal="true"
			aria-labelledby="modal-title"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
		>
			<div class="modal-header">
				<h2 id="modal-title">{employee.name}</h2>
				<button type="button" class="modal-close" onclick={closeModal} aria-label="Close">
					<Icon icon="mdi:close" />
				</button>
			</div>
			<div class="modal-body">
				<dl class="modal-details">
					<div class="modal-row">
						<dt>Name</dt>
						<dd>{employee.name}</dd>
					</div>
					<div class="modal-row">
						<dt>Job Title</dt>
						<dd>{employee.jobtitle}</dd>
					</div>
					<div class="modal-row">
						<dt>Campus</dt>
						<dd>{employee.location}</dd>
					</div>
					<div class="modal-row">
						<dt>Year</dt>
						<dd>{employee.year}</dd>
					</div>
					<div class="modal-row">
						<dt>Gross Pay</dt>
						<dd class="currency highlight">{formatCurrency(employee.grosspay)}</dd>
					</div>
					<div class="modal-row">
						<dt>Base Pay</dt>
						<dd class="currency">{formatCurrency(employee.basePay)}</dd>
					</div>
					<div class="modal-row">
						<dt>Overtime Pay</dt>
						<dd class="currency">{formatCurrency(employee.overtimePay)}</dd>
					</div>
					<div class="modal-row">
						<dt>Other Pay</dt>
						<dd class="currency">{formatCurrency(employee.otherPay)}</dd>
					</div>
				</dl>
			</div>
			<div class="modal-footer">
				<button type="button" class="modal-button primary" onclick={() => downloadJSON(employee)}>
					<Icon icon="mdi:code-json" />
					Download JSON
				</button>
				<button type="button" class="modal-button primary" onclick={() => downloadCSV(employee)}>
					<Icon icon="mdi:file-delimited" />
					Download CSV
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.search-page {
		min-height: 100vh;
		background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
		padding: 2rem 0;
		overflow-x: hidden;
	}

	.container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 0 2rem;
		width: 100%;
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
		max-width: 100%;
		min-width: 0;
	}

	.results-table-container {
		overflow-x: auto;
		-webkit-overflow-scrolling: touch;
		max-width: 100%;
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

	.sort-button {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.25rem;
		width: 100%;
		background: none;
		border: none;
		font: inherit;
		color: inherit;
		cursor: pointer;
		padding: 0;
	}

	.sort-button:hover span:first-child {
		color: #1d4ed8;
	}

	.results-table th.currency .sort-button {
		justify-content: flex-end;
	}

	.sort-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 1.5rem;
		height: 1.5rem;
		color: #94a3b8;
		transition: color 0.2s ease;
	}

	.sort-icon.active {
		color: #1d4ed8;
	}

	:global(.sort-arrow) {
		font-size: 1rem;
	}

	.results-table td {
		padding: 1rem;
		border-bottom: 1px solid #f1f5f9;
	}

	.results-table td.currency {
		text-align: right;
		font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
	}

	.employee-row {
		cursor: pointer;
	}

	.employee-row:hover {
		background: #f8fafc;
	}

	.employee-row:focus-visible {
		outline: 2px solid #3b82f6;
		outline-offset: -2px;
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

	.loading-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 2rem;
		background: white;
		border-radius: 16px;
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
		margin-bottom: 2rem;
	}

	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #e5e7eb;
		border-top: 4px solid #3b82f6;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	.loading-text {
		color: #6b7280;
		font-size: 1rem;
		font-weight: 500;
		margin: 0;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
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

	.pagination.top {
		margin-bottom: 1rem;
	}

	.pagination.bottom {
		margin-top: 1rem;
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

	.modal-backdrop {
		position: fixed;
		inset: 0;
		background: rgba(15, 23, 42, 0.55);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1.5rem;
		z-index: 1000;
		animation: fadeIn 0.15s ease-out;
		cursor: default;
	}

	.modal {
		background: white;
		border-radius: 16px;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.35);
		width: 100%;
		max-width: 520px;
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		animation: modalIn 0.2s ease-out;
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1.25rem 1.5rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.modal-header h2 {
		font-size: 1.25rem;
		font-weight: 700;
		color: #1f2937;
		margin: 0;
		line-height: 1.3;
	}

	.modal-close {
		background: none;
		border: none;
		padding: 0.375rem;
		cursor: pointer;
		border-radius: 0.5rem;
		color: #6b7280;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.25rem;
		transition: background 0.15s, color 0.15s;
	}

	.modal-close:hover {
		background: #f3f4f6;
		color: #111827;
	}

	.modal-body {
		padding: 1.25rem 1.5rem;
		overflow-y: auto;
	}

	.modal-details {
		display: flex;
		flex-direction: column;
		gap: 0;
		margin: 0;
	}

	.modal-row {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: 1rem;
		padding: 0.75rem 0;
		border-bottom: 1px solid #f1f5f9;
	}

	.modal-row:last-child {
		border-bottom: none;
	}

	.modal-row dt {
		font-size: 0.8125rem;
		font-weight: 600;
		color: #6b7280;
		text-transform: uppercase;
		letter-spacing: 0.03em;
		flex-shrink: 0;
	}

	.modal-row dd {
		margin: 0;
		color: #1f2937;
		font-weight: 500;
		text-align: right;
		word-break: break-word;
	}

	.modal-row dd.currency {
		font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', monospace;
	}

	.modal-row dd.highlight {
		color: #059669;
		font-weight: 700;
		font-size: 1.125rem;
	}

	.modal-footer {
		display: flex;
		gap: 0.75rem;
		padding: 1rem 1.5rem;
		border-top: 1px solid #e5e7eb;
		background: #f9fafb;
	}

	.modal-button {
		flex: 1;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.625rem 1rem;
		border-radius: 8px;
		font-weight: 600;
		font-size: 0.9rem;
		cursor: pointer;
		border: 2px solid transparent;
		transition: all 0.15s;
	}

	.modal-button.primary {
		background: linear-gradient(135deg, #3b82f6, #1d4ed8);
		color: white;
	}

	.modal-button.primary:hover {
		transform: translateY(-1px);
		box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
	}

	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	@keyframes modalIn {
		from { opacity: 0; transform: translateY(10px) scale(0.98); }
		to { opacity: 1; transform: translateY(0) scale(1); }
	}

	@media (max-width: 768px) {
		.search-page {
			padding: 1rem 0;
		}

		.container {
			padding: 0 1rem;
		}

		.page-header {
			margin-bottom: 1.5rem;
		}

		.page-header h1 {
			font-size: 1.75rem;
			margin-bottom: 0.5rem;
		}

		.page-header p {
			font-size: 1rem;
		}

		.search-form-container {
			padding: 1.25rem;
			border-radius: 12px;
		}

		.search-form {
			gap: 1rem;
		}

		.search-row {
			grid-template-columns: 1fr;
			gap: 1rem;
		}

		.search-input,
		.search-select {
			padding: 0.625rem 0.75rem;
			font-size: 0.95rem;
		}

		.search-actions {
			flex-direction: column;
			gap: 0.75rem;
		}

		.search-button {
			width: 100%;
			padding: 0.75rem 1rem;
		}

		.results-container {
			border-radius: 12px;
			margin-bottom: 1rem;
		}

		.pagination {
			flex-direction: column;
			gap: 0.75rem;
			padding: 0.75rem 1rem;
		}

		.pagination-controls {
			flex-wrap: wrap;
			justify-content: center;
		}

		.pagination-button {
			padding: 0.4rem 0.75rem;
			font-size: 0.875rem;
		}

		.results-table {
			font-size: 0.8125rem;
		}

		.results-table th,
		.results-table td {
			padding: 0.625rem 0.5rem;
		}

		.job-title {
			max-width: 140px;
		}

		.modal-backdrop {
			padding: 0.75rem;
		}

		.modal-header {
			padding: 1rem 1.25rem;
		}

		.modal-body,
		.modal-footer {
			padding: 1rem 1.25rem;
		}

		.modal-footer {
			flex-direction: column;
		}

		.modal-row {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.25rem;
		}

		.modal-row dd {
			text-align: left;
		}
	}

	@media (max-width: 480px) {
		.container {
			padding: 0 0.75rem;
		}

		.page-header h1 {
			font-size: 1.5rem;
		}

		.search-form-container {
			padding: 1rem;
		}
	}
</style>
