<template>
	<label class="floating-label my-4">
		<span>{{ label }}</span>
		<input :value="modelValue" @input="$emit('update:modelValue', $event.target.value)" :type="type"
			:placeholder="label" class="input w-full" :class="{ 'input-error': error != null }" />
		<p v-if="error != null && showMsg" class="text-error mb-4 text-sm font-bold mt-1">{{ error }}</p>
	</label>
</template>

<script setup>
	const props = defineProps({
		modelValue: {type: String, default: ''},
		label: {type: String, required: true},
		error: {type: String, default: null},
		type: {type: String, default: 'text'},
		showMsg: {type: Boolean, default: true}
	});

	const emit = defineEmits(['update:modelValue']);
</script>
