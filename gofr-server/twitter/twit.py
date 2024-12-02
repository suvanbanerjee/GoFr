from requests_oauthlib import OAuth1Session
import os
import json
import sys
import requests

# Set your consumer key and secret
consumer_key = os.environ.get("CONSUMER_KEY")
consumer_secret = os.environ.get("CONSUMER_SECRET")

# File to store access tokens
TOKEN_FILE = "twitter_tokens.json"

# Function to save tokens to a file
def save_tokens(access_token, access_token_secret):
    with open(TOKEN_FILE, "w") as file:
        json.dump({"access_token": access_token, "access_token_secret": access_token_secret}, file)

# Function to load tokens from a file
def load_tokens():
    if os.path.exists(TOKEN_FILE):
        with open(TOKEN_FILE, "r") as file:
            tokens = json.load(file)
            return tokens["access_token"], tokens["access_token_secret"]
    return None, None

# Function to perform manual PIN-based flow
def perform_oauth_flow():
    # Get request token
    request_token_url = "https://api.twitter.com/oauth/request_token?oauth_callback=oob&x_auth_access_type=write"
    oauth = OAuth1Session(consumer_key, client_secret=consumer_secret)

    try:
        fetch_response = oauth.fetch_request_token(request_token_url)
    except ValueError:
        print("There may have been an issue with the consumer_key or consumer_secret you entered.")
        return None, None

    resource_owner_key = fetch_response.get("oauth_token")
    resource_owner_secret = fetch_response.get("oauth_token_secret")
    print("Got OAuth token: %s" % resource_owner_key)

    # Get authorization
    base_authorization_url = "https://api.twitter.com/oauth/authorize"
    authorization_url = oauth.authorization_url(base_authorization_url)
    print("Please go here and authorize: %s" % authorization_url)
    verifier = input("Paste the PIN here: ")

    # Get the access token
    access_token_url = "https://api.twitter.com/oauth/access_token"
    oauth = OAuth1Session(
        consumer_key,
        client_secret=consumer_secret,
        resource_owner_key=resource_owner_key,
        resource_owner_secret=resource_owner_secret,
        verifier=verifier,
    )
    oauth_tokens = oauth.fetch_access_token(access_token_url)
    return oauth_tokens["oauth_token"], oauth_tokens["oauth_token_secret"]

# Load saved tokens or perform OAuth flow
access_token, access_token_secret = load_tokens()
if not access_token or not access_token_secret:
    access_token, access_token_secret = perform_oauth_flow()
    if access_token and access_token_secret:
        save_tokens(access_token, access_token_secret)
    else:
        print("OAuth flow failed.")
        exit(1)


# Check if text is passed as an argument
if len(sys.argv) < 2:
    print("Usage: python twit.py <text>")
    exit(1)

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
prompt = sys.argv[1]

api_url = "https://jsonplaceholder.typicode.com/posts"
response = requests.post(api_url, json={"prompt": prompt, "template": template})

if response.status_code == 201:
    post = response.json()
    tweet_text = post["response"]
else:
    exit(1)


oauth = OAuth1Session(
    consumer_key,
    client_secret=consumer_secret,
    resource_owner_key=access_token,
    resource_owner_secret=access_token_secret,
)

payload = {"text": tweet_text}

# Making the request
response = oauth.post(
    "https://api.twitter.com/2/tweets",
    json=payload,
)

if response.status_code != 201:
    raise Exception(
        "Request returned an error: {} {}".format(response.status_code, response.text)
    )

print("Response code: {}".format(response.status_code))

# Saving the response as JSON
json_response = response.json()
print(json.dumps(json_response, indent=4, sort_keys=True))