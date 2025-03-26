<template>
	<div v-if="session == null" class="flex w-full h-full items-center justify-center">
		<p class="text-error text-xl font-bold">Error: Session not found</p>
	</div>
	<div v-else-if="done" class="flex w-full h-full items-center justify-center flex-col p-8">
		<h1 class="font-bold text-2xl mb-8">{{session.workout.name}} completed</h1>
		<p class="text-center text-5xl font-extrabold">
			{{ currentTime }}
		</p>
		<Teleport to="body">
			<div class="fixed px-8 bottom-24 w-full">
				<button class="w-full btn btn-primary shadow" @click="home">
					<p>Done</p>
				</button>
			</div>
		</Teleport>

	</div>

	<template v-else>
		<div class="p-8">
			<h1 class="text-center text-3xl font-bold mb-4">
				{{ session.workout.name }}
			</h1>

			<transition name="slide-fade" mode="out-in">
				<h2 class="text-center text-2xl font-semibold" :key="currentExercise.exercise.name">
					{{ currentExercise.exercise.name }}
				</h2>
			</transition>

			<p class="text-center text-xl mb-4 opacity-60">
				{{ currentExercise.sets_done }} / {{ currentExercise.exercise.sets }} x
				{{currentExercise.exercise.reps}}
			</p>
			<p class="text-center text-xl mb-12 opacity-60">
				{{ currentExercise.exercise.rir }} RIR
			</p>

			<p class="text-center text-5xl font-extrabold">
				{{ currentTime }}
			</p>
		</div>


		<SessionDock @next="next" @skip="skip" />
	</template>
</template>

<script setup>
	import SessionDock from "@/components/SessionDock.vue";
	import {ref, onMounted, onUnmounted} from "vue";
	import {useRouter} from "vue-router";
	import {useLocalStorage} from '@vueuse/core'

	const router = useRouter();

	const currentTime = ref("00:00:00");
	let timer;

	const session = ref(null);
	const done = ref(false);
	const hasSession = useLocalStorage('session', false)
	const currentExercise = ref(null);

	function home() {
		router.push({name: "home"});
	}

	async function next() {
		currentExercise.value.sets_done++;
		if (currentExercise.value.sets_done > currentExercise.value.exercise.sets) {
			currentExercise.value.sets_done = currentExercise.value.exercise.sets;
			currentExercise.value.completed = true;
			currentExercise.value.active = false;
			currentExercise.value.skiped = false;

			const nextExercise = grabNext(session.value.exercise_sessions);

			// Completed
			if (nextExercise == null) {
				session.value.active = false;
				hasSession.value = false;
				done.value = true;
				clearInterval(timer);
				currentTime.value = "Grabbing time"
				await update(true);
				updateTime(Date.parse(session.value.endend_at) - Date.parse(session.value.started_at));
				return;
			} else {
				currentExercise.value = nextExercise;
				currentExercise.value.active = true;
				currentExercise.value.sets_done = 1;
				currentExercise.value.skiped = false;
			}
		}
		await update(false);
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
		const token = localStorage.getItem("token");
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session`;

		try {
			const response = await fetch(url, {
				method: "GET",
				headers: {
					"Content-Type": "application/json",
					Authorization: "Bearer " + token,
				},
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}
			session.value = await response.json();
			updateTime(Date.now() - Date.parse(session.value.started_at));
			// Grab the active exercise
			currentExercise.value = session.value.exercise_sessions.find(
				(exercise) => exercise.active
			);
			// If no active, grab the next
			if (currentExercise.value == null) {
				currentExercise.value = grabNext(session.value.exercise_sessions);
				currentExercise.value.sets_done = 1;
				currentExercise.value.active = true;
				await update(false);
			}
		} catch (error) {
			console.log(error);
		}
	}

	function grabNext(list) {
		list.sort((a, b) => a.exercise.order - b.exercise.order);
		let result = list.filter((exercise) => !exercise.skiped && !exercise.completed);
		if (result.length === 0) {
			result = list.filter((exercise) => !exercise.completed);
		}
		return result.at(0);
	}

	async function update(override) {
		const token = localStorage.getItem("token");
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${session.value.id}`;

		try {
			const response = await fetch(url, {
				method: "PATCH",
				headers: {
					"Content-Type": "application/json",
					Authorization: "Bearer " + token,
				},
				body: JSON.stringify(session.value),
			});
			if (!response.ok) {
				throw new Error((await response.json()).error);
			}
			if (override) {
				session.value = await response.json();
			}
			else {
				console.log(await response.json());
			}
		} catch (error) {
			console.log(error);
		}
	}

	function updateTime(time) {
		const totalSeconds = Math.floor(time / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		currentTime.value = `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
			2,
			"0"
		)}:${String(seconds).padStart(2, "0")}`;
	}

	onMounted(async () => {
		await grabSession();
		timer = setInterval(() => {
			updateTime(Date.now() - Date.parse(session.value.started_at));
		}, 750);
	});

	onUnmounted(() => {
		clearInterval(timer);
	});
</script>

<style scoped>
	.slide-fade-enter-active {
		transition: all 0.5s ease;
	}

	.slide-fade-enter-from {
		opacity: 0;
		transform: translateX(20px);
	}
</style>
