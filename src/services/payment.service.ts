import {
  MidtransCallbackRequest,
  MidtransGeneratePaymentLinkRequest,
} from "../DTO/midtrans.dto";
import { Order } from "../models/order.model";
import { Transaction } from "../models/transaction.model";
import { UserMembership } from "../models/userMembership.model";
import { GeneralLogRepository } from "../repo/log_repo";
import { MidtransRepository } from "../repo/midtrans/midtrans";
import { OrderRepository } from "../repo/order_repo";
import { PackageRepository } from "../repo/package_repo";
import { TransactionRepository } from "../repo/transaction_repo";
import { UserRepository } from "../repo/user_repo";
import { UserMembershipRepository } from "../repo/userMembership_repo";
import {
  MIDTRANS_STATUS_CAPTURE,
  MIDTRANS_STATUS_PENDING,
  MIDTRANS_STATUS_SETTLEMENT,
  STATUS_ORDER_FAILED,
  STATUS_ORDER_PENDING,
  STATUS_ORDER_SUCCESS,
} from "../utils/constant";

import { Logger } from "../utils/logger";
import { DefaultResponse } from "../utils/response";
import crypto from "crypto";

export class PaymentService {
  private packageRepo = new PackageRepository();
  private userRepo = new UserRepository();
  private midtransRepo = new MidtransRepository();
  private orderRepo = new OrderRepository();
  private transactionRepo = new TransactionRepository();
  private userMembershipRepo = new UserMembershipRepository();
  private generalLogRepo = new GeneralLogRepository();

  // ======================================================
  // GENERATE PAYMENT LINK
  // ======================================================
  async generatePaymentLink(
    package_id: number,
    user_id: number,
    email: string
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const packageData = await this.packageRepo.getPackageByID(package_id);
      if (!packageData || !packageData.id) {
        return { code: 404, message: "not found", error: "package not found" };
      }

      const user = await this.userRepo.getUserById(user_id);
      if (!user) {
        return { code: 404, message: "not found", error: "user not found" };
      }

      // Fallback nomor HP
      if (!user.phone_number) {
        const randomNum = Math.floor(Math.random() * 10000000000);
        user.phone_number = `08${randomNum.toString().padStart(10, "0")}`;
      }

      const paymentCode = this.generatePaymentCode(user.phone_number);
      const grossAmount = packageData.discountedPrice ?? 0;

      const midtransReq: MidtransGeneratePaymentLinkRequest = {
        transaction_details: {
          order_id: paymentCode,
          gross_amount: grossAmount,
        },
        item_details: [
          {
            name: `Package ${packageData.durationInMonth} month`,
            price: grossAmount,
            quantity: 1,
          },
        ],
        customer_details: {
          first_name: user.username,
          email,
          phone: user.phone_number,
        },
        usage_limit: 1,
        expiry: {
          start_time: new Date().toISOString(),
          duration: 1,
          unit: "days",
        },
        enabled_payments: [
          "credit_card",
          "bank_transfer",
          "gopay",
          "shopeepay",
          "qris",
        ],
      };

      const paymentURL = await this.midtransRepo.generatePaymentLink(
        midtransReq
      );
      resp.data = paymentURL;

      if (!user.user_id) {
        throw new Error("Invalid user_id");
      }

      if (!paymentURL.payment_url) {
        throw new Error("Invalid payment URL");
      }

      const orderReq: Order = {
        user_id: user.user_id,
        package_id: packageData.id!,
        status: STATUS_ORDER_PENDING,
        payment_code: paymentCode,
        payment_url: paymentURL.payment_url,
        created_by: email,
      };

      await this.orderRepo.createOrder(orderReq);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "bad request";
      resp.error = err.message;
      Logger.error(`[service][generatePaymentLink] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // PRIVATE: GENERATE PAYMENT CODE
  // ======================================================
  private generatePaymentCode(phoneNumber: string): string {
    const randomBytes = crypto.randomBytes(5).toString("hex").toUpperCase();
    const phoneLast5 = phoneNumber.slice(-5);
    const now = new Date();
    const secondsInDay =
      now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds();
    const timeComponent = secondsInDay.toString().padStart(5, "0");
    return `MT-${randomBytes.substring(0, 3)}${phoneLast5.substring(
      0,
      2
    )}${timeComponent}`;
  }

  // ======================================================
  // MIDTRANS CALLBACK HANDLER
  // ======================================================
  async midtransPaymentLinkNotification(
    req: MidtransCallbackRequest
  ): Promise<Error | null> {
    try {
      if (!req.order_id) throw new Error("order id is required");

      const splitted = req.order_id.split("-");
      if (splitted.length < 2) throw new Error("invalid order id");

      const orderCode = `${splitted[0]}-${splitted[1]}`;
      const order = await this.orderRepo.getOrderByPaymentCode(orderCode);
      if (!order || !order.id) throw new Error("order not found");

      await this.generalLogRepo.createGeneralLog({
        user_id: order.user_id,
        raw_body: `MidtransPaymentLinkNotification - ${JSON.stringify(req)}`,
        created_by: "SYSTEM",
      });

      if (order.status === STATUS_ORDER_SUCCESS) return null;

      const orderTransaction =
        await this.midtransRepo.midtransCheckStatusTransaction(req.order_id);

      if (orderTransaction.transaction_status !== req.transaction_status) {
        throw new Error("transaction status not match");
      }

      const amount = parseFloat(req.gross_amount ?? "0");

      if (
        [MIDTRANS_STATUS_SETTLEMENT, MIDTRANS_STATUS_CAPTURE].includes(
          req.transaction_status
        )
      ) {
        order.status = STATUS_ORDER_SUCCESS;

        const transactionReq: Transaction = {
          order_id: order.id!,
          webhook_id: req.transaction_id,
          amount,
          settlement_date: req.settlement_time || req.transaction_time,
          created_by: "SYSTEM",
        };

        await this.transactionRepo.createTransaction(transactionReq);

        const packageData = await this.packageRepo.getPackageByID(
          order.package_id
        );
        const userMembership =
          await this.userMembershipRepo.getUserMembershipByUserID(
            order.user_id
          );

        const addMonths = packageData?.durationInMonth ?? 1;

        if (!userMembership || !userMembership.id) {
          const expiredAt = new Date();
          expiredAt.setMonth(expiredAt.getMonth() + addMonths);

          await this.userMembershipRepo.createUserMembership({
            user_id: order.user_id,
            package_id: order.package_id,
            expired_at: expiredAt.toISOString(),
            is_membership_active: true,
            created_by: "SYSTEM",
          } as UserMembership);
        } else {
          const currExpiredAt = new Date(userMembership.expired_at);
          const newExpired = new Date(currExpiredAt);
          newExpired.setMonth(newExpired.getMonth() + addMonths);

          await this.userMembershipRepo.updateUserMembershipByUserID({
            user_id: order.user_id,
            package_id: order.package_id,
            expired_at: newExpired.toISOString(),
            is_membership_active: true,
            created_by: "SYSTEM",
          } as UserMembership);
        }

        Logger.info(
          `[service][midtransPaymentLinkNotification] order ${order.payment_code} success`
        );
      } else if (req.transaction_status !== MIDTRANS_STATUS_PENDING) {
        order.status = STATUS_ORDER_FAILED;
      }

      if (order.status !== STATUS_ORDER_PENDING) {
        await this.orderRepo.updateOrderStatus(order);
        Logger.info(
          `[service][midtransPaymentLinkNotification] order ${order.payment_code} updated to ${order.status}`
        );
      }

      return null;
    } catch (err: any) {
      Logger.error(`[service][midtransPaymentLinkNotification] ${err.message}`);
      return err;
    }
  }
}
