<template>
	<router-view />
	<Dock />
	<div class="h-8 w-full opacity-0 pointer-events-none">
	</div>
</template>

<script setup>
	import Dock from "./components/Dock.vue";
</script>
