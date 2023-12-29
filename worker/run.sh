#!/bin/sh
set -e

# Downloading model
if [ -d "$MODELPATH" ]; then
  echo "Model files already exist"
else
  echo "Downloading model"

  git lfs install

  git clone https://huggingface.co/jzli/epiCPhotoGasm-last-unicorn "$MODELPATH"

  rm $MODELPATH/*.ckpt

  echo "Downloaded model"
fi

# Downloading embeddings
if [ -d "/cspg_model/embeddings" ]; then
  echo "Embeddings already exist"
else
  echo "Downloading embeddings"

  git lfs install

  git clone https://huggingface.co/embed/EasyNegative "/cspg_model/embeddings/EasyNegative"
  echo "Downloaded embeddings"
fi

# Waiting for connection ability to rabbit
export AMQP_URI="amqp://${RABBITMQ_USER}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}/"
python /startup/rabbitmq.py -u "$AMQP_URI"

# Waiting for connection ability to s3
python /startup/s3.py

# Running app
gunicorn app.main:app \
  -w "${WORKERS:-1}" \
  -k uvicorn.workers.UvicornWorker \
  -b 0.0.0.0:5000 \
  --timeout 600

