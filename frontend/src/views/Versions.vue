<script setup lang="ts">
import { ref, computed } from "vue";
import { useDialog, useMessage } from "naive-ui";

const dialog = useDialog();
const message = useMessage();

const versions = ref([
  {
    version: "v1.2.0",
    date: "2026-06-20",
    size: "24.3 MB",
    notes: [
      "优化启动速度，冷启动耗时降低 40%",
      "修复了在低内存设备上的崩溃问题",
      "新增增量更新机制，减少下载流量",
      "改进日志系统，支持远程上报",
    ],
  },
  {
    version: "v1.1.3",
    date: "2026-06-12",
    size: "23.8 MB",
    notes: [
      "修复网络断开时无限重连的问题",
      "优化内存占用，降低约 15%",
      "修复文件校验失败后未回滚的 bug",
    ],
  },
  {
    version: "v1.1.2",
    date: "2026-06-03",
    size: "23.5 MB",
    notes: [
      "新增多线程下载支持",
      "修复进度条偶尔显示不准确的问题",
      "兼容旧版配置文件格式",
    ],
  },
  {
    version: "v1.1.1",
    date: "2026-05-25",
    size: "22.9 MB",
    notes: [
      "修复 Windows 7 下无法启动的问题",
      "优化磁盘写入性能",
    ],
  },
  {
    version: "v1.1.0",
    date: "2026-05-18",
    size: "22.4 MB",
    notes: [
      "全新 UI 界面重构",
      "支持断点续传",
      "新增更新回滚功能",
      "支持自定义更新源",
    ],
  },
  {
    version: "v1.0.5",
    date: "2026-05-10",
    size: "21.7 MB",
    notes: [
      "修复安装路径含空格时无法启动的问题",
      "优化更新检测逻辑",
    ],
  },
  {
    version: "v1.0.4",
    date: "2026-05-02",
    size: "21.2 MB",
    notes: [
      "新增静默更新模式",
      "修复代理环境下无法连接服务器的问题",
    ],
  },
  {
    version: "v1.0.3",
    date: "2026-04-24",
    size: "20.8 MB",
    notes: [
      "修复大文件下载中断后无法恢复的问题",
      "优化磁盘空间检测",
    ],
  },
  {
    version: "v1.0.2",
    date: "2026-04-15",
    size: "20.1 MB",
    notes: [
      "修复版本比较逻辑错误导致重复更新",
      "新增更新失败自动重试",
    ],
  },
  {
    version: "v1.0.1",
    date: "2026-04-08",
    size: "19.5 MB",
    notes: [
      "修复启动时偶现白屏的问题",
      "优化首屏加载速度",
    ],
  },
  {
    version: "v1.0.0",
    date: "2026-04-01",
    size: "18.9 MB",
    notes: [
      "首个正式版本发布",
      "支持自动更新与手动更新",
      "支持多平台安装包管理",
    ],
  },
]);

const pageSize = 5;
const currentPage = ref(1);

const pagedVersions = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return versions.value.slice(start, start + pageSize);
});

const showNewVersionModal = ref(false);
const newVersion = ref({ version: "", notes: "" });
const creating = ref(false);

function handleCreate() {
  if (!newVersion.value.version.trim()) {
    message.warning("请输入版本号");
    return;
  }
  creating.value = true;
  setTimeout(() => {
    versions.value.unshift({
      version: newVersion.value.version.trim(),
      date: new Date().toISOString().slice(0, 10),
      size: "0 MB",
      notes: newVersion.value.notes
        .split("\n")
        .map((n) => n.trim())
        .filter(Boolean),
    });
    currentPage.value = 1;
    creating.value = false;
    showNewVersionModal.value = false;
    newVersion.value = { version: "", notes: "" };
    message.success("版本创建成功");
  }, 1000);
}

function handleDelete(version: string) {
  dialog.warning({
    title: "确认删除",
    content: `确定要删除 ${version} 吗？此操作不可撤销。`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: () => {
      versions.value = versions.value.filter((v) => v.version !== version);
      const totalPages = Math.ceil(versions.value.length / pageSize);
      if (currentPage.value > totalPages && totalPages > 0) {
        currentPage.value = totalPages;
      }
      message.success("已删除");
    },
  });
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex justify-end">
      <n-button type="primary" @click="showNewVersionModal = true">
        打包新版本
      </n-button>
    </div>

    <n-card v-for="item in pagedVersions" :key="item.version">
      <template #header>
        <div class="flex items-center justify-between w-full">
          <div class="flex items-center gap-3">
            <span class="font-semibold">{{ item.version }}</span>
            <n-tag v-if="item === versions[0]" size="small" type="success">
              最新
            </n-tag>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-xs text-gray-400">{{ item.date }}</span>
            <n-button
              v-if="item === versions[0]"
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

    <div class="flex justify-center">
      <n-pagination
        v-model:page="currentPage"
        :page-count="Math.ceil(versions.length / pageSize)"
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
