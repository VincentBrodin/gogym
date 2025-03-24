<template>
	<div class="dock dock-lg">
		<router-link to="/" :class="{ 'dock-active': (page == 'home') }">
			<button class="flex flex-col justify-between items-center">
				<i class="bi bi-house text-xl"></i>
				<span class="dock-label">Home</span>
			</button>
		</router-link>

		<router-link to="/session" :class="[{ 'disabled': !session }, { 'dock-active': (page == 'session') }]">
			<button class="flex flex-col justify-between items-center" :disabled="!session">
				<i class="bi bi-activity text-xl"></i>
				<span class="dock-label">Session</span>
			</button>
		</router-link>

		<router-link to="/">
			<button class="flex flex-col justify-between items-center ">
				<i class="bi bi-clipboard-data text-xl"></i>
				<span class="dock-label">Stats</span>
			</button>
		</router-link>


		<router-link to="/">
			<button class="flex flex-col justify-between items-center">
				<i class="bi bi-gear text-xl"></i>
				<span class="dock-label">Settings</span>
			</button>
		</router-link>

	</div>
</template>

<script setup>
	import {ref, onMounted, watch} from "vue";
	import {useLocalStorage} from '@vueuse/core';
	import {useRoute} from 'vue-router';

	const session = useLocalStorage('session', false);
	const route = useRoute();
	const page = ref("home");
	watch(
		() => route.name,
		(newPage, oldPage) => {
			console.log('Route changed from', oldPage, 'to', newPage);
			page.value = newPage;
		}
	);

	async function grabSession() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session`;

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
			session.value = true;
		}
		catch (error) {
			session.value = false;
			console.log(error)
		}
	}

	onMounted(grabSession);
</script>
