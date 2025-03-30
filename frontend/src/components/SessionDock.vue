<template>
	<Teleport to="body">
		<div class="fixed px-8 bottom-24 w-full gap-4 flex flex-row">
			<button class="btn shadow grow" @click="promptSkip">
				<p>Skip</p>
			</button>
			<button class="btn btn-primary shadow grow" @click="next">
				<p>Next</p>
			</button>
		</div>
	</Teleport>

	<ConfirmationModal ref="confirmModal" promptText="Are you sure?" :detailText="`Are you sure that you want to skip?`"
		confirmText="Yes" @confirmed="skip" />
</template>

<script setup>
	import {ref} from 'vue';
	import ConfirmationModal from '@/components/ConfirmationModal.vue';
	const emit = defineEmits(["next", "skip"]);
	const confirmModal = ref(null);

	function promptSkip() {
		confirmModal.value.open();
	}

	function skip() {
		emit("skip", {});
	}

	function next() {
		emit("next", {});
	}
</script>
