import time
import json
import random
import asyncio
import aiohttp


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
    ["🚗", "🚕", "🚙", "🚌", "🚎", "🏎", "🚓", "🚑"]
]

car_emojis = ["🚗", "🚕", "🚙"]
emojis_flat = [emoji for sublist in emojis for emoji in sublist]

base_url = "http://localhost:8080"
endpoint = "/reaction"
full_url = base_url + endpoint
headers = {'Content-Type': 'application/json'}

NUM_SIMULATED_CLIENTS = 10_000
NUM_EMOJI_REACTIONS = 4

async def simulate_client(session, user_id: int, emoji_list: list) -> None:
    for _ in range(NUM_EMOJI_REACTIONS):
        emoji = random.choice(emoji_list)
        timestamp = int(time.time())
        payload = {"user_id": user_id, "emoji": emoji, "timestamp": timestamp}
        async with session.post(full_url, json=payload) as response:
            print(f'Client {user_id} selected emoji: {emoji}')
            print(f"Response code: {response.status}, message: await {response.text()}")



async def main():
    async with aiohttp.ClientSession() as session:
        tasks = [asyncio.create_task(simulate_client(session, i, emojis_flat)) for i in range(NUM_SIMULATED_CLIENTS)]
        await asyncio.gather(*tasks)

asyncio.run(main())

