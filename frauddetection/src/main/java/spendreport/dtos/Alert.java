package spendreport.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class Alert {
    @JsonProperty("transaction")
    private CardTransaction transaction;
    @JsonProperty("reason")
    private String reason;
}
