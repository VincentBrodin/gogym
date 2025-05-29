<template>
	<div class="card card-border bg-base-100 w-96 shadow-xl">
		<div class="card-body">
			<div class="flex flex-row justify-between">
				<div class="flex flex-col">
					<p class="text-xl text-left font-bold">{{ session.workout.name }}</p>
					<p class="opacity-60">
						<i class="bi bi-calendar"></i>
						{{ niceDate }}
					</p>
				</div>
				<div class="flex flex-col">
					<p class="text-xl font-bold text-blue-400 text-end">{{ time.toFixed(0) }}</p>
					<p class="opacity-60 text-end text-sm">min</p>
				</div>
			</div>

			<div class="flex flex-row justify-evenly mt-4">
				<div class="flex flex-col justify-center items-center">
					<p class="text-center text-xl font-bold">{{session.exercise_sessions.length}}</p>
					<p class="text-center opacity-60 text-sm">Exercises</p>
				</div>

				<div class="flex flex-col justify-center items-center">
					<p class="text-center text-xl font-bold">{{totalSets}}</p>
					<p class="text-center opacity-60 text-sm">Sets</p>
				</div>

				<div class="flex flex-col justify-center items-center">
					<p class="text-center text-xl font-bold">{{totalReps}}</p>
					<p class="text-center opacity-60 text-sm">Reps</p>
				</div>
			</div>
			<div class="card-actions flex flex-row gap-4">
				<button class="btn btn-primary grow" @click.stop="emitView">
					<i class="bi bi-eye text-lg"></i>
					View
				</button>
				<button class="btn btn-circle btn-ghost" @click.stop="promptRemove">
					<i class="bi bi-trash text-lg text-error"></i>
				</button>
			</div>

		</div>
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

	const emits = defineEmits(['view', 'remove']);

	const emitRemove = () => {
		emits('remove', props.session);
	};

	const emitView = () => {
		emits('view', props.session);
	};

	function promptRemove() {
		confirmModal.value.open()
	}


	const timeAgo = useTimeAgo(new Date(props.session.endend_at), {
		updateInterval: 60000, // One min 
		fullDateFormatter: date => date.toLocaleDateString()
	});

	const time = (new Date(props.session.endend_at) - new Date(props.session.started_at)) / 1000 / 60
	const niceDate = new Date(props.session.endend_at).toLocaleDateString('en-US', {
		weekday: 'long', // "Saturday"
		month: 'short',  // "Jan"
		day: 'numeric',  // "20"
	});

	const totalSets = props.session.exercise_sessions.reduce((sum, exerciseSession) => {
		return sum + (exerciseSession.sets_done || 0);
	}, 0);
	const totalReps = props.session.exercise_sessions.reduce((sum, exSession) => {
		const sets = exSession.sets_done || 0;
		const reps = exSession.exercise.reps || 0;
		return sum + sets * reps;
	}, 0);
</script>
