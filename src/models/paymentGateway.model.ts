export interface TransactionDetails {
  order_id: string;
  gross_amount: number;
  // payment_link_id?: string; // dikomentari di Go
}

export interface Expiry {
  start_time: string;
  duration: number;
  unit: string;
}

export interface ItemDetails {
  // id?: string;
  name: string;
  price: number;
  quantity: number;
  // brand?: string;
  // category?: string;
  // merchant_name?: string;
}

export interface CustomerDetails {
  first_name: string;
  // last_name?: string;
  email: string;
  phone: string;
  // notes?: string;
}

export interface MidtransGeneratePaymentLinkRequest {
  transaction_details: TransactionDetails;
  usage_limit: number;
  expiry: Expiry;
  enabled_payments: string[];
  item_details: ItemDetails[];
  customer_details: CustomerDetails;
}

export interface MidtransGeneratePaymentLinkResponse {
  order_id?: string;
  payment_url?: string;
  error_messages?: any;
}

// --- Callback ---
export interface MidtransCalbackMetadata {
  extra_info: {
    payment_link_id: string;
    payment_link_type: string;
  };
}

export interface MidtransCallbackVaNumbers {
  bank: string;
  va_number: string;
}

export interface MidtransCallbackPaymentAmount {
  paid_at: string;
  amount: string;
}

export interface MidtransCallbackRequest {
  currency?: string;
  custom_field1?: string;
  custom_field2?: string;
  custom_field3?: string;
  expiry_time?: string;
  fraud_status?: string;
  gross_amount?: string;
  merchant_id?: string;
  metadata?: MidtransCalbackMetadata;
  order_id?: string;
  payment_amounts?: MidtransCallbackPaymentAmount[];
  payment_type?: string;
  settlement_time?: string;
  signature_key?: string;
  status_code?: string;
  status_message?: string;
  transaction_id?: string;
  transaction_status?: string;
  transaction_time?: string;
  va_numbers?: MidtransCallbackVaNumbers[];

  // additional for credit card
  masked_card?: string;
  eci?: string;
  channel_response_message?: string;
  channel_response_code?: string;
  card_type?: string;
  bank?: string;
  approval_code?: string;

  // additional for QRIS
  issuer?: string;
  acquirer?: string;

  // additional for Permata VA
  permata_va_number?: string;
}

export interface MidtransCheckStatusResponse {
  masked_card: string;
  approval_code: string;
  bank: string;
  eci: string;
  channel_response_code: string;
  channel_response_message: string;
  transaction_time: string;
  gross_amount: string;
  currency: string;
  order_id: string;
  payment_type: string;
  signature_key: string;
  status_code: string;
  transaction_id: string;
  transaction_status: string;
  fraud_status: string;
  settlement_time: string;
  status_message: string;
  merchant_id: string;
  card_type: string;
  three_ds_version: string;
  challenge_completion: boolean;
}
