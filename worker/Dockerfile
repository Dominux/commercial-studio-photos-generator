FROM python:3.11-slim as base

# Building stage
FROM base AS builder

ENV PYTHONFAULTHANDLER=1 \
  PYTHONUNBUFFERED=1 \
  PYTHONHASHSEED=random \
  PIP_NO_CACHE_DIR=off \
  PIP_DISABLE_PIP_VERSION_CHECK=on \
  PIP_DEFAULT_TIMEOUT=100 \
  POETRY_NO_INTERACTION=1 \
  POETRY_VIRTUALENVS_CREATE=false \
  PATH="$PATH:/runtime/bin" \
  PYTHONPATH="$PYTHONPATH:/runtime/lib/python3.11/site-packages" \
  POETRY_VERSION=1.4.2

RUN apt update && apt install -y build-essential unzip wget
RUN pip install "poetry==$POETRY_VERSION"

LABEL com.nvidia.volumes.needed="nvidia_driver"
ENV CUDA_LAUNCH_BLOCKING 1
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES compute,utility

WORKDIR /src
COPY pyproject.toml /src/
RUN poetry export --without-hashes --no-interaction --no-ansi -f requirements.txt -o requirements.txt
RUN pip install --prefix=/runtime --force-reinstall -r requirements.txt
RUN pip install --prefix=/runtime xformers==0.0.22.post7 --index-url https://download.pytorch.org/whl/cu118
RUN pip install --prefix=/runtime triton==2.1.0 --index-url https://download.pytorch.org/whl/cu118
RUN pip install --prefix=/runtime --force-reinstall torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
RUN pip install --prefix=/runtime typing-extensions==4.8.0

# Runtime stage
FROM base AS runtime
COPY --from=builder /runtime /usr/local
RUN apt update && apt install -y curl git-lfs

WORKDIR /
COPY ./app ./app
COPY ./startup ./startup
COPY ./run.sh ./
RUN chmod +x run.sh

CMD ["sh", "./run.sh"]
