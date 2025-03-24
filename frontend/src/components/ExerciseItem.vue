<template>
	<div class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 gap-8 border">
		<div class="grow">
			<!--INFO-->
			<fieldset class="fieldset">
				<legend class="fieldset-legend text-left">Name</legend>
				<input v-model="exercise.name" type="text" class="input w-full" placeholder="Name" />
			</fieldset>
			<fieldset class="fieldset">
				<legend class="fieldset-legend text-left">Note</legend>
				<input v-model="exercise.note" type="text" class="input w-full" placeholder="Note" />
				<p class="fieldset-label">Not required</p>
			</fieldset>


			<!--SETS-->
			<fieldset class="fieldset w-full">
				<legend class="fieldset-legend text-left">Sets: {{exercise.sets}}</legend>
				<input v-model.number="exercise.sets" type="range" min="1" max="8" class="range" step="1" />
			</fieldset>

			<!--REPS-->
			<fieldset class="fieldset w-full">
				<legend class="fieldset-legend text-left">Reps: {{exercise.reps}}</legend>
				<input v-model.number="exercise.reps" type="range" min="1" max="20" class="range" step="1" />
			</fieldset>

			<!--RIR-->
			<fieldset class="fieldset w-full">
				<legend class="fieldset-legend text-left">RIR: {{exercise.rir}}</legend>
				<input v-model.number="exercise.rir" type="range" min="0" max="4" class="range" step="1" />
			</fieldset>

		</div>
		<div class="flex flex-col justify-between">
			<div class="flex flex-col gap-8">
				<!--UP-->
				<button v-if="exercise.order != 0" class="btn btn-square" @click.stop="moveUp">
					<i class="bi bi-arrow-up text-xl"></i>
				</button>

				<!--DOWN-->
				<button v-if="exercise.order != total - 1" class="btn btn-square" @click.stop="moveDown">
					<i class="bi bi-arrow-down text-xl"></i>
				</button>
			</div>
			<!--DELETE-->
			<button class="btn btn-square" @click.stop="promptRemove">
				<i class="bi bi-trash text-xl text-error"></i>
			</button>

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
