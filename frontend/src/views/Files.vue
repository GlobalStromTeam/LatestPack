<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useDialog, useMessage } from "naive-ui";
import { api } from "../api";
import {
  FolderOutlined,
  InsertDriveFileOutlined,
  CloudUploadOutlined,
  CreateNewFolderOutlined,
  EditOutlined,
  DeleteOutlined,
} from "@vicons/material";

const dialog = useDialog();
const message = useMessage();

interface FileItem {
  name: string;
  type: "folder" | "file";
  size: string;
  date: string;
}

const currentPath = ref<string[]>([]);
const pathBreadcrumbs = computed(() => {
  const crumbs = [{ label: "根目录", path: "" }];
  currentPath.value.forEach((segment, i) => {
    crumbs.push({
      label: segment,
      path: currentPath.value.slice(0, i + 1).join("/"),
    });
  });
  return crumbs;
});

const files = ref<FileItem[]>([]);
const loading = ref(false);

const sortedFiles = computed(() => {
  return [...files.value].sort((a, b) => {
    if (a.type === "folder" && b.type !== "folder") return -1;
    if (a.type !== "folder" && b.type === "folder") return 1;
    return a.name.localeCompare(b.name);
  });
});

async function fetchFiles() {
  loading.value = true;
  try {
    const res = await api.get("/files", {
      params: { path: currentPath.value.join("/") },
    });
    files.value = res.data.items;
  } catch (err: unknown) {
    if ((err as { response?: { status?: number } })?.response?.status !== 401) {
      message.error("加载文件列表失败");
    }
  } finally {
    loading.value = false;
  }
}

function navigateTo(path: string) {
  currentPath.value = path ? path.split("/") : [];
  fetchFiles();
}

function openItem(item: FileItem) {
  if (item.type === "folder") {
    currentPath.value.push(item.name);
    fetchFiles();
  }
}

function goUp() {
  currentPath.value.pop();
  fetchFiles();
}

const showNewFolderModal = ref(false);
const newFolderName = ref("");

async function handleCreateFolder() {
  const name = newFolderName.value.trim();
  if (!name) {
    message.warning("请输入文件夹名称");
    return;
  }
  try {
    await api.post("/files/folder", {
      path: currentPath.value.join("/"),
      name,
    });
    showNewFolderModal.value = false;
    newFolderName.value = "";
    await fetchFiles();
    message.success("文件夹已创建");
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "创建文件夹失败";
    message.error(msg);
  }
}

const showRenameModal = ref(false);
const renameTarget = ref<FileItem | null>(null);
const renameValue = ref("");

function openRename(item: FileItem) {
  renameTarget.value = item;
  renameValue.value = item.name;
  showRenameModal.value = true;
}

async function handleRename() {
  const name = renameValue.value.trim();
  if (!name) {
    message.warning("请输入新名称");
    return;
  }
  if (!renameTarget.value) return;
  try {
    await api.put("/files/rename", {
      path: currentPath.value.join("/"),
      oldName: renameTarget.value.name,
      newName: name,
    });
    showRenameModal.value = false;
    renameTarget.value = null;
    renameValue.value = "";
    await fetchFiles();
    message.success("已重命名");
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "重命名失败";
    message.error(msg);
  }
}

function handleDelete(item: FileItem) {
  dialog.warning({
    title: "确认删除",
    content: `确定要删除 ${item.name} 吗？${item.type === "folder" ? "文件夹内所有内容也将被删除。" : ""}`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        await api.delete("/files", {
          data: { path: currentPath.value.join("/"), name: item.name },
        });
        await fetchFiles();
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

const uploading = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

function triggerUpload() {
  fileInput.value?.click();
}

async function handleFileSelected(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const formData = new FormData();
    formData.append("path", currentPath.value.join("/"));
    formData.append("file", file);
    await api.post("/files/upload", formData);
    await fetchFiles();
    message.success("上传成功");
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "上传失败";
    message.error(msg);
  } finally {
    uploading.value = false;
    input.value = "";
  }
}

onMounted(fetchFiles);
</script>

<template>
  <div class="space-y-4">
    <input ref="fileInput" type="file" style="display: none" @change="handleFileSelected" />

    <div class="flex items-center justify-between">
      <n-breadcrumb>
        <n-breadcrumb-item
          v-for="crumb in pathBreadcrumbs"
          :key="crumb.path"
          @click="navigateTo(crumb.path)"
        >
          {{ crumb.label }}
        </n-breadcrumb-item>
      </n-breadcrumb>
      <div class="flex gap-2">
        <n-button :loading="uploading" @click="triggerUpload">
          <template #icon><n-icon :component="CloudUploadOutlined" /></template>
          上传
        </n-button>
        <n-button @click="showNewFolderModal = true">
          <template #icon><n-icon :component="CreateNewFolderOutlined" /></template>
          新建文件夹
        </n-button>
      </div>
    </div>

    <n-spin :show="loading">
      <n-card :bordered="false" class="p-0">
        <div
          v-if="currentPath.length > 0"
          class="flex items-center gap-2 px-4 py-2 border-b border-gray-100 cursor-pointer hover:bg-gray-50 text-sm text-gray-500"
          @click="goUp"
        >
          <n-icon :component="FolderOutlined" />
          <span>..</span>
        </div>

        <div v-if="sortedFiles.length === 0 && !loading" class="py-8">
          <n-empty description="此文件夹为空" />
        </div>

        <div
          v-for="item in sortedFiles"
          :key="item.name"
          class="flex items-center justify-between px-4 py-2 border-b border-gray-50 hover:bg-gray-50"
        >
          <div
            class="flex items-center gap-2 flex-1 min-w-0 cursor-pointer"
            @click="openItem(item)"
          >
            <n-icon
              :component="item.type === 'folder' ? FolderOutlined : InsertDriveFileOutlined"
              :color="item.type === 'folder' ? '#0f7b6c' : '#9ca3af'"
              size="20"
            />
            <span class="truncate">{{ item.name }}</span>
          </div>
          <div class="flex items-center gap-4 text-sm text-gray-400 shrink-0">
            <span class="w-20 text-right">{{ item.size }}</span>
            <span class="w-24">{{ item.date }}</span>
            <div class="flex gap-1">
              <n-button size="tiny" quaternary @click="openRename(item)">
                <template #icon><n-icon :component="EditOutlined" /></template>
              </n-button>
              <n-button size="tiny" quaternary type="error" @click="handleDelete(item)">
                <template #icon><n-icon :component="DeleteOutlined" /></template>
              </n-button>
            </div>
          </div>
        </div>
      </n-card>
    </n-spin>

    <n-modal
      v-model:show="showNewFolderModal"
      preset="card"
      title="新建文件夹"
      style="width: 400px"
    >
      <n-input v-model:value="newFolderName" placeholder="文件夹名称" @keydown.enter="handleCreateFolder" />
      <template #footer>
        <div class="flex justify-end gap-3">
          <n-button @click="showNewFolderModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateFolder">创建</n-button>
        </div>
      </template>
    </n-modal>

    <n-modal
      v-model:show="showRenameModal"
      preset="card"
      title="重命名"
      style="width: 400px"
    >
      <n-input v-model:value="renameValue" placeholder="新名称" @keydown.enter="handleRename" />
      <template #footer>
        <div class="flex justify-end gap-3">
          <n-button @click="showRenameModal = false">取消</n-button>
          <n-button type="primary" @click="handleRename">确认</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
