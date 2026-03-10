<script lang="ts">
	import type { PageProps } from './$types';
	import { Input } from '$lib/components/ui/input/index.js';
	import Button from '$lib/components/ui/button/button.svelte';

	let { data }: PageProps = $props();
</script>

<div class="space-y-4">
	<form class="space-y-2">
		<Input type="username" placeholder="Enter reg. number" class="max-w-xs" />
		<Input type="password" placeholder="Enter password" class="max-w-xs" />

		<Button type="submit">Begin Scrapping</Button>
	</form>
</div>
