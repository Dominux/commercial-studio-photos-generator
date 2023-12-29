import gc
import io
import random

import numpy as np
from diffusers import StableDiffusionPipeline, DPMSolverMultistepScheduler  # type: ignore
import torch

from app.common.schemas import SDRequest
from app.common.config import config


device = "cuda" if torch.cuda.is_available() else "cpu"


def clear_memory():
    gc.collect()
    torch.cuda.empty_cache()


class StableDiffusionGenerator:
    def __init__(self, model_name: str) -> None:
        self._model = self._init_model(model_name)

    def _init_model(self, model_name):
        model = StableDiffusionPipeline.from_pretrained(
            model_name, 
            torch_dtype=torch.float16, 
            use_safetensors=True, 
            safety_checker=None,
        )
        model = model.to(device)
        model.enable_freeu(s1=0.9, s2=0.2, b1=1.2, b2=1.4)

        # setting scheduler (sampling method/sampler)
        scheduler = DPMSolverMultistepScheduler.from_config(model.scheduler.config)
        scheduler.use_karras_sigmas = True
        model.scheduler = scheduler

        # setting embeddings
        for embedding in ("EasyNegative",):
            path = f"/cspg_model/embeddings/{embedding}/{embedding}.safetensors"
            model.load_textual_inversion(path)

        clear_memory()
        return model

    def txt2img(self, schema: SDRequest) -> bytes:
        seed = (
            schema.seed
            if schema.seed is not None
            else random.randint(0, np.iinfo(np.int32).max)
        )
        torch.manual_seed(seed)
        
        prompt = config.sd_prompt_mask.format(schema.product)

        with torch.no_grad():
            result = self._model(
                prompt=prompt,
                negative_prompt=config.sd_negative_prompt,
                width=schema.width,  # type: ignore
                height=schema.height,  # type: ignore
                guidance_scale=config.sd_cfg,
                num_inference_steps=schema.steps,
                num_images_per_prompt=schema.num_images,
                output_type="pil",
            ).images  # type: ignore
        clear_memory()

        buf = io.BytesIO()
        result[0].save(buf, format="PNG")
        return buf.getvalue()
