<template>
	<div class="flex flex-row justify-center items-center w-full h-full px-8">
		<form class="container-sm border border-base-300 rounded-lg w-full bg-base-200 p-8" @submit.prevent="submit">
			<h1 class="w-full text-center text-2xl font-bold">Create account</h1>

			<ValidationInput label="Username" v-model="username" type="text" :error="usernameErr" />
			<ValidationInput label="Email" v-model="email" type="email" :error="emailErr" />
			<ValidationInput label="Password" v-model="password" type="password" :error="passwordErr" />
			<ValidationInput label="Repeat Password" v-model="repeatPassword" type="password"
				:error="repeatPasswordErr" />

			<p class="mb-4 text-sm opacity-60">Already have an account? <router-link to="/login"
					class="underline text-primary">Sign in</router-link></p>
			<button type="submit" class="btn btn-primary w-full" :disabled="loading">Register</button>
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
	const usernameErr = ref(null);

	const email = ref("");
	const emailErr = ref(null);

	const password = ref("");
	const passwordErr = ref(null);

	const repeatPassword = ref("");
	const repeatPasswordErr = ref(null);

	const error = ref(null);


	const loading = ref(false);

	const submit = async () => {
		error.value = null;

		const url = `${import.meta.env.VITE_API_URL}api/register`;

		try {
			loading.value = true;
			if (password.value != repeatPassword.value) {
				repeatPasswordErr.value = "Passwords does not match";
				loading.value = false;
				return;
			}

			const response = await fetch(url, {
				method: "POST",
				headers: {"Content-Type": "application/json"},
				body: JSON.stringify({
					username: username.value,
					email: email.value,
					password: password.value
				})
			});

			if (response.ok) {
				const {token} = await response.json();
				localStorage.setItem('token', token);
				const redirectPath = route.query.redirect || '/';
				console.log(`Registered in moving user to ${redirectPath}`);
				router.push(redirectPath);
			} else {
				const output = await response.json();
				if (output.username) {
					usernameErr.value = output.username;
				}
				if (output.email) {
					emailErr.value = output.email;
				}
				if (output.password) {
					passwordErr.value = output.password;
				}
				console.log("Failed to register");
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
