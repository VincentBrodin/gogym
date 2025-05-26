<template>
	<div class="card card-border bg-base-100 shadow-xl w-full">
		<div class="card-body flex flex-row justify-between gap-8">
			<div class="grow flex flex-col gap-4">
				<!--INFO-->

				<input v-model="exercise.name" type="text" class="input input-ghost px-0 text-xl w-full mb"
					placeholder="Exercise Name" />

				<input v-model="exercise.note" type="text" class="input input-ghost w-full px-0 opacity-60"
					placeholder="Exercise Note" />

				<!--SETS-->
				<fieldset class="fieldset w-full">
					<div class="flex flex-row items-center justify-between w-full">
						<div class="flex flex-row items-center">
							<div
								class="bg-blue-200 p-1.5 mr-2 w-5 h-5 text-center rounded-md flex justify-center items-center">
								<i class="bi bi-hash text-blue-500 text-sm"></i>
							</div>
							<legend class="fieldset-legend text-left">Sets: </legend>
						</div>
						<p class="font-bold text-blue-500 text-xl text-end">{{exercise.sets}}</p>
					</div>

					<input v-model.number="exercise.sets" type="range" min="1" max="8" class="range w-full" step="1" />
					<div class="w-full flex flex-row justify-between">
						<p class="text-start opacity-60">1</p>
						<p class="text-end opacity-60">8</p>
					</div>
				</fieldset>

				<!--REPS-->
				<fieldset class="fieldset w-full">
					<div class="flex flex-row items-center justify-between w-full">
						<div class="flex flex-row items-center">
							<div
								class="bg-green-200 p-1.5 mr-2 w-5 h-5 text-center rounded-md flex justify-center items-center">
								<i class="bi bi-arrow-repeat text-green-500 text-sm"></i>
							</div>
							<legend class="fieldset-legend text-left">Reps:</legend>
						</div>
						<p class="font-bold text-green-500 text-xl text-end">{{exercise.reps}}</p>
					</div>
					<input v-model.number="exercise.reps" type="range" min="1" max="20" class="range w-full" step="1" />
					<div class="w-full flex flex-row justify-between">
						<p class="text-start opacity-60">1</p>
						<p class="text-end opacity-60">20</p>
					</div>
				</fieldset>

				<!--RIR-->
				<fieldset class="fieldset w-full">

					<div class="flex flex-row items-center justify-between w-full">
						<div class="flex flex-row items-center">
							<div
								class="bg-orange-200 p-1.5 mr-2 w-5 h-5 text-center rounded-md flex justify-center items-center">
								<i class="bi bi-bag text-orange-500 text-sm"></i>
							</div>
							<legend class="fieldset-legend text-left">RIR:</legend>
						</div>
						<p class="font-bold text-orange-500 text-xl text-end">{{exercise.rir}}</p>
					</div>
					<input v-model.number="exercise.rir" type="range" min="0" max="4" class="range w-full" step="1" />
					<div class="w-full flex flex-row justify-between">
						<p class="text-start opacity-60">0</p>
						<p class="text-end opacity-60">4</p>
					</div>
				</fieldset>

			</div>
			<div class="flex flex-col justify-between">
				<button class="btn btn-circle btn-ghost" @click.stop="promptRemove">
					<i class="bi bi-trash text-lg text-error"></i>
				</button>
				<div class="flex flex-col gap-4">
					<!--UP-->
					<button v-if="exercise.order != 0" class="btn btn-circle btn-ghost" @click.stop="moveUp">
						<i class="bi bi-arrow-up text-lg"></i>
					</button>

					<!--DOWN-->
					<button v-if="exercise.order != total - 1" class="btn btn-circle btn-ghost" @click.stop="moveDown">
						<i class="bi bi-arrow-down text-lg"></i>
					</button>
				</div>
			</div>
		</div>
		<ConfirmationModal ref="confirmModal" promptText="Are you sure?"
			:detailText="`Are you sure that you want to delete ${exercise.name}?`" confirmText="Yes"
			@confirmed="remove" />

	</div>

</template>


<script setup>
	import {watch, ref} from 'vue';
	import ConfirmationModal from './ConfirmationModal.vue';

	const props = defineProps({
		exercise: Object,
		total: Number
	});
	const emit = defineEmits(["move-up", "move-down", "remove", "change"]);
	const confirmModal = ref(null)


	function moveUp() {
		emit("move-up", props.exercise);
	}

	function moveDown() {
		emit("move-down", props.exercise);
	}
	function remove() {
		emit("remove", props.exercise);
	}
	function promptRemove() {
		confirmModal.value.open()
	}


	watch(
		() => ({
			name: props.exercise.name,
			note: props.exercise.note,
			sets: props.exercise.sets,
			reps: props.exercise.reps,
			rir: props.exercise.rir,
		}),
		(newValue, _) => {
			props.exercise.sets = parseInt(newValue.sets, 10);
			props.exercise.reps = parseInt(newValue.reps, 10);
			props.exercise.rir = parseInt(newValue.rir, 10);
			emit('change', {...props.exercise});
		},
		{deep: true}
	);

</script>
