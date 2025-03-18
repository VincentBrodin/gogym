<template>
	<div class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 border">
		<div>
			<h1 class="text-xl text-left">{{ workout.name }}</h1>
			<p class="text-xs text-left opacity-60">{{ workout.note }}</p>
			<p class="text-xs text-left opacity-60">Last did: {{timeAgo}}</p>

		</div>
		<div class="flex flex-row items-center justify-end gap-4">
			<button class="btn btn-circle btn-ghost" @click.stop="emitEdit">
				<i class="bi bi-pencil text-xl"></i>
			</button>
			<button class="btn btn-circle btn-ghost" @click.stop="emitRemove">
				<i class="bi bi-trash text-xl text-error"></i>
			</button>
		</div>
	</div>
</template>

<script setup>
	import {defineProps, defineEmits} from 'vue';
	import {useTimeAgo} from '@vueuse/core';

	const props = defineProps({
		workout: {
			type: Object,
			required: true
		}
	});

	const emits = defineEmits(['edit', 'remove']);

	const emitEdit = () => {
		emits('edit', props.workout);
	};

	const emitRemove = () => {
		emits('remove', props.workout);
	};


	const timeAgo = useTimeAgo(new Date(props.workout.last_done), {
		updateInterval: 60000, // One min 
		fullDateFormatter: date => date.toLocaleDateString()
	});
</script>
