import threading
import time
import requests
import json
import random


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

NUM_EMOJI_REACTIONS = 1
def simulate_client(user_id: int, emoji_list: list) -> None:
    for _ in range(NUM_EMOJI_REACTIONS):
        emoji = random.choice(emoji_list)
        timestamp = int(time.time())
        payload = {"user_id": user_id, "emoji": emoji, "timestamp": timestamp}
        print(f'Client {user_id} selected emoji: {emoji}')
        response = requests.post(full_url, headers=headers, data=json.dumps(payload))
        print(f"Response code: {response.status_code}, message: {response.text}")
        time.sleep(random.uniform(0.5, 2))

NUM_SIMULATED_CLIENTS = 1
threads = []
for i in range(NUM_SIMULATED_CLIENTS):
    t = threading.Thread(target=simulate_client, args=(i, emojis_flat))
    threads.append(t)
    t.start()

for t in threads:
    t.join()

