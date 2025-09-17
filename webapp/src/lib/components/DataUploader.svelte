<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import { ucCampuses, type UCCampus, type UploadProgress } from '$lib/types/wages';

	let selectedFile: File | null = null;
	let selectedLocation: UCCampus | '' = '';
	let selectedYear: number = new Date().getFullYear();
	let uploading = false;
	let uploadProgress: UploadProgress[] = [];
	let dragOver = false;

	// Generate year options (2010-2024)
	const yearOptions = Array.from({ length: 15 }, (_, i) => 2010 + i);

	onMount(() => {
		fetchUploadProgress();
	});

	async function fetchUploadProgress() {
		try {
			const response = await fetch('/api/upload');
			if (response.ok) {
				uploadProgress = await response.json();
			}
		} catch (error) {
			console.error('Error fetching upload progress:', error);
		}
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files[0]) {
			selectedFile = target.files[0];

			// Try to parse location and year from filename
			const filename = selectedFile.name;
			const match = filename.match(/wages_(\d{4})\.json$/);
			if (match) {
				selectedYear = parseInt(match[1]);
			}
		}
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		dragOver = false;

		if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
			selectedFile = event.dataTransfer.files[0];

			// Try to parse year from filename
			const filename = selectedFile.name;
			const match = filename.match(/wages_(\d{4})\.json$/);
			if (match) {
				selectedYear = parseInt(match[1]);
			}
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		dragOver = true;
	}

	function handleDragLeave() {
		dragOver = false;
	}

	async function uploadFile() {
		if (!selectedFile || !selectedLocation || !selectedYear) {
			alert('Please select a file, location, and year');
			return;
		}

		uploading = true;

		try {
			const formData = new FormData();
			formData.append('file', selectedFile);
			formData.append('location', selectedLocation);
			formData.append('year', selectedYear.toString());

			const response = await fetch('/api/upload', {
				method: 'POST',
				body: formData
			});

			const result = await response.json();

			if (response.ok) {
				alert(`Successfully uploaded ${result.uploadedRecords} records!`);
				selectedFile = null;
				selectedLocation = '';
				await fetchUploadProgress();
			} else {
				alert(`Upload failed: ${result.error}`);
			}
		} catch (error) {
			console.error('Upload error:', error);
			alert('Upload failed due to network error');
		} finally {
			uploading = false;
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes === 0) return '0 Bytes';
		const k = 1024;
		const sizes = ['Bytes', 'KB', 'MB', 'GB'];
		const i = Math.floor(Math.log(bytes) / Math.log(k));
		return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
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
			case 'completed': return 'text-green-600';
			case 'processing': return 'text-blue-600';
			case 'failed': return 'text-red-600';
			default: return 'text-gray-600';
		}
	}
</script>

<div class="uploader-container">
	<div class="upload-section">
		<h3>Upload Wage Data</h3>

		<div
			class="drop-zone"
			class:drag-over={dragOver}
			on:drop={handleDrop}
			on:dragover={handleDragOver}
			on:dragleave={handleDragLeave}
			role="button"
			tabindex="0"
		>
			<Icon icon="mdi:cloud-upload" class="upload-icon" />
			<p>Drop a JSON wage file here or click to select</p>
			<input
				type="file"
				accept=".json"
				on:change={handleFileSelect}
				class="file-input"
			/>
		</div>

		{#if selectedFile}
			<div class="file-info">
				<Icon icon="mdi:file-document" class="file-icon" />
				<div class="file-details">
					<span class="file-name">{selectedFile.name}</span>
					<span class="file-size">{formatFileSize(selectedFile.size)}</span>
				</div>
			</div>
		{/if}

		<div class="form-grid">
			<div class="form-group">
				<label for="location">Campus Location</label>
				<select id="location" bind:value={selectedLocation} disabled={uploading}>
					<option value="">Select Location</option>
					{#each ucCampuses as campus}
						<option value={campus}>{campus}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="year">Year</label>
				<select id="year" bind:value={selectedYear} disabled={uploading}>
					{#each yearOptions as year}
						<option value={year}>{year}</option>
					{/each}
				</select>
			</div>
		</div>

		<button
			class="upload-button"
			on:click={uploadFile}
			disabled={!selectedFile || !selectedLocation || !selectedYear || uploading}
		>
			{#if uploading}
				<Icon icon="mdi:loading" class="button-icon spinning" />
				Uploading...
			{:else}
				<Icon icon="mdi:upload" class="button-icon" />
				Upload File
			{/if}
		</button>
	</div>

	<div class="progress-section">
		<h3>Upload History</h3>

		{#if uploadProgress.length === 0}
			<p class="no-data">No uploads yet</p>
		{:else}
			<div class="progress-list">
				{#each uploadProgress as progress}
					<div class="progress-item">
						<div class="progress-header">
							<Icon
								icon={getStatusIcon(progress.status)}
								class="status-icon {getStatusColor(progress.status)}"
							/>
							<div class="progress-info">
								<span class="location-year">{progress.location} {progress.year}</span>
								<span class="status">{progress.status}</span>
							</div>
							<div class="progress-stats">
								{#if progress.totalRecords}
									<span class="records">
										{progress.uploadedRecords?.toLocaleString()} / {progress.totalRecords.toLocaleString()} records
									</span>
								{/if}
							</div>
						</div>

						{#if progress.status === 'processing' && progress.totalRecords}
							<div class="progress-bar">
								<div
									class="progress-fill"
									style="width: {((progress.uploadedRecords || 0) / progress.totalRecords) * 100}%"
								></div>
							</div>
						{/if}

						{#if progress.errorMessage}
							<div class="error-message">
								{progress.errorMessage}
							</div>
						{/if}
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	.uploader-container {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
		max-width: 1200px;
		margin: 0 auto;
	}

	.upload-section, .progress-section {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
		border: 1px solid var(--border);
	}

	h3 {
		margin: 0 0 1.5rem 0;
		color: var(--text-primary);
		font-size: 1.25rem;
		font-weight: 600;
	}

	.drop-zone {
		border: 2px dashed var(--border);
		border-radius: 0.75rem;
		padding: 3rem 2rem;
		text-align: center;
		cursor: pointer;
		transition: all 0.3s ease;
		position: relative;
		margin-bottom: 1.5rem;
	}

	.drop-zone:hover, .drop-zone.drag-over {
		border-color: var(--pri);
		background: var(--bg-secondary);
	}

	:global(.upload-icon) {
		font-size: 3rem;
		color: var(--text-secondary);
		margin-bottom: 1rem;
	}

	.drop-zone p {
		margin: 0;
		color: var(--text-secondary);
		font-size: 1rem;
	}

	.file-input {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}

	.file-info {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 1rem;
		background: var(--bg-secondary);
		border-radius: 0.5rem;
		margin-bottom: 1.5rem;
	}

	:global(.file-icon) {
		font-size: 1.5rem;
		color: var(--pri);
	}

	.file-details {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.file-name {
		font-weight: 500;
		color: var(--text-primary);
	}

	.file-size {
		font-size: 0.875rem;
		color: var(--text-secondary);
	}

	.form-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
		margin-bottom: 2rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-weight: 500;
		color: var(--text-primary);
		font-size: 0.875rem;
	}

	.form-group select {
		padding: 0.75rem;
		border: 1px solid var(--border);
		border-radius: 0.5rem;
		background: white;
		color: var(--text-primary);
		cursor: pointer;
	}

	.form-group select:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.upload-button {
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 0.875rem 1.5rem;
		background: linear-gradient(135deg, var(--pri), var(--sec));
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.upload-button:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
	}

	.upload-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	:global(.button-icon) {
		font-size: 1.125rem;
	}

	:global(.spinning) {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.no-data {
		text-align: center;
		color: var(--text-secondary);
		padding: 2rem;
		font-style: italic;
	}

	.progress-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.progress-item {
		border: 1px solid var(--border);
		border-radius: 0.5rem;
		padding: 1rem;
		background: var(--bg-secondary);
	}

	.progress-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		margin-bottom: 0.5rem;
	}

	:global(.status-icon) {
		font-size: 1.25rem;
	}

	.progress-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.location-year {
		font-weight: 500;
		color: var(--text-primary);
	}

	.status {
		font-size: 0.875rem;
		color: var(--text-secondary);
		text-transform: capitalize;
	}

	.progress-stats {
		text-align: right;
	}

	.records {
		font-size: 0.875rem;
		color: var(--text-secondary);
	}

	.progress-bar {
		height: 0.5rem;
		background: var(--border);
		border-radius: 0.25rem;
		overflow: hidden;
		margin-top: 0.5rem;
	}

	.progress-fill {
		height: 100%;
		background: linear-gradient(135deg, var(--pri), var(--sec));
		transition: width 0.3s ease;
	}

	.error-message {
		margin-top: 0.5rem;
		padding: 0.5rem;
		background: rgb(254 242 242);
		border: 1px solid rgb(252 165 165);
		border-radius: 0.25rem;
		color: rgb(153 27 27);
		font-size: 0.875rem;
	}

	@media (max-width: 768px) {
		.uploader-container {
			grid-template-columns: 1fr;
			gap: 1.5rem;
		}

		.form-grid {
			grid-template-columns: 1fr;
		}

		.upload-section, .progress-section {
			padding: 1.5rem;
		}

		.drop-zone {
			padding: 2rem 1rem;
		}
	}
</style>