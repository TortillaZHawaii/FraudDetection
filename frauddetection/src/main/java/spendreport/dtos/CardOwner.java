import com.fasterxml.jackson.annotation.JsonProperty;

@Data
public class CardOwner {
    @JsonProperty("id")
    private int id;
    @JsonProperty("first_name")
    private String firstName;
    @JsonProperty("last_name")
    private String lastName;
}
