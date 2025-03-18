<template>
	<div class="fixed px-8 bottom-24 w-full">
		<button class="w-full btn btn-primary shadow-xl" onclick="addWorkout.showModal()">
			<i class="bi bi-plus-square-dotted text-2xl"></i>
		</button>
	</div>

	<dialog id="addWorkout" class="modal modal-bottom sm:modal-middle">
		<div class="modal-box">
			<h3 class="text-lg font-bold">New workout</h3>
			<form method="dialog" class="w-full" @submit="submit">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Name</legend>
					<input v-model="workout.name" type="text" class="input w-full" placeholder="Name" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Note</legend>
					<input v-model="workout.note" type="text" class="input w-full" placeholder="Note" />
					<p class="fieldset-label">Not required</p>
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
	const emit = defineEmits(["add-workout"]);

	const workout = ref({
		name: "",
		note: "",
	});

	function submit() {
		emit("add-workout", {...workout.value});
		workout.value.name = "";
		workout.value.note = "";

	}
</script>
