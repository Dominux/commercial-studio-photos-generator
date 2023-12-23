from pathlib import Path

from pydantic import BaseModel


class SDImagePath(BaseModel):
    path: Path
