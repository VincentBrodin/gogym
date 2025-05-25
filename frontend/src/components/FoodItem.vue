<template>
	<div class="card card-border bg-base-100 w-full shadow-xl">
		<div class="card-body">
			<div class="flex flex-row justify-between items-start">
				<input type="text" v-model="foodItem.name" class="input input-ghost font-bold text-xl px-0"
					placeholder="Food Item">
				<div class="flex flex-row gap-4">
					<div class="flex flex-col">
						<p class="text-end text-xl font-bold text-success">{{ kcals.toFixed(1)}}</p>
						<p class="text-end text-sm opacity-60">kcal</p>
					</div>
					<button class="btn btn-ghost btn-error btn-circle" @click="emitRemove">
						<i class="bi bi-trash text-error"></i>
					</button>
				</div>
			</div>
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-4 align-middle w-full">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Calories (kcal)</legend>
					<input type="number" class="input" placeholder="165" v-model="foodItem.kcal" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Per (grams)</legend>
					<input type="number" class="input" placeholder="100" v-model="foodItem.per" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Weight (grams)</legend>
					<input type="number" class="input" placeholder="200" v-model="foodItem.weight" />
				</fieldset>
			</div>

			<div class="card card-border bg-base-100 w-full">
				<div class="card-body">
					<p>Calcultions:</p>
					<p class="text-sm font-mono opacity-90">({{foodItem.kcal}} kcal ÷ {{foodItem.per}}) × {{foodItem.weight}}g = {{kcals}} kcal</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
	import {defineProps, defineEmits, computed} from 'vue';

	const emits = defineEmits(['remove']);

	const props = defineProps({
		foodItem: {
			type: Object,
			required: true
		}
	});

	const kcals = computed(() => {
		return (props.foodItem.kcal / props.foodItem.per) * props.foodItem.weight;
	});

	const emitRemove = () => {
		emits('remove', props.foodItem);
	};
</script>
