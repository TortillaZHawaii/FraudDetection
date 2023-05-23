import { CardTransaction } from "./cardTransaction";

export class Alert {
    transaction: CardTransaction;
    reason: string;
  
    constructor(transaction: CardTransaction, reason: string) {
      this.transaction = transaction;
      this.reason = reason;
    }
}