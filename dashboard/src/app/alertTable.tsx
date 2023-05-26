import { useEffect, useState } from 'react';
import { Alert } from '../lib/alert';
import { Kafka } from 'kafkajs';

const kafka = new Kafka({
    clientId: 'dashboard',
    brokers: ['kafka:9092']
})

const consumer = kafka.consumer({ groupId: 'dashboard' })

export default function AlertTable() {
  const [alerts, setAlerts] = useState<Alert[]>([]);

  useEffect(() => {
    async function fetchData() {
      await consumer.connect();
      await consumer.subscribe({ topic: 'alerts', fromBeginning: true });

      await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
          if (!message.value) return;
          const alert = JSON.parse(message.value.toString());
          setAlerts((prevAlerts) => [...prevAlerts, alert]);
        },
      });
    }

    fetchData();
  }, []);

  return (
    <div>
      <table>
        <thead>
          <tr>
            <th>Transaction ID</th>
            <th>Amount</th>
            <th>Reason</th>
            <th>Owner</th>
          </tr>
        </thead>
        <tbody>
          {alerts.map((alert) => (
            <tr key={alert.transaction.utc.toUTCString()}>
              <td>{alert.transaction.utc.toUTCString()}</td>
              <td>{alert.transaction.amount}</td>
              <td>{alert.reason}</td>
              <td>{alert.transaction.owner.firstName + alert.transaction.owner.lastName}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

