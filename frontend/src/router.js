import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';
import Login from './views/Login.vue';
import Edit from './views/Edit.vue';
import Register from './views/Register.vue';
import Session from './views/Session.vue';
import Settings from './views/Settings.vue';
import Copy from './views/Copy.vue';
import Stats from './views/Stats.vue';
import Food from './views/Food.vue';
import ViewSession from './views/ViewSession.vue';
import { jwtDecode } from 'jwt-decode';

const routes = [
	{ path: '/login', name: 'login', component: Login },
	{ path: '/register', name: 'register', component: Register },
	{ path: '/', name: 'home', component: Home, meta: { requiresAuth: true } },
	{ path: '/edit', name: 'edit', component: Edit, meta: { requiresAuth: true } },
	{ path: '/session', name: 'session', component: Session, meta: { requiresAuth: true } },
	{ path: '/stats/view', name: 'view', component: ViewSession, meta: { requiresAuth: true } },
	{ path: '/settings', name: 'settings', component: Settings, meta: { requiresAuth: true } },
	{ path: '/copy', name: 'copy', component: Copy, meta: { requiresAuth: true } },
	{ path: '/stats', name: 'stats', component: Stats, meta: { requiresAuth: true } },
	{ path: '/food', name: 'food', component: Food },
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
