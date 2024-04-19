from datetime import datetime
import json
import os
import random
import uuid
from flask import Flask
import redis 
from confluent_kafka import Producer, Consumer, KafkaError


def get_coordinates():
    lat, lon = generate_random_coordinates()
    return {
        'url': os.getenv("WEBHOOK_ADDRESS", ""),
        'webhookId': uuid.uuid4().hex,
        'data': {
            'id': uuid.uuid4().hex,
            'lat': lat,
            'long': lon,
            'status': random.choice(["inside", "outside", "canceled"]),
            'created': datetime.now().strftime("%d/%m/%Y, %H:%M:%S"),
        }
    }

app = Flask(__name__)

redis_address = os.getenv("REDIS_ADDRESS", "localhost:6379")  

host, port = redis_address.split(":")  
port = int(port)  
# Create a connection to the Redis server  
redis_connection = redis.StrictRedis(host=host, port=port) 

topic_name = os.getenv("KAFKA_TOPIC", "geofence")

bootstrap_servers = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "172.17.112.242:9092")
producer = Producer({'bootstrap.servers': bootstrap_servers})

def generate_random_coordinates():
    lat = random.uniform(-90, 90)
    lon = random.uniform(-180, 180)
    return lat, lon



def delivery_report(err, msg):
    if err is not None:
        print(f'Message delivery failed: {err}')
    else:
        print(f'Message delivered to {msg.topic()} [{msg.partition()}]')

@app.route('/coordinates')
def coordinates():
    coordinates_data = get_coordinates()
    webhook_payload_json = json.dumps(coordinates_data)

   
    producer.produce(topic_name, key=str(uuid.uuid4()), value=webhook_payload_json, callback=delivery_report)
    producer.flush() 
    
    redis_connection.publish('coordinates', webhook_payload_json) 

    return webhook_payload_json


consumer_config = {
    'bootstrap.servers': bootstrap_servers,
    'group.id': 'webhook-group',
    'auto.offset.reset': 'earliest',
    'enable.auto.commit': True,
}

def consume_from_kafka():
    consumer = Consumer(consumer_config)
    consumer.subscribe([topic_name])

    try:
        while True:
            msg = consumer.poll(1.0)

            if msg is None:
                continue
            if msg.error():
                if msg.error().code() == KafkaError._PARTITION_EOF:
                    continue
                else:
                    print(msg.error())
                    break

            print(f"Received message: {msg.value().decode('utf-8')}")

    finally:
        consumer.close()

if __name__ == '__main__':
    # Start consuming messages from Kafka in a separate thread
    import threading
    threading.Thread(target=consume_from_kafka, daemon=True).start()

    app.run(host='0.0.0.0', port=8000)
