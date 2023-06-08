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

package spendreport.detectors;


import org.apache.flink.api.common.state.ValueState;
import org.apache.flink.api.common.state.ValueStateDescriptor;
import org.apache.flink.api.common.typeinfo.Types;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spendreport.dtos.Alert;
import spendreport.dtos.CardTransaction;

/**
 * Skeleton code for implementing a fraud detector.
 */
public class NormalDistributionDetector extends KeyedProcessFunction<String, CardTransaction, Alert> {

	private static final Logger LOG = LoggerFactory.getLogger(NormalDistributionDetector.class);

	private static final long serialVersionUID = 5L;

	// Fault-tolerant map-like state
	private transient ValueState<Double> averageState;
	private transient ValueState<Double> varianceState;
	private transient ValueState<Integer> countState;

	// Initialize the state in operator creation time
	// ValueState is fault-tolerant and stored in the state backend
	@Override
	public void open(Configuration parameters) {
		ValueStateDescriptor<Double> averageDescriptor = new ValueStateDescriptor<>(
			"average",
			Types.DOUBLE);
		averageState = getRuntimeContext().getState(averageDescriptor);

		ValueStateDescriptor<Double> varianceDescriptor = new ValueStateDescriptor<>(
			"variance",
			Types.DOUBLE);
		varianceState = getRuntimeContext().getState(varianceDescriptor);

		ValueStateDescriptor<Integer> countDescriptor = new ValueStateDescriptor<>(
			"count",
			Types.INT);
		countState = getRuntimeContext().getState(countDescriptor);
	}

	@Override
	public void processElement(CardTransaction transaction, KeyedProcessFunction<String, CardTransaction, Alert>.Context context, Collector<Alert> collector) throws Exception {
		// Get the current state for the current key
		Double average = averageState.value();
		Double variance = varianceState.value();
		Integer count = countState.value();

		if (count == null) {
			count = 1;
		} else {
			count++;
		}

		countState.update(count);

		double previousAverage;
		if (average == null) {
			average = (double) transaction.getAmount();
			previousAverage = average;
		} else {
			previousAverage = average;
			average = average + (transaction.getAmount() - average) / count;
		}

		averageState.update(average);

		double previousStandardDeviation = 0.0;
		if (count >= 2) {
			if (variance == null) {
				variance = 0.0;
			}
			previousStandardDeviation = Math.sqrt(variance);
			variance = variance + (transaction.getAmount() - previousAverage) * (transaction.getAmount() - average);
		}

		varianceState.update(variance);
		double standardDeviation = variance == null ? 0.0 : Math.sqrt(variance);

		LOG.info("Transaction amount: " + transaction.getAmount());
		LOG.info("Average: " + average + ", Variance: " + variance + ", Standard Deviation: " + standardDeviation +  ", Count: " + count);

		if (count > 10 && Math.abs(transaction.getAmount() - previousAverage) > previousStandardDeviation) {
			Alert alert = new Alert();
			alert.setTransaction(transaction);
			alert.setReason("Transaction is out of average range by more than standard deviation");
			collector.collect(alert);
		}
	}
}
