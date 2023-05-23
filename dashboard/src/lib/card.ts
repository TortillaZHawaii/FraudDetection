export class Card {
    cardNumber: string;
    cardType: string;
    expMonth: number;
    expYear: number;
    cvv: string;

    constructor({ card_number, card_type, exp_month, exp_year, cvv }: { card_number: string; card_type: string; exp_month: number; exp_year: number; cvv: string }) {
        this.cardNumber = card_number;
        this.cardType = card_type;
        this.expMonth = exp_month;
        this.expYear = exp_year;
        this.cvv = cvv;
    }
}