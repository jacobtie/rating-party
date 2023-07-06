<script setup lang="ts">
import { useSession } from '@/composables/session';
import router from '@/router';
import { createGame, getAllGames, type Game } from '@/services/game-service';
import { ref } from 'vue';

const { getUser, deleteUser } = useSession();

const user = getUser();

if (!user || !user.isAdmin || !user.jwt) {
  deleteUser();
  router.push('/');
}

const games = ref<Game[]>([]);

(async () => {
  try {
    const allGames = await getAllGames(user.jwt);
    if (allGames === false) {
      deleteUser();
      router.push('/');
      return;
    }
    games.value = allGames;
  } catch (err) {
    console.error(err);
  }
})();

const newPartyName = ref('');

const createParty = async () => {
  if (!newPartyName.value) return;
  const newGame = await createGame(user.jwt, newPartyName.value);
  router.push(`/admin/games/${newGame.gameId}`);
};

const logout = () => {
  deleteUser();
  router.push('/');
};
</script>

<template>
  <div class="full-height">
    <h1 class="main-title">Admin Page</h1>
    <div class="block">
      <v-text-field v-model="newPartyName" label="Party Name" variant="outlined" @keyup.enter="createParty"></v-text-field>
      <v-btn size="x-large" variant="tonal" @click="createParty">Create New Party</v-btn>
    </div>
    <div v-if="games && games.length > 0" class="block">
      <v-table>
        <tbody>
          <tr v-for="game in games" :key="game.gameId" class="game-row">
            <td>
              <router-link class="game-list-item" :to="`/admin/games/${game.gameId}`">
                <div>{{ game.gameName }}</div>
              </router-link>
            </td>
          </tr>
        </tbody>
      </v-table>
    </div>
    <div class="block">
      <v-btn variant="tonal" @click="logout">Logout</v-btn>
    </div>
  </div>
</template>

<style scoped>
.full-height {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.game-row {
  cursor: pointer;
}

.game-row:hover, .game-row:active, .game-row:focus {
  background-color: #d4cebc;
}

.game-list-item {
  text-decoration: none;
  color: black;
}
</style>
