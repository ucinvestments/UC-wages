<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { ucCampuses, type UploadProgress } from '$lib/types/wages';

	let uploadProgress: UploadProgress[] = [];
	let groupedProgress: Map<string, UploadProgress[]> = new Map();

	onMount(() => {
		fetchUploadProgress();
	});

	async function fetchUploadProgress() {
		try {
			const response = await fetch('/api/upload');
			if (response.ok) {
				uploadProgress = await response.json();
				groupProgressByCampus();
			}
		} catch (error) {
			console.error('Error fetching upload progress:', error);
		}
	}

	function groupProgressByCampus() {
		groupedProgress.clear();

		// Initialize all campuses
		ucCampuses.forEach(campus => {
			groupedProgress.set(campus, []);
		});

		// Group progress by campus
		uploadProgress.forEach(progress => {
			const campusProgress = groupedProgress.get(progress.location) || [];
			campusProgress.push(progress);
			groupedProgress.set(progress.location, campusProgress);
		});

		// Sort years within each campus
		groupedProgress.forEach((progresses, campus) => {
			progresses.sort((a, b) => b.year - a.year);
		});

		// Trigger reactivity
		groupedProgress = groupedProgress;
	}

	function getStatusIcon(status: string): string {
		switch (status) {
			case 'completed': return 'mdi:check-circle';
			case 'processing': return 'mdi:loading';
			case 'failed': return 'mdi:alert-circle';
			default: return 'mdi:clock';
		}
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'completed': return 'var(--success)';
			case 'processing': return 'var(--founder)';
			case 'failed': return 'var(--error)';
			default: return 'var(--text-secondary)';
		}
	}

	function getCampusTheme(campus: string): { primary: string; secondary: string } {
		// Campus-specific color themes
		const themes: Record<string, { primary: string; secondary: string }> = {
			'Berkeley': { primary: '#003262', secondary: '#FDB515' },
			'Los Angeles': { primary: '#2774AE', secondary: '#FFD100' },
			'San Diego': { primary: '#182B49', secondary: '#C69214' },
			'Davis': { primary: '#002855', secondary: '#FFBF00' },
			'Irvine': { primary: '#0064A4', secondary: '#FFD200' },
			'Santa Barbara': { primary: '#003660', secondary: '#FEBC11' },
			'Santa Cruz': { primary: '#003C6C', secondary: '#FDC700' },
			'Riverside': { primary: '#003DA5', secondary: '#FFB81C' },
			'Merced': { primary: '#002856', secondary: '#F1B82D' },
			'San Francisco': { primary: '#052049', secondary: '#FFC72C' },
			'UCOP': { primary: '#003B5C', secondary: '#F7931E' },
			'ASUCLA': { primary: '#2774AE', secondary: '#FFD100' },
			'UC SF Law': { primary: '#052049', secondary: '#FFC72C' }
		};

		return themes[campus] || { primary: 'var(--pri)', secondary: 'var(--sec)' };
	}

	function formatDate(date: string | undefined): string {
		if (!date) return 'N/A';
		return new Date(date).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="upload-history-container">
	<div class="header">
		<h2>Data Upload History</h2>
		<p>Overview of all wage data uploads by campus location</p>
	</div>

	{#if uploadProgress.length === 0}
		<div class="empty-state">
			<Icon icon="mdi:database-off" class="empty-icon" />
			<p>No upload history available</p>
		</div>
	{:else}
		<div class="campus-grid">
			{#each Array.from(groupedProgress.entries()) as [campus, progresses]}
				{@const theme = getCampusTheme(campus)}
				<div class="campus-card" style="--campus-primary: {theme.primary}; --campus-secondary: {theme.secondary};">
					<div class="campus-header">
						<h3>{campus}</h3>
						<span class="record-count">
							{progresses.length} {progresses.length === 1 ? 'year' : 'years'}
						</span>
					</div>

					{#if progresses.length === 0}
						<div class="no-uploads">
							<Icon icon="mdi:inbox" class="no-uploads-icon" />
							<p>No uploads for this campus</p>
						</div>
					{:else}
						<div class="year-list">
							{#each progresses as progress}
								<div class="year-item">
									<div class="year-header">
										<span class="year-label">{progress.year}</span>
										<Icon
											icon={getStatusIcon(progress.status)}
											style="color: {getStatusColor(progress.status)}; font-size: 1.25rem;"
										/>
									</div>

									<div class="year-details">
										{#if progress.totalRecords}
											<div class="detail-row">
												<Icon icon="mdi:account-group" class="detail-icon" />
												<span>{progress.totalRecords.toLocaleString()} records</span>
											</div>
										{/if}

										{#if progress.uploadedAt}
											<div class="detail-row">
												<Icon icon="mdi:calendar-clock" class="detail-icon" />
												<span>{formatDate(progress.uploadedAt)}</span>
											</div>
										{/if}

										{#if progress.status === 'processing' && progress.totalRecords}
											<div class="progress-bar">
												<div
													class="progress-fill"
													style="width: {((progress.uploadedRecords || 0) / progress.totalRecords) * 100}%"
												></div>
											</div>
											<span class="progress-text">
												{progress.uploadedRecords?.toLocaleString()} / {progress.totalRecords.toLocaleString()}
											</span>
										{/if}

										{#if progress.errorMessage}
											<div class="error-message">
												<Icon icon="mdi:alert" class="error-icon" />
												{progress.errorMessage}
											</div>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.upload-history-container {
		max-width: 1400px;
		margin: 0 auto;
	}

	.header {
		text-align: center;
		margin-bottom: 3rem;
	}

	.header h2 {
		font-size: 2rem;
		font-weight: 700;
		color: var(--pri);
		margin-bottom: 0.5rem;
	}

	.header p {
		color: var(--text-secondary);
		font-size: 1.125rem;
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		background: var(--bg-secondary);
		border-radius: 12px;
		border: 1px solid var(--border);
	}

	:global(.empty-icon) {
		font-size: 4rem;
		color: var(--text-secondary);
		margin-bottom: 1rem;
	}

	.empty-state p {
		color: var(--text-secondary);
		font-size: 1.125rem;
	}

	.campus-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
		gap: 1.5rem;
	}

	.campus-card {
		background: white;
		border-radius: 12px;
		overflow: hidden;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
		border: 1px solid var(--border);
		transition: transform 0.3s ease, box-shadow 0.3s ease;
	}

	.campus-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 10px 20px -5px rgb(0 0 0 / 0.15);
	}

	.campus-header {
		padding: 1.5rem;
		background: linear-gradient(135deg, var(--campus-primary), var(--campus-secondary));
		color: white;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.campus-header h3 {
		margin: 0;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.record-count {
		background: rgba(255, 255, 255, 0.2);
		padding: 0.25rem 0.75rem;
		border-radius: 1rem;
		font-size: 0.875rem;
		backdrop-filter: blur(10px);
	}

	.no-uploads {
		padding: 2rem;
		text-align: center;
		color: var(--text-secondary);
	}

	:global(.no-uploads-icon) {
		font-size: 2.5rem;
		margin-bottom: 0.5rem;
		opacity: 0.5;
	}

	.no-uploads p {
		margin: 0;
		font-size: 0.925rem;
	}

	.year-list {
		padding: 1rem;
		max-height: 400px;
		overflow-y: auto;
	}

	.year-list::-webkit-scrollbar {
		width: 6px;
	}

	.year-list::-webkit-scrollbar-track {
		background: var(--bg-secondary);
		border-radius: 3px;
	}

	.year-list::-webkit-scrollbar-thumb {
		background: var(--border);
		border-radius: 3px;
	}

	.year-list::-webkit-scrollbar-thumb:hover {
		background: var(--text-secondary);
	}

	.year-item {
		padding: 1rem;
		background: var(--bg-secondary);
		border-radius: 8px;
		margin-bottom: 0.75rem;
		border: 1px solid var(--border);
		transition: all 0.2s ease;
	}

	.year-item:last-child {
		margin-bottom: 0;
	}

	.year-item:hover {
		background: white;
		border-color: var(--campus-primary);
	}

	.year-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.75rem;
	}

	.year-label {
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--text-primary);
	}

	.year-details {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.detail-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	:global(.detail-icon) {
		font-size: 1rem;
		opacity: 0.7;
	}

	.progress-bar {
		height: 0.5rem;
		background: var(--border);
		border-radius: 0.25rem;
		overflow: hidden;
		margin-top: 0.25rem;
	}

	.progress-fill {
		height: 100%;
		background: linear-gradient(135deg, var(--campus-primary), var(--campus-secondary));
		transition: width 0.3s ease;
	}

	.progress-text {
		font-size: 0.75rem;
		color: var(--text-secondary);
		text-align: center;
		display: block;
		margin-top: 0.25rem;
	}

	.error-message {
		display: flex;
		align-items: flex-start;
		gap: 0.5rem;
		padding: 0.5rem;
		background: rgb(254 242 242);
		border: 1px solid rgb(252 165 165);
		border-radius: 0.25rem;
		color: rgb(153 27 27);
		font-size: 0.75rem;
		margin-top: 0.5rem;
	}

	:global(.error-icon) {
		font-size: 1rem;
		flex-shrink: 0;
		margin-top: 0.1rem;
	}

	@media (max-width: 1200px) {
		.campus-grid {
			grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
		}
	}

	@media (max-width: 768px) {
		.campus-grid {
			grid-template-columns: 1fr;
		}

		.header h2 {
			font-size: 1.5rem;
		}

		.header p {
			font-size: 1rem;
		}
	}
</style>