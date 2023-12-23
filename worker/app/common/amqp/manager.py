import asyncio

import aio_pika
from aiormq import AMQPConnectionError
from loguru import logger


class AMQPConnectionManager:
    def __init__(self) -> None:
        self._conn = None

    async def connect(self, rabbitmq_uri: str):
        while self._conn is None:
            try:
                self._conn = await aio_pika.connect_robust(rabbitmq_uri)
            except AMQPConnectionError:
                logger.info("Trying to establish connection with amqp")
                await asyncio.sleep(3)

    async def disconnect(self):
        if self._conn and not self._conn.is_closed:
            await self._conn.close()
            logger.info("Disconnected from AMQP")


amqp_conn_manager = AMQPConnectionManager()
