export interface DiscordAccount {
  id?: number; // int64 -> number
  uuid?: string; // string
  userId?: number; // int64 -> number
  discordAccountId?: any; // string
  username?: string; // string
  email?: string; // string
  createdAt?: string; // string timestamp
  updatedAt?: string; // string timestamp
}
