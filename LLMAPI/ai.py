from langchain.vectorstores import FAISS
from langchain.embeddings import HuggingFaceEmbeddings
from langchain.text_splitter import CharacterTextSplitter

# Load the text file
with open("C:\\Users\\hsanj\\OneDrive\\Desktop\\Go_backend\\Go idk\\output.txt", "r") as file:
    data = file.read()

# Split the text into smaller chunks
text_splitter = CharacterTextSplitter(chunk_size=500, chunk_overlap=50)
chunks = text_splitter.split_text(data)

# Create embeddings for chunks
embedding_model = HuggingFaceEmbeddings(model_name="TinyLlama/TinyLlama-1.1B-Chat-v1.0")
vector_store = FAISS.from_texts(chunks, embedding_model)