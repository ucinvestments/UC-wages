<script lang="ts">
	let { label, caption, columns, rows }: {
		label: string;
		caption: string;
		columns: string[];
		rows: string[][];
	} = $props();
</script>

<details>
	<summary>{label}</summary>
	{#if rows.length}
		<!-- svelte-ignore a11y_no_noninteractive_tabindex (keyboard users need to scroll wide tables) -->
		<div class="table-scroll" role="region" aria-label={caption} tabindex="0">
			<table>
				<caption>{caption}</caption>
				<thead>
					<tr>
						{#each columns as column}
							<th scope="col">{column}</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each rows as [heading, ...values]}
						<tr>
							<th scope="row">{heading}</th>
							{#each values as value}
								<td>{value}</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<p>No payroll data is available for this view.</p>
	{/if}
</details>

<style>
	details {
		border: 1px solid var(--border);
		border-radius: 12px;
		background: white;
	}

	summary {
		padding: 0.75rem 1rem;
		cursor: pointer;
		font-weight: 600;
		color: var(--pri);
	}

	summary:focus-visible,
	.table-scroll:focus-visible {
		outline: 2px solid var(--founder);
		outline-offset: -2px;
	}

	.table-scroll {
		max-height: 24rem;
		overflow: auto;
	}

	table {
		width: 100%;
		font-variant-numeric: tabular-nums;
	}

	caption,
	p {
		padding: 0.75rem 1rem;
		text-align: left;
	}

	th,
	td {
		padding: 0.5rem 1rem;
		border-top: 1px solid var(--border);
		text-align: left;
	}

	thead th {
		position: sticky;
		top: 0;
		background: var(--bg-secondary);
	}

	td {
		white-space: nowrap;
	}
</style>
