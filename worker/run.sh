#!/bin/sh
set -e

# Downloading model
if [ -d "$MODELPATH" ]; then
  echo "Model files already exists"
else
  mkdir -p /tts_model
  echo "Downloading model"

  git lfs install

  git clone https://huggingface.co/jzli/epiCPhotoGasm-last-unicorn "$MODELPATH"
  echo "Downloaded model"
fi

# Waiting for connection ability to rabbit
export AMQP_URI="amqp://${RABBITMQ_USER}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}/"
python /startup/rabbitmq.py -u "$AMQP_URI"

# Waiting for connection ability to s3
python /startup/s3.py

# Running app
gunicorn app.main:app \
  -w "${WORKERS:-4}" \
  -k uvicorn.workers.UvicornWorker \
  -b 0.0.0.0:5000 \
  --timeout 600

