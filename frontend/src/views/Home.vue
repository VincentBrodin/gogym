<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full h-full p-8">
		<div class="mb-8">
			<h1 class="text-left text-2xl font-bold">Hello {{name}}! </h1>
			<p class="text-left text-sm opacity-60">Are you ready to workout?</p>
		</div>
		<transition-group name="workout" tag="div" class="w-full">
			<WorkoutItem v-for="workout in workouts" :key="workout.id" :workout="workout" @click="startSession"
				@edit="editWorkout" @remove="removeWorkout" />
		</transition-group>
	</div>
	<AddWorkoutModal @add-workout="addWorkout" />
</template>

<script setup>
	import {ref, onMounted} from "vue";
	import {useRouter} from 'vue-router'
	import {jwtDecode} from 'jwt-decode';
	import {useLocalStorage} from '@vueuse/core'


	import AddWorkoutModal from "@/components/AddWorkoutModal.vue";
	import WorkoutItem from "@/components/WorkoutItem.vue";

	const router = useRouter()

	const workouts = ref([]);
	const loading = ref(true);
	const name = ref("John Doe");
	const session = useLocalStorage('session', false)


	async function loadWorkouts() {
		const token = localStorage.getItem('token');
		try {
			const decoded = jwtDecode(token);
			name.value = decoded.uname;
		}
		catch {
			console.log("Could not get username");
		}

		const url = `${import.meta.env.VITE_API_URL}api/restricted/workouts`;

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
			workouts.value = await response.json();
			workouts.value.sort((a, b) => Date.parse(a.last_done) - Date.parse(b.last_done));
			console.log(workouts.value)
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	onMounted(loadWorkouts);

	async function startSession(workout) {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${workout.id}`;
		try {
			const response = await fetch(url, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}

			router.push({
				name: 'session',
			});

			console.log("Started");
			session.value = true
		}
		catch (err) {
			console.log(err)
		}
	}

	async function addWorkout(workout) {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/workout`;
		try {
			const response = await fetch(url, {
				method: "PUT",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
				body: JSON.stringify(workout),
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}

			workouts.value.push(await response.json());
			console.log("Added");
		}
		catch (err) {
			console.log(err)
		}
	}

	function editWorkout(workout) {
		router.push({
			name: 'edit',
			query: {id: workout.id}
		});
	}

	async function removeWorkout(workout) {
		workouts.value.splice(workouts.value.indexOf(workout), 1);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/workout/${workout.id}`;

		try {

			const response = await fetch(url, {
				method: "DELETE",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token
				},
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

</script>
