<template>
	<router-view />
	<Dock />
	<div class="h-8 w-full bg-red-500">
	</div>
</template>

<script setup>
	import Dock from "./components/Dock.vue";
</script>
