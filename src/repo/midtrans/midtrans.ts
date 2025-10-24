// src/repos/MidtransRepo.ts
import axios from "axios";
import * as dotenv from "dotenv";
import {
  Expiry,
  MidtransCheckStatusResponse,
  MidtransGeneratePaymentLinkRequest,
  MidtransGeneratePaymentLinkResponse,
} from "DTO/midtrans.dto";

dotenv.config();

export class MidtransRepository {
  private baseUrl: string;
  private clientSecret: string;

  constructor() {
    this.baseUrl = process.env.MIDTRANS_BASE_URL || "";
    this.clientSecret = process.env.MIDTRANS_CLIENT_SECRET || "";
  }

  private generateTokenAuthorization(): string {
    const encoded = Buffer.from(this.clientSecret + ":").toString("base64");
    return `Basic ${encoded}`;
  }

  // --- Generate Payment Link ---
  async generatePaymentLink(
    req: MidtransGeneratePaymentLinkRequest
  ): Promise<MidtransGeneratePaymentLinkResponse> {
    const url = `${this.baseUrl}/v1/payment-links`;

    const requestPayload: MidtransGeneratePaymentLinkRequest = {
      transaction_details: req.transaction_details,
      usage_limit: 1,
      expiry: {
        start_time: new Date().toISOString().slice(0, 16).replace("T", " "), // format simplified
        duration: 1,
        unit: "days",
      } as Expiry,
      enabled_payments: [
        "credit_card",
        "bca_klikbca",
        "gopay",
        "permata_va",
        "bca_va",
        "bri_va",
        "bni_va",
        "indomaret",
        "shopeepay",
      ],
      item_details: req.item_details,
      customer_details: req.customer_details,
    };

    try {
      const res = await axios.post<MidtransGeneratePaymentLinkResponse>(
        url,
        requestPayload,
        {
          headers: {
            Authorization: this.generateTokenAuthorization(),
            "Content-Type": "application/json",
          },
        }
      );

      if (res.data.error_messages) {
        throw new Error(
          `Failed to generate payment link: ${res.data.error_messages}`
        );
      }

      return res.data;
    } catch (err: any) {
      throw new Error(
        `[Midtrans][generatePaymentLink] error: ${err.message || err}`
      );
    }
  }

  // --- Check Status Transaction ---
  async midtransCheckStatusTransaction(
    orderID: string
  ): Promise<MidtransCheckStatusResponse> {
    const url = `${this.baseUrl}/v2/${orderID}/status`;

    try {
      const res = await axios.get<MidtransCheckStatusResponse>(url, {
        headers: {
          Authorization: this.generateTokenAuthorization(),
          "Content-Type": "application/json",
        },
      });

      return res.data;
    } catch (err: any) {
      throw new Error(
        `[Midtrans][midtransCheckStatusTransaction] error: ${
          err.message || err
        }`
      );
    }
  }
}
