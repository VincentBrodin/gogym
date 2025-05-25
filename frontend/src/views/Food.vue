<template>
	<div class="w-full h-full p-8">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Food</h1>
		</div>
		<div class="w-full grid grid-cols-1 lg:grid-cols-3 gap-4 mb-8">
			<div class="flex flex-row p-4 bg-green-50 rounded-2xl border border-green-200 gap-2 shadow-md">
				<div
					class="bg-green-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-bullseye text-green-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Food Items</p>
					<p class="text-2xl font-bold">{{foodItems.length}}</p>
				</div>
			</div>
			<div class="flex flex-row p-4 bg-orange-50 rounded-2xl border border-orange-200 gap-2 shadow-md">
				<div
					class="bg-orange-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-graph-up-arrow text-orange-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Total Calories</p>
					<p class="text-2xl font-bold">{{totalKcals.toFixed(1)}}</p>
				</div>
			</div>
			<div class="flex flex-row p-4 bg-purple-50 rounded-2xl border border-purple-200 gap-2 shadow-md">
				<div
					class="bg-purple-200 p-1.5 mr-2 h-full aspect-square text-center rounded-2xl flex justify-center items-center">
					<i class="bi bi-calculator text-purple-500 text-2xl"></i>
				</div>
				<div>
					<p class="opacity-60">Avg per Item</p>
					<p v-if="foodItems.length != 0" class="text-2xl font-bold">
						{{(totalKcals/foodItems.length).toFixed(1)}}</p>
					<p v-else class="text-2xl font-bold">0.0</p>
				</div>
			</div>
		</div>
		<div class="flex flex-row justify-between">
			<h2 class="font-bold text-2xl mb-4">Food Items</h2>
			<p class="opacity-60"><span class="font-bold">{{totalKcals}} kcal</span> total</p>
		</div>
		<div v-if="foodItems.length != 0" class="w-full flex flex-col gap-4 pb-40">
			<FoodItem v-for="foodItem in foodItems" :key="foodItem.id" :foodItem="foodItem" @remove="remove" />

			<div class="card card-border w-full bg-gradient-to-r from-emerald-600 to-teal-600 shadow-xl">
				<div class="card-body">
					<p class="text-2xl text-center text-white font-bold">
						<i class="bi bi-calculator"></i>
						Total Meal Calories
					</p>
					<p class="text-2xl text-center text-white font-bold">{{totalKcals.toFixed(1)}}</p>
					<p class="text-sm opacity-60 text-white text-center">kcal</p>
					<p class="text-sm opacity-60 text-white text-center">From <span
							class="opacity-100 font-bold">{{foodItems.length}}</span> Food Items</p>
				</div>
			</div>
		</div>
		<div v-else>
			<div class="card card-border w-full shadow-xl py-6">
				<div class="card-body">
					<p class="text-2xl text-center font-bold">
						No food items yet
					</p>
					<p class="text-sm opacity-60 text-center">Add your first food item to start tracking calories</p>
				</div>
			</div>
		</div>
	</div>
	<AddFoodItemModal @add-food-item="addFoodItem" />
</template>

<script setup>
	import {useLocalStorage} from '@vueuse/core'
	import {computed} from 'vue';
	import FoodItem from '@/components/FoodItem.vue';
	import AddFoodItemModal from '@/components/AddFoodItemModal.vue';

	const foodItems = useLocalStorage("foodItems", []);

	function addFoodItem(foodItem) {
		foodItems.value.push(foodItem);
	}

	function remove(foodItem) {
		const index = foodItems.value.indexOf(foodItem);
		if (index !== -1) {
			foodItems.value.splice(index, 1);
		}
	}

	const totalKcals = computed(() => {
		let value = 0;
		for (let foodItem of foodItems.value) {
			value += (foodItem.kcal / foodItem.per) * foodItem.weight;
		}
		return value
	})
</script>
