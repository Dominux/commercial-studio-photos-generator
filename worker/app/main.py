import asyncio
import uuid
from pathlib import Path
import traceback
from contextlib import asynccontextmanager

from fastapi import FastAPI
from loguru import logger
from pydantic import BaseModel
from aio_pika import Message

from app.common.amqp.consumer import listen_queue
from app.common.app_runtime import app_runtime
from app.common.config import config
from app.common.s3 import get_s3
from app.common.schemas import SDRequest, SDResponse


@asynccontextmanager
async def lifespan(app: FastAPI):
    await app_runtime.init(config.rabbitmq_uri, config.modelpath)
    asyncio.create_task(generate_for_queue())

    yield

    await app_runtime.on_down()


async def generate_for_queue():
    s3_gen = get_s3()
    s3 = await anext(s3_gen)

    channel = app_runtime.channel
    sd_generator = app_runtime.sd_generator

    async for msg in listen_queue(
        channel, config.rabbitmq_request_queue, SDRequest  # type: ignore
    ):
        try:
            image = sd_generator.txt2img(msg)

            path = Path(config.s3_outputs_path) / str(msg.id).replace("-", "")
            path = path.with_suffix(".png")
            await s3.put_object(path, image)

            msg_back = SDResponse(id=msg.id, status="done")
            msg_back = Message(body=msg_back.model_dump_json().encode())
            await channel.default_exchange.publish(
                msg_back, routing_key=config.rabbitmq_response_queue
            )

        except Exception:
            logger.error(traceback.format_exc())


app = FastAPI(title="Commercial studio photots generator", docs_url=None, lifespan=lifespan)

class HealthCheckSchema(BaseModel):
    msg: str


@app.get("/healthcheck", response_model=HealthCheckSchema)
async def healthcheck():
    """
    Healthcheck endpoint.

    ### Output
    * {"msg": "ok"}
    """
    return {"msg": "ok"}
