import { ref, type Ref } from 'vue';

export type User = {
  jwt: string;
  isAdmin?: boolean;
  gameId?: string;
}

let jwt: Ref<string> | undefined;
let isAdmin: Ref<boolean> | undefined;
let gameId: Ref<string | null> | undefined;

export function useSession() {
  if (!jwt) {
    jwt = ref(localStorage.getItem('jwt') ?? '');
    isAdmin = ref(localStorage.getItem('isAdmin') === 'true');
    gameId = ref(localStorage.getItem('gameId') ?? null);
  }
  return {
    getUser(): User {
      const user: User = { jwt: jwt?.value ?? '' };
      if (isAdmin?.value) {
        user.isAdmin = isAdmin.value;
      }
      if (gameId?.value) {
        user.gameId = gameId.value;
      }
      return user;
    },
    setUser(user: User) {
      jwt!.value = user.jwt;
      isAdmin!.value = user.isAdmin ?? false;
      gameId!.value = user.gameId ?? '';
      localStorage.setItem('jwt', user.jwt);
      localStorage.setItem('isAdmin', String(user.isAdmin ?? false));
      if (user.gameId) {
        localStorage.setItem('gameId', user.gameId);
      }
    },
    deleteUser() {
      jwt!.value = '';
      isAdmin!.value = false;
      gameId!.value = null;
      localStorage.removeItem('jwt');
      localStorage.removeItem('isAdmin');
      localStorage.removeItem('gameId');
    },
  };
}
