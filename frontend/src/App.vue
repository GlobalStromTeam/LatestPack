<script setup lang="ts">
import themeOverrides from "./theme";
</script>

<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <router-view />
  </n-config-provider>
</template>
