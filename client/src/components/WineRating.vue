<script setup lang="ts">
import type { Rating } from '@/services/rating-service';
import type { Wine } from '@/services/wine-service';

const props = defineProps<{
  rating: Rating
  wines: Wine[]
}>();

const wine = props.wines.find((w) => w.wineId === props.rating.wineId)!;

const validation = {
  isBetweenInclusive: (min: number, max: number) => (v: number) => v >= min && v <= max,
};
</script>

<template>
  <div>
    <h3>Wine {{ wine.wineCode }}</h3>
    <v-text-field v-model.number="rating.sightRating" :rules="[validation.isBetweenInclusive(0, 4)]" label="Sight (0-4)"></v-text-field>
    <v-text-field v-model.number="rating.aromaRating" :rules="[validation.isBetweenInclusive(0, 6)]" label="Aroma (0-6)"></v-text-field>
    <v-text-field v-model.number="rating.tasteRating" :rules="[validation.isBetweenInclusive(0, 6)]" label="Taste (0-6)"></v-text-field>
    <v-text-field v-model.number="rating.overallRating" :rules="[validation.isBetweenInclusive(0, 4)]" label="Overall (0-4)"></v-text-field>
    <v-textarea v-model="rating.comments" label="Comments"></v-textarea>
  </div>
</template>
