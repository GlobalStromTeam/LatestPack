import { ref, computed } from "vue";
import { defineStore } from "pinia";
import { api } from "../api";

export interface ChannelConfig {
  endpoint?: string;
  bucket?: string;
  region?: string;
  accessKey?: string;
  secretKey?: string;
  path?: string;
}

export interface Channel {
  id: string;
  name: string;
  type: "local" | "webdav" | "s3";
  enabled: boolean;
  weight: number;
  config: ChannelConfig;
}

export const useChannelStore = defineStore("channels", () => {
  const channels = ref<Channel[]>([]);
  const loading = ref(false);

  const enabledChannels = computed(() => channels.value.filter((c) => c.enabled));
  const localChannel = computed(() => channels.value.find((c) => c.type === "local"));

  async function fetchChannels() {
    loading.value = true;
    try {
      const res = await api.get("/channels");
      channels.value = res.data;
    } catch {
      channels.value = [];
    } finally {
      loading.value = false;
    }
  }

  async function createChannel(data: Omit<Channel, "id">) {
    const res = await api.post("/channels", data);
    channels.value.push(res.data);
    return res.data as Channel;
  }

  async function updateChannel(id: string, data: Partial<Channel>) {
    const res = await api.put(`/channels/${id}`, data);
    const idx = channels.value.findIndex((c) => c.id === id);
    if (idx !== -1) channels.value[idx] = res.data;
    return res.data as Channel;
  }

  async function deleteChannel(id: string) {
    await api.delete(`/channels/${id}`);
    channels.value = channels.value.filter((c) => c.id !== id);
  }

  async function toggleChannel(id: string, enabled: boolean) {
    return updateChannel(id, { enabled });
  }

  return {
    channels,
    loading,
    enabledChannels,
    localChannel,
    fetchChannels,
    createChannel,
    updateChannel,
    deleteChannel,
    toggleChannel,
  };
});
