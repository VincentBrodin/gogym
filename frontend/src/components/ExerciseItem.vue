<template>
	<button class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 gap-8 border border-primary">
		<div class="grow">
			<!--INFO-->
			<h1 class="text-xl text-left">{{exercise.name}}</h1>
			<p class="text-xs text-left opacity-60">{{exercise.note}}</p>

			<!--SETS-->
			<fieldset class="fieldset w-full mt-4">
				<legend class="fieldset-legend text-left">Sets: {{exercise.sets}}</legend>
				<input v-model="exercise.sets" type="range" min="1" max="8" class="range" step="1" />
			</fieldset>

			<!--REPS-->
			<fieldset class="fieldset w-full mt-4">
				<legend class="fieldset-legend text-left">Reps: {{exercise.reps}}</legend>
				<input v-model="exercise.reps" type="range" min="1" max="20" class="range" step="1" />
			</fieldset>

		</div>
		<div class="flex flex-col justify-between">
			<div class="flex flex-col gap-4">
				<!--UP-->
				<button v-if="exercise.order != 0" class="btn btn-square btn-ghost" @click.stop="moveUp(exercise)">
					<i class="bi bi-arrow-up text-xl"></i>
				</button>

				<!--DOWN-->
				<button v-if="exercise.order != total - 1" class="btn btn-square btn-ghost"
					@click.stop="moveDown(exercise)">
					<i class="bi bi-arrow-down text-xl"></i>
				</button>
			</div>
			<!--DELETE-->
			<button class="btn btn-square btn-ghost" @click.stop="remove(exercise)">
				<i class="bi bi-trash text-xl text-error"></i>
			</button>

		</div>

	</button>

</template>


<script setup>
	import {watch} from 'vue';

	const props = defineProps({
		exercise: Object,
		total: Number
	});
	const emit = defineEmits(["move-up", "move-down", "remove", "change"]);

	function moveUp() {
		emit("move-up", props.exercise);
	}

	function moveDown() {
		emit("move-down", props.exercise);
	}
	function remove() {
		emit("remove", props.exercise);
	}

	watch(
		() => props.exercise,
		(newValue, _) => {
			newValue.sets = parseInt(newValue.sets)
			newValue.reps = parseInt(newValue.reps)
			emit('change', newValue);
		},
		{deep: true}
	);


</script>
