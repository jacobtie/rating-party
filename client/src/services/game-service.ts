import { baseUrl } from './utils';

export type Game = {
  gameId: string
  gameCode: string
  gameName: string
  isRunning: boolean
  areResultsShared: boolean
}

export async function getAllGames(jwt: string): Promise<Game[] | false> {
  const response = await fetch(`${baseUrl}/games`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  if (!response.ok) {
    return false;
  }
  const games: Game[] = await response.json();
  return games;
}

export async function getGame(jwt: string, gameId: string): Promise<Game | false> {
  const response = await fetch(`${baseUrl}/games/${gameId}`, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  if (!response.ok) {
    return false;
  }
  const game: Game = await response.json();
  return game;
}

export async function createGame(jwt: string, gameName: string): Promise<Game> {
  const response = await fetch(`${baseUrl}/games`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
    body: JSON.stringify({ gameName }),
  });
  const game: Game = await response.json();
  return game;
}

export async function updateGame(jwt: string, gameId: string, gameName: string, isRunning: boolean, areResultsShared = false): Promise<void> {
  await fetch(`${baseUrl}/games/${gameId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
    body: JSON.stringify({
      gameName,
      isRunning,
      areResultsShared,
    }),
  });
}

export async function deleteGame(jwt: string, gameId: string): Promise<Game> {
  const response = await fetch(`${baseUrl}/games/${gameId}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwt}`,
    },
  });
  const game: Game = await response.json();
  return game;
}