<template>
	<div class="w-full">
		<button class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 border" @click="emitStart">
			<div>
				<h1 class="text-xl text-left font-semibold">{{ session.workout.name }}</h1>
				<p class="text-xs text-left opacity-60">{{ session.workout.note }}</p>
				<p class="text-xs text-left opacity-60 mb-4">{{timeAgo}}</p>
				<p class="text-xs text-left opacity-60">{{duration}}</p>

			</div>
			<div class="flex flex-row items-center justify-end gap-4">
				<button class="btn btn-square btn-error" @click.stop="promptRemove">
					<i class="bi bi-trash text-xl"></i>
				</button>
			</div>
		</button>
		<ConfirmationModal ref="confirmModal" promptText="Are you sure?"
			:detailText="`Are you sure that you want to delete this session?`" confirmText="Yes"
			@confirmed="emitRemove" />
	</div>

</template>

<script setup>
	import {ref, defineProps, defineEmits} from 'vue';
	import {useTimeAgo} from '@vueuse/core';
	import ConfirmationModal from './ConfirmationModal.vue';

	const props = defineProps({
		session: {
			type: Object,
			required: true
		}
	});

	const confirmModal = ref(null);
	const duration = updateTime(Date.parse(props.session.endend_at) - Date.parse(props.session.started_at));

	function updateTime(time) {
		const totalSeconds = Math.floor(time / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
			2,
			"0"
		)}:${String(seconds).padStart(2, "0")}`;
	}

	const emits = defineEmits(['remove']);

	const emitRemove = () => {
		emits('remove', props.session);
	};

	function promptRemove() {
		confirmModal.value.open()
	}


	const timeAgo = useTimeAgo(new Date(props.session.endend_at), {
		updateInterval: 60000, // One min 
		fullDateFormatter: date => date.toLocaleDateString()
	});
</script>
