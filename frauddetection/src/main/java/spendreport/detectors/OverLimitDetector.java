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


import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import spendreport.dtos.Alert;
import spendreport.dtos.CardTransaction;

/**
 * Skeleton code for implementing a fraud detector.
 */
public class OverLimitDetector extends KeyedProcessFunction<String, CardTransaction, Alert> {

	private static final long serialVersionUID = 2L;

	@Override
	public void processElement(CardTransaction transaction, KeyedProcessFunction<String, CardTransaction, Alert>.Context context, Collector<Alert> collector) {
		if (transaction.getAmount() > transaction.getLimitLeft()) {
			Alert alert = new Alert();
			alert.setReason("Transaction amount is over the limit");
			alert.setTransaction(transaction);
			collector.collect(alert);
		}
	}
}
