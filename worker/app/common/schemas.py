from uuid import UUID

from pydantic import BaseModel


class SDRequest(BaseModel):
    id: UUID
    product: str | None = "a bulgarian pepper"
    seed: int | None = None
    width: int | None = 512
    height: int | None = 512
    steps: int | None = 20
    num_images: int | None = 1


class SDResponse(BaseModel):
    id: UUID
    status: str
