<template>
	<div v-if="session == null" class="flex h-screen items-center justify-center bg-base-200">
		<p class="text-error text-xl font-bold">Error: Session not found</p>
	</div>

	<template v-else>
		<div class="p-8">
			<h1 class="text-center text-3xl font-bold mb-4">
				{{ session.workout.name }}
			</h1>
			<h2 class="text-center text-2xl font-semibold mb-4">
				{{ currentExercise.exercise.name }}
			</h2>
			<p class="text-center text-xl mb-12">
				{{ currentExercise.sets_done }} / {{ currentExercise.exercise.sets }}
			</p>
			<p class="text-center text-5xl font-extrabold ">
				{{ currentTime }}
			</p>

		</div>

		<SessionDock @next="next" @skip="skip" />
	</template>
</template>

<script setup>
	import SessionDock from "@/components/SessionDock.vue";
	import {ref, onMounted, onUnmounted} from "vue";
	import {useRouter} from 'vue-router'

	const router = useRouter()

	const currentTime = ref("00:00:00");
	let timer;

	const session = ref(null)
	const currentExercise = ref(null)

	async function next() {
		currentExercise.value.sets_done++;
		if (currentExercise.value.sets_done > currentExercise.value.exercise.sets) {
			currentExercise.value.sets_done = currentExercise.value.exercise.sets;
			currentExercise.value.completed = true;
			currentExercise.value.active = false;
			currentExercise.value.skiped = false;
			const next = grabNext(session.value.exercise_sessions)
			if (next == null) {
				session.value.active = false;
				router.push({
					name: 'home',
				});

			}
			else {
				currentExercise.value = next;
				currentExercise.value.active = true;
				currentExercise.value.sets_done = 1;
				currentExercise.value.skiped = false;
			}
		}
		await update();
	}

	async function skip() {
		currentExercise.value.sets_done = 0;
		currentExercise.value.active = false;
		currentExercise.value.skiped = true;

		currentExercise.value = grabNext(session.value.exercise_sessions);

		currentExercise.value.active = true;
		currentExercise.value.skiped = false;
		currentExercise.value.sets_done = 1;
		await update();
	}

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
			updateTime();
			// Grab the active exercise
			currentExercise.value = session.value.exercise_sessions.find((exercise) => exercise.active);
			// If no active
			if (currentExercise.value == null) {
				currentExercise.value = grabNext(session.value.exercise_sessions);
				currentExercise.value.sets_done = 1
				currentExercise.value.active = true
				await update();
			}
			//console.log(session.value)
		}
		catch (error) {
			console.log(error)
		}
	}

	function grabNext(list) {
		list.sort((a, b) => a.exercise.order - b.exercise.order);
		let result = list.filter((exercise) => !exercise.skiped && !exercise.completed);
		if (result.length == 0) {
			result = list.filter((exercise) => !exercise.completed);
		}
		return result.at(0);
	}

	async function update() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${session.value.id}`;

		try {
			const response = await fetch(url, {
				method: "PATCH",
				headers: {
					"Content-Type": "application/json",
					"Authorization": "Bearer " + token,
				},
				body: JSON.stringify(session.value)
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}
			console.log(await response.json());
		}
		catch (error) {
			console.log(error)
		}

	}

	function updateTime() {
		const time = Date.now() - Date.parse(session.value.started_at);
		const totalSeconds = Math.floor(time / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		currentTime.value = `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
	}

	onMounted(async () => {
		await grabSession();
		timer = setInterval(() => {
			updateTime();
		}, 750);
	});

	onUnmounted(() => {
		clearInterval(timer);
	});
</script>
