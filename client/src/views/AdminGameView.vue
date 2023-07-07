<script setup lang="ts">
import { useSession } from '@/composables/session';
import router from '@/router';
import { deleteGame, getGame, updateGame, type Game } from '@/services/game-service';
import { createWine, deleteWine, getAllWines, type Wine } from '@/services/wine-service';
import { ref } from 'vue';

const { getUser, deleteUser } = useSession();

const user = getUser();

if (!user || !user.isAdmin || !user.jwt) {
  deleteUser();
  router.push('/');
}

const game = ref<Game | null>(null);

const gameId = `${router.currentRoute.value.params.gameId}`;

(async () => {
  try {
    const gameFromServer = await getGame(user.jwt, gameId);
    if (gameFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    game.value = gameFromServer;
  } catch (err) {
    console.error(err);
  }
})();

const switchGameStatus = async () => {
  try {
    await updateGame(user.jwt, gameId, game.value!.gameName, !game.value!.isRunning);
    game.value!.isRunning = !game.value!.isRunning;
  } catch (err) {
    console.error(err);
  }
};

const removeGame = async () => {
  if (!window.confirm('Are you sure you want to delete this party?')) return;
  try {
    await deleteGame(user.jwt, gameId);
    router.push('/admin');
  } catch (err) {
    console.error(err);
  }
};

const wines = ref<Wine[]>([]);

(async () => {
  try {
    const winesFromServer = await getAllWines(user.jwt, gameId);
    if (winesFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    wines.value = winesFromServer;
  } catch (err) {
    console.error(err);
  }
})();

const newWineName = ref('');
const newWineCode = ref('');
const newWineYear = ref('2023');

const addWine = async () => {
  if (!newWineName.value || !newWineCode.value || !newWineYear.value || !Number.isInteger(newWineYear.value) || +newWineYear.value <= 0) return;
  try {
    const wine = await createWine(user.jwt, gameId, newWineName.value, newWineCode.value, +newWineYear.value);
    wines.value.push(wine);
    newWineName.value = '';
    newWineCode.value = '';
    newWineYear.value = '2023';
  } catch (err) {
    console.error(err);
  }
};

const removeWine = async (wineId: string) => {
  if (!window.confirm('Are you sure you want to delete this wine?')) return;
  try {
    await deleteWine(user.jwt, gameId, wineId);
    wines.value = wines.value.filter((wine) => wine.wineId !== wineId);
  } catch (err) {
    console.error(err);
  }
};

const goBack = () => {
  router.push('/admin');
};

const logout = () => {
  if (!window.confirm('Are you sure you want to log out?')) return;
  deleteUser();
  router.push('/');
};
</script>

<template>
  <div v-if="game" class="full-height">
    <h1 class="main-title">{{ game.gameName }}</h1>
    <h1 class="main-title">Code: {{ game.gameCode }}</h1>
    <v-btn size="x-large" :color="game.isRunning ? 'red' : 'green'" @click="switchGameStatus">{{ game.isRunning ? 'Stop Party' : 'Start Party' }}</v-btn>
    <div class="block">
      <v-text-field v-model="newWineName" label="Wine Name" variant="outlined" @keyup.enter="addWine"></v-text-field>
      <v-text-field v-model="newWineCode" label="Wine Code (ex. A)" variant="outlined" @keyup.enter="addWine"></v-text-field>
      <v-text-field v-model.number="newWineYear" label="Wine Code (ex. A)" variant="outlined" @keyup.enter="addWine"></v-text-field>
      <v-btn size="x-large" variant="tonal" @click="addWine">Add Wine</v-btn>
    </div>
    <div v-if="wines.length > 0" class="block">
      <h2>Wines</h2>
      <v-table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Code</th>
            <th>Year</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="wine in wines" :key="wine.wineId">
            <td>{{ wine.wineName }}</td>
            <td>{{ wine.wineCode }}</td>
            <td>{{ wine.wineYear }}</td>
            <td><v-btn color="red" size="small" @click="removeWine(wine.wineId)">Delete</v-btn></td>
          </tr>
        </tbody>
      </v-table>
    </div>
    <div class="block">
      <v-btn variant="tonal" @click="goBack">Back</v-btn>
    </div>
    <div class="block">
      <v-btn variant="tonal" @click="removeGame">Delete Party</v-btn>
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

.main-title:first-of-type {
  margin-bottom: 0;
}

.block:first-of-type {
  margin-top: 16px;
}
</style>
