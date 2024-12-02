from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from rag import db, llama, run
import uvicorn

# Initialize FastAPI app
app = FastAPI()

# Pydantic model for API input
class QueryRequest(BaseModel):
    question: str

@app.post("/generate_content/")
async def generate_content(request: QueryRequest):
    try:
        # Email generation template
        template = """You are an AI email content generator. Generate a professional email based on the following context.
        
        The response MUST be in this exact format:
        Subject: [Your generated subject line]
        Body: [Your generated email body]
        
        Keep the email professional, clear, and concise.
        Include a proper greeting and sign-off in the body.
        Do not include any content outside of what's provided in the context.
       
        """
        
        # Format the template with the user's question
        formatted_template = template.format(context=request.question)
        
        # Generate the response
        response = run(formatted_template, request.question)

        # Validate response format
        if not ("Subject:" in response and "Body:" in response):
            # If response doesn't have the correct format, restructure it
            response = f"Subject: Re: {request.question}\nBody: {response}"

        return {"status": "success", "response": response}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/generate_post/")
async def generate_post(request: QueryRequest, prompt=None, template=None):
    try:
        # template = request.question
        # prompt = "Generate a social media post based on the following context."
        # # Template for generating posts
        # template = """You are an AI content generator tasked with creating engaging, concise, and professional social media posts.

        # Use only the information provided in the context.
        # Do not add any extra details, assumptions, or speculative content.
        # Maintain a tone suitable for [platform: e.g., LinkedIn, Twitter, Instagram, etc.].
        # The post should be clear, concise, and adhere to any specified character limits or formatting guidelines.
        # If the content is technical or professional, ensure the language is precise and jargon-free (if possible). For creative or casual posts, keep the tone friendly and approachable.

        # Deliverables:
        # Post text: A short, engaging post based on the context.
        # Hashtags (if required): Relevant hashtags derived from the context.
        # Call-to-action (optional): If appropriate, include a call-to-action to increase engagement.
        # Note:
        # Do not generate any content outside the context provided. If the context is insufficient, indicate that more information is required."""

        response=run(template,prompt)

        return {"status": "success", "response": response}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

