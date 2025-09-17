<script lang="ts">
	import '../app.css';
	import { dev } from '$app/environment';
	import { inject } from '@vercel/analytics';
	import { fade } from 'svelte/transition';
	import { page } from '$app/stores';
	import Icon from '@iconify/svelte';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

	onMount(() => {
		if (browser) {
			// Initialize Vercel Analytics
			inject({ mode: dev ? 'development' : 'production' });
		}
	});

	function toggleDonationInfo() {
		const donationInfo = document.getElementById('donationInfo');
		if (donationInfo) {
			donationInfo.style.display = donationInfo.style.display === 'none' ? 'block' : 'none';
		}
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav class="navbar">
	<div class="nav-container">
		<a href="/" class="logo-link">
			<div class="logo">
				<Icon icon="mdi:currency-usd" class="logo-icon" />
				<span class="logo-text">UC Wages</span>
			</div>
		</a>

		<div class="nav-links">
			<a href="/" class="nav-link" class:active={$page.url.pathname === '/'}>
				<Icon icon="mdi:chart-line" class="nav-icon" />
				Explorer
			</a>
			<a href="/about" class="nav-link" class:active={$page.url.pathname === '/about'}>
				<Icon icon="mdi:information" class="nav-icon" />
				About
			</a>
			<a href="/data" class="nav-link" class:active={$page.url.pathname === '/data'}>
				<Icon icon="mdi:database" class="nav-icon" />
				Data
			</a>
			<a href="/methodology" class="nav-link" class:active={$page.url.pathname === '/methodology'}>
				<Icon icon="mdi:book-open-page-variant" class="nav-icon" />
				Methodology
			</a>
			<a
				href="https://github.com/ucinvestments/UC-wages"
				target="_blank"
				rel="noopener noreferrer"
				class="nav-link external"
			>
				<Icon icon="mdi:github" class="nav-icon" />
				GitHub
				<Icon icon="mdi:open-in-new" class="external-icon" />
			</a>
		</div>
	</div>
</nav>

<main in:fade={{ duration: 300 }}>
	{@render children?.()}
</main>

<footer class="footer">
	<div class="footer-content">
		<div class="footer-section">
			<h4 class="footer-title">UC Wage Explorer</h4>
			<p class="footer-text">
				Transparency in University of California employee compensation data.
			</p>
		</div>

		<div class="footer-section">
			<h4 class="footer-title">Quick Links</h4>
			<div class="footer-links">
				<a href="/" class="footer-link">Explorer</a>
				<a href="/about" class="footer-link">About</a>
				<a href="/data" class="footer-link">Data</a>
				<a href="/methodology" class="footer-link">Methodology</a>
				<a
					href="https://ucannualwage.ucop.edu/wage/"
					target="_blank"
					rel="noopener noreferrer"
					class="footer-link"
				>
					UC Annual Wage
					<Icon icon="mdi:open-in-new" class="footer-external" />
				</a>
			</div>
		</div>

		<div class="footer-section">
			<h4 class="footer-title">Contact</h4>
			<div class="footer-links">
				<a href="mailto:admin@ucinvestments.info" class="footer-link">Contact Us</a>
				<a href="mailto:press@ucinvestments.info" class="footer-link">Press Inquiries</a>
				<a href="mailto:dev@ucinvestments.info" class="footer-link">Development</a>
			</div>
		</div>

		<div class="footer-section">
			<h4 class="footer-title">Support This Project</h4>
			<p class="footer-text">
				This project is self-funded. Donations help cover hosting and development costs.
			</p>
			<div class="donation-links">
				<button class="donate-button" on:click={toggleDonationInfo}>
					<Icon icon="mdi:heart" class="donate-icon" />
					Donate
				</button>
				<div class="donation-info" id="donationInfo" style="display: none;">
					<div class="crypto-address">
						<strong>ETH:</strong>
						<code class="address">0x623c7559ddC51BAf15Cc81bf5bc13c0B0EA14c01</code>
					</div>
					<div class="crypto-address">
						<strong>XMR:</strong>
						<code class="address"
							>44bvXALNkxUgSkGChKQPnj79v6JwkeYEkGijgKyp2zRq3EiuL6oewAv5u2c7FN7jbN1z7uj1rrPfL77bbsJ3cC8U2ADFoTj</code
						>
					</div>
					<p class="alt-contact">
						Or contact <a href="mailto:admin@ucinvestments.info">Admin</a> for alternatives.
					</p>
				</div>
			</div>
		</div>

		<div class="footer-section">
			<h4 class="footer-title">Data Sources</h4>
			<p class="footer-text">
				Last updated: 2024<br />
				UC Annual Wage Database
			</p>
		</div>
	</div>

	<div class="footer-bottom">
		<p>&copy; 2024 UC Wage Explorer. Educational purposes only.</p>
	</div>
</footer>

<style>
	:global(html) {
		scroll-behavior: smooth;
	}

	.navbar {
		position: sticky;
		top: 0;
		z-index: 100;
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(20px);
		border-bottom: 1px solid var(--border);
		box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
	}

	.nav-container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 1rem 2rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.logo-link {
		text-decoration: none;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		transition: transform 0.3s ease;
	}

	.logo:hover {
		transform: scale(1.05);
	}

	:global(.logo-icon) {
		font-size: 2rem;
		color: var(--founder);
	}

	.logo-text {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--pri);
		letter-spacing: -0.01em;
	}

	.nav-links {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.nav-link {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.625rem 1.25rem;
		color: var(--text-secondary);
		text-decoration: none;
		font-weight: 500;
		border-radius: 0.75rem;
		transition: all 0.2s ease;
		position: relative;
	}

	:global(.nav-icon) {
		font-size: 1.125rem;
	}

	.nav-link:hover {
		background: var(--bg-secondary);
		color: var(--pri);
		transform: translateY(-1px);
	}

	.nav-link.active {
		background: linear-gradient(135deg, var(--founder), var(--pri));
		color: white;
	}

	.nav-link.external {
		border: 2px solid var(--border);
	}

	:global(.external-icon) {
		font-size: 0.875rem;
		margin-left: -0.125rem;
	}

	main {
		min-height: calc(100vh - 400px);
	}

	.footer {
		background: linear-gradient(180deg, var(--bg-secondary) 0%, white 100%);
		border-top: 1px solid var(--border);
		margin-top: 4rem;
		padding: 3rem 0 1.5rem;
	}

	.footer-content {
		max-width: 1400px;
		margin: 0 auto;
		padding: 0 2rem;
		display: grid;
		grid-template-columns: 2fr 1fr 1fr 1fr 1fr;
		gap: 2rem;
		margin-bottom: 2rem;
	}

	.footer-section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.footer-title {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 1.125rem;
		font-weight: 600;
		color: var(--pri);
		margin: 0;
	}

	.footer-text {
		color: var(--text-secondary);
		line-height: 1.6;
		font-size: 0.925rem;
		margin: 0;
	}

	.footer-links {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.footer-link {
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
		color: var(--text-secondary);
		text-decoration: none;
		font-size: 0.925rem;
		transition: all 0.2s ease;
		width: fit-content;
	}

	.footer-link:hover {
		color: var(--founder);
		transform: translateX(4px);
	}

	:global(.footer-external) {
		font-size: 0.75rem;
	}

	.footer-bottom {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem 2rem 0;
		border-top: 1px solid var(--border);
		text-align: center;
	}

	.footer-bottom p {
		color: var(--text-secondary);
		font-size: 0.875rem;
		margin: 0;
	}

	.donation-links {
		margin-top: 1rem;
	}

	.donate-button {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.625rem 1rem;
		background: linear-gradient(135deg, var(--golden-gate), var(--sec));
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.3s ease;
		font-size: 0.875rem;
	}

	.donate-button:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
	}

	:global(.donate-icon) {
		font-size: 1rem;
		color: white;
	}

	.donation-info {
		margin-top: 1rem;
		padding: 1rem;
		background: var(--bg-secondary);
		border-radius: 0.5rem;
		border: 1px solid var(--border);
	}

	.crypto-address {
		margin-bottom: 0.75rem;
	}

	.crypto-address strong {
		display: block;
		color: var(--pri);
		font-size: 0.875rem;
		margin-bottom: 0.25rem;
	}

	.address {
		display: block;
		background: white;
		padding: 0.5rem;
		border-radius: 0.25rem;
		font-family: 'JetBrains Mono', monospace;
		font-size: 0.75rem;
		word-break: break-all;
		border: 1px solid var(--border);
		color: var(--text-primary);
		cursor: text;
		user-select: all;
	}

	.alt-contact {
		font-size: 0.875rem;
		color: var(--text-secondary);
		margin: 0.75rem 0 0;
	}

	.alt-contact a {
		color: var(--founder);
		text-decoration: none;
		font-weight: 500;
	}

	.alt-contact a:hover {
		color: var(--pri);
		text-decoration: underline;
	}

	@media (max-width: 768px) {
		.nav-container {
			flex-direction: column;
			gap: 1rem;
			padding: 1rem;
		}

		.nav-links {
			width: 100%;
			justify-content: center;
		}

		.nav-link {
			padding: 0.5rem 1rem;
			font-size: 0.875rem;
		}

		:global(.nav-icon) {
			font-size: 1rem;
		}

		.logo-text {
			font-size: 1.125rem;
		}

		.footer-content {
			grid-template-columns: 1fr;
			gap: 2rem;
			padding: 0 1.5rem;
		}

		.address {
			font-size: 0.7rem;
		}

		.footer-section {
			text-align: center;
		}

		.footer-links {
			align-items: center;
		}
	}

	@media (max-width: 480px) {
		.nav-links {
			flex-direction: column;
			width: 100%;
		}

		.nav-link {
			width: 100%;
			justify-content: center;
		}
	}
</style>
