import { ref, computed } from "vue";
import { defineStore } from "pinia";
import axios from "axios";

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(null);
  const username = ref<string | null>(null);

  const isAuthenticated = computed(() => token.value !== null);

  function initialize() {
    try {
      const raw = localStorage.getItem("auth");
      if (raw) {
        const data = JSON.parse(raw);
        token.value = data.token;
        username.value = data.username;
      }
    } catch {
      localStorage.removeItem("auth");
    }
  }

  async function login(user: string, password: string) {
    const res = await axios.post("/api/auth/login", {
      username: user,
      password,
    });
    token.value = res.data.token;
    username.value = res.data.username;
    localStorage.setItem(
      "auth",
      JSON.stringify({ token: token.value, username: username.value }),
    );
  }

  function logout() {
    token.value = null;
    username.value = null;
    localStorage.removeItem("auth");
  }

  initialize();

  return { token, username, isAuthenticated, login, logout };
});
