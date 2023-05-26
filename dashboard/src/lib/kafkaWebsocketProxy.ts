import WebSocket from 'ws';
import { Kafka } from 'kafkajs';

const kafka = new Kafka({
  clientId: 'my-app',
  brokers: ['localhost:9092'],
});

const consumer = kafka.consumer({ groupId: 'alert-group' });

async function consumeAlerts() {
  await consumer.connect();
  await consumer.subscribe({ topic: 'alerts', fromBeginning: true });

  await consumer.run({
    eachMessage: async ({ message }) => {
      if (message.value !== null && message.value !== undefined) {
        const alert = JSON.parse(message.value.toString());
        // send the alert to all connected clients
        wss.clients.forEach((client) => {
          if (client.readyState === WebSocket.OPEN) {
            client.send(JSON.stringify(alert));
          }
        });
      }
    },
  });
}

const wss = new WebSocket.Server({ port: 3000 });

wss.on('connection', (ws) => {
  console.log('Client connected');
});

consumeAlerts();