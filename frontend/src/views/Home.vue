<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full h-full p-8">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Hello {{name}}! </h1>
			<p class="text-left text-sm opacity-60">Are you ready to workout?</p>
		</div>
		<div class="w-full grid grid-cols-1 lg:grid-cols-3 gap-4 mb-16">
			<div class="flex flex-row p-4 bg-blue-50 rounded-2xl border border-blue-200 gap-2 shadow-md">
				<div
					class="bg-blue-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-bullseye text-blue-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Total Workouts</p>
					<p class="text-2xl font-bold">{{workouts.length}}</p>
				</div>
			</div>
			<div class="flex flex-row p-4 bg-green-50 rounded-2xl border border-green-200 gap-2 shadow-md">
				<div
					class="bg-green-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-graph-up-arrow text-green-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Completed</p>
					<p class="text-2xl font-bold">{{workouts.length}}</p>
				</div>
			</div>
			<div class="flex flex-row p-4 bg-purple-50 rounded-2xl border border-purple-200 gap-2 shadow-md">
				<div
					class="bg-purple-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-clock text-purple-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Total Time</p>
					<p class="text-2xl font-bold">{{workouts.length}}h</p>
				</div>
			</div>
		</div>
		<h2 class="font-bold text-2xl mb-4">Your workouts</h2>
		<transition-group name="workout" tag="div"
			class="w-full grid grid-cols-1 xl:grid-cols-3 lg:grid-cols-2 place-items-center gap-4 pb-40">
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
