<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useDialog, useMessage } from "naive-ui";
import { api } from "../api";

interface VersionItem {
  version: string;
  date: string;
  size: string;
  notes: string[];
}

const dialog = useDialog();
const message = useMessage();

const versions = ref<VersionItem[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = 5;
const total = ref(0);

const showNewVersionModal = ref(false);
const newVersion = ref({ version: "", notes: "" });
const creating = ref(false);

function isLatest(version: string) {
  return currentPage.value === 1 && versions.value.length > 0 && versions.value[0].version === version;
}

async function fetchVersions() {
  loading.value = true;
  try {
    const res = await api.get("/versions", {
      params: { page: currentPage.value, pageSize },
    });
    versions.value = res.data.items;
    total.value = res.data.total;
  } catch (err: unknown) {
    if ((err as { response?: { status?: number } })?.response?.status !== 401) {
      message.error("加载版本列表失败");
    }
  } finally {
    loading.value = false;
  }
}

function onPageChange(page: number) {
  currentPage.value = page;
  fetchVersions();
}

async function handleCreate() {
  if (!newVersion.value.version.trim()) {
    message.warning("请输入版本号");
    return;
  }
  creating.value = true;
  try {
    const notes = newVersion.value.notes
      .split("\n")
      .map((n) => n.trim())
      .filter(Boolean);
    await api.post("/versions", {
      version: newVersion.value.version.trim(),
      notes,
    });
    showNewVersionModal.value = false;
    newVersion.value = { version: "", notes: "" };
    currentPage.value = 1;
    await fetchVersions();
    message.success("版本创建成功");
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "创建版本失败";
    message.error(msg);
  } finally {
    creating.value = false;
  }
}

function handleDelete(version: string) {
  dialog.warning({
    title: "确认删除",
    content: `确定要删除 ${version} 吗？此操作不可撤销。`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        await api.delete(`/versions/${version}`);
        await fetchVersions();
        message.success("已删除");
      } catch (err: unknown) {
        const msg =
          (err as { response?: { data?: { error?: string } } })?.response?.data
            ?.error ?? "删除失败";
        message.error(msg);
      }
    },
  });
}

onMounted(fetchVersions);
</script>

<template>
  <div class="space-y-4">
    <div class="flex justify-end">
      <n-button type="primary" @click="showNewVersionModal = true">
        打包新版本
      </n-button>
    </div>

    <n-spin :show="loading">
      <template v-if="!loading && versions.length === 0">
        <n-empty description="暂无版本" />
      </template>
      <n-card v-for="item in versions" :key="item.version">
        <template #header>
          <div class="flex items-center justify-between w-full">
            <div class="flex items-center gap-3">
              <span class="font-semibold">{{ item.version }}</span>
              <n-tag v-if="isLatest(item.version)" size="small" type="success">
                最新
              </n-tag>
            </div>
            <div class="flex items-center gap-3">
              <span class="text-xs text-gray-400">{{ item.date }}</span>
              <n-button
                v-if="isLatest(item.version)"
                size="tiny"
                type="error"
                quaternary
                @click="handleDelete(item.version)"
              >
                删除
              </n-button>
            </div>
          </div>
        </template>

        <div class="text-sm text-gray-500 mb-3">
          大小: {{ item.size }}
        </div>

        <ul class="list-disc pl-5 space-y-1.5 text-sm text-gray-600 m-0">
          <li v-for="note in item.notes" :key="note">{{ note }}</li>
        </ul>
      </n-card>
    </n-spin>

    <div v-if="total > 0" class="flex justify-center">
      <n-pagination
        v-model:page="currentPage"
        :page-count="Math.ceil(total / pageSize)"
        @update:page="onPageChange"
      />
    </div>

    <n-modal
      v-model:show="showNewVersionModal"
      preset="card"
      title="打包新版本"
      style="width: 480px"
      :mask-closable="false"
    >
      <n-form label-placement="top">
        <n-form-item label="版本号">
          <n-input v-model:value="newVersion.version" placeholder="例如 v1.3.0" />
        </n-form-item>
        <n-form-item label="更新日志">
          <n-input
            v-model:value="newVersion.notes"
            type="textarea"
            placeholder="每行一条更新内容"
            :rows="5"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <n-button @click="showNewVersionModal = false">取消</n-button>
          <n-button type="primary" :loading="creating" @click="handleCreate">
            创建
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
