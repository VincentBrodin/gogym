<template>
	<div class="flex flex-row justify-center items-center w-full h-full px-8">
		<form class="container-sm border border-base-300 rounded-lg w-full bg-base-200 p-8" @submit.prevent="submit">
			<h1 class="w-full text-center text-2xl font-bold">Login</h1>

			<ValidationInput label="Username" v-model="username" type="text" :error="error" :showMsg="false" />
			<ValidationInput label="Password" v-model="password" type="password" :error="error" :showMsg="false" />

			<p v-if="error != null" class="text-error mb-4 text-sm font-bold">{{ error }}</p>

			<button type="submit" class="btn btn-primary btn-outline w-full" :disabled="loading">Login</button>
		</form>
	</div>
</template>

<script setup>
	import {ref} from "vue";
	import {useRouter, useRoute} from 'vue-router'
	import ValidationInput from '../components/ValidationInput.vue'

	const router = useRouter()
	const route = useRoute()

	const username = ref("");
	const password = ref("");
	const error = ref(null);
	const loading = ref(false);

	const submit = async () => {
		error.value = null;

		const url = `${import.meta.env.VITE_API_URL}/api/login`;

		try {
			loading.value = true;
			const response = await fetch(url, {
				method: "POST",
				headers: {"Content-Type": "application/json"},
				body: JSON.stringify({
					username: username.value,
					password: password.value
				})
			});

			if (response.ok) {
				const {token} = await response.json();
				localStorage.setItem('token', token);
				const redirectPath = route.query.redirect || '/';
				console.log(`Logged in moving user to ${redirectPath}`);
				router.push(redirectPath);
			} else {
				console.log("Failed to login");
				error.value = "Username or password is incorrect!";
			}
		} catch (err) {
			console.log(err)
			error.value = "An error occurred. Please try again.";
		}
		finally {
			loading.value = false;
		}
	};
</script>
