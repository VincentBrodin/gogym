<template>
	<Teleport to="body">
		<div class="fixed px-8 bottom-24 w-full">
			<button class="w-full btn btn-primary" onclick="addExercise.showModal()">
				<i class="bi bi-plus text-2xl"></i>
				Add Exercise
			</button>
		</div>
	</Teleport>

	<dialog id="addExercise" class="modal modal-bottom sm:modal-middle">
		<div class="modal-box">
			<h3 class="text-lg font-bold">New exercise</h3>
			<form method="dialog" class="w-full" @submit="submit">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Name</legend>
					<input v-model="exercise.name" type="text" class="input w-full" placeholder="Name" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Note</legend>
					<input v-model="exercise.note" type="text" class="input w-full" placeholder="Note" />
					<p class="fieldset-label">Not required</p>
				</fieldset>

				<!--SETS-->
				<fieldset class="fieldset w-full mt-4">
					<legend class="fieldset-legend text-left">Sets: {{exercise.sets}}</legend>
					<input v-model.number="exercise.sets" type="range" min="1" max="8" class="range w-full" step="1" />
				</fieldset>

				<!--REPS-->
				<fieldset class="fieldset w-full mt-4">
					<legend class="fieldset-legend text-left">Reps: {{exercise.reps}}</legend>
					<input v-model.number="exercise.reps" type="range" min="1" max="20" class="range w-full" step="1" />
				</fieldset>

				<!--RIR-->
				<fieldset class="fieldset w-full mt-4">
					<legend class="fieldset-legend text-left">RIR: {{exercise.rir}}</legend>
					<input v-model.number="exercise.rir" type="range" min="0" max="4" class="range w-full" step="1" />
				</fieldset>


				<button type="submit" class="btn btn-primary w-full mt-4">Add</button>
			</form>
		</div>
		<form method="dialog" class="modal-backdrop">
			<button>close</button>
		</form>

	</dialog>
</template>

<script setup>
	import {ref} from "vue";
	const emit = defineEmits(["add-exercise"]);

	const exercise = ref({
		name: "",
		note: "",
		sets: 3,
		reps: 8,
		rir: 0,
	});

	function submit() {
		exercise.value.sets = parseInt(exercise.value.sets);
		exercise.value.reps = parseInt(exercise.value.reps);
		exercise.value.rir = parseInt(exercise.value.rir);
		emit("add-exercise", {...exercise.value});
		exercise.value.name = "";
		exercise.value.note = "";
		exercise.value.sets = 3;
		exercise.value.reps = 8;
		exercise.value.rir = 0;

	}
</script>
