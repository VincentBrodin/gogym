<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full h-full">
		<div class="w-full p-8">
			<div class="mb-8">
				<h1 class="text-left text-2xl ">Hello {{name}}! </h1>
				<p class="text-left text-sm opacity-60">Are you ready to workout?</p>
			</div>
			<button class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 border border-primary"
				v-for="workout in workouts" :key="workout.id" @click="startWorkout(workout)">
				<div>
					<h1 class="text-xl text-left">{{workout.name}}</h1>
					<p class="text-xs text-left opacity-60">{{workout.note}}</p>
				</div>
				<div>
					<button class="btn btn-square btn-ghost" @click.stop="editWorkout(workout)">
						<i class="bi bi-pencil text-xl"></i>
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
	import {ref, onMounted} from "vue";
	import {useRouter} from 'vue-router'
	import {jwtDecode} from 'jwt-decode';

	const router = useRouter()

	const workouts = ref([]);
	const loading = ref(true);
	const name = ref("John Doe");



	async function loadWorkouts() {
		const token = localStorage.getItem('token');
		try {
			const decoded = jwtDecode(token);
			name.value = decoded.uname;
		}
		catch {
			console.log("Could not get username");
		}

		const url = `${import.meta.env.VITE_API_URL}/api/restricted/workouts`;

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
			console.log(workouts.value)
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	onMounted(loadWorkouts);

	function startWorkout(workout) {
		console.log("Start")
		console.log(workout)
	}

	function editWorkout(workout) {
		router.push({
			name: 'edit',
			query: {id: workout.id}
		});
	}
</script>
