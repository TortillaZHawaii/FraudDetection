import { Card } from "./card";
import { CardOwner } from "./cardOwner";

export class CardTransaction {
    amount: number;
    limitLeft: number;
    currency: string;
    latitude: number;
    longitude: number;
    card: Card;
    owner: CardOwner;
    utc: Date;
  
    constructor({ amount, limit_left, currency, latitude, longitude, card, owner, utc }: { amount: number; limit_left: number; currency: string; latitude: number; longitude: number; card: any; owner: any; utc: string }) {
        this.amount = amount;
        this.limitLeft = limit_left;
        this.currency = currency;
        this.latitude = latitude;
        this.longitude = longitude;
        this.card = new Card(card);
        this.owner = new CardOwner(owner);
        this.utc = new Date(utc);
    }
}