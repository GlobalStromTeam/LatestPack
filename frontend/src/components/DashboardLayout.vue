<script setup lang="ts">
import { ref, h, computed } from "vue";
import type { Component } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import { useAuthStore } from "../stores/auth";
import {
  HomeOutlined,
  SettingsOutlined,
  ExitToAppOutlined,
  AssessmentOutlined,
  SystemUpdateAltOutlined,
  FolderOutlined,
  CloudOutlined,
} from "@vicons/material";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const collapsed = ref(false);

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

const menuOptions = [
  {
    label: "概述",
    key: "dashboard",
    icon: renderIcon(HomeOutlined),
  },
  {
    label: "版本管理",
    key: "versions",
    icon: renderIcon(SystemUpdateAltOutlined),
  },
  {
    label: "文件管理",
    key: "files",
    icon: renderIcon(FolderOutlined),
  },
  {
    label: "渠道管理",
    key: "channels",
    icon: renderIcon(CloudOutlined),
  },
  {
    label: "数据分析",
    key: "analytics",
    icon: renderIcon(AssessmentOutlined),
  },
  {
    type: "divider",
    key: "d1",
  },
  {
    label: "设置",
    key: "settings",
    icon: renderIcon(SettingsOutlined),
  },
  {
    label: "退出登录",
    key: "logout",
    icon: renderIcon(ExitToAppOutlined),
  },
];

const activeKey = computed(() => {
  const path = route.path;
  if (path === "/dashboard") return "dashboard";
  if (path.startsWith("/dashboard/versions")) return "versions";
  if (path.startsWith("/dashboard/files")) return "files";
  if (path.startsWith("/dashboard/channels")) return "channels";
  if (path.startsWith("/dashboard/analytics")) return "analytics";
  if (path.startsWith("/dashboard/settings")) return "settings";
  return "dashboard";
});

function handleMenuSelect(key: string) {
  if (key === "logout") {
    authStore.logout();
    router.push("/login");
    return;
  }
  const routeMap: Record<string, string> = {
    dashboard: "/dashboard",
    versions: "/dashboard/versions",
    files: "/dashboard/files",
    channels: "/dashboard/channels",
    analytics: "/dashboard/analytics",
    settings: "/dashboard/settings",
  };
  if (routeMap[key]) {
    router.push(routeMap[key]);
  }
}
</script>

<template>
  <n-layout has-sider class="h-screen">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="220"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      :native-scrollbar="false"
      class="flex flex-col"
    >
      <div class="flex items-center justify-center h-16 px-4 border-b border-gray-100">
        <span
          v-if="!collapsed"
          class="text-lg font-bold tracking-wide whitespace-nowrap"
        >
          LatestPack
        </span>
        <span v-else class="text-lg font-bold">L</span>
      </div>

      <div class="flex-1 py-2">
        <n-menu
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
          :value="activeKey"
          @update:value="handleMenuSelect"
        />
      </div>
    </n-layout-sider>

    <n-layout class="flex-1">
      <n-layout-header bordered class="h-16 flex items-center px-6">
        <h2 class="text-base font-semibold m-0">
          {{ route.meta.title || "概述" }}
        </h2>
      </n-layout-header>
      <n-layout-content
        content-style="padding: 24px;"
        :native-scrollbar="false"
        class="h-[calc(100vh-64px)]"
      >
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>
