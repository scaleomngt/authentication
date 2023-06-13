import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from '../views/Home.vue';
import qrcode from '../views/qrcode.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },{
    path: '/qrcode',
    name: 'qrcode',
    component: qrcode,
  },
];

const router = new VueRouter({
  routes,
});

export default router;
