export interface EconomicCalendarItem {
  calendarId: string; // CalendarId
  date: string; // Date
  country: string; // Country
  category: string; // Category
  event: string; // Event
  reference: string; // Reference
  referenceDate: string; // ReferenceDate
  source: string; // Source
  sourceURL: string; // SourceURL
  actual: string; // Actual
  previous: string; // Previous
  forecast: string; // Forecast
  teForecast: string; // TEForecast
  url: string; // URL
  dateSpan: string; // DateSpan
  importance: number; // Importance (int → number)
  lastUpdate: string; // LastUpdate
  revised: string; // Revised
  currency: string; // Currency
  unit: string; // Unit
  ticker: string; // Ticker
  symbol: string; // Symbol
}
