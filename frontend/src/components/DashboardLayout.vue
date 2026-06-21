<script setup lang="ts">
import { ref, h, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import {
  HomeOutlined,
  SettingsOutlined,
  ExitToAppOutlined,
  AssessmentOutlined,
  SystemUpdateAltOutlined,
  FolderOutlined,
} from "@vicons/material";

const router = useRouter();
const route = useRoute();
const collapsed = ref(false);

function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) });
}

type Component = ReturnType<typeof h> extends infer R
  ? R extends (...args: unknown[]) => infer V
    ? V
    : never
  : never;

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
    label: "数据分析",
    key: "analytics",
    icon: renderIcon(AssessmentOutlined),
  },
];

const bottomMenuOptions = [
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
  if (path.startsWith("/dashboard/analytics")) return "analytics";
  if (path.startsWith("/dashboard/settings")) return "settings";
  return "dashboard";
});

function handleMenuSelect(key: string) {
  if (key === "logout") {
    router.push("/login");
    return;
  }
  const routeMap: Record<string, string> = {
    dashboard: "/dashboard",
    versions: "/dashboard/versions",
    files: "/dashboard/files",
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
      <div class="flex items-center h-16 px-4 border-b border-gray-100">
        <span
          v-if="!collapsed"
          class="text-lg font-bold tracking-wide whitespace-nowrap"
        >
          LatestPack
        </span>
        <span v-else class="text-lg font-bold mx-auto">L</span>
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

      <div class="border-t border-gray-100 py-2">
        <n-menu
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="bottomMenuOptions"
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
