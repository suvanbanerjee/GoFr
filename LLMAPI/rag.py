import os
from langchain.document_loaders import TextLoader
from langchain_community.embeddings import OllamaEmbeddings
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_community.vectorstores import Chroma
from langchain.prompts import ChatPromptTemplate, PromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_community.chat_models import ChatOllama
from langchain_core.runnables import RunnablePassthrough
from langchain.retrievers.multi_query import MultiQueryRetriever
def db():
    local_path = "GoFrServer/output.txt"

    # Local PDF file uploads
    if local_path:
        loader = TextLoader(local_path, encoding="latin1") 
        data = loader.load()

    # Split and chunk 
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=100)
    chunks = text_splitter.split_documents(data)
    # Split and chunk 
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=100)
    chunks = text_splitter.split_documents(data)

    current_dir = os.getcwd()
    persistent_directory = os.path.join(current_dir, "db", "chroma_db_for_GitHub")
    embedding_function = OllamaEmbeddings(model="nomic-embed-text", show_progress=True)

    if os.path.exists(persistent_directory):
        vector_db = Chroma(
            persist_directory=persistent_directory, 
            embedding_function=embedding_function,
            collection_name="local-rag"
        )
        print("Loaded existing Chroma vector store.")
    else:
        vector_db = Chroma.from_documents(
            documents=chunks, 
            embedding=OllamaEmbeddings(model="nomic-embed-text", show_progress=True),
            collection_name="local-rag",
            persist_directory=persistent_directory
        )
        vector_db.persist()
    return vector_db

def llama(vector_db,template):
    # LLM from Ollama
    local_model = "llama3.2"
    llm = ChatOllama(model=local_model)

    QUERY_PROMPT = PromptTemplate(
        input_variables=["question"],
        template=template,
    )

    retriever = MultiQueryRetriever.from_llm(
        vector_db.as_retriever(), 
        llm,
        prompt=QUERY_PROMPT
    )

    template = """ Answer the question based ONLY on the following context:
    {context}
    Question: {question}
    """

    prompt = ChatPromptTemplate.from_template(template)

    chain = (
        {"context": retriever, "question": RunnablePassthrough()}
        | prompt
        | llm
        | StrOutputParser()
    )
    return chain

template = """You are an AI content generator tasked with creating engaging, concise, and professional social media posts.

    Use only the information provided in the context.
    Do not add any extra details, assumptions, or speculative content.
    Maintain a tone suitable for [platform: e.g., LinkedIn, Twitter, Instagram, etc.].
    The post should be clear, concise, and adhere to any specified character limits or formatting guidelines.
    If the content is technical or professional, ensure the language is precise and jargon-free (if possible). For creative or casual posts, keep the tone friendly and approachable.

    Deliverables:
    Post text: A short, engaging post based on the context.
    Hashtags (if required): Relevant hashtags derived from the context.
    Call-to-action (optional): If appropriate, include a call-to-action to increase engagement.
    Note:
    Do not generate any content outside the context provided. If the context is insufficient, indicate that more information is required."""
def run(template,prompt):
    current_dir = os.getcwd()
    persistent_directory = os.path.join(current_dir, "db", "chroma_db_for_GitHub")
    embedding_function = OllamaEmbeddings(model="nomic-embed-text", show_progress=True)

    if os.path.exists(persistent_directory):
        vector_db = Chroma(
            persist_directory=persistent_directory, 
            embedding_function=embedding_function,
            collection_name="local-rag"
        )
    else:
        vector_db = db()

    chain=llama(vector_db,template)

    response = chain.invoke(prompt)

    return response
    
run(template,"Create a linkedin post on Circuit Breaker in HTTP Communication")