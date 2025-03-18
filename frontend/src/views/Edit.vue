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
				<!--INFO-->
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
			<!--EXERCISE-->
			<transition-group name="exercise" tag="div" class="w-full">
				<ExerciseItem v-for="exercise in workout.exercises" :key="exercise.id" :exercise="exercise"
					:total="workout.exercises.length" @move-up="moveUp" @move-down="moveDown" @remove="remove"
					@change="triggerUpdateExercise(exercise)" />
			</transition-group>
			<div class="w-full h-32">

			</div>
		</div>
		<AddExerciseModal @add-exercise="addExercise" />
	</div>
</template>

<script setup>
	const optemistic = true;
	import {ref, onMounted, nextTick, watch} from "vue";
	import {useRoute} from 'vue-router'
	import AddExerciseModal from '@/components/AddExerciseModal.vue'
	import ExerciseItem from '@/components/ExerciseItem.vue'

	const route = useRoute()

	const workout = ref({name: "", note: ""});
	const loading = ref(true);
	const error = ref(false);

	const id = route.query.id;

	let updateExerciseTimeout = null
	let updateWorkoutTimeout = null


	const debouncedUpdates = new Map();
	function triggerUpdateExercise(exercise) {
		const {id} = exercise;
		if (!debouncedUpdates.has(id)) {
			debouncedUpdates.set(id, debounce(id, updateExercise, 750));
		}
		debouncedUpdates.get(id)(exercise);
	}

	const timers = new Map();
	function debounce(id, func, delay) {
		return function (...args) {
			if (timers.has(id)) {
				clearTimeout(timers.get(id));
			}
			const timer = setTimeout(() => {
				func(...args);
				timers.delete(id);
			}, delay);
			timers.set(id, timer);
		};
	}




	async function moveUp(exercise) {
		const start = exercise.order;
		const above = workout.value.exercises.find(e => e.order == start - 1);
		if (above == undefined) return;
		exercise.order = above.order;
		above.order = start;
		workout.value.exercises.sort((a, b) => a.order - b.order);

		await nextTick();
		await triggerUpdateAllExercises();

	}

	async function moveDown(exercise) {
		const start = exercise.order;
		const below = workout.value.exercises.find(e => e.order == start + 1);
		if (below == undefined) return;
		exercise.order = below.order;
		below.order = start;
		workout.value.exercises.sort((a, b) => a.order - b.order);

		await nextTick();
		await triggerUpdateAllExercises();
	}

	async function remove(exercise) {
		workout.value.exercises.splice(workout.value.exercises.indexOf(exercise), 1);
		workout.value.exercises.sort((a, b) => a.order - b.order);

		await nextTick();

		for (let i in workout.value.exercises) {
			workout.value.exercises[i].order = i;
		}
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}/api/restricted/exercise/${exercise.id}`;
		try {
			const response = await fetch(url, {
				method: "DELETE",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				}
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}

			if (!optemistic) {
				workout.value.exercises = await response.json();
				workout.value.exercises.sort((a, b) => a.order - b.order);
				console.log(workout.value.exercises)
			}
			console.log("Removed")
		}
		catch (err) {
			console.log(err)
		}
	}

	async function updateExercise(exercise) {
		// Force numbers 
		exercise.sets = parseInt(exercise.sets, 10);
		exercise.reps = parseInt(exercise.reps, 10);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}/api/restricted/exercise/${exercise.id}`;
		try {
			const response = await fetch(url, {
				method: "PATCH",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
				body: JSON.stringify(exercise)
			});

			if (!response.ok) {
				throw new Error((await response.json()).error);
			}

			console.log("Update " + exercise.id)
		}
		catch (err) {
			console.log(err)
		}
	}

	async function triggerUpdateAllExercises() {
		if (updateExerciseTimeout) clearTimeout(updateExerciseTimeout)
		updateExerciseTimeout = setTimeout(async () => {
			// Force numbers 
			for (let exercise of workout.value.exercises) {
				exercise.sets = parseInt(exercise.sets, 10);
				exercise.reps = parseInt(exercise.reps, 10);
			}
			const token = localStorage.getItem('token');
			const url = `${import.meta.env.VITE_API_URL}/api/restricted/exercises/${workout.value.id}`;
			try {
				const response = await fetch(url, {
					method: "PATCH",
					headers: {
						"Content-Type": "application/json",
						"Authorization": "Bearer " + token
					},
					body: JSON.stringify(workout.value.exercises)
				});
				if (!response.ok) {
					throw new Error((await response.json()).error);
				}

				if (!optemistic) {
					workout.value.exercises = await response.json();
					workout.value.exercises.sort((a, b) => a.order - b.order);
					console.log(workout.value.exercises)
				}
				console.log("Update")
			}
			catch (err) {
				console.log(err)
			}
		}, 750)
	}

	async function triggerUpdateWorkout() {
		if (updateWorkoutTimeout) clearTimeout(updateWorkoutTimeout)
		updateWorkoutTimeout = setTimeout(async () => {
			// Force numbers 
			for (let exercise of workout.value.exercises) {
				exercise.sets = parseInt(exercise.sets, 10);
				exercise.reps = parseInt(exercise.reps, 10);
			}
			const token = localStorage.getItem('token');
			const url = `${import.meta.env.VITE_API_URL}/api/restricted/workout/${workout.value.id}`;
			try {
				const response = await fetch(url, {
					method: "PATCH",
					headers: {
						"Content-Type": "application/json",
						"Authorization": "Bearer " + token
					},
					body: JSON.stringify(workout.value)
				});
				if (!response.ok) {
					throw new Error((await response.json()).error);
				}

				console.log("Update")
			}
			catch (err) {
				console.log(err)
			}
		}, 750)
	}




	async function addExercise(exercise) {
		exercise.workout_id = workout.value.id;
		// Force numbers 
		exercise.sets = parseInt(exercise.sets, 10);
		exercise.reps = parseInt(exercise.reps, 10);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}/api/restricted/exercise`;
		try {
			const response = await fetch(url, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
				body: JSON.stringify(exercise),
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}
			workout.value.exercises.push(await response.json());
			workout.value.exercises.sort((a, b) => a.order - b.order);
			console.log("Added")
		}
		catch (err) {
			console.log(err)
		}
	}

	async function loadWorkout() {
		const token = localStorage.getItem('token');
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
				throw new Error((await response.json()).error);
			}
			workout.value = await response.json();
			workout.value.exercises.sort((a, b) => a.order - b.order);
			console.log(workout.value)
		}
		catch (err) {
			console.log(err)
			error.value = true;
		} finally {
			loading.value = false;
		}
	}

	onMounted(loadWorkout);

	watch(
		() => ({
			name: workout.value.name,
			note: workout.value.note,
		}),
		(_, oldValue) => {
			if (oldValue.name == "") {
				return;
			}
			triggerUpdateWorkout();
		},
		{deep: true}
	);


</script>
