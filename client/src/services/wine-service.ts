import { baseUrl } from './utils';

export type Wine = {
  wineId: string
  wineName: string
  wineCode: string
  wineYear: number
}

export async function getAllWines(jwt: string, gameId: string): Promise<Wine[] | false> {
  const response = await fetch(`${baseUrl}/games/${gameId}/wines`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  if (!response.ok) {
    return false;
  }
  const wines = await response.json();
  return wines;
}

export async function createWine(jwt: string, gameId: string, wineName: string, wineCode: string, wineYear: number): Promise<Wine> {
  const response = await fetch(`${baseUrl}/games/${gameId}/wines`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
    body: JSON.stringify({
      wineName,
      wineCode,
      wineYear,
    }),
  });
  const wine = await response.json();
  return wine;
}

export async function deleteWine(jwt: string, gameId: string, wineId: string): Promise<Wine> {
  const response = await fetch(`${baseUrl}/games/${gameId}/wines/${wineId}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  const wine = await response.json();
  return wine;
}