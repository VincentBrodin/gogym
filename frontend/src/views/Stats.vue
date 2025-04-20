<template>
	<div v-if="loading" class="w-full h-full flex justify-center items-center">
		<p>Loading...</p>
	</div>
	<div v-else class="w-full min-h-full p-8">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Stats</h1>
			<h1 class="text-left text-xl font-bold text-success mt-2">
				<i class="bi bi-arrow-up"></i>
				{{Math.round(stats.totalPercentGain)}}%
			</h1>
		</div>

		<transition-group name="session" tag="div" class="w-full">
			<StatItem v-for="session in sessions" :key="session.id" :session="session" @remove="removeSession" />
		</transition-group>
	</div>
</template>

<script setup>
	import StatItem from '@/components/StatItem.vue';
	import {ref, onMounted} from 'vue';

	const sessions = ref(null);
	const stats = ref(null);
	const loading = ref(true);

	async function loadSessions() {
		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/sessions`;
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
			sessions.value = await response.json();
			sessions.value.sort((a, b) => Date.parse(b.endend_at) - Date.parse(a.endend_at));
			stats.value = computeStats(sessions.value)
			grabImprovemt(sessions.value)
			console.log(sessions.value)
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	async function removeSession(session) {
		sessions.value.splice(sessions.value.indexOf(session), 1);

		const token = localStorage.getItem('token');
		const url = `${import.meta.env.VITE_API_URL}api/restricted/session/${session.id}`;

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
			console.log(await response.json());
		}
		catch (error) {
			console.error('Error loading workouts:', error);
		} finally {
			loading.value = false;
		}
	}

	function computeStats(sessions) {
		const byExercise = {};
		sessions.forEach(session => {
			const date = new Date(session.started_at);
			session.exercise_sessions.forEach(es => {
				es.exercise_weights.forEach(w => {
					if (w.weight === 0) return;            // skip zeros
					const id = es.exercise.id;
					if (!byExercise[id]) {
						byExercise[id] = {
							name: es.exercise.name,
							records: []
						};
					}
					byExercise[id].records.push({date, weight: w.weight});
				});
			});
		});

		const exerciseStats = {};
		Object.entries(byExercise).forEach(([id, {name, records}]) => {
			if (records.length < 2) return;           // not enough data
			records.sort((a, b) => a.date - b.date);
			const first = records[0];
			const last = records[records.length - 1];
			const kgGain = last.weight - first.weight;
			const percentGain = first.weight > 0
				? (kgGain / first.weight) * 100
				: 0;

			exerciseStats[id] = {
				name,
				firstDate: first.date,
				lastDate: last.date,
				firstWeight: first.weight,
				lastWeight: last.weight,
				kgGain,
				percentGain
			};
		});

		const statsArr = Object.values(exerciseStats);
		if (statsArr.length === 0) {
			return {averageKgPerWeek: 0, totalPercentGain: 0, exercises: {}};
		}
		const firstDates = statsArr.map(s => s.firstDate.getTime());
		const lastDates = statsArr.map(s => s.lastDate.getTime());
		const earliest = Math.min(...firstDates);
		const latest = Math.max(...lastDates);
		const weeks = (latest - earliest) / (1000 * 60 * 60 * 24 * 7) || 1;

		const totalKgGain = statsArr.reduce((sum, s) => sum + s.kgGain, 0);
		const totalPercentGain = statsArr.reduce((sum, s) => sum + s.percentGain, 0)
			/ statsArr.length;

		return {
			averageKgPerWeek: totalKgGain / weeks,
			totalPercentGain,
			exercises: exerciseStats
		};
	}

	function grabImprovemt(list) {
		console.log(list)
		const groupBy = (arr, keyFn) =>
			arr.reduce((acc, item) => {
				const key = keyFn(item);
				(acc[key] ??= []).push(item);
				return acc;
			}, {});

		const result = groupBy(list, ({workout}) => workout.id);
		console.log(result)
		for (let id in result) {
			const order = result[id].sort((a, b) => Date.parse(a.endend_at) - Date.parse(b.endend_at))
			let last = null;
			for (let session of order) {
				if (last != null) {
					let gainTotal = 0;
					let count = 0;
					for (let lex of last.exercise_sessions) {
						for (let cex of session.exercise_sessions) {
							if (lex.exercise.id == cex.exercise.id) {
								const lmax = getMaxWeight(lex);
								const cmax = getMaxWeight(cex);
								count++;
								if (lmax === 0 || cmax === 0 || lmax === cmax) {
									continue;
								}
								gainTotal += ((cmax - lmax) / lmax) * 100;
							}
						}
					}
					if (count > 0) {
						session.gain = Math.round(gainTotal / count);
					} else {
						session.gain = 0;
					}
				}
				else {
					session.gain = 0;
				}
				last = session
			}
		}
	}

	function getMaxWeight(exercise) {
		const weights = exercise.exercise_weights
			.map(item => Number(item.weight))
			.filter(w => !isNaN(w));

		return weights.length
			? Math.max(...weights)
			: 0;  // or null, depending on how you want to handle “no data”
	}
	onMounted(loadSessions);

</script>
