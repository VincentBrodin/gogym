<template>
	<Teleport to="body">
		<div class="fixed px-8 bottom-24 w-full">
			<button class="w-full btn btn-primary" onclick="addFoodItem.showModal()">
				<i class="bi bi-plus text-2xl"></i>
				Add Food Item
			</button>
		</div>
	</Teleport>

	<dialog id="addFoodItem" class="modal modal-bottom sm:modal-middle">
		<div class="modal-box">
			<h3 class="text-lg font-bold">New Food Item</h3>
			<form method="dialog" class="w-full" @submit="submit">
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Name</legend>
					<input v-model="foodItem.name" type="text" class="input w-full" placeholder="Food Item" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Calories (kcal)</legend>
					<input v-model="foodItem.kcal" type="number" class="input w-full" placeholder="165" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Per (grams)</legend>
					<input v-model="foodItem.per" type="number" class="input w-full" placeholder="165" />
				</fieldset>
				<fieldset class="fieldset">
					<legend class="fieldset-legend">Weight (grams)</legend>
					<input v-model="foodItem.weight" type="number" class="input w-full" placeholder="165" />
				</fieldset>


				<button type="submit" class="btn btn-primary w-full mt-4">Add</button>
			</form>
		</div>
		<form method="dialog" class="modal-backdrop">
			<button>close</button>
		</form>

	</dialog>
</template>

<script setup>
	import {ref} from "vue";
	import {uuid} from "vue-uuid";
	const emit = defineEmits(["add-food-item"]);

	const foodItem = ref({
		name: "",
		kcal: 0,
		per: 100,
		weight: 0,
		id: uuid,
	});

	function submit() {
		emit("add-food-item", {...foodItem.value});
		foodItem.value.name = "";
		foodItem.value.weight = 0;
		foodItem.value.per = 0;
		foodItem.value.kcal = 0;
	}
</script>
