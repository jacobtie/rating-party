import { baseUrl } from './utils';

export type Rating = {
  ratingId?: string
  gameId: string
  participantId?: string
  username: string
  wineId: string
  sightRating: number
  aromaRating: number
  tasteRating: number
  overallRating: number
  comments: string
}

export async function getAllRatings(jwt: string, gameId: string): Promise<Rating[] | false> {
  const response = await fetch(`${baseUrl}/games/${gameId}/ratings`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  if (!response.ok) {
    return false;
  }
  const ratings: Rating[] = await response.json();
  return ratings;
}

export async function putRating(jwt: string, rating: Rating): Promise<void> {
  await fetch(`${baseUrl}/games/${rating.gameId}/wines/${rating.wineId}/ratings`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
    body: JSON.stringify(rating),
  });
}

export async function getResults(jwt: string, gameId: string): Promise<Record<string, unknown>[] | false> {
  const response = await fetch(`${baseUrl}/games/${gameId}/ratings/results`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  if (!response.ok) {
    return false;
  }
  const results: Record<string, unknown>[] = await response.json();
  return results;
}