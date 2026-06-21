<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useMessage, useDialog } from "naive-ui";
import { useChannelStore } from "../stores/channels";
import type { Channel, ChannelConfig } from "../stores/channels";
import {
  AddOutlined,
  DeleteOutlined,
  EditOutlined,
  StorageOutlined,
  CloudSyncOutlined,
  FolderOutlined,
} from "@vicons/material";

const message = useMessage();
const dialog = useDialog();
const store = useChannelStore();

const showModal = ref(false);
const editingChannel = ref<Channel | null>(null);
const saving = ref(false);

const formType = ref<"local" | "webdav" | "s3">("webdav");
const formName = ref("");
const formEnabled = ref(true);
const formWeight = ref(50);

const webdavConfig = ref<Pick<ChannelConfig, "endpoint" | "path" | "accessKey" | "secretKey">>({
  endpoint: "",
  path: "",
  accessKey: "",
  secretKey: "",
});

const s3Config = ref<Pick<ChannelConfig, "endpoint" | "bucket" | "region" | "accessKey" | "secretKey" | "path">>({
  endpoint: "",
  bucket: "",
  region: "",
  accessKey: "",
  secretKey: "",
  path: "",
});

const typeOptions = [
  { label: "WebDAV", value: "webdav" as const, icon: CloudSyncOutlined },
  { label: "S3", value: "s3" as const, icon: StorageOutlined },
];

function getTypeIcon(type: string) {
  if (type === "local") return FolderOutlined;
  if (type === "webdav") return CloudSyncOutlined;
  return StorageOutlined;
}

function getTypeLabel(type: string) {
  if (type === "local") return "本地";
  if (type === "webdav") return "WebDAV";
  return "S3";
}

function openCreate() {
  editingChannel.value = null;
  formType.value = "webdav";
  formName.value = "";
  formEnabled.value = true;
  formWeight.value = 50;
  webdavConfig.value = { endpoint: "", path: "", accessKey: "", secretKey: "" };
  s3Config.value = { endpoint: "", bucket: "", region: "", accessKey: "", secretKey: "", path: "" };
  showModal.value = true;
}

function openEdit(channel: Channel) {
  editingChannel.value = channel;
  formType.value = channel.type;
  formName.value = channel.name;
  formEnabled.value = channel.enabled;
  formWeight.value = channel.weight;
  if (channel.type === "webdav") {
    webdavConfig.value = {
      endpoint: channel.config.endpoint ?? "",
      path: channel.config.path ?? "",
      accessKey: channel.config.accessKey ?? "",
      secretKey: channel.config.secretKey ?? "",
    };
  } else if (channel.type === "s3") {
    s3Config.value = {
      endpoint: channel.config.endpoint ?? "",
      bucket: channel.config.bucket ?? "",
      region: channel.config.region ?? "",
      accessKey: channel.config.accessKey ?? "",
      secretKey: channel.config.secretKey ?? "",
      path: channel.config.path ?? "",
    };
  }
  showModal.value = true;
}

function buildConfig(): ChannelConfig {
  if (formType.value === "webdav") {
    return {
      endpoint: webdavConfig.value.endpoint,
      path: webdavConfig.value.path,
      accessKey: webdavConfig.value.accessKey,
      secretKey: webdavConfig.value.secretKey,
    };
  }
  if (formType.value === "s3") {
    return {
      endpoint: s3Config.value.endpoint,
      bucket: s3Config.value.bucket,
      region: s3Config.value.region,
      accessKey: s3Config.value.accessKey,
      secretKey: s3Config.value.secretKey,
      path: s3Config.value.path,
    };
  }
  return {};
}

async function handleSave() {
  const name = formName.value.trim();
  if (!name) {
    message.warning("请输入渠道名称");
    return;
  }
  if (formType.value !== "local" && formWeight.value < 1) {
    message.warning("权重不能小于 1");
    return;
  }

  saving.value = true;
  try {
    const config = buildConfig();
    if (editingChannel.value) {
      await store.updateChannel(editingChannel.value.id, {
        name,
        type: formType.value,
        enabled: formEnabled.value,
        weight: formType.value === "local" ? 0 : formWeight.value,
        config,
      });
      message.success("渠道已更新");
    } else {
      await store.createChannel({
        name,
        type: formType.value,
        enabled: formEnabled.value,
        weight: formType.value === "local" ? 0 : formWeight.value,
        config,
      });
      message.success("渠道已创建");
    }
    showModal.value = false;
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "操作失败";
    message.error(msg);
  } finally {
    saving.value = false;
  }
}

function handleDelete(channel: Channel) {
  if (channel.type === "local") {
    message.warning("本地渠道不可删除");
    return;
  }
  dialog.warning({
    title: "确认删除",
    content: `确定要删除渠道「${channel.name}」吗？`,
    positiveText: "删除",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        await store.deleteChannel(channel.id);
        message.success("渠道已删除");
      } catch (err: unknown) {
        const msg =
          (err as { response?: { data?: { error?: string } } })?.response?.data
            ?.error ?? "删除失败";
        message.error(msg);
      }
    },
  });
}

async function handleToggle(channel: Channel, enabled: boolean) {
  try {
    await store.toggleChannel(channel.id, enabled);
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "操作失败";
    message.error(msg);
  }
}

onMounted(() => {
  store.fetchChannels();
});
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div class="text-sm text-gray-500">
        共 {{ store.channels.length }} 个渠道，{{ store.enabledChannels.length }} 个已启用
      </div>
      <n-button type="primary" @click="openCreate">
        <template #icon><n-icon :component="AddOutlined" /></template>
        添加渠道
      </n-button>
    </div>

    <n-spin :show="store.loading">
      <div v-if="store.channels.length === 0 && !store.loading" class="py-8">
        <n-empty description="暂无渠道" />
      </div>

      <div class="space-y-3">
        <n-card
          v-for="channel in store.channels"
          :key="channel.id"
          size="small"
          :class="{ 'opacity-50': !channel.enabled }"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3 min-w-0">
              <n-icon
                :component="getTypeIcon(channel.type)"
                :color="channel.enabled ? '#0f7b6c' : '#9ca3af'"
                size="22"
              />
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-medium truncate">{{ channel.name }}</span>
                  <n-tag :bordered="false" size="small" :type="channel.type === 'local' ? 'default' : 'info'">
                    {{ getTypeLabel(channel.type) }}
                  </n-tag>
                  <n-tag v-if="channel.type === 'local'" :bordered="false" size="small" type="warning">
                    内置
                  </n-tag>
                </div>
                <div class="text-xs text-gray-400 mt-0.5">
                  <template v-if="channel.type === 'webdav'">
                    {{ channel.config.endpoint ?? '' }}
                  </template>
                  <template v-else-if="channel.type === 's3'">
                    {{ channel.config.bucket ?? '' }}{{ channel.config.endpoint ? ` · ${channel.config.endpoint}` : '' }}
                  </template>
                  <template v-else>
                    本地文件系统
                  </template>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3 shrink-0">
              <template v-if="channel.type !== 'local'">
                <div class="text-xs text-gray-400 w-14 text-right">
                  权重 {{ channel.weight }}
                </div>
              </template>
              <n-switch
                :value="channel.enabled"
                @update:value="(v: boolean) => handleToggle(channel, v)"
              />
              <n-button
                v-if="channel.type !== 'local'"
                size="tiny"
                quaternary
                @click="openEdit(channel)"
              >
                <template #icon><n-icon :component="EditOutlined" /></template>
              </n-button>
              <n-button
                v-if="channel.type !== 'local'"
                size="tiny"
                quaternary
                type="error"
                @click="handleDelete(channel)"
              >
                <template #icon><n-icon :component="DeleteOutlined" /></template>
              </n-button>
            </div>
          </div>
        </n-card>
      </div>
    </n-spin>

    <n-modal
      v-model:show="showModal"
      preset="card"
      :title="editingChannel ? '编辑渠道' : '添加渠道'"
      style="width: 520px"
    >
      <div class="space-y-4">
        <div v-if="!editingChannel">
          <div class="text-sm text-gray-500 mb-2">渠道类型</div>
          <div class="flex gap-3">
            <div
              v-for="opt in typeOptions"
              :key="opt.value"
              class="flex-1 border rounded-md p-3 cursor-pointer text-center transition-colors"
              :class="formType === opt.value
                ? 'border-[#0f7b6c] bg-[#0f7b6c]/5'
                : 'border-gray-200 hover:border-gray-300'"
              @click="formType = opt.value"
            >
              <n-icon :component="opt.icon" size="24" :color="formType === opt.value ? '#0f7b6c' : '#9ca3af'" />
              <div class="text-sm mt-1" :class="formType === opt.value ? 'text-[#0f7b6c] font-medium' : 'text-gray-500'">
                {{ opt.label }}
              </div>
            </div>
          </div>
        </div>

        <n-input v-model:value="formName" placeholder="渠道名称" />

        <template v-if="formType === 'webdav'">
          <n-input v-model:value="webdavConfig.endpoint" placeholder="WebDAV 地址（如 https://dav.example.com/path）" />
          <n-input v-model:value="webdavConfig.path" placeholder="远程路径（可选）" />
          <n-input v-model:value="webdavConfig.accessKey" placeholder="用户名" />
          <n-input v-model:value="webdavConfig.secretKey" type="password" show-password-on="click" placeholder="密码" />
        </template>

        <template v-if="formType === 's3'">
          <n-input v-model:value="s3Config.endpoint" placeholder="Endpoint（如 https://s3.amazonaws.com）" />
          <n-input v-model:value="s3Config.bucket" placeholder="Bucket 名称" />
          <n-input v-model:value="s3Config.region" placeholder="Region（如 us-east-1）" />
          <n-input v-model:value="s3Config.path" placeholder="路径前缀（可选）" />
          <n-input v-model:value="s3Config.accessKey" placeholder="Access Key" />
          <n-input v-model:value="s3Config.secretKey" type="password" show-password-on="click" placeholder="Secret Key" />
        </template>

        <template v-if="formType !== 'local'">
          <div>
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm text-gray-500">权重</span>
              <span class="text-sm font-medium">{{ formWeight }}</span>
            </div>
            <n-slider v-model:value="formWeight" :min="1" :max="100" :step="1" />
            <div class="flex justify-between text-xs text-gray-400 mt-1">
              <span>低优先级</span>
              <span>高优先级</span>
            </div>
          </div>
        </template>

        <div class="flex items-center gap-2">
          <n-switch v-model:value="formEnabled" />
          <span class="text-sm">{{ formEnabled ? '已启用' : '已禁用' }}</span>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="saving" @click="handleSave">
            {{ editingChannel ? '保存' : '创建' }}
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>
