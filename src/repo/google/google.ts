// src/repos/GoogleRepo.ts
import { OAuth2Client } from "google-auth-library";
import axios from "axios";

export interface GoogleCallbackRequest {
  state: string;
  code: string;
}

export interface GoogleCallbackResponse {
  id: string;
  email: string;
  name: string;
  picture: string;
  verified_email: boolean;
}

export class GoogleRepository {
  private oauth2Client: OAuth2Client;

  constructor() {
    this.oauth2Client = new OAuth2Client({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      redirectUri: process.env.GOOGLE_REDIRECT_URI!,
    });
  }

  // --- LOGIN (Generate OAuth2 URL) ---
  public async login(): Promise<string> {
    const url = this.oauth2Client.generateAuthUrl({
      access_type: "offline",
      scope: [
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/userinfo.profile",
      ],
      state: "state-token",
    });
    return url;
  }

  // --- CALLBACK (Exchange Code & Get User Info) ---
  public async callback(
    req: GoogleCallbackRequest
  ): Promise<GoogleCallbackResponse> {
    if (req.state !== "state-token") {
      throw new Error("invalid credentials");
    }

    // Tukar auth code dengan access token
    const { tokens } = await this.oauth2Client.getToken(req.code);
    this.oauth2Client.setCredentials(tokens);

    if (!tokens.access_token) {
      throw new Error("No access token received from Google");
    }

    // Panggil Google userinfo endpoint
    const resp = await axios.get(
      "https://www.googleapis.com/oauth2/v2/userinfo",
      {
        headers: { Authorization: `Bearer ${tokens.access_token}` },
      }
    );

    return resp.data as GoogleCallbackResponse;
  }
}
