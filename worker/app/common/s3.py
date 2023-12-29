from pathlib import Path

from aiobotocore.client import AioBaseClient
from aiobotocore.session import get_session

from app.common.config import config


async def get_s3():
    session = get_session()
    async with session.create_client(
        "s3",
        endpoint_url=config.s3_uri,
        aws_secret_access_key=config.s3_secret_access_key,
        aws_access_key_id=config.s3_access_key_id,
    ) as client:
        yield S3(client)


class S3:
    def __init__(self, client: AioBaseClient) -> None:
        self._client = client

    async def put_object(self, path: Path, data: bytes):
        await self._client.put_object(Bucket=config.s3_bucket, Key=str(path), Body=data)  # type: ignore
