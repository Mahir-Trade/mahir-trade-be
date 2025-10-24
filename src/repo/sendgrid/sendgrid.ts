import sgMail from "@sendgrid/mail";
import dotenv from "dotenv";

dotenv.config();

export interface SendgridAttachment {
  content: string;
  type: string;
  filename: string;
  disposition: string;
  contentId: string;
}

export interface SendgridSendEmailRequest {
  from: string;
  to: string;
  subject: string;
  body: string;
  senderName: string;
  recepientName: string;
  attach?: SendgridAttachment;
}

export class SendgridRepository {
  constructor() {
    const apiKey =
      process.env.SENDGRID_API_KEY ||
      "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";
    if (!apiKey) throw new Error("SENDGRID_API_KEY is not set");
    sgMail.setApiKey(apiKey);
  }

  // --- SEND EMAIL ---
  public async sendEmail(req: SendgridSendEmailRequest): Promise<void> {
    const msg: any = {
      to: {
        email: req.to,
        name: req.recepientName,
      },
      from: {
        email: req.from,
        name: req.senderName,
      },
      subject: req.subject,
      html: req.body,
    };

    if (req.attach && req.attach.content) {
      msg.attachments = [
        {
          content: req.attach.content,
          type: req.attach.type,
          filename: req.attach.filename,
          disposition: req.attach.disposition,
          content_id: req.attach.contentId,
        },
      ];
    }

    const response = await sgMail.send(msg);

    if (response[0].statusCode !== 202) {
      throw new Error(
        `Sendgrid failed with status code ${response[0].statusCode}`
      );
    }
  }
}
