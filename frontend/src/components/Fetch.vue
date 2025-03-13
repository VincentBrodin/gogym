<script>
	const API_URL = `http://localhost:8080`

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
				const json = await (await fetch(url)).json();
				this.msg_body = json.message;
				this.msg_name = json.name;

				document.title = `${this.msg_body} ${this.msg_name}`
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
