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
			<transition-group name="workout" tag="div" class="w-full">
				<button
					class="w-full flex flex-row justify-between bg-base-200 rounded-xl p-4 mb-4 border border-primary"
					v-for="workout in workouts" :key="workout.id" @click="startWorkout(workout)">
					<div>
						<h1 class="text-xl text-left">{{workout.name}}</h1>
						<p class="text-xs text-left opacity-60">{{workout.note}}</p>
					</div>
					<div>
						<button class="btn btn-square btn-ghost" @click.stop="editWorkout(workout)">
							<i class="bi bi-pencil text-xl"></i>
						</button>
						<button class="btn btn-square btn-ghost" @click.stop="removeWorkout(workout)">
							<i class="bi bi-trash text-xl text-error"></i>
						</button>

					</div>
				</button>
			</transition-group>
		</div>
	</div>
	<AddWorkoutModal @add-workout="addWorkout" />
</template>

<script setup>
	import {ref, onMounted} from "vue";
	import {useRouter} from 'vue-router'
	import {jwtDecode} from 'jwt-decode';

	import AddWorkoutModal from "@/components/AddWorkoutModal.vue";

	const router = useRouter()

	const workouts = ref([]);
	const loading = ref(true);
	const name = ref("John Doe");

	async function addWorkout(workout) {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}/api/restricted/workout`;
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

	async function removeWorkout(workout) {
		workouts.value.splice(workouts.value.indexOf(workout), 1);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}/api/restricted/workout/${workout.id}`;

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
