<template>
	<div class="w-full h-full p-8 flex flex-col gap-4">
		<!--<div class="dropdown w-full mb-4">
			<div tabindex="0" role="button" class="btn m-1 w-full">
				Theme
				<svg width="12px" height="12px" class="inline-block h-2 w-2 fill-current opacity-60"
					xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2048 2048">
					<path d="M1799 349l242 241-1017 1017L7 590l242-241 775 775 775-775z"></path>
				</svg>
			</div>
			<ul tabindex="0" class="dropdown-content bg-base-300 rounded-box z-1 p-2 shadow-2xl w-full">
				<li v-for="theme in themes" :key="theme" class="my-2">
					<input type="radio" name="theme-dropdown"
						class="theme-controller btn btn-block justify-start" :aria-label="theme"
						:value="theme" />
				</li>
			</ul>
		</div>-->
		<fieldset class="fieldset">
			<legend class="fieldset-legend">Username</legend>
			<input v-model="username" type="text" class="input w-full" placeholder="Name" />
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
			loading.value = false;
		}
		catch {
			console.log("Could not get username");
		}
	}
	async function triggerUpdate(newUsername) {
		console.log(newUsername)
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
					body: JSON.stringify({"username": newUsername})
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


	onMounted(getUsername);

	watch(
		() => ({
			username: username.value,
		}),
		(newValue, oldValue) => {
			if (oldValue.username == "") {
				return;
			}
			if (newValue.username !== "") {
				triggerUpdate(newValue.username);
			}
		},
	);

</script>
