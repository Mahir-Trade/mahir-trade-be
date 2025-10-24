export interface DiscordRoleRequest {
  user_id: string;
}

export interface GetDiscordUserRequest {
  username: string;
}

export interface DiscordUser {
  id: string;
  email: string;
  username: string;
}
