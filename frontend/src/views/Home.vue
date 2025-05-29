<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full h-full p-8">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Hello {{name}}! </h1>
			<p class="text-left text-sm opacity-60">Are you ready to workout?</p>
		</div>
		<div class="w-full grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
			<div
				class="flex flex-row p-4 bg-gradient-to-br from-blue-500/5 to-indigo-500/5 rounded-2xl border border-blue-200 gap-2 shadow-md">
				<div
					class="bg-blue-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-bullseye text-blue-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Total Workouts</p>
					<p class="text-2xl font-bold">{{workouts.length}}</p>
				</div>
			</div>
			<div
				class="flex flex-row p-4 bg-gradient-to-br from-emerald-500/5 to-teal-500/5 rounded-2xl border border-green-200 gap-2 shadow-md">
				<div
					class="bg-green-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-graph-up-arrow text-green-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Completed</p>
					<p class="text-2xl font-bold">{{workouts.length}}</p>
				</div>
			</div>
			<div
				class="flex flex-row p-4 bg-gradient-to-br from-purple-500/5 to-pink-500/5 rounded-2xl border border-purple-200 gap-2 shadow-md">
				<div
					class="bg-purple-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-clock text-purple-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Total Time</p>
					<p v-if="totalTime == null" class="text-2xl font-bold">Loading</p>
					<p v-else class="text-2xl font-bold">{{totalTime.toFixed(1)}} h</p>
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
	const totalTime = ref(null);
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

	async function loadSessions() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/sessions`;
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
			console.log("Hello world")
			const allSessions = await response.json();
			const grouped = allSessions.reduce((acc, session) => {
				const workoutId = session.workout.id;
				if (!acc[workoutId]) {
					acc[workoutId] = [];
				}
				acc[workoutId].push(session);
				return acc;
			}, {});

			let tt = 0;
			for (let workout of workouts.value) {
				const sessions = grouped[workout.id];
				console.log(sessions);
				if (sessions == undefined || sessions == null) {
					workout.time = -1;
				} else {
					let time = 0;
					for (let session of sessions) {
						const start = new Date(session.started_at)
						const end = new Date(session.endend_at)
						const diffInMs = end - start;
						const diffInMinutes = diffInMs / 1000 / 60;
						time += diffInMinutes;
						console.log(diffInMinutes)
					}
					workout.time = time / sessions.length;
					tt += workout.time;
				}
			}
			if (tt <= 0) {
				totalTime.value = -1
			} else {
				totalTime.value = tt / 60
			}
		}
		catch (error) {
			console.error('Error loading sessions:', error);
		}
	}


	onMounted(async () => {
		await loadWorkouts()
		await loadSessions()
	});

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
