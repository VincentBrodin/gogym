<template>
	<div class="dock dock-lg">
		<router-link to="/">
			<button class="flex flex-col justify-between items-center"
				:class="{ 'dock-active': ($route.name == 'Home') }">
				<i class="bi bi-house text-xl"></i>
				<span class="dock-label">Home</span>
			</button>
		</router-link>

		<router-link to="/session" :class="{ 'disabled': !session }">
			<button class="flex flex-col justify-between items-center" :disabled="!session">
				<i class="bi bi-activity text-xl"></i>
				<span class="dock-label">Session</span>
			</button>
		</router-link>

		<router-link to="/">
			<button class="flex flex-col justify-between items-center">
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
	import {ref, onMounted} from "vue";
	const session = ref(false);

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
