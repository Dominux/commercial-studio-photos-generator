import argparse
import asyncio

import aio_pika
from aiormq import AMQPConnectionError
from loguru import logger


async def wait_for_connection(uri: str, interval: int):
    while True:
        try:
            await aio_pika.connect_robust(uri)
        except AMQPConnectionError:
            await asyncio.sleep(interval)
            logger.info("Attempt to establish connection with rabbitmq")
        else:
            logger.info("Established connection with rabbitmq")
            return


parser = argparse.ArgumentParser()
parser.add_argument("-u", "--uri", required=True)
parser.add_argument("-i", "--interval", default=3)


if __name__ == "__main__":
    args = parser.parse_args()
    awaitable = wait_for_connection(args.uri, args.interval)
    asyncio.run(awaitable)
