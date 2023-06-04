import asyncio
import logging
from aiokafka import AIOKafkaConsumer
from websockets import serve, broadcast, ConnectionClosed

# Kafka consumer configuration
bootstrap_servers = 'kafkac:9092'
topic_name = 'alerts'

# Websocket server configuration
websocket_host = '0.0.0.0'
websocket_port = 12000

# Set of connected WebSocket clients
clients = set()

async def kafka_listener():
    logging.info('Starting Kafka consumer')
    consumer = AIOKafkaConsumer(
        topic_name,
        bootstrap_servers=bootstrap_servers,
        group_id='alerts-notification-api',
    )

    await consumer.start()
    try:
        async for msg in consumer:
            logging.info('Received message from Kafka %s', msg.value)
            # Broadcast message to all connected WebSocket clients
            broadcast(clients, msg.value.decode())
    finally:
        await consumer.stop()


async def websocket_handler(websocket):
    logging.info('Client connected')
    clients.add(websocket)
    try:
        await websocket.wait_closed()
    except ConnectionClosed:
        pass
    finally:
        logging.info('Client disconnected')
        clients.remove(websocket)


async def main():
    websocket_task = serve(websocket_handler, websocket_host, websocket_port)
    kafka_task = asyncio.create_task(kafka_listener())

    await asyncio.gather(websocket_task, kafka_task)

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    logging.info('Starting alerts-notification-api')
    asyncio.run(main())
