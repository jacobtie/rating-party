<script setup lang="ts">
import WineRating from '@/components/WineRating.vue';
import { useSession } from '@/composables/session';
import router from '@/router';
import { getGame, type Game } from '@/services/game-service';
import { getAllRatings, getResults, putRating, type Rating } from '@/services/rating-service';
import { getAllWines, type Wine } from '@/services/wine-service';
import { ref } from 'vue';

const { getUser, deleteUser } = useSession();

const user = getUser();

if (!user.jwt) {
  router.push('/');
}

const game = ref<Game | null>(null);
(async () => {
  const { gameId } = user;
  if (!gameId) {
    deleteUser();
    router.push('/');
    return;
  }
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

const wines = ref<Wine[]>([]);
const ratings = ref<Rating[]>([]);
const results = ref<Record<string, unknown>[]>([]);
(async () => {
  const { gameId } = user;
  if (!gameId) {
    deleteUser();
    router.push('/');
    return;
  }
  try {
    const winesFromServer = await getAllWines(user.jwt, gameId);
    if (winesFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    wines.value = winesFromServer;
    const ratingsFromServer = await getAllRatings(user.jwt, gameId);
    if (ratingsFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    ratings.value = ratingsFromServer;
    const wineIds = wines.value.map((wine) => wine.wineId);
    ratings.value = ratings.value.filter((rating) => wineIds.includes(rating.wineId));
    for (const wine of wines.value ) {
      if (ratings.value.some((rating) => rating.wineId === wine.wineId)) continue;
      ratings.value.push({
        wineId: wine.wineId,
        gameId: gameId,
        sightRating: 0,
        aromaRating: 0,
        tasteRating: 0,
        overallRating: 0,
        comments: '',
      });
    }
    ratings.value.sort((a, b) => {
      const wineA = wines.value.find((wine) => wine.wineId === a.wineId)!;
      const wineB = wines.value.find((wine) => wine.wineId === b.wineId)!;
      return wineA.wineCode.localeCompare(wineB.wineCode);
    });
    if (!game.value?.areResultsShared) return;
    const resultsFromServer = await getResults(user.jwt, user.gameId!);
    if (resultsFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    results.value = resultsFromServer;
  } catch (err) {
    console.error(err);
  }
})();
const saveRatings = async () => {
  try {
    const gameFromServer = await getGame(user.jwt, user.gameId!);
    if (gameFromServer === false) {
      deleteUser();
      router.push('/');
      return;
    }
    game.value = gameFromServer;
    if (!game.value.isRunning) {
      alert('The game was already ended');
      return;
    }
    for (const rating of ratings.value) {
      if (rating.sightRating < 0 || rating.sightRating > 4) return showValidationErrorSnackbar();
      if (rating.aromaRating < 0 || rating.aromaRating > 6) return showValidationErrorSnackbar();
      if (rating.tasteRating < 0 || rating.tasteRating > 6) return showValidationErrorSnackbar();
      if (rating.overallRating < 0 || rating.overallRating > 4) return showValidationErrorSnackbar();
    }
    await Promise.all(ratings.value.map((rating) => putRating(user.jwt, rating)));
    showSuccessSnackbar();
  } catch (err) {
    console.log(err);
  }
};

const isErrorSnackbarActive = ref(false);
const showValidationErrorSnackbar = () => {
  isErrorSnackbarActive.value = true;
  setTimeout(() => {
    isErrorSnackbarActive.value = false;
  }, 3000);
};

const isSuccessSnackbarActive = ref(false);
const showSuccessSnackbar = () => {
  isSuccessSnackbarActive.value = true;
  setTimeout(() => {
    isSuccessSnackbarActive.value = false;
  }, 3000);
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
    <div v-if="game.isRunning && wines.length > 0">
      <div v-for="rating of ratings" :key="rating.wineId" class="block">
        <WineRating :rating="rating" :wines="wines" />
      </div>
      <div class="block">
        <v-btn color="green" @click="saveRatings">Save</v-btn>
      </div>
    </div>
    <div v-else-if="!game.areResultsShared" class="block">
      <p>Party is stopped. Please refresh when the host announces that the party has started or has shared the results.</p>
    </div>
    <div v-else class="block">
      <h2>Results</h2>
      <v-table class="results-table">
        <thead>
          <tr>
            <th>Wine Name</th>
            <th>Wine Code</th>
            <th>Wine Year</th>
            <th>Average</th>
            <th>Rank</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="res of results" :key="(res.wineId as string)">
            <td>{{ res.wineName }}</td>
            <td>{{ res.wineCode }}</td>
            <td>{{ res.wineYear }}</td>
            <td>{{ res.avg }}</td>
            <td>{{ res.rank }}</td>
          </tr>
        </tbody>
      </v-table>
    </div>
    <div class="block">
      <v-btn variant="tonal" @click="logout">Logout</v-btn>
    </div>
    <v-snackbar v-model="isErrorSnackbarActive" color="red">
      Please enter valid ratings
    </v-snackbar>
    <v-snackbar v-model="isSuccessSnackbarActive" color="green">
      Ratings saved
    </v-snackbar>
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
</style>
