<script lang="ts">
	import Container from '$lib/components/Container.svelte';
	import Section from '$lib/components/Section.svelte';
	import DataUploader from '$lib/components/DataUploader.svelte';
</script>

<svelte:head>
	<title>Data Upload - UC Wage Explorer</title>
	<meta name="description" content="Upload UC wage data files to the database" />
</svelte:head>

<Container>
	<Section
		title="Data Management"
		subtitle="Upload and manage UC wage data files"
	>
		<DataUploader />
	</Section>

	<Section title="Upload Instructions" subtitle="How to upload wage data files">
		<div class="instructions">
			<div class="instruction-card">
				<h4>File Format</h4>
				<p>Upload JSON files with the following structure:</p>
				<pre><code>{JSON.stringify({
					location: "Berkeley",
					year: 2024,
					scraped_at: "2024-01-01T00:00:00Z",
					total_records: 1000,
					records: [
						{
							id: 1,
							firstname: "John",
							lastname: "Doe",
							title: "Professor",
							basepay: "100000.00",
							overtimepay: "5000.00",
							adjustpay: "2000.00",
							grosspay: "107000.00"
						}
					]
				}, null, 2)}</code></pre>
			</div>

			<div class="instruction-card">
				<h4>Supported Locations</h4>
				<ul>
					<li>ASUCLA</li>
					<li>Berkeley</li>
					<li>Davis</li>
					<li>UC SF Law</li>
					<li>Irvine</li>
					<li>Los Angeles</li>
					<li>Merced</li>
					<li>Riverside</li>
					<li>San Diego</li>
					<li>San Francisco</li>
					<li>Santa Barbara</li>
					<li>Santa Cruz</li>
					<li>UCOP</li>
				</ul>
			</div>

			<div class="instruction-card">
				<h4>Data Processing</h4>
				<ul>
					<li>Files are processed in batches of 1,000 records</li>
					<li>Duplicate records are automatically updated</li>
					<li>Upload progress is tracked in real-time</li>
					<li>Failed uploads will show error messages</li>
				</ul>
			</div>
		</div>
	</Section>
</Container>

<style>
	.instructions {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
		gap: 2rem;
		margin-top: 2rem;
	}

	.instruction-card {
		background: var(--bg-secondary);
		border-radius: 0.75rem;
		padding: 1.5rem;
		border: 1px solid var(--border);
	}

	.instruction-card h4 {
		margin: 0 0 1rem 0;
		color: var(--pri);
		font-size: 1.125rem;
		font-weight: 600;
	}

	.instruction-card p {
		margin: 0 0 1rem 0;
		color: var(--text-secondary);
		line-height: 1.6;
	}

	.instruction-card ul {
		margin: 0;
		padding-left: 1.25rem;
		color: var(--text-secondary);
		line-height: 1.6;
	}

	.instruction-card li {
		margin-bottom: 0.5rem;
	}

	pre {
		background: white;
		border: 1px solid var(--border);
		border-radius: 0.5rem;
		padding: 1rem;
		overflow-x: auto;
		font-size: 0.875rem;
		margin: 0;
	}

	code {
		font-family: 'JetBrains Mono', 'Courier New', monospace;
		color: var(--text-primary);
	}

	@media (max-width: 768px) {
		.instructions {
			grid-template-columns: 1fr;
			gap: 1.5rem;
		}

		.instruction-card {
			padding: 1.25rem;
		}

		pre {
			font-size: 0.8rem;
		}
	}
</style>