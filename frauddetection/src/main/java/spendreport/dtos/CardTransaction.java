import com.fasterxml.jackson.annotation.JsonProperty;

@Data
public class CardTransaction {
    @JsonProperty("amount")
    private int amount;
    @JsonProperty("limit_left")
    private int limitLeft;
    @JsonProperty("currency")
    private String currency;
    @JsonProperty("latitude")
    private double latitude;
    @JsonProperty("longitude")
    private double longitude;
    @JsonProperty("card")
    private Card card;
    @JsonProperty("owner")
    private CardOwner owner;
    @JsonProperty("utc")
    private Date utc;
} 
