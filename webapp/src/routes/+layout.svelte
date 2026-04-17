<script lang="ts">
  import "../app.css";
  import { dev } from "$app/environment";
  import { fade } from "svelte/transition";
  import { injectAnalytics } from "@vercel/analytics/sveltekit";
  import Navbar from "$lib/components/Navbar.svelte";
  import Footer from "$lib/components/Footer.svelte";

  injectAnalytics({ mode: dev ? "development" : "production" });

  let { children } = $props();

  const navLinks = [
    { href: "/", label: "Explorer", icon: "mdi:chart-line" },
    { href: "/search", label: "Search", icon: "mdi:account-search" },
    { href: "/about", label: "About", icon: "mdi:information" },
    { href: "/data", label: "Data", icon: "mdi:database" },
    { href: "/methodology", label: "Methodology", icon: "mdi:book-open-page-variant" },
    {
      href: "https://github.com/ucinvestments/UC-wages",
      label: "GitHub",
      icon: "mdi:github",
      external: true,
    },
  ];

  const footerLinks = [
    { href: "/", label: "Explorer" },
    { href: "/search", label: "Search" },
    { href: "/about", label: "About" },
    { href: "/data", label: "Data" },
    { href: "/methodology", label: "Methodology" },
    {
      href: "https://ucannualwage.ucop.edu/wage/",
      label: "UC Annual Wage",
      external: true,
    },
  ];

  const cryptoAddresses = [
    { label: "ETH", address: "0x623c7559ddC51BAf15Cc81bf5bc13c0B0EA14c01" },
    {
      label: "XMR",
      address:
        "44bvXALNkxUgSkGChKQPnj79v6JwkeYEkGijgKyp2zRq3EiuL6oewAv5u2c7FN7jbN1z7uj1rrPfL77bbsJ3cC8U2ADFoTj",
    },
  ];
</script>

<Navbar appName="UC Wages" logoIcon="mdi:currency-usd" links={navLinks} />

<main in:fade={{ duration: 300 }}>
  {@render children?.()}
</main>

<Footer
  appName="UC Wage Explorer"
  description="Transparency in University of California employee compensation data."
  links={footerLinks}
  contactEmail="admin@ucinvestments.info"
  altContactEmail="admin@ucinvestments.info"
  {cryptoAddresses}
  dataSources={"UC Annual Wage Database<br />Source: ucannualwage.ucop.edu"}
/>

<style>
  :global(html) {
    scroll-behavior: smooth;
  }

  main {
    min-height: calc(100vh - 400px);
  }
</style>
