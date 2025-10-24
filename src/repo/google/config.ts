export class GoogleCfg {
  clientId: string;
  clientSecret: string;
  redirectUrl: string;
  serviceAccountFilePath: string;
  videoBucketName: string;
  imageBucketName: string;
  fileBucketName: string;

  constructor() {
    this.clientId = process.env.GOOGLE_CLIENT_ID || "";
    this.clientSecret = process.env.GOOGLE_CLIENT_SECRET || "";
    this.redirectUrl = process.env.GOOGLE_REDIRECT_URL || "";
    this.serviceAccountFilePath =
      process.env.GOOGLE_SERVICE_ACCOUNT_FILE_PATH || "";
    this.videoBucketName = process.env.GOOGLE_VIDEO_BUCKET_NAME || "";
    this.imageBucketName = process.env.GOOGLE_IMAGE_BUCKET_NAME || "";
    this.fileBucketName = process.env.GOOGLE_FILE_BUCKET_NAME || "";
  }
}
