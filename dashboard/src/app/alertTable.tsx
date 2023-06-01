'use client';

import { useEffect, useState } from 'react';
import { Alert } from '../lib/alert';

export default function AlertTable() {
  const [alerts, setAlerts] = useState<Alert[]>([]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:12000/');

    console.log('Connecting to websocket');

    ws.onopen = () => {
      console.log('Connected to websocket');
    };

    ws.onmessage = (event) => {
      console.log('Received message', event.data);
      const alert = JSON.parse(event.data);
      setAlerts((alerts) => [...alerts, alert]);
    };

    return () => {
      console.log('Closing websocket');
      ws.close();
    };
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
            <tr key={alert.transaction.card.cardNumber}>
              <td>{alert.reason}</td>
            </tr>
            // <tr key={alert.transaction.utc.toUTCString()}>
            //   <td>{alert.transaction.utc.toUTCString()}</td>
            //   <td>{alert.transaction.amount}</td>
            //   <td>{alert.reason}</td>
            //   <td>{alert.transaction.owner.firstName + alert.transaction.owner.lastName}</td>
            // </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

