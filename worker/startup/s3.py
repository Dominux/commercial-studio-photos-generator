import os
import asyncio

import aioboto3
from aiobotocore.session import get_session
from loguru import logger


host = os.environ["S3_HOST"]
port = os.environ["S3_PORT"]
secret = os.environ["S3_SECRET_ACCESS_KEY"]
key = os.environ["S3_ACCESS_KEY_ID"]
BUCKET = os.environ["S3_BUCKET"]

URL = f"http://{host}:{port}"


async def wait_for_connection(interval: int):
    while True:
        try:
            session = get_session()
            async with session.create_client(
                "s3",
                endpoint_url=URL,
                aws_secret_access_key=secret,
                aws_access_key_id=key,
            ):
                logger.info("Established connection with s3")
                return
        except Exception:
            await asyncio.sleep(interval)
            logger.info("Attempt to establish connection with s3")


async def create_bucket():
    session = aioboto3.Session()
    async with session.resource(
        "s3",
        endpoint_url=URL,
        aws_secret_access_key=secret,
        aws_access_key_id=key,
    ) as s3:
        try:
            await s3.create_bucket(Bucket=BUCKET)
        except s3.meta.client.exceptions.BucketAlreadyOwnedByYou:
            logger.info("Bucket already exists")
        else:
            logger.info("Created bucket")


async def main(interval: int):
    await wait_for_connection(interval)
    await create_bucket()


if __name__ == "__main__":
    asyncio.run(main(3))
