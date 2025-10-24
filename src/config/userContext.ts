import { AsyncLocalStorage } from "async_hooks";

interface UserContextData {
  email: string;
  username: string;
}

const userContext = new AsyncLocalStorage<UserContextData>();

export const UserContext = {
  run: (data: UserContextData, callback: () => void) => {
    userContext.run(data, callback);
  },
  get: (): UserContextData | undefined => {
    return userContext.getStore();
  },
};
