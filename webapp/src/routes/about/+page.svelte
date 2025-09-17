<script lang="ts">
	import Icon from '@iconify/svelte';
</script>

<svelte:head>
	<title>About - UC Wage Explorer</title>
	<meta
		name="description"
		content="Learn about the UC Wage Explorer project and its mission to provide transparency in University of California employee compensation data."
	/>
</svelte:head>

<div class="page-container">
	<section class="hero-section">
		<h1 class="page-title">About UC Wage Explorer</h1>
		<p class="hero-subtitle">
			Promoting transparency and accountability through open access to University of California
			compensation data
		</p>
	</section>

	<section class="content-section">
		<div class="content-grid">
			<div class="card">
				<div class="card-icon">
					<Icon icon="mdi:target" />
				</div>
				<h2>Our Mission</h2>
				<p>
					The UC Wage Explorer provides accessible, interactive visualizations of University of
					California employee compensation data. We believe that transparency in public employee
					wages fosters accountability and enables informed discussions about resource allocation in
					higher education.
				</p>
			</div>

			<div class="card">
				<div class="card-icon">
					<Icon icon="mdi:database-search" />
				</div>
				<h2>Data Source</h2>
				<p>
					All data is sourced from the official UC Annual Wage website, which publishes compensation
					information for UC employees as required by California state law. We aggregate and
					visualize this public data to make it more accessible and understandable.
				</p>
			</div>

			<div class="card">
				<div class="card-icon">
					<Icon icon="mdi:code-braces" />
				</div>
				<h2>Open Source</h2>
				<p>
					This project is fully open source and available on GitHub. We encourage contributions,
					suggestions, and feedback from the community. The codebase includes both the data scraper
					and this web application.
				</p>
				<a
					href="https://github.com/ucinvestments/UC-wages"
					target="_blank"
					rel="noopener noreferrer"
					class="github-link"
				>
					<Icon icon="mdi:github" />
					View on GitHub
				</a>
			</div>

			<div class="card">
				<div class="card-icon">
					<Icon icon="mdi:shield-check" />
				</div>
				<h2>Privacy & Ethics</h2>
				<p>
					We only display publicly available information. No private or sensitive data is collected
					or displayed. This tool is designed for educational and research purposes to promote
					transparency in public institutions.
				</p>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">About the Developer</h2>
		<div class="developer-card">
			<div class="developer-info">
				<h3>Stephen Okita</h3>
				<p class="developer-bio">
					I am, as of now the only dev on this project. I would like help please.
				</p>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Contact Us</h2>
		<div class="contact-grid">
			<a href="mailto:dev@ucinvestments.info" class="contact-card">
				<Icon icon="mdi:code-tags" class="contact-icon" />
				<h3>Development</h3>
				<p>Technical inquiries and contributions</p>
				<span class="email">dev@ucinvestments.info</span>
			</a>

			<a href="mailto:admin@ucinvestments.info" class="contact-card">
				<Icon icon="mdi:account-cog" class="contact-icon" />
				<h3>Administration</h3>
				<p>General inquiries and support</p>
				<span class="email">admin@ucinvestments.info</span>
			</a>

			<a href="mailto:press@ucinvestments.info" class="contact-card">
				<Icon icon="mdi:newspaper-variant" class="contact-icon" />
				<h3>Press</h3>
				<p>Media and press inquiries</p>
				<span class="email">press@ucinvestments.info</span>
			</a>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Technology Stack</h2>
		<div class="tech-grid">
			<div class="tech-item">
				<Icon icon="mdi:language-go" class="tech-icon go" />
				<span>Go</span>
			</div>
			<div class="tech-item">
				<Icon icon="simple-icons:svelte" class="tech-icon svelte" />
				<span>SvelteKit</span>
			</div>
			<div class="tech-item">
				<Icon icon="mdi:language-typescript" class="tech-icon typescript" />
				<span>TypeScript</span>
			</div>
			<div class="tech-item">
				<Icon icon="simple-icons:postgresql" class="tech-icon postgres" />
				<span>PostgreSQL</span>
			</div>
			<div class="tech-item">
				<Icon icon="simple-icons:d3dotjs" class="tech-icon d3" />
				<span>D3.js</span>
			</div>
			<div class="tech-item">
				<Icon icon="simple-icons:tailwindcss" class="tech-icon tailwind" />
				<span>Tailwind CSS</span>
			</div>
		</div>
	</section>

	<section class="content-section">
		<h2 class="section-title">Support This Project</h2>
		<div class="support-card">
			<p>
				UC Wage Explorer is a self-funded project. If you find this tool valuable, please consider
				supporting its continued development and hosting costs.
			</p>
			<div class="support-buttons">
				<a
					href="https://github.com/ucinvestments/UC-wages"
					target="_blank"
					rel="noopener noreferrer"
					class="support-button"
				>
					<Icon icon="mdi:star" />
					Star on GitHub
				</a>
				<a href="mailto:admin@ucinvestments.info" class="support-button secondary">
					<Icon icon="mdi:email" />
					Contact for Support Options
				</a>
			</div>
		</div>
	</section>
</div>

<style>
	.page-container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	:global(body) {
		overflow-x: hidden;
	}

	.hero-section {
		padding: 1.5rem 1rem;
		background: linear-gradient(135deg, var(--pri) 0%, var(--founder) 100%);
		position: relative;
		overflow: hidden;
		margin-bottom: 4rem;
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
		max-width: 600px;
		margin: 0 auto;
		line-height: 1.6;
		text-align: center;
		position: relative;
		z-index: 1;
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

	.content-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 2rem;
		margin-bottom: 3rem;
	}

	.card {
		background: white;
		border-radius: 16px;
		padding: 2rem;
		box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
		transition:
			transform 0.3s ease,
			box-shadow 0.3s ease;
	}

	.card:hover {
		transform: translateY(-4px);
		box-shadow: 0 20px 40px -5px rgba(0, 0, 0, 0.15);
	}

	.card-icon {
		width: 48px;
		height: 48px;
		background: linear-gradient(135deg, var(--founder), var(--pri));
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 1.5rem;
		color: white;
		font-size: 24px;
	}

	.card h2 {
		font-size: 1.5rem;
		margin-bottom: 1rem;
		color: var(--pri);
	}

	.card p {
		color: var(--text-secondary);
		line-height: 1.6;
	}

	.github-link {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		margin-top: 1rem;
		padding: 0.75rem 1.5rem;
		background: var(--pri);
		color: white;
		text-decoration: none;
		border-radius: 8px;
		font-weight: 500;
		transition: all 0.3s ease;
	}

	.github-link:hover {
		background: var(--founder);
		transform: translateY(-2px);
	}

	.developer-card {
		background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
		border-radius: 16px;
		padding: 3rem;
		margin-bottom: 2rem;
	}

	.developer-info h3 {
		font-size: 1.75rem;
		color: var(--pri);
		margin-bottom: 0.5rem;
	}

	.developer-role {
		font-size: 1.125rem;
		color: var(--founder);
		font-weight: 500;
		margin-bottom: 1rem;
	}

	.developer-bio {
		color: var(--text-secondary);
		line-height: 1.6;
		max-width: 800px;
	}

	.contact-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
		gap: 2rem;
		margin-bottom: 2rem;
	}

	.contact-card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		text-decoration: none;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
		transition: all 0.3s ease;
		border: 2px solid transparent;
		text-align: center;
	}

	.contact-card:hover {
		transform: translateY(-4px);
		border-color: var(--founder);
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
	}

	.contact-icon {
		font-size: 2.5rem;
		color: var(--founder);
		margin-bottom: 1rem;
	}

	.contact-card h3 {
		font-size: 1.25rem;
		color: var(--pri);
		margin-bottom: 0.5rem;
	}

	.contact-card p {
		color: var(--text-secondary);
		font-size: 0.925rem;
		margin-bottom: 1rem;
	}

	.email {
		display: block;
		color: var(--founder);
		font-weight: 500;
		font-size: 0.925rem;
	}

	.tech-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 2rem;
		justify-content: center;
		margin-top: 2rem;
	}

	.tech-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}

	.tech-icon {
		font-size: 3rem;
		transition: transform 0.3s ease;
	}

	.tech-item:hover .tech-icon {
		transform: scale(1.1);
	}

	.tech-icon.go {
		color: #00add8;
	}

	.tech-icon.svelte {
		color: #ff3e00;
	}

	.tech-icon.typescript {
		color: #3178c6;
	}

	.tech-icon.postgres {
		color: #4169e1;
	}

	.tech-icon.d3 {
		color: #f68e56;
	}

	.tech-icon.tailwind {
		color: #06b6d4;
	}

	.tech-item span {
		color: var(--text-secondary);
		font-weight: 500;
	}

	.support-card {
		background: linear-gradient(135deg, #fef3c7, #fde68a);
		border-radius: 16px;
		padding: 2.5rem;
		text-align: center;
	}

	.support-card p {
		color: var(--text-primary);
		font-size: 1.125rem;
		margin-bottom: 2rem;
		max-width: 600px;
		margin-left: auto;
		margin-right: auto;
	}

	.support-buttons {
		display: flex;
		gap: 1rem;
		justify-content: center;
		flex-wrap: wrap;
	}

	.support-button {
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

	.support-button.secondary {
		background: white;
		color: var(--pri);
		border: 2px solid var(--pri);
	}

	.support-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.support-button.secondary:hover {
		background: var(--pri);
		color: white;
	}

	@media (max-width: 768px) {
		.page-title {
			font-size: 2rem;
		}

		.hero-subtitle {
			font-size: 1rem;
		}

		.content-grid {
			grid-template-columns: 1fr;
		}

		.contact-grid {
			grid-template-columns: 1fr;
		}

		.developer-card {
			padding: 2rem;
		}

		.support-buttons {
			flex-direction: column;
		}

		.support-button {
			width: 100%;
			justify-content: center;
		}
	}
</style>
