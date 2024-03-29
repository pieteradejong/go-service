import requests
import json


url = "http://yourapi.com/post"
headers = {'Content-Type': 'application/json'}
payload = {"emoji": "😀"}  # This string contains a UTF-8 encoded emoji

response = requests.post(url, headers=headers, data=json.dumps(payload))
print(response.status_code)

