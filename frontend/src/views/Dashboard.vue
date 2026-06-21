<script setup lang="ts">
import { ref } from "vue";

const stats = ref([
  { label: "今日启动", value: "372", change: "+8.1%" },
  { label: "今日更新", value: "56", change: "+3.4%" },
  { label: "流量 (MB)", value: "2,847", change: "+15.7%" },
]);

const latestVersion = ref({
  version: "v1.2.0",
  date: "2026-06-20",
  notes: [
    "优化启动速度，冷启动耗时降低 40%",
    "修复了在低内存设备上的崩溃问题",
    "新增增量更新机制，减少下载流量",
    "改进日志系统，支持远程上报",
  ],
});
</script>

<template>
  <div class="space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <n-card v-for="stat in stats" :key="stat.label" size="small">
        <div class="text-sm text-gray-500 mb-1">{{ stat.label }}</div>
        <div class="text-2xl font-bold">{{ stat.value }}</div>
        <div class="text-sm text-emerald-600 mt-1">{{ stat.change }}</div>
      </n-card>
    </div>

    <n-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span>更新公告</span>
          <n-tag size="small" type="info">{{ latestVersion.version }}</n-tag>
        </div>
      </template>
      <div class="text-xs text-gray-400 mb-3">{{ latestVersion.date }}</div>
      <ul class="list-disc pl-5 space-y-2 text-sm text-gray-600 m-0">
        <li v-for="note in latestVersion.notes" :key="note">{{ note }}</li>
      </ul>
    </n-card>
  </div>
</template>
