import { MidtransCallbackRequest } from "../DTO/midtrans.dto";
import { PaymentService } from "../services/payment.service";
import { DefaultResponse } from "../utils/response";
import { Request, Response } from "express";

export class PaymentController {
  private paymentService: PaymentService;

  constructor() {
    this.paymentService = new PaymentService();

    this.createPayment = this.createPayment.bind(this);
    this.paymentLinkCallback = this.paymentLinkCallback.bind(this);
  }

  // ======================================================
  // CREATE PAYMENT
  // ======================================================
  async createPayment(req: Request, res: Response): Promise<Response> {
    try {
      const { package_id, email, user_id } = req.body;

      if (!package_id || !user_id || !email) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "package_id, user_id, and email are required",
        } as DefaultResponse);
      }

      const resp = await this.paymentService.generatePaymentLink(
        package_id,
        user_id,
        email
      );

      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[PaymentController][createPayment] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // ======================================================
  // MIDTRANS CALLBACK
  // ======================================================
  async paymentLinkCallback(req: Request, res: Response): Promise<Response> {
    try {
      const callbackReq = req.body as MidtransCallbackRequest;
      const err = await this.paymentService.midtransPaymentLinkNotification(
        callbackReq
      );

      if (err) {
        console.error("[PaymentController][paymentLinkCallback]", err);
        return res.status(200).json({ message: "success" });
      }

      return res.status(200).json({ message: "success" });
    } catch (err: any) {
      console.error("[PaymentController][paymentLinkCallback] error:", err);
      return res.status(200).json({ message: "success" });
    }
  }
}
