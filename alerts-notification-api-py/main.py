import asyncio
import json
import logging
from kafka import KafkaConsumer
from websockets import serve, ConnectionClosed

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
    consumer = KafkaConsumer(
        topic_name,
        bootstrap_servers=bootstrap_servers,
        api_version=(3, 3, 0),
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )
    
    logging.info('Listening for Kafka messages')
    for message in consumer:
        logging.info('Received Kafka message: %s', message.value)
        alert_message = message.value
        logging.info('Broadcasting alert message to connected clients')
        await broadcast_message(alert_message)
    logging.info('Kafka consumer stopped')

async def broadcast_message(message):
    if not clients:
        logging.info('No connected clients')
        return
    
    serialized_message = json.dumps(message)
    for client in clients:
        try:
            logging.info('Sending message to client')
            await client.send(serialized_message)
            logging.info('Message sent to client')
        except ConnectionClosed:
            logging.info('Client disconnected')
            clients.remove(client)

async def websocket_handler(websocket):
    logging.info('Client connected')
    clients.add(websocket)
    try:
        async for message in websocket:
            # Do something with the received message, if required
            pass
    except ConnectionClosed:
        pass
    finally:
        logging.info('Client disconnected')
        clients.remove(websocket)

async def main():
    kafka_task = asyncio.create_task(kafka_listener())
    websocket_task = serve(websocket_handler, websocket_host, websocket_port)
    
    await asyncio.gather(kafka_task, websocket_task)

async def echo(websocket):
    async for message in websocket:
        await websocket.send(message)

async def test():
    kafka_task = asyncio.create_task(kafka_listener())

    async with serve(echo, websocket_host, websocket_port):
        await kafka_task  # run forever

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    logging.info('Starting alerts-notification-api')
    # asyncio.run(main())
    asyncio.run(test())
