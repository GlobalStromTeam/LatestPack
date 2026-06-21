<script setup lang="ts">
import { ref, reactive } from "vue";
import type { FormInst, FormRules } from "naive-ui";

const formRef = ref<FormInst | null>(null);
const loading = ref(false);

const model = reactive({
  username: "",
  password: "",
});

const rules: FormRules = {
  username: {
    required: true,
    message: "请输入用户名",
    trigger: "blur",
  },
  password: {
    required: true,
    message: "请输入密码",
    trigger: "blur",
  },
};

const handleLogin = () => {
  formRef.value?.validate((errors) => {
    if (!errors) {
      loading.value = true;
      setTimeout(() => {
        loading.value = false;
      }, 1500);
    }
  });
};
</script>

<template>
  <!-- 背景动画 -->
  <div class="flex justify-center items-center min-h-screen relative overflow-hidden">
    <div class="ocean-bg">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
      <div class="blob blob-3"></div>
      <div class="blob blob-4"></div>
      <div class="blob blob-5"></div>
    </div>
    
    <n-card :bordered="false" style="width: 400px" class="relative z-10">
      <div class="text-center mb-6">
        <h1 class="text-[28px] font-bold m-0 tracking-wide">LatestPack</h1>
      </div>
      <n-form
        ref="formRef"
        :model="model"
        :rules="rules"
        label-placement="top"
      >
        <n-form-item path="username" label="用户名">
          <n-input
            v-model:value="model.username"
            placeholder="请输入用户名"
          />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="model.password"
            type="password"
            show-password-on="mousedown"
            placeholder="请输入密码"
            @keydown.enter="handleLogin"
          />
        </n-form-item>
        <n-button
          type="primary"
          block
          :loading="loading"
          @click="handleLogin"
          round
        >
          登录
        </n-button>
      </n-form>
    </n-card>
  </div>
</template>

<style scoped>
.ocean-bg {
  position: absolute;
  inset: 0;
  filter: blur(40px);
  overflow: hidden;
  z-index: 0;
}

.blob {
  position: absolute;
  border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
}

.blob-1 {
  width: 50vmax;
  height: 50vmax;
  top: -20%;
  left: -10%;
  background: radial-gradient(circle, #0a4c6e, #062c43);
  animation: drift1 18s ease-in-out infinite alternate;
}

.blob-2 {
  width: 45vmax;
  height: 45vmax;
  top: 30%;
  right: -15%;
  background: radial-gradient(circle, #0e7490, #155e75);
  animation: drift2 22s ease-in-out infinite alternate;
}

.blob-3 {
  width: 40vmax;
  height: 40vmax;
  bottom: -10%;
  left: 20%;
  background: radial-gradient(circle, #1e3a5f, #0c2d48);
  animation: drift3 20s ease-in-out infinite alternate;
}

.blob-4 {
  width: 35vmax;
  height: 35vmax;
  top: 10%;
  left: 40%;
  background: radial-gradient(circle, #065f46, #047857);
  animation: drift4 25s ease-in-out infinite alternate;
}

.blob-5 {
  width: 30vmax;
  height: 30vmax;
  bottom: 20%;
  right: 10%;
  background: radial-gradient(circle, #312e81, #1e1b4b);
  animation: drift5 19s ease-in-out infinite alternate;
}

@keyframes drift1 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%; }
  50% { border-radius: 60% 40% 30% 70% / 50% 60% 40% 50%; }
  100% { transform: translate(8vw, 6vh) rotate(45deg) scale(1.1); border-radius: 40% 60% 60% 40% / 70% 30% 70% 30%; }
}

@keyframes drift2 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); border-radius: 60% 40% 30% 70% / 50% 60% 40% 50%; }
  50% { border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%; }
  100% { transform: translate(-6vw, 8vh) rotate(-30deg) scale(1.15); border-radius: 50% 50% 30% 70% / 60% 40% 60% 40%; }
}

@keyframes drift3 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); border-radius: 50% 50% 30% 70% / 60% 40% 60% 40%; }
  50% { border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%; }
  100% { transform: translate(5vw, -5vh) rotate(60deg) scale(0.95); border-radius: 40% 60% 60% 40% / 70% 30% 70% 30%; }
}

@keyframes drift4 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); border-radius: 40% 60% 60% 40% / 70% 30% 70% 30%; }
  50% { border-radius: 60% 40% 30% 70% / 50% 60% 40% 50%; }
  100% { transform: translate(-4vw, 7vh) rotate(-45deg) scale(1.08); border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%; }
}

@keyframes drift5 {
  0% { transform: translate(0, 0) rotate(0deg) scale(1); border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%; }
  50% { border-radius: 50% 50% 30% 70% / 60% 40% 60% 40%; }
  100% { transform: translate(7vw, -3vh) rotate(35deg) scale(1.12); border-radius: 60% 40% 30% 70% / 50% 60% 40% 50%; }
}
</style>
