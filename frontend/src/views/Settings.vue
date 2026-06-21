<script setup lang="ts">
import { ref } from "vue";
import { useMessage, useDialog } from "naive-ui";
import { api } from "../api";
import { useAuthStore } from "../stores/auth";
import {
  PersonOutlined,
  LockOutlined,
  SaveOutlined,
} from "@vicons/material";

const message = useMessage();
const dialog = useDialog();
const authStore = useAuthStore();

const currentUsername = ref(authStore.username ?? "");
const newUsername = ref("");
const savingUsername = ref(false);

const currentPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const savingPassword = ref(false);

async function handleUpdateUsername() {
  const name = newUsername.value.trim();
  if (!name) {
    message.warning("请输入新账户名");
    return;
  }
  if (name === currentUsername.value) {
    message.warning("新账户名与当前相同");
    return;
  }
  savingUsername.value = true;
  try {
    await api.put("/auth/username", { username: name });
    authStore.username = name;
    currentUsername.value = name;
    newUsername.value = "";
    message.success("账户名已更新");
  } catch (err: unknown) {
    const msg =
      (err as { response?: { data?: { error?: string } } })?.response?.data
        ?.error ?? "更新账户名失败";
    message.error(msg);
  } finally {
    savingUsername.value = false;
  }
}

async function handleUpdatePassword() {
  if (!currentPassword.value) {
    message.warning("请输入当前密码");
    return;
  }
  if (!newPassword.value) {
    message.warning("请输入新密码");
    return;
  }
  if (newPassword.value.length < 6) {
    message.warning("新密码至少 6 位");
    return;
  }
  if (newPassword.value !== confirmPassword.value) {
    message.warning("两次输入的新密码不一致");
    return;
  }

  dialog.warning({
    title: "确认修改密码",
    content: "修改密码后需要重新登录，确定继续？",
    positiveText: "确认修改",
    negativeText: "取消",
    onPositiveClick: async () => {
      savingPassword.value = true;
      try {
        await api.put("/auth/password", {
          currentPassword: currentPassword.value,
          newPassword: newPassword.value,
        });
        message.success("密码已修改，请重新登录");
        authStore.logout();
        window.location.href = "/login";
      } catch (err: unknown) {
        const msg =
          (err as { response?: { data?: { error?: string } } })?.response?.data
            ?.error ?? "修改密码失败";
        message.error(msg);
      } finally {
        savingPassword.value = false;
      }
    },
  });
}
</script>

<template>
  <div class="space-y-6 max-w-2xl">
    <n-card>
      <template #header>
        <div class="flex items-center gap-2">
          <n-icon :component="PersonOutlined" />
          <span>修改账户名</span>
        </div>
      </template>
      <div class="space-y-4">
        <div>
          <div class="text-sm text-gray-500 mb-1">当前账户名</div>
          <div class="text-base font-medium">{{ currentUsername }}</div>
        </div>
        <n-input
          v-model:value="newUsername"
          placeholder="输入新账户名"
          @keydown.enter="handleUpdateUsername"
        />
        <div class="flex justify-end">
          <n-button
            type="primary"
            :loading="savingUsername"
            @click="handleUpdateUsername"
          >
            <template #icon><n-icon :component="SaveOutlined" /></template>
            保存
          </n-button>
        </div>
      </div>
    </n-card>

    <n-card>
      <template #header>
        <div class="flex items-center gap-2">
          <n-icon :component="LockOutlined" />
          <span>修改密码</span>
        </div>
      </template>
      <div class="space-y-4">
        <n-input
          v-model:value="currentPassword"
          type="password"
          show-password-on="click"
          placeholder="当前密码"
        />
        <n-input
          v-model:value="newPassword"
          type="password"
          show-password-on="click"
          placeholder="新密码（至少 6 位）"
        />
        <n-input
          v-model:value="confirmPassword"
          type="password"
          show-password-on="click"
          placeholder="确认新密码"
          @keydown.enter="handleUpdatePassword"
        />
        <div class="flex justify-end">
          <n-button
            type="primary"
            :loading="savingPassword"
            @click="handleUpdatePassword"
          >
            <template #icon><n-icon :component="SaveOutlined" /></template>
            修改密码
          </n-button>
        </div>
      </div>
    </n-card>
  </div>
</template>
