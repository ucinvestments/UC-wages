<script lang="ts">
	import DataUploader from '$lib/components/DataUploader.svelte';
	import Icon from '@iconify/svelte';
</script>

<svelte:head>
	<title>Data Overview - UC Wage Explorer</title>
	<meta name="description" content="View the status of UC wage data uploads and available datasets" />
</svelte:head>

<div class="page-container">
	<section class="hero-section">
		<h1 class="page-title">Data Overview</h1>
		<p class="hero-subtitle">
			Comprehensive wage data from all University of California campuses
		</p>
	</section>

	<div class="data-summary">
		<div class="summary-card">
			<Icon icon="mdi:school" class="summary-icon" />
			<div class="summary-content">
				<span class="summary-value">13</span>
				<span class="summary-label">UC Campuses</span>
			</div>
		</div>

		<div class="summary-card">
			<Icon icon="mdi:calendar-range" class="summary-icon" />
			<div class="summary-content">
				<span class="summary-value">15</span>
				<span class="summary-label">Years of Data</span>
			</div>
		</div>

		<div class="summary-card">
			<Icon icon="mdi:database" class="summary-icon" />
			<div class="summary-content">
				<span class="summary-value">2M+</span>
				<span class="summary-label">Total Records</span>
			</div>
		</div>

		<div class="summary-card">
			<Icon icon="mdi:update" class="summary-icon" />
			<div class="summary-content">
				<span class="summary-value">2024</span>
				<span class="summary-label">Latest Update</span>
			</div>
		</div>
	</div>

	<section class="data-section">
		<DataUploader />
	</section>

	<section class="info-section">
		<h2 class="section-title">About the Dataset</h2>

		<div class="info-grid">
			<div class="info-card">
				<Icon icon="mdi:source-branch" class="info-icon" />
				<h3>Data Source</h3>
				<p>
					All data is sourced directly from the UC Annual Wage website, which publishes
					employee compensation information as mandated by California state law.
				</p>
			</div>

			<div class="info-card">
				<Icon icon="mdi:file-document-outline" class="info-icon" />
				<h3>Data Fields</h3>
				<ul>
					<li>Employee Name (anonymized where applicable)</li>
					<li>Job Title</li>
					<li>Campus Location</li>
					<li>Base Pay</li>
					<li>Overtime Pay</li>
					<li>Other Pay & Adjustments</li>
					<li>Gross Pay Total</li>
				</ul>
			</div>

			<div class="info-card">
				<Icon icon="mdi:refresh" class="info-icon" />
				<h3>Update Schedule</h3>
				<p>
					Data is updated annually when new information is released by the UC system,
					typically in the first quarter of each year for the previous year's data.
				</p>
			</div>

			<div class="info-card">
				<Icon icon="mdi:shield-check" class="info-icon" />
				<h3>Data Privacy</h3>
				<p>
					All displayed data is publicly available. Names may be redacted for certain
					employees as per UC's privacy policies and state regulations.
				</p>
			</div>
		</div>
	</section>

	<section class="download-section">
		<div class="download-card">
			<Icon icon="mdi:download" class="download-icon" />
			<div class="download-content">
				<h3>Access Raw Data</h3>
				<p>
					The complete dataset is available in our GitHub repository. You can download
					individual campus data files or clone the entire repository for analysis.
				</p>
				<a href="https://github.com/ucinvestments/UC-wages/tree/main/data" target="_blank" rel="noopener noreferrer" class="download-button">
					<Icon icon="mdi:github" />
					View on GitHub
				</a>
			</div>
		</div>
	</section>
</div>

<style>
	:global(body) {
		overflow-x: hidden;
	}

	.page-container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	.hero-section {
		padding: 1.5rem 1rem;
		background: linear-gradient(135deg, var(--pri) 0%, var(--founder) 100%);
		position: relative;
		overflow: hidden;
		margin-bottom: 3rem;
		border-radius: 16px;
	}

	.hero-section::before {
		content: "";
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
	}

	.page-title {
		font-family: "Space Grotesk", sans-serif;
		font-size: 2.25rem;
		font-weight: 700;
		color: white;
		margin-bottom: 1rem;
		letter-spacing: -0.02em;
		text-align: center;
		position: relative;
		z-index: 1;
	}

	.hero-subtitle {
		font-size: 1.25rem;
		color: white;
		opacity: 0.9;
		text-align: center;
		position: relative;
		z-index: 1;
		max-width: 600px;
		margin: 0 auto;
	}

	.data-summary {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 1.5rem;
		margin-bottom: 3rem;
	}

	.summary-card {
		background: white;
		border-radius: 12px;
		padding: 1.5rem;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
		border: 1px solid var(--border);
		display: flex;
		align-items: center;
		gap: 1.5rem;
		transition: transform 0.3s ease, box-shadow 0.3s ease;
	}

	.summary-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 20px -5px rgb(0 0 0 / 0.15);
	}

	:global(.summary-icon) {
		font-size: 2.5rem;
		color: var(--founder);
	}

	.summary-content {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.summary-value {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--pri);
	}

	.summary-label {
		font-size: 0.925rem;
		color: var(--text-secondary);
	}

	.data-section {
		margin-bottom: 4rem;
	}

	.info-section {
		margin-bottom: 4rem;
	}

	.section-title {
		font-size: 2rem;
		font-weight: 700;
		text-align: center;
		margin-bottom: 2rem;
		color: var(--pri);
	}

	.info-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 2rem;
	}

	.info-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
		border: 1px solid var(--border);
		transition: transform 0.3s ease, box-shadow 0.3s ease;
	}

	.info-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 10px 20px -5px rgb(0 0 0 / 0.15);
	}

	:global(.info-icon) {
		font-size: 2rem;
		color: var(--founder);
		margin-bottom: 1rem;
	}

	.info-card h3 {
		margin: 0 0 1rem 0;
		color: var(--pri);
		font-size: 1.25rem;
	}

	.info-card p {
		color: var(--text-secondary);
		line-height: 1.6;
		margin: 0;
	}

	.info-card ul {
		list-style: none;
		padding: 0;
		margin: 0;
		color: var(--text-secondary);
		line-height: 1.8;
	}

	.info-card li {
		padding-left: 1.5rem;
		position: relative;
		margin-bottom: 0.5rem;
	}

	.info-card li::before {
		content: "•";
		position: absolute;
		left: 0;
		color: var(--founder);
		font-weight: bold;
	}

	.download-section {
		margin-top: 4rem;
	}

	.download-card {
		background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
		border-radius: 16px;
		padding: 3rem;
		display: flex;
		align-items: center;
		gap: 2rem;
		box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
	}

	:global(.download-icon) {
		font-size: 4rem;
		color: var(--founder);
	}

	.download-content {
		flex: 1;
	}

	.download-content h3 {
		margin: 0 0 1rem 0;
		color: var(--pri);
		font-size: 1.5rem;
	}

	.download-content p {
		color: var(--text-secondary);
		line-height: 1.6;
		margin-bottom: 1.5rem;
	}

	.download-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.875rem 1.75rem;
		background: var(--pri);
		color: white;
		text-decoration: none;
		border-radius: 8px;
		font-weight: 500;
		transition: all 0.3s ease;
	}

	.download-button:hover {
		background: var(--founder);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	@media (max-width: 768px) {
		.page-container {
			padding: 1rem;
		}

		.page-title {
			font-size: 1.75rem;
		}

		.hero-subtitle {
			font-size: 1rem;
		}

		.data-summary {
			grid-template-columns: repeat(2, 1fr);
			gap: 1rem;
		}

		.summary-card {
			flex-direction: column;
			text-align: center;
			padding: 1.25rem;
		}

		.info-grid {
			grid-template-columns: 1fr;
		}

		.download-card {
			flex-direction: column;
			text-align: center;
			padding: 2rem;
		}
	}

	@media (max-width: 480px) {
		.data-summary {
			grid-template-columns: 1fr;
		}
	}
</style>