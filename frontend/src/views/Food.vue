<template>
	<!--<StreamBarcodeReader @decode="onDecode" @loaded="onLoaded"></StreamBarcodeReader>-->

	<div class="w-full h-full p-8">
		<div class="w-full flex flex-row justify-center items-center">
			<div class="radial-progress" :style="{ '--value': progress }" :aria-valuenow="progress" role="progressbar">
				{{ progress }}%
			</div>
		</div>

		<div>
			<fieldset class="fieldset">
				<legend class="fieldset-legend">kcals per 100g</legend>
				<input v-model="kcals" type="number" class="input w-full" placeholder="kcals" />
			</fieldset>

			<fieldset class="fieldset">
				<legend class="fieldset-legend">Weight (g)</legend>
				<input v-model="weight" type="number" class="input w-full" placeholder="Weight" />
			</fieldset>
		</div>
		<p>{{(kcals/100)*weight}}kcals</p>
	</div>


</template>

<script setup>
	import {ref, onMounted} from 'vue'

	const progress = ref(0);
	const back = ref(false);
	const kcals = ref(0);
	const weight = ref(0);

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
