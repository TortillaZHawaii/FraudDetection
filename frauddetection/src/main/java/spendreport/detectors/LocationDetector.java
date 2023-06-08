package spendreport.detectors;

import org.apache.flink.streaming.api.functions.windowing.ProcessWindowFunction;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;
import org.apache.flink.util.Collector;
import spendreport.dtos.Alert;
import spendreport.dtos.CardTransaction;

public class LocationDetector extends ProcessWindowFunction<CardTransaction, Alert, String, TimeWindow> {
    static final double MAX_DISTANCE = 0.02;

    @Override
    public void process(String s, ProcessWindowFunction<CardTransaction, Alert, String, TimeWindow>.Context context, Iterable<CardTransaction> iterable, Collector<Alert> collector) {
        double latSum = 0;
        double longSum = 0;
        int count = 0;

        for (CardTransaction transaction : iterable) {
            latSum += transaction.getLatitude();
            longSum += transaction.getLongitude();
            count++;
        }

        double latMean = latSum / count;
        double longMean = longSum / count;

        for (CardTransaction transaction : iterable) {
            double distance = Math.pow(transaction.getLatitude() - latMean, 2)
                    + Math.pow(transaction.getLongitude() - longMean, 2);
            if (distance > MAX_DISTANCE * MAX_DISTANCE) {
                Alert alert = new Alert();
                alert.setReason("Too far from average location in window");
                alert.setTransaction(transaction);
                collector.collect(alert);
            }
        }
    }
}
