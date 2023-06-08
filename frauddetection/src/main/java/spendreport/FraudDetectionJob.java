/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package spendreport;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.datastream.KeyedStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import spendreport.detectors.ExpiredCardDetector;
import spendreport.detectors.NormalDistributionDetector;
import spendreport.detectors.OverLimitDetector;
import spendreport.detectors.SmallThenLargeDetector;
import spendreport.dtos.Alert;
import spendreport.dtos.CardTransaction;

/**
 * Skeleton code for the data stream walkthrough
 */
public class FraudDetectionJob {
	public static void main(String[] args) throws Exception {
		StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

		String bootstrapServers = "kafkac:9092";

		KafkaSource<String> kafkaSource = KafkaSource.<String>builder()
				.setBootstrapServers(bootstrapServers)
				.setTopics("transactions")
				.setGroupId("fraud-detector")
				.setValueOnlyDeserializer(new SimpleStringSchema())
				.build();

		DataStream<String> kafkaStream = env.fromSource(kafkaSource, WatermarkStrategy.noWatermarks(), "Kafka Source");

		DataStream<CardTransaction> transactions = kafkaStream
				.map(CardTransaction::fromString)
				.name("transactions");

		KeyedStream<CardTransaction, String> keyedTransactions = transactions
			// partition the stream according to the id of the account
			// it allows to have all the transactions of the same account in the same operator instance
			// and prevents race conditions when updating the state
			.keyBy(CardTransaction::getAccountId);

		DataStream<Alert> overLimitAlerts = keyedTransactions
			.process(new OverLimitDetector())
			.name("over-limit-alerts");

		DataStream<Alert> smallThenLargeAlerts = keyedTransactions
			.process(new SmallThenLargeDetector())
			.name("small-then-large-alerts");

		DataStream<Alert> expiredCardAlerts = keyedTransactions
			.process(new ExpiredCardDetector())
			.name("expired-card-alerts");

		DataStream<Alert> normalDistributionAlerts = keyedTransactions
			.process(new NormalDistributionDetector())
			.name("normal-distribution-alerts");

		DataStream<Alert> alerts = overLimitAlerts.union(smallThenLargeAlerts).union(expiredCardAlerts)
				.union(normalDistributionAlerts);

		KafkaSink<String> alertKafkaSink = KafkaSink.<String>builder()
				.setBootstrapServers(bootstrapServers)
				.setRecordSerializer(KafkaRecordSerializationSchema.builder()
						.setTopic("alerts")
						.setValueSerializationSchema(new SimpleStringSchema())
						.build()
				).build();

		alerts
			.map(Alert::toJson)
			.sinkTo(alertKafkaSink)
			.name("kafka-alerts");

		alerts
			.addSink(new AlertSink())
			.name("send-alerts");

		env.execute("Fraud Detection");
	}
}
