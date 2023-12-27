import asyncio
import uuid
from pathlib import Path
import traceback
from contextlib import asynccontextmanager

from fastapi import FastAPI
from loguru import logger
from pydantic import BaseModel

from app.common.amqp.consumer import listen_queue
from app.common.app_runtime import app_runtime
from app.common.config import config
from app.common.s3 import get_s3
from app.common.schemas import SDRequest


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
        channel, config.rabbitmq_queue, SDRequest  # type: ignore
    ):
        try:
            image = sd_generator.txt2img(msg)

            img_id = uuid.uuid4()
            path = Path(config.s3_outputs_path) / str(img_id).replace("-", "")
            path = path.with_suffix(".png")
            await s3.put_object(path, image)

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
