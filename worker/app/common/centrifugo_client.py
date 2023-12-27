import httpx


class CentrifugoClient:
    def __init__(self, url: str, api_key: str) -> None:
        self._url = url
        self._headers = {"Authorization": api_key}

    async def publish(self, channel: str, data: dict) -> None:
        json = {"channel": channel, "data": data}
        async with httpx.AsyncClient() as client:
            await client.post(self._url, json=str(json), headers=self._headers)
