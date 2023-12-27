from aio_pika import Channel

from app.common.amqp import amqp_conn_manager
from app.stable_diffusion import StableDiffusionGenerator


class AppRuntime:
    def __init__(self):
        self._ch_gen = None
        self._ch = None

        self._sd_generator = None

    @property
    def channel(self) -> Channel:
        assert self._ch
        return self._ch  # type: ignore

    @property
    def sd_generator(self) -> StableDiffusionGenerator:
        assert self._sd_generator
        return self._sd_generator

    async def init(self, rabbitmq_uri: str | None, modelpath: str):
        if rabbitmq_uri is not None:
            await amqp_conn_manager.connect(rabbitmq_uri)

            self._ch_gen = self._channels()
            self._ch = await anext(self._ch_gen)

        self._sd_generator = StableDiffusionGenerator(model_name=modelpath)

    async def on_down(self):
        await amqp_conn_manager.disconnect()

    @staticmethod
    async def _channels():
        assert amqp_conn_manager._conn and not amqp_conn_manager._conn.is_closed

        async with amqp_conn_manager._conn:
            yield await amqp_conn_manager._conn.channel()


app_runtime = AppRuntime()
