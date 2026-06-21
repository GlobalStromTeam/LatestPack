<script setup lang="ts">
import { ref, onMounted } from "vue";
import { api } from "../api";
import { useMessage } from "naive-ui";

const message = useMessage();
const loading = ref(true);
const error = ref<string | null>(null);

const stats = ref<Array<{ label: string; value: string; change: string }>>([]);
const latestVersion = ref<{
  version: string;
  date: string;
  notes: string[];
}>({ version: "", date: "", notes: [] });

async function fetchData() {
  loading.value = true;
  error.value = null;
  try {
    const [statsRes, versionRes] = await Promise.all([
      api.get("/dashboard/stats"),
      api.get("/dashboard/latest-version"),
    ]);
    const d = statsRes.data;
    stats.value = [
      {
        label: "今日启动",
        value: String(d.launches.value),
        change: `+${d.launches.change}%`,
      },
      {
        label: "今日更新",
        value: String(d.updates.value),
        change: `+${d.updates.change}%`,
      },
      {
        label: `流量 (${d.traffic.unit})`,
        value: String(d.traffic.value),
        change: `+${d.traffic.change}%`,
      },
    ];
    latestVersion.value = versionRes.data;
  } catch (err: unknown) {
    if ((err as { response?: { status?: number } })?.response?.status !== 401) {
      error.value = "加载数据失败";
      message.error(error.value);
    }
  } finally {
    loading.value = false;
  }
}

onMounted(fetchData);
</script>

<template>
  <n-spin :show="loading">
    <template v-if="error">
      <n-result status="error" title="加载失败" :description="error">
        <template #footer>
          <n-button @click="fetchData">重试</n-button>
        </template>
      </n-result>
    </template>
    <template v-else>
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
  </n-spin>
</template>
