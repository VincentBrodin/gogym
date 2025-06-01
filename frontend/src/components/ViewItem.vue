<template>
	<div class="card card-border bg-base-100 w-full shadow-xl">
		<div class="card-body">
			<div class="flex flex-row items-center">
				<div class="bg-blue-200 p-1.5 mr-2 w-7 h-7 text-center rounded-md flex justify-center items-center">
					<i class="bi bi-bullseye text-blue-500 text-lg"></i>
				</div>
				<p class="text-lg font-bold">{{props.exercise.exercise.name}}</p>
			</div>
			<div class="overflow-x-auto">
				<table class="table text-center">
					<!-- head -->
					<thead>
						<tr>
							<th>Set</th>
							<th>Weight</th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="weight in props.exercise.exercise_weights">
							<th>{{weight.set}}</th>
							<td>{{weight.weight}} {{type}}</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>
</template>

<script setup>
	const props = defineProps({
		exercise: {
			type: Object,
			required: true
		},
		imperial: {
			type: Boolean,
			required: true
		}
	});
	console.log(props.exercise)

	const type = props.imperial ? "lbs" : "kg"
</script>
