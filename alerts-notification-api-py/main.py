# Import KafkaConsumer from Kafka library
import time
from kafka import KafkaConsumer

# Import sys module
import sys

if __name__ == "__main__":
    print("Starting consumer")

    # Define server with port
    bootstrap_servers = 'localhost:9092'
    print("Server: " + bootstrap_servers)

    # Define topic name from where the message will recieve
    topicName = 'alerts'
    print("Topic: " + topicName)

    # Initialize consumer variable
    consumer = KafkaConsumer (
        topicName,
        group_id ='group1',
        bootstrap_servers = bootstrap_servers,
        auto_offset_reset = 'latest'
    )

    print("Consumer: " + str(consumer))

    # Read and print message from consumer
    for msg in consumer:
        print("Topic Name=%s,Message=%s"%(msg.topic,msg.value))

    # Terminate the script
    sys.exit()
