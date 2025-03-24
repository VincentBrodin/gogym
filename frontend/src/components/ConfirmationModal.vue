<template>
	<Teleport to="body">
		<dialog id="confirmModal" ref="dialog" class="modal modal-bottom sm:modal-middle">
			<div class="modal-box">
				<h3 class="text-lg font-bold">{{promptText}}</h3>
				<p class="opacity-60">{{detailText}}</p>
				<form method="dialog" class="w-full" @submit.prevent="confirmAction">
					<button type="submit" class="btn btn-error w-full mt-4">{{confirmText}}</button>
				</form>
			</div>
			<form method="dialog" class="modal-backdrop">
				<button>close</button>
			</form>

		</dialog>
	</Teleport>
</template>

<script setup>
	import {ref} from 'vue'
	const props = defineProps({
		promptText: {type: String, required: true},
		detailText: {type: String, required: true},
		confirmText: {type: String, required: true},
	})
	const emit = defineEmits(['confirmed'])
	const dialog = ref(null)

	function open() {
		if (dialog.value) {
			dialog.value.showModal();
		}
	}

	function close() {
		if (dialog.value) {
			dialog.value.close();
		}
	}

	function confirmAction() {
		close();
		emit('confirmed')
	}
	defineExpose({open})
</script>
