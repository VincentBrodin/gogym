<template>
	<p v-if="session == null">Not found</p>
	<p v-else>Found</p>
</template>

<script setup>
	import {ref, onMounted} from "vue";
	const session = ref(null)

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
			session.value = await response.json();
			console.log(session.value)
		}
		catch (error) {
			console.log(error)
		}
	}

	onMounted(grabSession);
</script>
