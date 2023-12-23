from typing import TypeVar

from aio_pika import RobustChannel
from loguru import logger
from pydantic import BaseModel


Model = TypeVar("Model", bound=BaseModel)


async def listen_queue(channel: RobustChannel, queue_name: str, model: type[Model]):
    # Will take no more than 1 messages in advance
    await channel.set_qos(prefetch_count=1)

    # Declaring queue
    queue = await channel.declare_queue(queue_name, auto_delete=True)

    async with queue.iterator() as queue_iter:
        async for message in queue_iter:
            try:
                async with message.process():
                    msg_body = message.body.decode()
                    logger.debug(f'received msg "{msg_body}" from queue "{queue_name}"')
                    yield model.model_validate_json(msg_body)

            except Exception as e:
                logger.error(e)
