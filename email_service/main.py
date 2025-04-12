import constants
import redis
import logging
import sender

redis_client = redis.StrictRedis(
    host=constants.REDIS_HOST,
    port=constants.REDIS_PORT,
    password=constants.REDIS_PASSWORD,
    decode_responses=True
)

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def process_message(topic, message):
    data = eval(message)
    if topic == constants.EMAIL_CODE_TOPIC:
        sender.send_code(data["email"], data["code"])
    elif topic == constants.PASSWORD_TOPIC:
        sender.send_password(data["email"], data["password"])
    else:
        logger.warning(f"Unknown topic: {topic}")


def process_redis_messages():
    logger.info("Listening for messages...")
    pubsub = redis_client.pubsub()
    pubsub.subscribe(constants.REDIS_TOPICS)

    for message in pubsub.listen():
        if message['type'] == 'message':
            topic = message['channel']
            message_data = message['data']
            logger.info(f"Received message from {topic}: {message_data}")
            process_message(topic, message_data)


if __name__ == "__main__":
    process_redis_messages()
