import logging
import sys

from loguru import logger as LOGGER
from pydantic_settings import BaseSettings


class Config(BaseSettings):
    rabbitmq_host: str
    rabbitmq_port: int
    rabbitmq_user: str
    rabbitmq_password: str
    rabbitmq_request_queue: str
    rabbitmq_response_queue: str

    @property
    def rabbitmq_uri(self) -> str:
        return "amqp://{}:{}@{}:{}/".format(
            self.rabbitmq_user,
            self.rabbitmq_password,
            self.rabbitmq_host,
            self.rabbitmq_port,
        )

    s3_host: str
    s3_port: int
    s3_access_key_id: str
    s3_secret_access_key: str
    s3_bucket: str
    s3_outputs_path: str

    @property
    def s3_uri(self) -> str:
        return "http://{}:{}".format(
            self.s3_host,
            self.s3_port,
        )

    modelpath: str
    sd_prompt_mask: str
    sd_negative_prompt: str
    sd_cfg: float

    LOG_LEVEL: str = "INFO"
    LOG_FORMAT: str = (
        "<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level: <8}</level> "
        "| <level>{message}</level>"
    )


config = Config()  # type: ignore


logging.basicConfig(level=logging.INFO)


class InterceptHandler(logging.Handler):
    loglevel_mapping = {
        50: "CRITICAL",
        40: "ERROR",
        30: "WARNING",
        20: "INFO",
        10: "DEBUG",
        0: "NOTSET",
    }

    def emit(self, record: logging.LogRecord):
        try:
            level = LOGGER.level(record.levelname).name
        except (AttributeError, ValueError):
            level = str(record.levelno)

        frame, depth = logging.currentframe(), 2
        while frame.f_code.co_filename == logging.__file__:
            if frame.f_back:
                frame = frame.f_back
            depth += 1

        LOGGER.opt(depth=depth, exception=record.exc_info).log(
            level, record.getMessage()
        )


class CustomLogger:
    @classmethod
    def make_logger(cls, level, logs_format):
        _logger = cls.customize_logging(level=level, logs_format=logs_format)
        return _logger

    @classmethod
    def customize_logging(cls, level: str, logs_format: str):
        intercept_handler = InterceptHandler()

        LOGGER.remove()
        LOGGER.add(
            sink=sys.stdout,
            enqueue=True,
            backtrace=True,
            level=level.upper(),
            format=logs_format,
        )

        lognames = [
            "asyncio",
            "aio_pika",
            "fastapi",
            "uvicorn",
            "uvicorn.access",
            "uvicorn.error",
        ]

        for _log in lognames:
            _logger = logging.getLogger(_log)
            _logger.handlers = [intercept_handler]
            _logger.propagate = False

        return LOGGER.bind(request_id=None, method=None)


CustomLogger.make_logger(level=config.LOG_LEVEL, logs_format=config.LOG_FORMAT)
