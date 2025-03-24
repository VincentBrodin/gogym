<template>
	<div class="w-full">
		<button class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 border" @click="emitStart">
			<div>
				<h1 class="text-xl text-left font-semibold">{{ workout.name }}</h1>
				<p class="text-xs text-left opacity-60">{{ workout.note }}</p>
				<p class="text-xs text-left opacity-60">Last did: {{timeAgo}}</p>

			</div>
			<div class="flex flex-row items-center justify-end gap-4">
				<button class="btn btn-square " @click.stop="emitEdit">
					<i class="bi bi-pencil text-xl"></i>
				</button>
				<button class="btn btn-square" @click.stop="promptRemove">
					<i class="bi bi-trash text-xl text-error"></i>
				</button>
			</div>
		</button>
		<ConfirmationModal ref="confirmModal" promptText="Are you sure?"
			:detailText="`Are you sure that you want to delete ${workout.name}?`" confirmText="Yes"
			@confirmed="emitRemove" />

	</div>

</template>

<script setup>
	import {ref, defineProps, defineEmits} from 'vue';
	import {useTimeAgo} from '@vueuse/core';
	import ConfirmationModal from './ConfirmationModal.vue';

	const props = defineProps({
		workout: {
			type: Object,
			required: true
		}
	});

	const confirmModal = ref(null)

	const emits = defineEmits(['click', 'edit', 'remove']);

	const emitStart = () => {
		emits('click', props.workout);
	};


	const emitEdit = () => {
		emits('edit', props.workout);
	};

	const emitRemove = () => {
		emits('remove', props.workout);
	};

	function promptRemove() {
		confirmModal.value.open()
	}


	const timeAgo = useTimeAgo(new Date(props.workout.last_done), {
		updateInterval: 60000, // One min 
		fullDateFormatter: date => date.toLocaleDateString()
	});
</script>
