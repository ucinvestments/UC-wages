<script lang="ts">
	import Icon from '@iconify/svelte';
</script>

<svelte:head>
	<title>Methodology - UC Wage Explorer</title>
	<meta name="description" content="Learn about the data collection and processing methodology used by UC Wage Explorer" />
</svelte:head>

<div class="page-container">
	<section class="hero-section">
		<h1 class="page-title">Data Collection Methodology</h1>
		<p class="hero-subtitle">
			Transparent and systematic approach to gathering University of California wage data
		</p>
	</section>

	<section class="content-section">
		<h2 class="section-title">Data Collection Process</h2>

		<div class="process-grid">
			<div class="process-card">
				<div class="process-number">1</div>
				<h3>Source Identification</h3>
				<p>
					Data is collected directly from the official UC Annual Wage website
					(ucannualwage.ucop.edu), which publishes employee compensation data as
					mandated by California state law. This ensures all data is publicly
					available and legally accessible.
				</p>
			</div>

			<div class="process-card">
				<div class="process-number">2</div>
				<h3>Automated Collection</h3>
				<p>
					We use a high-performance Go-based scraper with concurrent workers to
					efficiently collect data. The scraper respects rate limits and implements
					retry logic to ensure reliable data collection without overloading the source.
				</p>
			</div>

			<div class="process-card">
				<div class="process-number">3</div>
				<h3>Data Validation</h3>
				<p>
					Each data point is validated for completeness and accuracy. Records are
					checked for required fields including employee name (anonymized where necessary),
					title, location, year, and various pay components.
				</p>
			</div>

			<div class="process-card">
				<div class="process-number">4</div>
				<h3>Storage & Organization</h3>
				<p>
					Data is stored in a structured JSON format, organized by campus location
					and year. This hierarchical structure enables efficient querying and
					maintains data integrity across years of historical information.
				</p>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Technical Implementation</h2>

		<div class="tech-details">
			<div class="detail-card">
				<h3><Icon icon="mdi:code-braces" /> Scraper Architecture</h3>
				<div class="detail-content">
					<h4>Core Components:</h4>
					<ul>
						<li><strong>Language:</strong> Go (Golang) for high performance and concurrency</li>
						<li><strong>Worker Pool:</strong> Configurable concurrent workers (default: 3-10)</li>
						<li><strong>Rate Limiting:</strong> Built-in delays between requests (default: 1 second)</li>
						<li><strong>Retry Logic:</strong> Automatic retry on failures with exponential backoff</li>
						<li><strong>Progress Tracking:</strong> Resume capability for interrupted scraping sessions</li>
					</ul>

					<h4>API Interaction:</h4>
					<pre><code>{`POST https://ucannualwage.ucop.edu/wage/search
Content-Type: application/json

{
  "op": "search",
  "page": 1,
  "rows": 100,
  "year": "2024",
  "location": "Berkeley"
}`}</code></pre>
				</div>
			</div>

			<div class="detail-card">
				<h3><Icon icon="mdi:database" /> Data Structure</h3>
				<div class="detail-content">
					<h4>Storage Format:</h4>
					<pre><code>{`data/
├── [Campus_Name]/
│   ├── wages_2024.json
│   ├── wages_2023.json
│   └── ...
└── scrape_progress.json`}</code></pre>

					<h4>Record Schema:</h4>
					<pre><code>{`{
  "location": "Berkeley",
  "year": 2024,
  "scraped_at": "2025-09-13T19:16:37Z",
  "total_records": 37078,
  "records": [
    {
      "firstname": "*****",
      "lastname": "*****",
      "title": "Professor",
      "location": "Berkeley",
      "year": "2024",
      "basepay": "150,000.00",
      "overtimepay": "0.00",
      "adjustpay": "5,000.00",
      "grosspay": "155,000.00"
    }
  ]
}`}</code></pre>
				</div>
			</div>

			<div class="detail-card">
				<h3><Icon icon="mdi:chart-line" /> Data Processing</h3>
				<div class="detail-content">
					<h4>Aggregation Pipeline:</h4>
					<ul>
						<li><strong>Ingestion:</strong> JSON files are parsed and validated</li>
						<li><strong>Transformation:</strong> Pay values converted from strings to numbers</li>
						<li><strong>Aggregation:</strong> Calculate totals, averages, and counts by campus/year</li>
						<li><strong>Database Storage:</strong> PostgreSQL with Drizzle ORM for efficient queries</li>
						<li><strong>Caching:</strong> Aggregated results cached for performance</li>
					</ul>

					<h4>Key Metrics Calculated:</h4>
					<ul>
						<li>Total wages by campus and year</li>
						<li>Average wages per employee</li>
						<li>Employee count trends</li>
						<li>Pay distribution analysis</li>
						<li>Year-over-year growth rates</li>
					</ul>
				</div>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Data Coverage</h2>

		<div class="coverage-grid">
			<div class="coverage-card">
				<Icon icon="mdi:school" class="coverage-icon" />
				<h3>13 UC Locations</h3>
				<p>Complete coverage of all UC campuses and affiliated institutions</p>
				<ul class="campus-list">
					<li>UC Berkeley</li>
					<li>UC Davis</li>
					<li>UC Irvine</li>
					<li>UCLA</li>
					<li>UC Merced</li>
					<li>UC Riverside</li>
					<li>UC San Diego</li>
					<li>UC San Francisco</li>
					<li>UC Santa Barbara</li>
					<li>UC Santa Cruz</li>
					<li>UC Office of the President</li>
					<li>UC SF Law</li>
					<li>ASUCLA</li>
				</ul>
			</div>

			<div class="coverage-card">
				<Icon icon="mdi:calendar-range" class="coverage-icon" />
				<h3>15 Years of Data</h3>
				<p>Historical data from 2010 to 2024</p>
				<div class="stats">
					<div class="stat">
						<span class="stat-value">195</span>
						<span class="stat-label">Location-Year Combinations</span>
					</div>
					<div class="stat">
						<span class="stat-value">~2M+</span>
						<span class="stat-label">Individual Records</span>
					</div>
					<div class="stat">
						<span class="stat-value">100%</span>
						<span class="stat-label">Public Data</span>
					</div>
				</div>
			</div>

			<div class="coverage-card">
				<Icon icon="mdi:update" class="coverage-icon" />
				<h3>Update Frequency</h3>
				<p>Annual updates when new data becomes available</p>
				<ul>
					<li>Data typically released in Q1 for previous year</li>
					<li>Automated scraping process for updates</li>
					<li>Historical data preserved for trend analysis</li>
					<li>Version control for data changes</li>
				</ul>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Privacy & Anonymization</h2>

		<div class="privacy-card">
			<Icon icon="mdi:shield-lock" class="privacy-icon" />
			<div class="privacy-content">
				<p>
					While all data displayed is publicly available per California state law, we respect
					individual privacy. Names are anonymized in certain contexts (shown as "*****") while
					maintaining the ability to analyze compensation trends by job title, department, and campus.
				</p>
				<p>
					The UC Annual Wage website implements its own anonymization for employees earning below
					certain thresholds or in specific categories. We preserve these anonymizations in our dataset.
				</p>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Data Quality Assurance</h2>

		<div class="quality-grid">
			<div class="quality-item">
				<Icon icon="mdi:check-circle" class="quality-icon" />
				<h4>Completeness Checks</h4>
				<p>Verify all expected fields are present</p>
			</div>
			<div class="quality-item">
				<Icon icon="mdi:calculator" class="quality-icon" />
				<h4>Calculation Validation</h4>
				<p>Ensure gross pay equals sum of components</p>
			</div>
			<div class="quality-item">
				<Icon icon="mdi:compare" class="quality-icon" />
				<h4>Cross-Reference</h4>
				<p>Compare totals with official UC reports</p>
			</div>
			<div class="quality-item">
				<Icon icon="mdi:alert-circle" class="quality-icon" />
				<h4>Anomaly Detection</h4>
				<p>Flag unusual patterns for review</p>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Limitations & Disclaimers</h2>

		<div class="disclaimer-card">
			<Icon icon="mdi:information-outline" class="disclaimer-icon" />
			<div class="disclaimer-content">
				<ul>
					<li>Data reflects only base salary and standard compensation components</li>
					<li>Does not include benefits, retirement contributions, or other non-wage compensation</li>
					<li>Mid-year hires or departures may show partial year compensation</li>
					<li>Job titles and departments may change between years</li>
					<li>Some records may be anonymized at the source for privacy</li>
					<li>Data accuracy depends on source reporting to UC system</li>
				</ul>
			</div>
		</div>
	</section>

	<section class="content-section">
		<div class="cta-card">
			<h3>Questions About Our Methodology?</h3>
			<p>We're committed to transparency in our data collection and processing methods.</p>
			<div class="cta-buttons">
				<a href="https://github.com/ucinvestments/UC-wages" target="_blank" rel="noopener noreferrer" class="cta-button">
					<Icon icon="mdi:github" />
					View Source Code
				</a>
				<a href="mailto:dev@ucinvestments.info" class="cta-button secondary">
					<Icon icon="mdi:email" />
					Contact Technical Team
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
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	.hero-section {
		padding: 1.5rem 1rem;
		background: linear-gradient(135deg, var(--pri) 0%, var(--founder) 100%);
		position: relative;
		overflow: hidden;
		margin-bottom: 4rem;
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

	.content-section {
		margin-bottom: 4rem;
	}

	.section-title {
		font-size: 2rem;
		font-weight: 700;
		text-align: center;
		margin-bottom: 2rem;
		color: var(--pri);
	}

	.process-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 2rem;
		margin-bottom: 3rem;
	}

	.process-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
		position: relative;
		transition: transform 0.3s ease, box-shadow 0.3s ease;
	}

	.process-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
	}

	.process-number {
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, var(--founder), var(--pri));
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
		font-size: 1.25rem;
		margin-bottom: 1rem;
	}

	.process-card h3 {
		color: var(--pri);
		margin-bottom: 1rem;
		font-size: 1.25rem;
	}

	.process-card p {
		color: var(--text-secondary);
		line-height: 1.6;
	}

	.tech-details {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.detail-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
	}

	.detail-card h3 {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: var(--pri);
		margin-bottom: 1.5rem;
		font-size: 1.5rem;
	}

	.detail-content h4 {
		color: var(--founder);
		margin-top: 1.5rem;
		margin-bottom: 1rem;
	}

	.detail-content:first-child h4:first-child {
		margin-top: 0;
	}

	.detail-content ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.detail-content li {
		padding: 0.5rem 0;
		color: var(--text-secondary);
		padding-left: 1.5rem;
		position: relative;
	}

	.detail-content li::before {
		content: "•";
		position: absolute;
		left: 0;
		color: var(--founder);
		font-weight: bold;
	}

	pre {
		background: var(--bg-secondary);
		padding: 1rem;
		border-radius: 8px;
		overflow-x: auto;
		font-size: 0.875rem;
		line-height: 1.5;
		margin: 1rem 0;
	}

	code {
		font-family: 'JetBrains Mono', monospace;
		color: var(--text-primary);
	}

	.coverage-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: 2rem;
	}

	.coverage-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
		text-align: center;
	}

	.coverage-icon {
		font-size: 3rem;
		color: var(--founder);
		margin-bottom: 1rem;
	}

	.coverage-card h3 {
		color: var(--pri);
		margin-bottom: 1rem;
	}

	.coverage-card p {
		color: var(--text-secondary);
		margin-bottom: 1.5rem;
	}

	.campus-list {
		list-style: none;
		padding: 0;
		display: grid;
		grid-template-columns: 1fr;
		gap: 0.5rem;
		text-align: left;
		max-width: 250px;
		margin: 0 auto;
	}

	.campus-list li {
		padding: 0.25rem 0;
		color: var(--text-secondary);
		font-size: 0.925rem;
		border-left: 3px solid var(--founder);
		padding-left: 1rem;
	}

	.stats {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		margin-top: 1rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--founder);
	}

	.stat-label {
		font-size: 0.875rem;
		color: var(--text-secondary);
		text-align: center;
	}

	.privacy-card, .disclaimer-card {
		background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
		border-radius: 12px;
		padding: 2rem;
		display: flex;
		gap: 2rem;
		align-items: flex-start;
	}

	.privacy-icon, .disclaimer-icon {
		font-size: 2.5rem;
		color: var(--founder);
		flex-shrink: 0;
	}

	.privacy-content p, .disclaimer-content {
		color: var(--text-secondary);
		line-height: 1.6;
	}

	.privacy-content p:first-child {
		margin-bottom: 1rem;
	}

	.disclaimer-content ul {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.disclaimer-content li {
		padding: 0.5rem 0;
		padding-left: 1.5rem;
		position: relative;
	}

	.disclaimer-content li::before {
		content: "•";
		position: absolute;
		left: 0;
		color: var(--founder);
		font-weight: bold;
	}

	.quality-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 2rem;
	}

	.quality-item {
		text-align: center;
	}

	.quality-icon {
		font-size: 2.5rem;
		color: var(--founder);
		margin-bottom: 1rem;
	}

	.quality-item h4 {
		color: var(--pri);
		margin-bottom: 0.5rem;
		font-size: 1.125rem;
	}

	.quality-item p {
		color: var(--text-secondary);
		font-size: 0.925rem;
	}

	.cta-card {
		background: linear-gradient(135deg, var(--founder), var(--pri));
		color: white;
		border-radius: 16px;
		padding: 3rem;
		text-align: center;
	}

	.cta-card h3 {
		font-size: 1.75rem;
		margin-bottom: 1rem;
	}

	.cta-card p {
		font-size: 1.125rem;
		margin-bottom: 2rem;
		opacity: 0.9;
	}

	.cta-buttons {
		display: flex;
		gap: 1rem;
		justify-content: center;
		flex-wrap: wrap;
	}

	.cta-button {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.875rem 1.75rem;
		background: white;
		color: var(--pri);
		text-decoration: none;
		border-radius: 8px;
		font-weight: 500;
		transition: all 0.3s ease;
	}

	.cta-button.secondary {
		background: transparent;
		color: white;
		border: 2px solid white;
	}

	.cta-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	}

	.cta-button.secondary:hover {
		background: white;
		color: var(--pri);
	}

	@media (max-width: 768px) {
		.page-title {
			font-size: 1.75rem;
		}

		.hero-subtitle {
			font-size: 1rem;
		}

		.process-grid {
			grid-template-columns: 1fr;
		}

		.coverage-grid {
			grid-template-columns: 1fr;
		}

		.quality-grid {
			grid-template-columns: repeat(2, 1fr);
		}

		.privacy-card, .disclaimer-card {
			flex-direction: column;
			text-align: center;
		}

		.cta-buttons {
			flex-direction: column;
		}

		.cta-button {
			width: 100%;
			justify-content: center;
		}
	}
</style>