<template>
	<!--<StreamBarcodeReader @decode="onDecode"></StreamBarcodeReader>-->
	<div class="w-full h-full p-8">
		<div class="mb-4">
			<h1 class="text-left text-2xl font-bold">Food</h1>
		</div>
		<div>
			<div class="flex flex-row gap-4">
				<fieldset class="fieldset grow">
					<legend class="fieldset-legend">kcals</legend>
					<input v-model="kcals" type="number" min="1" class="input w-full" placeholder="kcals" />
				</fieldset>
				<fieldset class="fieldset grow">
					<legend class="fieldset-legend">per grams</legend>
					<input v-model="perWeight" type="number" min="1" class="input w-full" placeholder="per grams" />
				</fieldset>
			</div>

			<fieldset class="fieldset">
				<legend class="fieldset-legend">Weight (g)</legend>
				<input v-model="weight" type="number" min="1" class="input w-full" placeholder="Weight" />
			</fieldset>
		</div>
		<div class="flex flex-row gap-4 w-full my-8">
			<button class="btn" @click="promptClear">Clear</button>
			<button class="grow btn btn-primary" @click="add">Add: {{Math.round((kcals/Math.max(perWeight, 1))*weight)}}kcals</button>
		</div>
		<p class="text-xl">Total: {{totalKcals}}kcals</p>
	</div>

	<ConfirmationModal ref="confirmModal" promptText="Are you sure?"
		:detailText="`Are you sure that you want to clear the calories?`" confirmText="Yes" @confirmed="clear" />
</template>

<script setup>
	import ConfirmationModal from '@/components/ConfirmationModal.vue';
	import {ref, onMounted} from 'vue'

	import {useLocalStorage} from '@vueuse/core'

	const progress = ref(0);
	const back = ref(false);
	const kcals = ref(0);
	const perWeight = ref(100);
	const weight = ref(0);
	const confirmModal = ref(null);
	const totalKcals = useLocalStorage('kcals', 0)
	//function onDecode(value) {
	//	console.log(value)
	//	alert(value)
	//}

	function add() {
		totalKcals.value += Math.round((kcals.value / perWeight.value) * weight.value);
		kcals.value = 0;
		weight.value = 0;
		perWeight.value = 100;
	}

	function promptClear() {
		if (confirmModal) {
			confirmModal.value.open();
		}
	}

	function clear() {
		kcals.value = 0;
		weight.value = 0;
		perWeight.value = 100;
		totalKcals.value = 0;
	}


	onMounted(() => {
		setInterval(() => {
			if (back.value) {
				progress.value -= 1;
			}
			else {
				progress.value += 1;
			}

			if (progress.value >= 100 || progress.value <= 0) {
				back.value = !back.value;
			}
		}, 250)
	});
</script>
