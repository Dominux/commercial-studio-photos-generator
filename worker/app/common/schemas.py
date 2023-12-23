from pydantic import BaseModel


class SDRequest(BaseModel):
    product: str | None = "a bulgarian pepper"
    seed: int | None = None
    width: int | None = 512
    height: int | None = 512
    steps: int | None = 4
    num_images: int | None = 1
