export class CardOwner {
    id: number;
    firstName: string;
    lastName: string;

    constructor({ id, first_name, last_name }: { id: number; first_name: string; last_name: string }) {
        this.id = id;
        this.firstName = first_name;
        this.lastName = last_name;
    }
}