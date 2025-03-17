<template>
	<div class="dock dock-lg">
		<router-link to="/">
			<button class="flex flex-col justify-between items-center">
				<i class="bi bi-house text-xl"></i>
				<span class="dock-label">Home</span>
			</button>
		</router-link>

		<router-link to="/">
			<button class="flex flex-col justify-between items-center">
				<i class="bi bi-clipboard-data text-xl"></i>
				<span class="dock-label">Stats</span>
			</button>
		</router-link>

		<router-link to="/">
			<button class="flex flex-col justify-between items-center">
				<i class="bi bi-gear text-xl"></i>
				<span class="dock-label">Settings</span>
			</button>
		</router-link>
	</div>
</template>
