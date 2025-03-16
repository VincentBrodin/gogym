<script>
	const API_URL = `http://localhost:8080/api/restricted/account`
	const TOKEN = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywidW5hbWUiOiJ2aW50ZSIsImV4cCI6MTc0MjEyMzQxNH0.GAbRQGc6C1IUpc_2NXARFXuHBTIHeJvtlNotcg7HQYo`

	export default {
		data: () => ({
			msg_body: '',
			msg_name: '',
		}),

		created() {
			// fetch on init
			this.fetchData()
		},

		methods: {
			async fetchData() {
				const url = `${API_URL}`
				const json = await (await fetch(url, {
					method: "GET", 
					headers: {Authorization: `Bearer ${TOKEN}`}
				})).json();
				console.log(json)
			},
		}
	}
</script>

<template>
	<h1 class="text-3xl">{{msg_body}} {{msg_name}}</h1>
</template>

<style>
	a {
		text-decoration: none;
		color: #42b883;
	}

	li {
		line-height: 1.5em;
		margin-bottom: 20px;
	}

	.author,
	.date {
		font-weight: bold;
	}
</style>
