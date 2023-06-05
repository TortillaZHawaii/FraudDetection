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

import java.util.Calendar;
import java.util.Date;
import java.util.GregorianCalendar;

/**
 * Skeleton code for implementing a fraud detector.
 */
public class ExpiredCardDetector extends KeyedProcessFunction<String, CardTransaction, Alert> {

	private static final long serialVersionUID = 3L;

	@Override
	public void processElement(CardTransaction transaction, KeyedProcessFunction<String, CardTransaction, Alert>.Context context, Collector<Alert> collector) throws Exception {
		if (isExpired(transaction)) {
			Alert alert = new Alert();
			alert.setReason("Card is expired");
			alert.setTransaction(transaction);
			collector.collect(alert);
		}
	}

	private static Boolean isExpired(CardTransaction transaction)
	{
		final var card = transaction.getCard();
		final var expiryYear = card.getExpYear();
		final var expiryMonth = card.getExpMonth();
		// for some reason months in java start at 0
		// we need to find first expired date
		final var expiredDate = expiryMonth == 12 ? new GregorianCalendar(expiryYear + 1, Calendar.JANUARY, 1).getTime()
				: new GregorianCalendar(expiryYear, expiryMonth, 1).getTime();

		final var date = transaction.getUtc();
		return date.after(expiredDate);
	}
}
