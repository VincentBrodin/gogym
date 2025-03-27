<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full h-full p-8">
		<transition-group name="session" tag="div" class="w-full">
			<StatItem v-for="session in sessions" :key="session.id" :session="session" @remove="removeSession" />
		</transition-group>
	</div>

</template>

<script setup>
	import StatItem from '@/components/StatItem.vue';
	import {ref, onMounted} from 'vue';

	const sessions = ref(null);
	const loading = ref(true);

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
			sessions.value = await response.json();
			sessions.value.sort((a, b) => Date.parse(b.endend_at) - Date.parse(a.endend_at));
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	async function removeSession(session) {
		sessions.value.splice(sessions.value.indexOf(session), 1);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${session.id}`;

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
			console.log(await response.json());
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	onMounted(loadSessions);

</script>
