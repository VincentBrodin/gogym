<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else-if="session == null" class="flex w-full h-full items-center justify-center">
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
		<div class="p-4 flex flex-col gap-4 pb-30">
			<div class="card card-border bg-gradient-to-br from-blue-500/5 to-indigo-500/5 w-full shadow-xl">
				<div class="card-body">
					<p class="text-lg font-bold"><i class="bi bi-bullseye text-blue-300"></i> Current Exercise</p>
					<select class="select w-full font-semibold text-xl text-center" @change="switchExercise"
						ref="selectRef">
						<option v-for="item, i in session.exercise_sessions" :key="item.exercise.id"
							:selected="item.active" :value="i">
							{{item.exercise.name}}</option>
					</select>
				</div>
			</div>
			<div class="card card-border bg-gradient-to-br from-emerald-500/5 to-teal-500/5 w-full shadow-xl">
				<div class="card-body">
					<p class="text-center text-3xl font-bold" :key="currentExercise.sets_done">
						{{ currentExercise.sets_done }} / {{ currentExercise.exercise.sets }} x
						{{currentExercise.exercise.reps}}
					</p>
					<p class="text-center text-lg font-bold text-green-600" :key="currentExercise.exercise.rir">
						{{ currentExercise.exercise.rir }} RIR
					</p>
					<progress class="progress w-full" :value="progress" max="100"></progress>
				</div>
			</div>
			<div class="card card-border bg-gradient-to-br from-blue-500/5 to-indigo-500/5 w-full shadow-xl">
				<div class="card-body">
					<p class="text-lg font-bold text-center"><i class="bi bi-clock text-blue-600"></i> Session Time</p>
					<p class="text-center text-3xl font-bold text-blue-600">
						{{ currentTime }}
					</p>
				</div>
			</div>
			<div class="card card-border bg-gradient-to-br from-purple-500/5 to-pink-500/5 w-full shadow-xl">
				<div class="card-body">
					<p class="text-lg font-bold text-center"><i class="bi bi-pause-fill text-purple-600"></i> Rest Time
					</p>
					<p class="text-center text-3xl font-bold text-purple-600">
						{{ timeoutTime }}
					</p>
					<div class="w-full flex flex-row justify-evenly">
						<button class="btn text-purple-600 btn-outline" @click="updateTimeout(30)">
							+ 30s
						</button>
						<button class="btn text-purple-600 btn-outline" @click="updateTimeout(60)">
							+ 1m
						</button>
						<button class="btn text-purple-600 btn-outline" @click="updateTimeout(120)">
							+ 2m
						</button>
						<button class="btn text-purple-600 btn-outline" @click="updateTimeout(180)">
							+ 3m
						</button>
					</div>
				</div>
			</div>

			<div class="card card-border bg-gradient-to-br from-orange-500/5 to-red-500/5 w-full shadow-xl">
				<div class="card-body">
					<div class="flex items-center justify-center gap-3"><svg xmlns="http://www.w3.org/2000/svg"
							width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
							stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
							class="lucide lucide-weight h-6 w-6 text-orange-600">
							<circle cx="12" cy="5" r="3"></circle>
							<path
								d="M6.5 8a2 2 0 0 0-1.905 1.46L2.1 18.5A2 2 0 0 0 4 21h16a2 2 0 0 0 1.925-2.54L19.4 9.5A2 2 0 0 0 17.48 8Z">
							</path>
						</svg>
						<div class="text-center">
							<div class="text-4xl font-bold text-gray-900">
								{{currentExercise.exercise_weights[currentExercise.sets_done - 1].weight}}</div>
							<div v-if="imperial" class="text-lg text-orange-600 font-semibold">LB</div>
							<div v-else class="text-lg text-orange-600 font-semibold">KG</div>
						</div>
					</div>
					<div class="w-full grid grid-cols-4 align-middle gap-4 mb-4">
						<button v-for="weight in weights" :key="`primary-${weight}`"
							class="btn text-success btn-outline" @click="updateWeight(weight)">
							+{{ weight }}
						</button>
					</div>
					<div class="w-full grid grid-cols-4 align-middle gap-4">
						<button v-for="weight in weights" :key="`error-${weight}`" class="btn text-error btn-outline"
							@click="updateWeight(-weight)">
							-{{ weight }}
						</button>
					</div>
				</div>
			</div>
		</div>
		<SessionDock @next="next" @skip="skip" />
	</template>
</template>

<script setup>
	import SessionDock from "@/components/SessionDock.vue";
	import {ref, onMounted, onUnmounted} from "vue";
	import {jwtDecode} from 'jwt-decode';
	import {useRouter} from "vue-router";
	import {useLocalStorage} from '@vueuse/core'

	const router = useRouter();

	const currentTime = ref("00:00:00");
	const timeoutTime = ref("00:00:00");
	const timeout = useLocalStorage("timeout", Date.now());
	const weights = ref([1, 1.25, 2.5, 5, 10, 15, 20, 25]);

	let timer;


	const loading = ref(true);
	const session = ref(null);
	const done = ref(false);
	const progress = ref(0);
	const selectRef = ref(null)
	const hasSession = useLocalStorage('session', false)
	const currentExercise = ref(null);
	const imperial = ref(false);

	async function updateTimeout(time) {
		const done = (timeout.value - Date.now()) <= 0
		if (done) {
			timeout.value = Date.now() + time * 1000;
		}
		else {
			timeout.value += time * 1000;
		}

		timeoutTime.value = updateTime(timeout.value - Date.now());
	}

	async function updateWeight(weight) {
		currentExercise.value.exercise_weights[currentExercise.value.sets_done - 1].weight += weight
		if (currentExercise.value.exercise_weights[currentExercise.value.sets_done - 1].weight < 0) {
			currentExercise.value.exercise_weights[currentExercise.value.sets_done - 1].weight = 0
		}
		await update(false);
	}

	function home() {
		router.push({name: "home"});
	}

	async function switchExercise() {
		currentExercise.value.active = false;
		currentExercise.value.skiped = true;
		currentExercise.value.sets_done = 1
		const i = Number(selectRef.value.value)

		console.log(currentExercise.value.exercise.name)
		currentExercise.value = session.value.exercise_sessions[i]
		console.log(currentExercise.value.exercise.name)

		currentExercise.value.active = true
		currentExercise.value.skiped = false
		currentExercise.value.completed = false
		currentExercise.value.sets_done = 1
		updateProgress();
		await update(false);
	}

	async function next() {
		currentExercise.value.sets_done++;
		if (currentExercise.value.sets_done > currentExercise.value.exercise.sets) {
			currentExercise.value.sets_done = currentExercise.value.exercise.sets;
			currentExercise.value.completed = true;
			currentExercise.value.active = false;
			currentExercise.value.skiped = false;
			updateProgress();

			const nextExercise = grabNext(session.value.exercise_sessions);

			// Completed
			if (nextExercise == null) {
				session.value.active = false;
				hasSession.value = false;
				done.value = true;
				clearInterval(timer);
				currentTime.value = "Grabbing time"
				await update(true);
				currentTime.value = updateTime(Date.parse(session.value.endend_at) - Date.parse(session.value.started_at));
				return;
			} else {
				currentExercise.value = nextExercise;
				currentExercise.value.active = true;
				currentExercise.value.sets_done = 1;
				currentExercise.value.skiped = false;
			}
		}
		else if (currentExercise.value.sets_done != 1) {
			const lastWeight = currentExercise.value.exercise_weights[currentExercise.value.sets_done - 2].weight;
			currentExercise.value.exercise_weights[currentExercise.value.sets_done - 1].weight = lastWeight;
		}
		await update();
	}

	async function skip() {
		console.log(session.value.exercise_sessions)
		currentExercise.value.sets_done = 0;
		currentExercise.value.active = false;
		currentExercise.value.skiped = true;

		const old = currentExercise.value;
		currentExercise.value = grabNext(session.value.exercise_sessions);
		if (currentExercise.value == null || currentExercise.value == old) {
			console.log("All null")
			for (let exercise of session.value.exercise_sessions) {
				exercise.skiped = false;
			}
			old.skiped = true;
		}

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
			session.value.exercise_sessions.sort(function (a, b) {
				return a.exercise.order - b.exercise.order;
			});
			console.log(session.value)
			currentTime.value = updateTime(Date.now() - Date.parse(session.value.started_at));
			// Grab the active exercise
			currentExercise.value = session.value.exercise_sessions.find(
				(exercise) => exercise.active
			);

			for (let exercise_session of session.value.exercise_sessions) {
				exercise_session.exercise_weights.sort((a, b) => a.order - b.order);
				for (let exercise_weight of exercise_session.exercise_weights) {
					exercise_weight.init_weight = exercise_weight.weight;
				}
			}
			// If no active, grab the next
			if (currentExercise.value == null) {
				currentExercise.value = grabNext(session.value.exercise_sessions);
				currentExercise.value.sets_done = 1;
				currentExercise.value.active = true;

				await update(false);
			}
			updateProgress();
			loading.value = false;
		} catch (error) {
			console.log(error);
		}
	}

	function updateProgress() {
		const completed = session.value.exercise_sessions.filter(ex => ex.completed).length;
		progress.value = Math.round((completed / session.value.exercise_sessions.length) * 100)

	}

	function grabNext(list) {
		list.sort((a, b) => a.exercise.order - b.exercise.order);
		let result = list.filter((exercise) => !exercise.skiped && !exercise.completed);
		if (result.length === 0) {
			result = list.filter((exercise) => !exercise.completed);
		}
		if (result.length == 0) {
			return null;
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
		const totalSeconds = Math.max(Math.floor(time / 1000), 0);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(
			2,
			"0"
		)}:${String(seconds).padStart(2, "0")}`;
	}

	function getImperial() {
		const token = localStorage.getItem('token');
		try {
			const decoded = jwtDecode(token);
			imperial.value = decoded.imperial;
			if (imperial.value) {
				weights.value = [1, 2.5, 5, 10, 25, 35, 45, 55];
			}
		}
		catch {
			console.log("Could not get imperial");
		}
	}

	onMounted(async () => {
		await grabSession();
		currentTime.value = updateTime(Date.now() - Date.parse(session.value.started_at));
		timeoutTime.value = updateTime(timeout.value - Date.now());
		timer = setInterval(() => {
			currentTime.value = updateTime(Date.now() - Date.parse(session.value.started_at));
			timeoutTime.value = updateTime(timeout.value - Date.now());
		}, 500);
		getImperial();
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
