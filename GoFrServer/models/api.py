from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from models.rag import db, llama, run
import logging
import sys

# Initialize logger
logger = logging.getLogger('uvicorn.error')
logger.setLevel(logging.DEBUG)

# Initialize FastAPI app
app = FastAPI()

# Pydantic model for API input
class QueryRequest(BaseModel):
    Context: str

@app.post("/generate_post/")
async def generate_post(request: QueryRequest):
    try:
       
        # # Template for generating posts
        template = """You are an AI generator tasked with creating engaging, concise, and professional social media posts.

        Use only the information provided in the context.
        Do not add any extra details, assumptions, or speculative content.
        Maintain a tone suitable for [platform: e.g., LinkedIn, Twitter, Instagram, etc.].
        The post should be clear, concise, and adhere to any specified character limits or formatting guidelines.
        If the content is technical or professional, ensure the language is precise and jargon-free (if possible). For creative or casual posts, keep the tone friendly and approachable.
    """
        
        res=run(template,request.Context)
        logger.debug(res)
        return {"status": "success", "response": res}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
