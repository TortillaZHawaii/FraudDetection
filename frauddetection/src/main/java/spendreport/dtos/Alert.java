package spendreport.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.Data;

@Data
public class Alert {
    @JsonProperty("transaction")
    private CardTransaction transaction;
    @JsonProperty("reason")
    private String reason;

    public String toJson() {
        ObjectMapper mapper = new ObjectMapper();
        try {
            return mapper.writeValueAsString(this);
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }
}
