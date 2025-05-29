<template>
	<div class="w-full h-full p-8 flex flex-col gap-4">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Settings</h1>
		</div>
		<fieldset class="fieldset">
			<legend class="fieldset-legend">Username</legend>
			<input v-model="username" type="text" class="input w-full" placeholder="Name" />
		</fieldset>
		<fieldset class="fieldset">
			<legend class="fieldset-legend">Weight System</legend>
			<label class="label">
				<input type="checkbox" v-model="imperial" class="checkbox" />
				Imperial
			</label>
		</fieldset>


		<button class="btn w-full" @click="logout" :disabled="loading">Log out</button>
		<button class="btn btn-error w-full" @click="deleteAccount" :disabled="loading">Delete account</button>
	</div>
</template>

<script setup>
	import {ref, onMounted, watch} from 'vue';
	import {jwtDecode} from 'jwt-decode';
	import {useRouter} from 'vue-router'
	//const themes = ["silk", "dracula", "lofi", "black"]

	let updateTimeout = null;
	const username = ref("");
	const imperial = ref(false);
	const router = useRouter();
	const loading = ref(true);

	function logout() {
		loading.value = true;
		localStorage.removeItem('token');
		router.push({
			name: 'login',
			query: {redirect: '/settings'}
		});
	}

	async function deleteAccount() {
		loading.value = true;
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/account`;
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
			localStorage.removeItem('token');
			router.push({
				name: 'login',
				query: {redirect: '/settings'}
			});
		}
		catch (err) {
			console.log(err)
		}

	}

	function getUsername() {
		const token = localStorage.getItem('token');
		try {
			const decoded = jwtDecode(token);
			username.value = decoded.uname;
		}
		catch {
			console.log("Could not get username");
		}
	}

	function getImperial() {
		const token = localStorage.getItem('token');
		try {
			const decoded = jwtDecode(token);
			imperial.value = decoded.imperial;
		}
		catch {
			console.log("Could not get imperial");
		}
	}

	async function triggerUpdate(newUsername, newSystem) {
		console.log(newUsername)
		console.log(newSystem)
		if (updateTimeout) clearTimeout(updateTimeout)
		updateTimeout = setTimeout(async () => {
			const token = localStorage.getItem('token');
			const url = `${import.meta.env.VITE_API_URL}api/restricted/account`;
			try {
				const response = await fetch(url, {
					method: "PATCH",
					headers: {
						"Content-Type": "application/json",
						"Authorization": "Bearer " + token
					},
					body: JSON.stringify({"username": newUsername, "imperial": newSystem})
				});
				if (!response.ok) {
					throw new Error((await response.json()).error);
				}
				const json = await response.json();
				console.log(json);
				localStorage.setItem('token', json.token)
			}
			catch (err) {
				console.log(err)
			}
		}, 750)
	}


	onMounted(() => {
		getUsername()
		getImperial()
		loading.value = false
	});


	watch(
		() => username.value,
		(newValue, oldValue) => {
			if (oldValue === "") return;
			if (newValue !== "") {
				triggerUpdate(newValue, imperial.value);
			}
		}
	);

	watch(
		() => imperial.value,
		(newValue) => {
			console.log("Hello world");
			triggerUpdate(username.value, newValue);
		}
	);
</script>
