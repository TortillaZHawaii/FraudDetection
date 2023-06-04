export class Card {
    card_number: string;
    card_type: string;
    exp_month: number;
    exp_year: number;
    cvv: string;

    constructor({ card_number, card_type, exp_month, exp_year, cvv }: { card_number: string; card_type: string; exp_month: number; exp_year: number; cvv: string }) {
        this.card_number = card_number;
        this.card_type = card_type;
        this.exp_month = exp_month;
        this.exp_year = exp_year;
        this.cvv = cvv;
    }
}