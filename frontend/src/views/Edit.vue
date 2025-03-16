<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else-if="error" class="w-full h-full flex justify-center items-center">
		<p>Error</p>
	</div>
	<div v-else class="w-full h-full">
		<div class="w-full p-8">
			<div class="mb-8">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Name</legend>
					<input v-model="workout.name" type="text" class="input w-full" placeholder="Name" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Note</legend>
					<input v-model="workout.note" type="text" class="input w-full" placeholder="Note" />
					<p class="fieldset-label">Not required</p>
				</fieldset>
			</div>
			<button
				class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 gap-8 border border-primary"
				v-for="exercise in workout.exercises" :key="workout.id">
				<div class="grow">
					<h1 class="text-xl text-left">{{exercise.name}}</h1>
					<p class="text-xs text-left opacity-60">{{exercise.note}}</p>

					<fieldset class="fieldset w-full mt-4">
						<legend class="fieldset-legend text-left">Sets: {{exercise.sets}}</legend>
						<input v-model="exercise.sets" type="range" min="1" max="8" value="{{exercise.sets}}"
							class="range" step="1" />
					</fieldset>
					<fieldset class="fieldset w-full mt-4">
						<legend class="fieldset-legend text-left">Reps: {{exercise.reps}}</legend>
						<input v-model="exercise.reps" type="range" min="1" max="20" value="{{exercise.reps}}"
							class="range" step="1" />
					</fieldset>

				</div>
				<div class="flex flex-col justify-between">
					<div class="flex flex-col gap-4">
						<button class="btn btn-square btn-ghost">
							<i class="bi bi-arrow-up text-xl"></i>
						</button>
						<button class="btn btn-square btn-ghost">
							<i class="bi bi-arrow-down text-xl"></i>
						</button>
					</div>
					<button class="btn btn-square btn-ghost">
						<i class="bi bi-trash text-xl text-error"></i>
					</button>

				</div>

			</button>
		</div>

		<div class="fixed px-8 bottom-24 w-full">
			<button class="w-full btn btn-primary shadow-xl">
				<i class="bi bi-plus-square-dotted text-2xl"></i>
			</button>
		</div>
	</div>
</template>

<script setup>
	import {ref, onMounted, defineComponent} from "vue";
	import {jwtDecode} from 'jwt-decode';
	import {useRouter, useRoute} from 'vue-router'

	const router = useRouter()
	const route = useRoute()

	const workout = ref(null);
	const loading = ref(true);
	const error = ref(false);

	const id = route.query.id;

	async function loadWorkout() {
		const token = localStorage.getItem('token');
		try {
			//const decoded = jwtDecode(token);
			_ = jwtDecode(token);
		}
		catch {
			console.log("Could not get token");
		}

		const url = `${import.meta.env.VITE_API_URL}/api/restricted/workout/${id}`;

		try {

			const response = await fetch(url, {
				method: "GET",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
			});
			if (!response.ok) {
				throw new Error('Failed to fetch workouts');
			}
			workout.value = await response.json();
			console.log(workout.value)
		}
		catch (_) {
			error.value = true;
		} finally {
			loading.value = false;
		}
	}

	onMounted(loadWorkout);
</script>
