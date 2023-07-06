<script setup lang="ts">
import { useSession } from '@/composables/session';
import router from '@/router';
import { signin } from '@/services/session-service';
import { ref } from 'vue';

const { getUser, setUser } = useSession();

if (getUser().jwt && getUser().isAdmin) {
  router.push('/admin');
} else if (getUser().jwt) {
  router.push('/game');
}

const username = ref('');
const gameCode = ref('');
const errMsg = ref('');

const joinGame = async () => {
  if (!username.value || !gameCode.value) return;
  try {
    const res = await signin(username.value, gameCode.value);
    if (!res.jwt) {
      throw Error('No JWT');
    }
    setUser(res);
    if (res.isAdmin) {
      router.push('/admin');
    } else {
      router.push('/game');
    }
  } catch (err) {
    errMsg.value = 'Invalid party code';
  }
};
</script>

<template>
  <div class="full-height">
    <div class="login-box">
      <h2 class="login-header">Login</h2>
      <v-text-field v-model="username" class="login-field" label="Username" variant="outlined" @keyup.enter="joinGame"></v-text-field>
      <v-text-field v-model="gameCode" class="login-field" label="Party Code" variant="outlined" @keyup.enter="joinGame"></v-text-field>
      <p>{{ errMsg }}</p>
      <v-btn size="x-large" variant="tonal" @click="joinGame">Join</v-btn>
    </div>
  </div>
</template>

<style scoped>
.full-height {
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-box {
  background-color: #F0EAD6;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  width: 95vw;
  display: flex;
  flex-direction: column;
}

.login-header {
  margin-bottom: 16px;
}
</style>
