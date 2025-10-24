import { Order } from "@models/order.model";
import db from "../db/db.config";
import { Pool } from "pg";
import { OrderQueries } from "./queries/order_queries";

export class OrderRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE ORDER ---
  async createOrder(req: Order): Promise<number> {
    try {
      if (!req.status) {
        req.status = "pending";
      }

      const result = await this.pg.query<{ id: number }>(
        OrderQueries.CreateOrder,
        [
          req.user_id,
          req.package_id,
          req.status,
          req.payment_code,
          req.payment_url,
          req.created_by,
        ]
      );

      if (result.rowCount === 0) {
        throw new Error("failed to create order");
      }

      return result.rows[0].id;
    } catch (err: any) {
      console.error(`[orderRepo][createOrder] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET ORDER BY PAYMENT CODE ---
  async getOrderByPaymentCode(paymentCode: string): Promise<Order> {
    try {
      const result = await this.pg.query<Order>(
        OrderQueries.GetOrderByPaymentCode,
        [paymentCode]
      );

      if (result.rowCount === 0) {
        throw new Error(`order with payment code ${paymentCode} not found`);
      }

      return result.rows[0];
    } catch (err: any) {
      console.error(`[orderRepo][getOrderByPaymentCode] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE ORDER STATUS ---
  async updateOrderStatus(req: Order): Promise<boolean> {
    try {
      const result = await this.pg.query(OrderQueries.UpdateOrderStatus, [
        req.status,
        req.updated_by,
        req.id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[orderRepo][updateOrderStatus] error: ${err.message}`);
      throw err;
    }
  }
}
