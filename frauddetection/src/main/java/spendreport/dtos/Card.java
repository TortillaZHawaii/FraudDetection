import com.fasterxml.jackson.annotation.JsonProperty;

@Data
public class Card {
    @JsonProperty("card_number")
    private String cardNumber;
    @JsonProperty("card_type")
    private String cardType;
    @JsonProperty("exp_month")
    private int expMonth;
    @JsonProperty("exp_year")
    private int expYear;
    @JsonProperty("cvv")
    private String cvv;
}


