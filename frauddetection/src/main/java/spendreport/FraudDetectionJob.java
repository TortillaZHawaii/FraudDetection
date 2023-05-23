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
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import spendreport.dtos.Alert;
import spendreport.dtos.CardTransaction;

import java.util.Properties;

/**
 * Skeleton code for the data stream walkthrough
 */
public class FraudDetectionJob {
	public static void main(String[] args) throws Exception {
		StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

		KafkaSource<String> kafkaSource = KafkaSource.<String>builder()
				.setBootstrapServers("kafka:9092")
				.setTopics("transactions")
				.setGroupId("fraud-detector")
				.setValueOnlyDeserializer(new SimpleStringSchema())
				.build();

		DataStream<String> kafkaStream = env.fromSource(kafkaSource, WatermarkStrategy.noWatermarks(), "Kafka Source");

		DataStream<CardTransaction> transactions = kafkaStream
				.map(CardTransaction::fromString)
				.name("transactions");

		DataStream<Alert> alerts = transactions
			// partition the stream according to the id of the account
			// it allows to have all the transactions of the same account in the same operator instance
			// and prevents race conditions when updating the state
			.keyBy(CardTransaction::getAccountId)
			.process(new FraudDetector())
			.name("fraud-detector");

		alerts
			.addSink(new AlertSink())
			.name("send-alerts");

		env.execute("Fraud Detection");
	}
}
