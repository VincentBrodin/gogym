import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';
import About from './views/About.vue';
import Login from './views/Login.vue';
import Edit from './views/Edit.vue';
import Register from './views/Register.vue';
import { jwtDecode } from 'jwt-decode';

const routes = [
	{ path: '/', name: 'Home', component: Home, meta: { requiresAuth: true } },
	{ path: '/about', name: 'About', component: About },
	{ path: '/login', name: 'login', component: Login },
	{ path: '/register', name: 'register', component: Register },
	{ path: '/edit', name: 'edit', component: Edit }
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});


function hasValidToken() {
	const token = localStorage.getItem('token');
	if (!token) return false;

	try {
		const decoded = jwtDecode(token);
		const currentTime = Date.now() / 1000;
		return decoded.exp > currentTime;
	}
	catch (error) {
		console.error("Invalid token:", error)
		localStorage.removeItem("token")
		return false
	}
}


router.beforeEach((to, _, next) => {
	if (to.meta.requiresAuth && !hasValidToken()) {
		next({
			name: 'login',
			query: { redirect: to.fullPath }
		})
	} else {
		next()
	}
})

export default router;
