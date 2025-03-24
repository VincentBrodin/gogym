<template>
	<div class="fixed px-8 bottom-24 w-full">
		<button class="w-full btn shadow mb-4" @click="skip">
			<p>Skip</p>
		</button>
		<button class="w-full btn btn-primary shadow" @click="next">
			<p>Next</p>
		</button>
	</div>
</template>

<script setup>
	const emit = defineEmits(["next", "skip"]);

	function skip() {
		emit("skip", {});
	}

	function next() {
		emit("next", {});
	}
</script>
