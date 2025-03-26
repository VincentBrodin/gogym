<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else-if="error" class="w-full h-full flex justify-center items-center">
		<p>Error</p>
	</div>
	<div v-else class="w-full h-full p-8">
		<h1 class="text-2xl font-bold mb-8">Do you want to copy {{workout.name}}?</h1>
		<ul class="w-full my-8">
			<li v-for="exercise in workout.exercises">
				{{ exercise.name }}
				<span v-if="exercise.note != '' && exercise.note != null"> - 
					<span class="opacity-60">{{exercise.note}}</span>
				</span>
			</li>
		</ul>
		<button class="btn btn-primary w-full mb-4" @click="copyWorkout">Copy</button>
		<router-link to="/" class="btn w-full">No</router-link>
	</div>

</template>

<script setup>
	import {useRoute, useRouter} from 'vue-router';
	import {ref, onMounted} from 'vue';
	const route = useRoute()
	const router = useRouter()
	const id = route.query.id;

	const loading = ref(true);
	const error = ref(false);
	const workout = ref(null);

	async function loadWorkout() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/workout/${id}`;
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

	async function copyWorkout() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/workout/copy/${id}`;
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
			const json = await response.json();
			router.push({
				name: 'edit',
				query: {'id': json.id}
			});

			console.log(json)
		}
		catch (err) {
			console.log(err)
			error.value = true;
		} finally {
			loading.value = false;
		}
	}


	onMounted(loadWorkout);
</script>
