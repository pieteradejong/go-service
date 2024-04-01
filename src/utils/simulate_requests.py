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

# for _ in range(100):
#     emoji = random.choice(emojis_flat)
#     payload = {"message": emoji}
#     print(f'Selected emoji: {emoji}')
#     response = requests.post(full_url, headers=headers, data=json.dumps(payload))

for emoji in car_emojis:
    for _ in range(10):
        payload = {"message": emoji}
        print(f'Selected emoji: {emoji}')
        response = requests.post(full_url, headers=headers, data=json.dumps(payload))
