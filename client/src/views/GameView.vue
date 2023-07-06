<script setup lang="ts">
import { useSession } from '@/composables/session';
import router from '@/router';

const { getUser } = useSession();

if (!getUser().jwt) {
  router.push('/');
}

</script>

<template>
  <div>
    <h1>Game</h1>
  </div>
</template>
