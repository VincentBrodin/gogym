<template>
	<div class="p-8">
		<div v-if="loading" class="w-full h-full flex justify-center items-center">
			<p>Loading...</p>
		</div>
		<div v-else class="flex flex-col gap-4">
			<div class="card card-border bg-base-100 w-full shadow-xl">
				<div class="card-body">
					<div class="gap-4 grid grid-cols-1 lg:grid-cols-2">
						<div>
							<p class="text-xl font-bold mb-2">{{session.workout.name}}</p>
							<p class="opacity-60 mb-2">
								<i class="bi bi-calendar"></i>
								{{niceDate}}
							</p>
							<p class="opacity-60">
								<i class="bi bi-clock"></i>
								{{startTime}}
								-
								{{endTime}}
							</p>
						</div>
						<div class="grid grid-cols-2 lg:grid-cols-4 align-middle gap-4">
							<div class="text-center">
								<p class="text-xl font-bold text-blue-600">{{time}}</p>
								<p class="opacity-60">Duration</p>
							</div>
							<div class="text-center">
								<p class="text-xl font-bold text-green-600">10h</p>
								<p class="opacity-60">Total Volume</p>
							</div>
							<div class="text-center">
								<p class="text-xl font-bold text-purple-600">{{sets}}</p>
								<p class="opacity-60">Total Sets</p>
							</div>
							<div class="text-center">
								<p class="text-xl font-bold text-orange-600">{{reps}}</p>
								<p class="opacity-60">Total Reps</p>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="flex flex-row justify-between">
				<p class="text-lg font-bold">Exercises</p>
				<p class="opacity-60">{{session.exercise_sessions.length}} exercises</p>
			</div>
			<div class="flex flex-col gap-4 pb-10">
				<ViewItem v-for="exercise in session.exercise_sessions" :key="exercise.id" :exercise="exercise" :imperial="imperial" />
			</div>
		</div>
	</div>
</template>

<script setup>
	import ViewItem from '@/components/ViewItem.vue';
	import {jwtDecode} from 'jwt-decode';
	import {onMounted, ref} from 'vue';
	import {useRoute} from 'vue-router'
	const route = useRoute()
	const id = route.query.id;

	const session = ref(null)

	const niceDate = ref(null)
	const startTime = ref(null)
	const endTime = ref(null)

	const time = ref(null)

	const sets = ref(null)
	const reps = ref(null)

	const imperial = ref(false)

	const loading = ref(true)

	async function loadSession() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${id}`;
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
			session.value.exercise_sessions = session.value.exercise_sessions.sort((a, b) => a.exercise.order - b.exercise.order)
			niceDate.value = new Date(session.value.started_at).toLocaleDateString('en-US', {
				weekday: 'long', // "Saturday"
				month: 'short',  // "Jan"
				day: 'numeric',  // "20"
			});

			const start = new Date(session.value.started_at)
			const end = new Date(session.value.endend_at)

			startTime.value = start.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});
			endTime.value = end.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit', hour12: false});


			const diffMs = Math.abs(end - start);
			const diffMinutes = Math.floor(diffMs / (1000 * 60));

			const hours = Math.floor(diffMinutes / 60);
			const minutes = diffMinutes % 60;

			if (hours == 0) {
				time.value = `${minutes}m`;
			} else {
				time.value = `${hours}h ${minutes}m`;
			}

			sets.value = session.value.exercise_sessions.reduce((sum, exerciseSession) => {
				return sum + (exerciseSession.sets_done || 0);
			}, 0);
			reps.value = session.value.exercise_sessions.reduce((sum, exSession) => {
				const sets = exSession.sets_done || 0;
				const reps = exSession.exercise.reps || 0;
				return sum + sets * reps;
			}, 0);

		}
		catch (error) {
			console.error('Error loading session:', error);
		} finally {
			loading.value = false;
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


	onMounted(async () => {
		await loadSession();
		getImperial();
	})

</script>
