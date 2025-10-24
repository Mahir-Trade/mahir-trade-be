// src/repos/GoogleRepo.ts
import { Storage, GetSignedUrlConfig } from "@google-cloud/storage";
import { TranscoderServiceClient } from "@google-cloud/video-transcoder";
import { exec } from "child_process";
import { promisify } from "util";
import * as fs from "fs";
import * as path from "path";
import * as url from "url";
import { FileUpload } from "@dto/google.dto";

const execAsync = promisify(exec);

export interface URLParserResponse {
  url: string;
  host: string;
  path: string;
  bucketName: string;
}

export class BucketRepository {
  private storage: Storage;
  private transcoder: TranscoderServiceClient;
  private serviceAccountPath: string;

  constructor() {
    this.serviceAccountPath =
      process.env.GOOGLE_SERVICE_ACCOUNT_FILE_PATH || "";
    this.storage = new Storage({ keyFilename: this.serviceAccountPath });
    this.transcoder = new TranscoderServiceClient({
      keyFilename: this.serviceAccountPath,
    });
  }

  // --- PRESIGNED URL ---
  async presignedURL(bucketName: string, privateUrl: string): Promise<string> {
    const urlParsed = new URL(privateUrl);
    const decodePath = decodeURIComponent(urlParsed.pathname);
    const cleanUrl = `${urlParsed.protocol}//${urlParsed.host}${decodePath}`;

    const objectName = cleanUrl.replace(
      `https://${urlParsed.host}/${bucketName}/`,
      ""
    );
    const temporaryAccess = Date.now() + 60 * 60 * 1000; // 1h

    const jsonKey = fs.readFileSync(this.serviceAccountPath);
    const key = JSON.parse(jsonKey.toString());

    const options: GetSignedUrlConfig = {
      version: "v4",
      action: "read",
      expires: temporaryAccess,
    };

    const [signedUrl] = await this.storage
      .bucket(bucketName)
      .file(objectName)
      .getSignedUrl(options);
    return signedUrl;
  }

  // --- UPLOAD FILE ---
  async uploadFile(form: FileUpload): Promise<string> {
    const dest = `gs://${form.bucket_name}/${form.file_name}`;

    // Upload file ke bucket
    await execAsync(`gsutil cp ${form.localFilePath} ${dest}`);

    // Set Content-Type
    await execAsync(
      `gsutil setmeta -h "Content-Type:${form.file_content_type}" ${dest}`
    );

    // Generate signed URL
    const signed = await this.presignedURL(
      form.bucket_name,
      `https://storage.cloud.google.com/${form.bucket_name}/${form.file_name}`
    );
    return signed;
  }
  // --- URL PARSER ---
  urlParser(fileURL: string): URLParserResponse {
    const urlParsed = new URL(fileURL);
    const decodePath = decodeURIComponent(urlParsed.pathname);
    const cleanUrl = `${urlParsed.protocol}//${urlParsed.host}${decodePath}`;

    return {
      url: cleanUrl,
      host: urlParsed.host,
      path: decodePath,
      bucketName: decodePath.split("/")[1],
    };
  }

  // --- START TRANSCODING JOB ---
  async startTranscodingJob(
    bucketName: string,
    filename: string
  ): Promise<void> {
    const ctx = { timeout: 30 * 60 * 1000 }; // 30 minutes

    const inputURI = `gs://${bucketName}/${filename}`;
    const outputURI = `gs://${bucketName}_transcoded/`;

    const fileName240p = filename.replace(".mp4", "-240p");
    const fileName360p = filename.replace(".mp4", "-360p");
    const fileName480p = filename.replace(".mp4", "-480p");
    const fileName720p = filename.replace(".mp4", "-720p");

    const parent = "projects/mahir-trade-429013/locations/asia-southeast1";

    const job = {
      parent,
      job: {
        inputUri: inputURI,
        outputUri: outputURI,
        config: {
          elementaryStreams: [
            {
              key: "audio",
              audioStream: { codec: "aac", bitrateBps: 128000 },
            },
            {
              key: fileName240p,
              videoStream: {
                h264: {
                  bitrateBps: 500000,
                  frameRate: 30,
                  heightPixels: 240,
                  widthPixels: 426,
                },
              },
            },
            {
              key: fileName360p,
              videoStream: {
                h264: {
                  bitrateBps: 800000,
                  frameRate: 30,
                  heightPixels: 360,
                  widthPixels: 640,
                },
              },
            },
            {
              key: fileName480p,
              videoStream: {
                h264: {
                  bitrateBps: 1000000,
                  frameRate: 30,
                  heightPixels: 480,
                  widthPixels: 854,
                },
              },
            },
            {
              key: fileName720p,
              videoStream: {
                h264: {
                  bitrateBps: 2500000,
                  frameRate: 30,
                  heightPixels: 720,
                  widthPixels: 1280,
                },
              },
            },
          ],
          muxStreams: [
            {
              key: fileName240p,
              elementaryStreams: [fileName240p, "audio"],
              container: "mp4",
            },
            {
              key: fileName360p,
              elementaryStreams: [fileName360p, "audio"],
              container: "mp4",
            },
            {
              key: fileName480p,
              elementaryStreams: [fileName480p, "audio"],
              container: "mp4",
            },
            {
              key: fileName720p,
              elementaryStreams: [fileName720p, "audio"],
              container: "mp4",
            },
          ],
        },
      },
    };

    console.log(
      `[repo][google][StartTranscodingJob] starting job for file: ${filename}`
    );
    const [resp] = await this.transcoder.createJob(job as any);

    const jobName = resp.name!;
    while (true) {
      const [getJobResp] = await this.transcoder.getJob({ name: jobName });
      if (getJobResp.state === "SUCCEEDED") {
        console.log(
          `[repo][google][StartTranscodingJob] Job succeeded: ${jobName}`
        );
        break;
      }
      if (getJobResp.state === "FAILED") {
        console.error(
          `[repo][google][StartTranscodingJob] Job failed: ${getJobResp.error?.message}`
        );
        break;
      }
      await new Promise((r) => setTimeout(r, 5000));
    }
  }
}
