import constants
import requests
from Diploma import Diploma

HEADERS = {
    "Authorization": f"Bearer {constants.BEARER_TOKEN}",
    "Content-Type": "application/json"
}


def post_diploma(diploma: Diploma):
    url = f'http://{constants.API_HOST}:{constants.API_PORT}/api/v1/service/diploma'
    payload = {
        "user_id": diploma.user_id,
        "olympiad_id": diploma.olympiad_id,
        "class": diploma.diploma_class,
        "level": diploma.level
    }
    requests.post(url, json=payload, headers=HEADERS)