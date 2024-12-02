from requests_oauthlib import OAuth1Session
import os
import json

# Set your consumer key and secret
consumer_key = os.environ.get("CONSUMER_KEY")
consumer_secret = os.environ.get("CONSUMER_SECRET")

# File to store access tokens
TOKEN_FILE = "twitter_tokens.json"
LAST_TWEET_FILE = "last_tweet_id.txt"

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
    request_token_url = "https://api.twitter.com/oauth/request_token?oauth_callback=oob&x_auth_access_type=read"
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

# Load the last seen tweet ID
def load_last_tweet_id():
    if os.path.exists(LAST_TWEET_FILE):
        with open(LAST_TWEET_FILE, "r") as file:
            return file.read().strip()
    return None

# Save the latest tweet ID
def save_last_tweet_id(tweet_id):
    with open(LAST_TWEET_FILE, "w") as file:
        file.write(tweet_id)

# Define hashtags to search for
hashtags = [
    "#GolangNews",
    "#GoLangUpdates",
    "#GoProgrammingNews",
]

# Combine hashtags into a search query
search_query = " OR ".join(hashtags) + " -is:retweet"  # Exclude retweets
url = "https://api.twitter.com/2/tweets/search/recent"

params = {
    "query": search_query,
    "max_results": 1,  # Adjust as needed
    "tweet.fields": "created_at,text,author_id",
}

last_tweet_id = load_last_tweet_id()
if last_tweet_id:
    params["since_id"] = last_tweet_id

oauth = OAuth1Session(
    consumer_key,
    client_secret=consumer_secret,
    resource_owner_key=access_token,
    resource_owner_secret=access_token_secret,
)

response = oauth.get(url, params=params)

if response.status_code != 200:
    raise Exception(
        "Request returned an error: {} {}".format(response.status_code, response.text)
    )

# Parse and display the response
tweets = response.json()

if "data" in tweets:
    print(json.dumps(tweets, indent=4, sort_keys=True))
    latest_tweet_id = tweets["data"][0]["id"]  # Save the most recent tweet ID
    save_last_tweet_id(latest_tweet_id)
else:
    print("No new tweets found.")
