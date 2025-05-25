<template>
	<div class="card card-border bg-base-100 w-96 shadow-xl">
		<div class="card-body">
			<div class="flex justify-between">
				<h2 class="text-2xl font-bold">{{workout.name}}</h2>
				<div class="dropdown dropdown-end">
					<div tabindex="0" role="button" class="m-1">
						<i class="bi bi-three-dots-vertical"></i>
					</div>
					<ul tabindex="0" class="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm">
						<li>
							<a @click="emitEdit">
								<i class="bi bi-pencil"></i>
								Edit
							</a>
						</li>
						<li>
							<a class="text-error" @click="promptRemove">
								<i class="bi bi-trash"></i>
								Delete
							</a>
						</li>
					</ul>
				</div>
			</div>
			<p>{{ workout.note }}</p>
			<div class="flex flex-row gap-4">
				<div class="flex flex-row items-center">
					<div class="bg-blue-200 p-1.5 mr-2 w-5 h-5 text-center rounded-md flex justify-center items-center">
						<i class="bi bi-bullseye text-blue-500 text-sm"></i>
					</div>
					<p class="text-sm opacity-60">{{workout.exercises.length}} exercises</p>
				</div>
				<div class="flex flex-row items-center">
					<div
						class="bg-purple-200 p-1.5 mr-2 w-5 h-5 text-center rounded-md flex justify-center items-center">
						<i class="bi bi-clock text-purple-500 text-sm"></i>
					</div>
					<p class="text-sm opacity-60">60 min</p>
				</div>
			</div>
			<p class="opacity-60">Last done:</p>
			<p class="font-bold capitalize">{{timeAgo}}</p>
			<p class="opacity-60">{{niceDate}}</p>
			<div class="card-actions mt-5">
				<button class="btn btn-primary w-full" @click="emitStart">
					<i class="bi bi-caret-right"></i>
					Start Workout
				</button>
			</div>
		</div>
	</div>
	<ConfirmationModal ref="confirmModal" promptText="Are you sure?"
		:detailText="`Are you sure that you want to delete ${workout.name}?`" confirmText="Yes"
		@confirmed="emitRemove" />

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

	const niceDate = new Date(props.workout.last_done).toLocaleDateString('en-US', {
		weekday: 'long', // "Saturday"
		month: 'short',  // "Jan"
		day: 'numeric',  // "20"
	});
</script>
