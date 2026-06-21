<script setup lang="ts">
import { ref, computed } from "vue";
import { useDialog, useMessage } from "naive-ui";
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

const files = ref<Record<string, FileItem[]>>({
  "": [
    { name: "releases", type: "folder", size: "-", date: "2026-06-20" },
    { name: "config", type: "folder", size: "-", date: "2026-06-15" },
    { name: "README.md", type: "file", size: "2.4 KB", date: "2026-06-01" },
    { name: "changelog.txt", type: "file", size: "8.1 KB", date: "2026-06-20" },
    { name: "launcher.exe", type: "file", size: "18.9 MB", date: "2026-04-01" },
  ],
  releases: [
    { name: "v1.2.0", type: "folder", size: "-", date: "2026-06-20" },
    { name: "v1.1.3", type: "folder", size: "-", date: "2026-06-12" },
  ],
  "releases/v1.2.0": [
    { name: "app-1.2.0.exe", type: "file", size: "24.3 MB", date: "2026-06-20" },
    { name: "patch-1.2.0.zip", type: "file", size: "3.1 MB", date: "2026-06-20" },
  ],
  "releases/v1.1.3": [
    { name: "app-1.1.3.exe", type: "file", size: "23.8 MB", date: "2026-06-12" },
  ],
  config: [
    { name: "update.json", type: "file", size: "0.5 KB", date: "2026-06-15" },
    { name: "mirror.json", type: "file", size: "0.3 KB", date: "2026-06-10" },
  ],
});

const currentKey = computed(() => currentPath.value.join("/"));
const currentFiles = computed(() => files.value[currentKey.value] || []);

const sortedFiles = computed(() => {
  return [...currentFiles.value].sort((a, b) => {
    if (a.type === "folder" && b.type !== "folder") return -1;
    if (a.type !== "folder" && b.type === "folder") return 1;
    return a.name.localeCompare(b.name);
  });
});

function navigateTo(path: string) {
  currentPath.value = path ? path.split("/") : [];
}

function openItem(item: FileItem) {
  if (item.type === "folder") {
    currentPath.value.push(item.name);
  }
}

const showNewFolderModal = ref(false);
const newFolderName = ref("");

function handleCreateFolder() {
  const name = newFolderName.value.trim();
  if (!name) {
    message.warning("请输入文件夹名称");
    return;
  }
  if (currentFiles.value.some((f) => f.name === name)) {
    message.error("同名项已存在");
    return;
  }
  const key = currentKey.value;
  if (!files.value[key]) files.value[key] = [];
  files.value[key].push({
    name,
    type: "folder",
    size: "-",
    date: new Date().toISOString().slice(0, 10),
  });
  const childKey = key ? `${key}/${name}` : name;
  if (!files.value[childKey]) files.value[childKey] = [];
  showNewFolderModal.value = false;
  newFolderName.value = "";
  message.success("文件夹已创建");
}

const showRenameModal = ref(false);
const renameTarget = ref<FileItem | null>(null);
const renameValue = ref("");

function openRename(item: FileItem) {
  renameTarget.value = item;
  renameValue.value = item.name;
  showRenameModal.value = true;
}

function handleRename() {
  const name = renameValue.value.trim();
  if (!name) {
    message.warning("请输入新名称");
    return;
  }
  if (!renameTarget.value) return;
  if (currentFiles.value.some((f) => f.name === name && f !== renameTarget.value)) {
    message.error("同名项已存在");
    return;
  }
  const oldName = renameTarget.value.name;
  const key = currentKey.value;
  const list = files.value[key];
  const idx = list.findIndex((f) => f.name === oldName);
  if (idx !== -1) {
    if (renameTarget.value.type === "folder") {
      const oldChildKey = key ? `${key}/${oldName}` : oldName;
      const newChildKey = key ? `${key}/${name}` : name;
      if (files.value[oldChildKey]) {
        files.value[newChildKey] = files.value[oldChildKey];
        delete files.value[oldChildKey];
      }
    }
    list[idx].name = name;
  }
  showRenameModal.value = false;
  renameTarget.value = null;
  renameValue.value = "";
  message.success("已重命名");
}

function handleDelete(item: FileItem) {
  dialog.warning({
    title: "确认删除",
    content: `确定要删除 ${item.name} 吗？${item.type === "folder" ? "文件夹内所有内容也将被删除。" : ""}`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: () => {
      const key = currentKey.value;
      files.value[key] = files.value[key].filter((f) => f.name !== item.name);
      if (item.type === "folder") {
        const childKey = key ? `${key}/${item.name}` : item.name;
        delete files.value[childKey];
      }
      message.success("已删除");
    },
  });
}

function handleUpload() {
  message.info("文件上传功能需对接后端 API");
}
</script>

<template>
  <div class="space-y-4">
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
        <n-button @click="handleUpload">
          <template #icon><n-icon :component="CloudUploadOutlined" /></template>
          上传
        </n-button>
        <n-button @click="showNewFolderModal = true">
          <template #icon><n-icon :component="CreateNewFolderOutlined" /></template>
          新建文件夹
        </n-button>
      </div>
    </div>

    <n-card :bordered="false" class="p-0">
      <div
        v-if="currentPath.length > 0"
        class="flex items-center gap-2 px-4 py-2 border-b border-gray-100 cursor-pointer hover:bg-gray-50 text-sm text-gray-500"
        @click="currentPath.pop()"
      >
        <n-icon :component="FolderOutlined" />
        <span>..</span>
      </div>

      <div v-if="sortedFiles.length === 0 && currentPath.length === 0" class="py-8">
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
