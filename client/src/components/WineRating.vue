<script setup lang="ts">
import type { Rating } from '@/services/rating-service';
import type { Wine } from '@/services/wine-service';
import { computed } from 'vue';

const props = defineProps<{
  rating: Rating
  wines: Wine[]
}>();

const wine = props.wines.find((w) => w.wineId === props.rating.wineId)!;

const validation = {
  isBetweenInclusive: (min: number, max: number) => (v: number) => v >= min && v <= max,
};

const totalRating = computed(() => {
  return (!Number.isNaN(Number(props.rating.sightRating)) ? Number(props.rating.sightRating) : 0)
    + (!Number.isNaN(Number(props.rating.aromaRating)) ? Number(props.rating.aromaRating) : 0)
    + (!Number.isNaN(Number(props.rating.tasteRating)) ? Number(props.rating.tasteRating) : 0)
    + (!Number.isNaN(Number(props.rating.overallRating)) ? Number(props.rating.overallRating) : 0);
});
</script>

<template>
  <div>
    <h3>Wine {{ wine.wineCode }}</h3>
    <v-text-field v-model.number="rating.sightRating" :rules="[validation.isBetweenInclusive(0, 4)]" label="Sight (0-4)"></v-text-field>
    <v-text-field v-model.number="rating.aromaRating" :rules="[validation.isBetweenInclusive(0, 6)]" label="Aroma (0-6)"></v-text-field>
    <v-text-field v-model.number="rating.tasteRating" :rules="[validation.isBetweenInclusive(0, 6)]" label="Taste (0-6)"></v-text-field>
    <v-text-field v-model.number="rating.overallRating" :rules="[validation.isBetweenInclusive(0, 4)]" label="Overall (0-4)"></v-text-field>
    <v-text-field v-model.number="totalRating" disabled label="Total"></v-text-field>
    <v-textarea v-model="rating.comments" label="Comments"></v-textarea>
  </div>
</template>
