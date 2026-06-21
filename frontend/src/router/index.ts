import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import DashboardLayout from "../components/DashboardLayout.vue";
import Dashboard from "../views/Dashboard.vue";
import Analytics from "../views/Analytics.vue";
import Versions from "../views/Versions.vue";
import Files from "../views/Files.vue";
import Settings from "../views/Settings.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "Login",
      component: Login,
    },
    {
      path: "/dashboard",
      component: DashboardLayout,
      meta: { title: "概述" },
      children: [
        {
          path: "",
          name: "Dashboard",
          component: Dashboard,
          meta: { title: "概述" },
        },
        {
          path: "versions",
          name: "Versions",
          component: Versions,
          meta: { title: "版本管理" },
        },
        {
          path: "files",
          name: "Files",
          component: Files,
          meta: { title: "文件管理" },
        },
        {
          path: "analytics",
          name: "Analytics",
          component: Analytics,
          meta: { title: "数据分析" },
        },
        {
          path: "settings",
          name: "Settings",
          component: Settings,
          meta: { title: "设置" },
        },
      ],
    },
    {
      path: "/",
      redirect: "/login",
    },
  ],
});

export default router;
