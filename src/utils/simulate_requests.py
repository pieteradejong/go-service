import time
import random
import asyncio
import aiohttp
import logging

emojis = [
    ["😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣"],
    ["🥰", "😍", "🤩", "😘", "😗", "😚", "😙"],
    ["😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓"],
    ["🤔", "🤭", "🤫", "🤥", "😶", "😐", "😑", "😬"],
    ["🙄", "😯", "😦", "😧", "😮", "🤯", "😲", "😳"],
    ["🥺", "😢", "😭", "😤", "😠", "😡", "🤬"],
    ["🤠", "😇", "🤗", "🤡", "🤥", "🤓", "😈", "👿"],
    ["👋", "🤚", "🖐", "✋", "🖖", "👌", "🤏"],
    ["🍏", "🍎", "🍐", "🍊", "🍋", "🍌", "🍉", "🍇"],
    ["🚗", "🚕", "🚙", "🚌", "🚎", "🏎", "🚓", "🚑"],
]

emotional_states = {
    "Happy/Positive": ["😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣"],
    "Love/Affection": ["😘", "😗", "😚", "😙", "🥰", "😍", "🤩"],
    "Excitement/Thrill": ["🤪", "😲", "🤯", "😮"],
    "Thoughtful/Contemplative": ["🤔", "🤨", "🧐", "🤓"],
    "Disappointment/Sadness": ["😢", "😭", "🥺", "😦", "😧"],
    "Anger/Frustration": ["😠", "😡", "🤬"],
    "Mischievous/Playful": ["😈", "👿", "🤡", "🤠", "😇"],
}

emojis_flat = [emoji for sublist in emojis for emoji in sublist]

base_url = "http://localhost:8080"
endpoint = "/reaction"
full_url = base_url + endpoint
headers = {"Content-Type": "application/json"}

NUM_SIMULATED_CLIENTS = 10_000
NUM_EMOJI_REACTIONS = 4
CHUNK_SIZE = 1000


def init():
    logging.basicConfig(
        level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
    )


async def simulate_client(session, user_id: int, emotional_states: dict) -> None:
    mood = random.choice(list(emotional_states.keys()))
    emoji_list = emotional_states[mood]
    for _ in range(NUM_EMOJI_REACTIONS):
        emoji = random.choice(emoji_list)
        timestamp = int(time.time())
        payload = {"user_id": user_id, "emoji": emoji, "timestamp": timestamp}
        try:
            async with session.post(full_url, json=payload) as response:
                logging.info(f"Client {user_id} in {mood} mood selected emoji: {emoji}")
                logging.info(f"Response code: {response.status}")
        except aiohttp.ClientError as e:
            logging.info(f"Request failed for user {user_id}: {str(e)}")


async def main():
    async with aiohttp.ClientSession() as session:
        for chunk in range(0, NUM_SIMULATED_CLIENTS, CHUNK_SIZE):
            tasks = [
                asyncio.create_task(simulate_client(session, i, emotional_states))
                for i in range(chunk, min(chunk + 1000, NUM_SIMULATED_CLIENTS))
            ]
            await asyncio.gather(*tasks)


def run():
    init()
    asyncio.run(main())


if __name__ == "__main__":
    run()
