export interface Attachment {
  content?: string; // omitempty → optional
  type?: string; // omitempty → optional
  filename?: string; // omitempty → optional
  disposition?: string; // omitempty → optional
  fileName?: string; // omitempty → optional
  contentId?: string; // omitempty → optional
}

export interface SendgridSendEmailRequest {
  from: string; // required, email
  to: string; // required, email
  senderName: string; // required
  recepientName: string; // required
  subject: string; // required
  body: string; // required
  attach: Attachment; // required (tapi fields di dalam optional)
}
